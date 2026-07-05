package down

import (
	"controller/internal/system"
	"controller/internal/ui"
	"fmt"
	"os"
	"os/exec"
)

func Run() error {
	ui.Header("Stop System")

	sys, err := system.LoadDefinition("repos.yaml")
	if err != nil {
		return fmt.Errorf("could not read system definition: %w", err)
	}

	infraDir := sys.GetOrchestrationDir()
	if _, err := os.Stat(infraDir); os.IsNotExist(err) {
		return fmt.Errorf("orchestration directory not found at %s", infraDir)
	}

	ui.Info("Using orchestration directory: %s", infraDir)
	ui.Info("Running: docker-compose down")
	cmd := exec.Command("docker-compose", "down")
	cmd.Dir = infraDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error stopping docker-compose: %w", err)
	}

	ui.Success("Workspace has been stopped")
	return nil
}
