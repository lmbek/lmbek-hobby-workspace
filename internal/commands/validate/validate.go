package validate

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"workspace-controller/internal/system"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func Run() {
	sys, err := system.LoadDefinition("system/system-definition.yaml")
	if err != nil {
		slog.Error("Could not read system definition", "error", err)
		os.Exit(1)
	}

	slog.Info("Validating system...")

	workspaceDir := getEnv("SERVICES_DIR", "../workspace/services")
	hasErrors := false

	for name, svc := range sys.Services {
		slog.Info("Checking service", "name", name)
		targetPath := filepath.Join(workspaceDir, name)

		if err := performStaticGitChecks(targetPath, svc.Version); err != nil {
			slog.Error("Git validation failed", "service", name, "error", err)
			hasErrors = true
		}

		// Health Check
		if svc.HealthCheck != "" {
			if err := performHealthCheck(svc.HealthCheck); err != nil {
				slog.Warn("Health check failed", "service", name, "url", svc.HealthCheck, "error", err)
				fmt.Printf("\nHint: Is the service running? Use '<cli> up' to start it.\n")
				system.PrintCLINote()
			} else {
				slog.Info("Service is healthy", "service", name)
			}
		}
	}

	infraDir := getEnv("INFRA_DIR", "../workspace/infrastructure")
	if sys.Infrastructure != nil {
		slog.Info("Checking infrastructure component")
		targetPath := infraDir

		// Infrastructure now only has Git versioning validation
		if err := performStaticGitChecks(targetPath, sys.Infrastructure.Version); err != nil {
			slog.Error("Git validation failed", "component", "infrastructure", "error", err)
			hasErrors = true
		} else {
			slog.Info("Infrastructure repository is consistent")
		}
	} else {
		slog.Info("No infrastructure defined, skipping")
	}

	if sys.Tools != nil {
		slog.Info("Checking tools component")
		targetPath := getEnv("TOOLS_DIR", "../workspace/tools")
		if err := performStaticGitChecks(targetPath, sys.Tools.Version); err != nil {
			slog.Error("Git validation failed", "component", "tools", "error", err)
			hasErrors = true
		} else {
			slog.Info("Tools repository is consistent")
		}
	} else {
		slog.Info("No tools defined, skipping")
	}

	if hasErrors {
		slog.Error("Validation status: FAILED")
		os.Exit(1)
	} else {
		slog.Info("Validation status: PASSED")
	}
}

func performStaticGitChecks(targetPath string, expectedVersion string) error {
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
		// Only warn about local changes in infrastructure to avoid blocking during development
		if filepath.Base(targetPath) == "infrastructure" {
			slog.Warn("Local changes detected in infrastructure", "path", targetPath)
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
