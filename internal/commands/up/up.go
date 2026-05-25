package up

import (
	"log/slog"
	"os"
	"os/exec"
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

	slog.Info("Starting workspace...")

	infraDir := getEnv("INFRA_DIR", "../workspace/infrastructure")
	if _, err := os.Stat(infraDir); os.IsNotExist(err) {
		slog.Error("Infrastructure directory not found. Please run 'sync' first.", "path", infraDir)
		os.Exit(1)
	}

	slog.Info("Running docker-compose up -d", "path", infraDir)
	// Note: Each component in infrastructure may have its own docker-compose or deployment logic.
	// Currently, we assume a central docker-compose.yaml exists in the infrastructure directory.
	cmd := exec.Command("docker-compose", "up", "-d")
	cmd.Dir = infraDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		slog.Error("Error starting docker-compose", "error", err)
		os.Exit(1)
	}

	// Post-Up Hooks
	if len(sys.Hooks.PostUp) > 0 {
		slog.Info("Running Post-Up Hooks...")
		for _, hook := range sys.Hooks.PostUp {
			slog.Info("Executing hook", "command", hook)
			if err := runHook(hook); err != nil {
				slog.Error("Hook failed", "command", hook, "error", err)
			}
		}
	}

	slog.Info("Workspace is up and running")
}

func runHook(command string) error {
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
