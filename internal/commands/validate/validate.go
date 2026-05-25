package validate

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

	for name, svc := range sys.Services {
		fmt.Printf("\nService: %s\n", name)
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
	}

	fmt.Println("\nValidation summary:")
	if hasErrors {
		fmt.Println("Status: FAILED (Check the issues above)")
		os.Exit(1)
	} else {
		fmt.Println("Status: PASSED")
	}
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
