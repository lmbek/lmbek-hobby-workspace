package up

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"workspace-controller/internal/system"
)

func Run() {
	sys, err := system.LoadDefinition("system/system-definition.yaml")
	if err != nil {
		fmt.Printf("Error: Could not read system definition: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("STARTING WORKSPACE")
	fmt.Println("==================")

	infraDir := "../workspace/infrastructure"
	if _, err := os.Stat(infraDir); os.IsNotExist(err) {
		fmt.Println("Error: infrastructure directory not found. Please run 'start' first.")
		os.Exit(1)
	}

	fmt.Println("Running: docker-compose up -d --build")
	cmd := exec.Command("docker-compose", "up", "-d", "--build")
	cmd.Dir = infraDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error starting docker-compose: %v\n", err)
		os.Exit(1)
	}

	// Post-Up Hooks
	if len(sys.Hooks.PostUp) > 0 {
		fmt.Println("\nRunning Post-Up Hooks:")
		for _, hook := range sys.Hooks.PostUp {
			fmt.Printf("  [HOOK] %s\n", hook)
			if err := runHook(hook); err != nil {
				fmt.Printf("  [ERROR] Hook failed: %v\n", err)
			}
		}
	}

	fmt.Println("\nWorkspace is up and running.")
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
