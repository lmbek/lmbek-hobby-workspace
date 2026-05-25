package down

import (
	"log/slog"
	"os"
	"os/exec"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func Run() {
	slog.Info("Stopping workspace...")

	infraDir := getEnv("INFRA_DIR", "../workspace/infrastructure")
	if _, err := os.Stat(infraDir); os.IsNotExist(err) {
		slog.Error("Infrastructure directory not found", "path", infraDir)
		os.Exit(1)
	}

	slog.Info("Running docker-compose down", "path", infraDir)
	cmd := exec.Command("docker-compose", "down")
	cmd.Dir = infraDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		slog.Error("Error stopping docker-compose", "error", err)
		os.Exit(1)
	}

	slog.Info("Workspace has been stopped")
}
