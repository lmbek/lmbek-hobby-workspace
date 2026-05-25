package sync

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

	// 1. Create directories
	workspaceDir := "workspace"
	infraDir := "infrastructure"

	ensureDir(workspaceDir)
	ensureDir(infraDir)

	fmt.Println("SYNC SYSTEM")
	fmt.Println("===========")

	// 2. Process services
	for name, svc := range sys.Services {
		targetPath := filepath.Join(workspaceDir, name)

		fmt.Printf("\nService: %s\n", name)

		// Check for placeholder URL
		if strings.Contains(svc.Repository, "@company") || svc.Repository == "" {
			fmt.Printf("  [SKIP] Skipping '%s' - Repository URL is a placeholder (%s)\n", name, svc.Repository)
			continue
		}

		if _, err := os.Stat(targetPath); os.IsNotExist(err) {
			// Clone if it doesn't exist
			fmt.Printf("  [CLONE] Fetching %s...\n", svc.Repository)
			if err := runGit("clone", svc.Repository, targetPath); err != nil {
				fmt.Printf("  [ERROR] Error cloning: %v\n", err)
				continue
			}
		} else {
			fmt.Printf("  [EXISTS] Directory already exists.\n")
		}

		// Checkout version
		fmt.Printf("  [CHECKOUT] Setting version to %s...\n", svc.Version)
		if err := runGitInDir(targetPath, "checkout", svc.Version); err != nil {
			fmt.Printf("  [ERROR] Error checking out %s: %v\n", svc.Version, err)
		}
	}

	fmt.Println("\nSync finished.")
}

func ensureDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("Creating directory: %s\n", path)
		os.MkdirAll(path, 0755)
	}
}

func runGit(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runGitInDir(dir string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
