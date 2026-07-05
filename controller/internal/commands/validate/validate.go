package validate

import (
	"controller/internal/system"
	"controller/internal/ui"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func Run() error {
	ui.Header("System Validation")

	sys, err := system.LoadDefinition("repos.yaml")
	if err != nil {
		return fmt.Errorf("could not read system definition: %w", err)
	}

	ui.Info("Validating components...")

	hasErrors := false

	categories := map[string]map[string]system.Component{
		"proxy":          sys.Proxy,
		"applications":   sys.Applications,
		"infrastructure": sys.Infrastructure,
		"orchestrator":   sys.Orchestrator,
		"platform":       sys.Platform,
		"tools":          sys.Tools,
		"docs":           sys.Docs,
	}

	for catName, components := range categories {
		catDir := system.GetCategoryDir(catName)

		for name, comp := range components {
			slog.Debug("Checking component", "category", catName, "name", name)
			targetPath := filepath.Join(catDir, name)

			displayName := name
			if displayName == "" {
				displayName = catName
			}

			if comp.Repository != "" {
				if err := performStaticGitChecks(targetPath, comp.Version, catName); err != nil {
					ui.Error("[%s] Git validation failed: %v", displayName, err)
					hasErrors = true
				} else {
					ui.Success("[%s] Git OK", displayName)
				}
			} else {
				ui.Success("[%s] Local component OK", displayName)
			}

			// Health Check
			if comp.HealthCheck != "" {
				if err := performHealthCheck(comp.HealthCheck); err != nil {
					ui.Warn("[%s] Health check failed: %v", displayName, err)
					ui.Info("  Hint: Is the component running? Use '<cli> up' to start it.")
				} else {
					ui.Success("[%s] Healthy", displayName)
				}
			}
		}
	}

	if hasErrors {
		return fmt.Errorf("validation status: FAILED")
	}

	ui.Success("Validation status: PASSED")
	return nil
}

func performStaticGitChecks(targetPath string, expectedVersion string, category string) error {
	// 1. Check if directory exists
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s. Please run '<cli> init'", targetPath)
	}

	// 2. Check current version (git tag/branch)
	currentVersion, err := getGitCurrentVersion(targetPath)
	if err != nil {
		return fmt.Errorf("could not determine current version: %v", err)
	}

	if currentVersion != expectedVersion {
		slog.Warn("Version mismatch", "path", targetPath, "expected", expectedVersion, "found", currentVersion)
		slog.Info("Hint: The controller is configured to stay on the current branch (main).", "path", targetPath)
		return nil
	}
	slog.Info("Version matches", "path", targetPath, "version", currentVersion)

	// 3. Check for local changes
	isDirty, err := isGitDirty(targetPath)
	if err != nil {
		return fmt.Errorf("could not check git status: %v", err)
	}
	if isDirty {
		// Only warn about local changes in infrastructure and orchestrator to avoid blocking during development
		if category == "infrastructure" || category == "orchestrator" {
			slog.Warn("Local changes detected in "+category, "path", targetPath)
			return nil
		}
		return fmt.Errorf("local changes detected in %s! Please commit or stash them", targetPath)
	}
	slog.Info("No local changes detected", "path", targetPath)

	return nil
}

func performHealthCheck(url string) error {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		client := http.Client{
			Timeout: 2 * time.Second,
		}
		resp, err := client.Get(url)
		if err != nil {
			return fmt.Errorf("HTTP check failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return fmt.Errorf("HTTP check returned status code: %d", resp.StatusCode)
		}
		return nil
	}

	if strings.HasPrefix(url, "tcp://") {
		address := strings.TrimPrefix(url, "tcp://")
		conn, err := net.DialTimeout("tcp", address, 2*time.Second)
		if err != nil {
			return fmt.Errorf("TCP check failed: %v", err)
		}
		conn.Close()
		return nil
	}

	return fmt.Errorf("unsupported health check protocol: %s", url)
}

func getGitCurrentVersion(dir string) (string, error) {
	// Try to get exact tag match first
	cmd := exec.Command("git", "describe", "--tags", "--exact-match")
	cmd.Dir = dir
	output, err := cmd.Output()
	if err == nil {
		return strings.TrimSpace(string(output)), nil
	}

	// If no tag, get branch name or hash
	cmd = exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = dir
	output, err = cmd.Output()
	if err != nil {
		return "", err
	}

	branch := strings.TrimSpace(string(output))
	if branch == "HEAD" {
		// Detached HEAD, get hash
		cmd = exec.Command("git", "rev-parse", "--short", "HEAD")
		cmd.Dir = dir
		output, err = cmd.Output()
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(output)), nil
	}

	return branch, nil
}

func isGitDirty(dir string) (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return len(strings.TrimSpace(string(output))) > 0, nil
}
