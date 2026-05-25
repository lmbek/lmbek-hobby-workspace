package validate

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"workspace-controller/internal/system"
)

func Run() {
	sys, err := system.LoadDefinition("system/system-definition.yaml")
	if err != nil {
		fmt.Printf("Error: Could not read system definition: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("VALIDATE SYSTEM")
	fmt.Println("===============")

	workspaceDir := "workspace"
	hasErrors := false

	fmt.Println("\nChecking Services:")
	for name, svc := range sys.Services {
		fmt.Printf("\n- Service: %s\n", name)
		targetPath := filepath.Join(workspaceDir, name)

		// 1. Check if directory exists
		if _, err := os.Stat(targetPath); os.IsNotExist(err) {
			fmt.Printf("  [MISSING] Directory does not exist: %s\n", targetPath)
			fmt.Printf("  [ACTION] Please run 'workspace-controller sync' to clone this service.\n")
			hasErrors = true
			continue
		}

		// 2. Check current version (git tag/branch)
		currentVersion, err := getGitCurrentVersion(targetPath)
		if err != nil {
			fmt.Printf("  [ERROR] Could not determine current version: %v\n", err)
			hasErrors = true
		} else {
			if currentVersion == svc.Version {
				fmt.Printf("  [OK] Version matches: %s\n", currentVersion)
			} else {
				fmt.Printf("  [MISMATCH] Version mismatch! Expected: %s, Found: %s\n", svc.Version, currentVersion)
				hasErrors = true
			}
		}

		// 3. Check for local changes
		isDirty, err := isGitDirty(targetPath)
		if err != nil {
			fmt.Printf("  [ERROR] Could not check git status: %v\n", err)
			hasErrors = true
		} else if isDirty {
			fmt.Printf("  [DIRTY] Local changes detected! Please commit or stash them.\n")
			hasErrors = true
		} else {
			fmt.Printf("  [CLEAN] No local changes detected.\n")
		}

		// 4. Health Check
		if svc.HealthCheck != "" {
			if err := performHealthCheck(svc.HealthCheck); err != nil {
				fmt.Printf("  [UNHEALTHY] %v\n", err)
				// We don't set hasErrors = true here because it might just not be running yet
				fmt.Printf("  [HINT] Is the service running? Use 'workspace-controller up' to start it.\n")
			} else {
				fmt.Printf("  [HEALTHY] Service is reachable and responding correctly.\n")
			}
		}
	}

	fmt.Println("\nChecking Infrastructure:")
	for name, infra := range sys.Infrastructure {
		fmt.Printf("\n- Component: %s\n", name)
		if infra.HealthCheck != "" {
			if err := performHealthCheck(infra.HealthCheck); err != nil {
				fmt.Printf("  [UNHEALTHY] %v\n", err)
				fmt.Printf("  [HINT] Is the infrastructure running? Use 'workspace-controller up' to start it.\n")
			} else {
				fmt.Printf("  [HEALTHY] Component is reachable and responding correctly.\n")
			}
		} else {
			fmt.Printf("  [SKIP] No health check defined.\n")
		}
	}

	fmt.Println("\nValidation summary:")
	if hasErrors {
		fmt.Println("Status: FAILED (Check the issues above)")
		os.Exit(1)
	} else {
		fmt.Println("Status: PASSED (Static checks okay)")
	}
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
