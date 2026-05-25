package sync

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
		fmt.Printf("Error: Could not read system definition: %v\n", err)
		os.Exit(1)
	}

	// 1. Create directories using environment variables or defaults
	workspaceDir := getEnv("SERVICES_DIR", "../workspace/services")
	infraDir := getEnv("INFRA_DIR", "../workspace/infrastructure")
	toolsDir := getEnv("TOOLS_DIR", "../workspace/tools")

	ensureDir(workspaceDir)
	ensureDir(infraDir)

	fmt.Println("SYNC SYSTEM")
	fmt.Println("===========")

	// 2. Process services
	for name, svc := range sys.Services {
		processGitComponent(workspaceDir, name, svc.Repository, svc.Version)
	}

	// 3. Process infrastructure
	if sys.Infrastructure != nil {
		processGitComponent(infraDir, "infrastructure", sys.Infrastructure.Repository, sys.Infrastructure.Version)
	}

	// 4. Process tools
	if sys.Tools != nil {
		processGitComponent(toolsDir, "tools", sys.Tools.Repository, sys.Tools.Version)
	}

	// 5. Post-Sync Hooks
	if len(sys.Hooks.PostSync) > 0 {
		fmt.Println("\nRunning Post-Sync Hooks:")
		for _, hook := range sys.Hooks.PostSync {
			fmt.Printf("  [HOOK] %s\n", hook)
			if err := runHook(hook); err != nil {
				fmt.Printf("  [ERROR] Hook failed: %v\n", err)
			}
		}
	}

	fmt.Println("\nSync finished.")
}

func runHook(command string) error {
	// Simple shell execution
	var cmd *exec.Cmd
	if strings.Contains(os.Getenv("OS"), "Windows") {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func processGitComponent(baseDir, name, repo, version string) {
	targetPath := filepath.Join(baseDir, name)
	// Special case: if name matches the end of baseDir, don't join
	if filepath.Base(baseDir) == name {
		targetPath = baseDir
	}

	fmt.Printf("\nComponent: %s\n", name)

	// Check for placeholder URL
	if strings.Contains(repo, "@company") || repo == "" {
		fmt.Printf("  [SKIP] Skipping '%s' - Repository URL is a placeholder (%s)\n", name, repo)
		return
	}

	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		// Clone if it doesn't exist
		fmt.Printf("  [CLONE] Fetching %s...\n", repo)
		if err := runGit("clone", repo, targetPath); err != nil {
			fmt.Printf("  [ERROR] Error cloning: %v\n", err)
			return
		}
	} else {
		fmt.Printf("  [EXISTS] Directory already exists.\n")
	}

	// Checkout version
	fmt.Printf("  [CHECKOUT] Setting version to %s...\n", version)
	if err := runGitInDir(targetPath, "checkout", version); err != nil {
		fmt.Printf("  [ERROR] Error checking out %s: %v\n", version, err)
	}
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
