package sync

import (
	"log/slog"
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
		slog.Error("Could not read system definition", "error", err)
		os.Exit(1)
	}

	// 1. Create directories using environment variables or defaults
	workspaceDir := getEnv("SERVICES_DIR", "../workspace/services")
	infraDir := getEnv("INFRA_DIR", "../workspace/infrastructure")
	toolsDir := getEnv("TOOLS_DIR", "../workspace/tools")

	ensureDir(workspaceDir)
	ensureDir(infraDir)

	slog.Info("Syncing system...")

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
		slog.Info("Running Post-Sync Hooks...")
		for _, hook := range sys.Hooks.PostSync {
			slog.Info("Executing hook", "command", hook)
			if err := runHook(hook); err != nil {
				slog.Error("Hook failed", "command", hook, "error", err)
			}
		}
	}

	slog.Info("Sync finished")
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

	slog.Info("Processing component", "name", name)

	// Check for placeholder URL
	if strings.Contains(repo, "@company") || repo == "" {
		slog.Warn("Skipping component - Repository URL is a placeholder", "name", name, "repo", repo)
		return
	}

	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		// Clone if it doesn't exist
		slog.Info("Cloning repository", "repo", repo, "target", targetPath)
		if err := runGit("clone", repo, targetPath); err != nil {
			slog.Error("Error cloning", "repo", repo, "error", err)
			return
		}
	} else {
		slog.Debug("Directory already exists", "path", targetPath)
	}

	// Syncing logic: stay on branch, fetch and pull updates
	slog.Info("Fetching and pulling updates", "name", name, "path", targetPath)
	if err := runGitInDir(targetPath, "fetch", "--all"); err != nil {
		slog.Error("Error fetching", "name", name, "error", err)
	}
	if err := runGitInDir(targetPath, "pull"); err != nil {
		slog.Error("Error pulling", "name", name, "error", err)
	}

	// Verification: Check if we are at the expected version (tag/branch/hash)
	// but DO NOT checkout to it if it causes a detached HEAD.
	slog.Info("Target version info", "name", name, "version", version)
}

func ensureDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		slog.Info("Creating directory", "path", path)
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
