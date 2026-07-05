package up

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"workspace/git-controller/internal/system"
	"workspace/git-controller/internal/ui"
)

func Run() error {
	ui.Header("Start System")

	sys, workspace, err := system.LoadDefinition("system-definition.yaml")
	if err != nil {
		return fmt.Errorf("could not read system definition: %w", err)
	}

	infraDir := sys.GetOrchestrationDir(workspace)
	if _, err := os.Stat(infraDir); os.IsNotExist(err) {
		return fmt.Errorf("orchestration directory not found at %s. Please run '<cli> init' first", infraDir)
	}

	ui.Step(4, "Starting Docker containers")
	ui.Info("Using orchestration directory: %s", infraDir)
	ui.Info("Running: docker-compose up -d")

	cmd := exec.Command("docker-compose", "up", "-d")
	cmd.Dir = infraDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error starting docker-compose: %w", err)
	}

	// Post-Up Hooks
	if len(sys.Hooks.PostUp) > 0 {
		ui.Step(5, "Running Post-Up Hooks")
		for _, hook := range sys.Hooks.PostUp {
			ui.Info("Running: %s", hook)
			if err := runHook(hook); err != nil {
				ui.Error("Hook failed: %s (%v)", hook, err)
			}
		}
	}

	ui.Success("Workspace is up and running")
	return nil
}

func runHook(command string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
