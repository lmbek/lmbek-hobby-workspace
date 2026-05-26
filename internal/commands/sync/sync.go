package sync

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"workspace-controller/internal/gitutil"
	"workspace-controller/internal/system"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func Run() {
	fmt.Println("SYSTEM SYNCHRONIZATION")
	fmt.Println("======================")

	sys, err := system.LoadDefinition("system/system-definition.yaml")
	if err != nil {
		slog.Error("Error loading system definition", "error", err)
		os.Exit(1)
	}

	workspaceDir := getEnv("SERVICES_DIR", "../workspace/services")
	infraDir := getEnv("INFRA_DIR", "../workspace/infrastructure")
	toolsDir := getEnv("TOOLS_DIR", "../workspace/tools")

	// Ensure absolute paths
	if abs, err := filepath.Abs(workspaceDir); err == nil {
		workspaceDir = abs
	}
	if abs, err := filepath.Abs(infraDir); err == nil {
		infraDir = abs
	}
	if abs, err := filepath.Abs(toolsDir); err == nil {
		toolsDir = abs
	}

	gitutil.EnsureDir(workspaceDir)
	gitutil.EnsureDir(infraDir)
	gitutil.EnsureDir(toolsDir)

	hasErrors := false

	// Sync Services
	for name, svc := range sys.Services {
		if err := gitutil.ProcessGitComponent(workspaceDir, name, svc.Repository, svc.Version); err != nil {
			gitutil.HandleGitError(err)
			hasErrors = true
		}
	}

	// Sync Infrastructure
	if sys.Infrastructure != nil {
		if err := gitutil.ProcessGitComponent(infraDir, "infrastructure", sys.Infrastructure.Repository, sys.Infrastructure.Version); err != nil {
			gitutil.HandleGitError(err)
			hasErrors = true
		}
	}

	// Sync Tools
	if sys.Tools != nil {
		if err := gitutil.ProcessGitComponent(toolsDir, "tools", sys.Tools.Repository, sys.Tools.Version); err != nil {
			gitutil.HandleGitError(err)
			hasErrors = true
		}
	}

	if hasErrors {
		fmt.Println("\nSynchronization failed with errors.")
		os.Exit(1)
	}

	// Running Hooks
	if len(sys.Hooks.PostSync) > 0 {
		fmt.Println("\n[3] Running Post-Sync Hooks:")
		for _, hook := range sys.Hooks.PostSync {
			if err := gitutil.RunHook(hook); err != nil {
				slog.Error("Hook failed", "command", hook, "error", err)
			}
		}
	}

	fmt.Println("\nSynchronization finished successfully!")
}
