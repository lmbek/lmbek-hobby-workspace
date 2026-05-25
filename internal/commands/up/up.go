package up

import (
	"fmt"
	"os"
	"os/exec"
)

func Run() {
	fmt.Println("STARTING WORKSPACE")
	fmt.Println("==================")

	infraDir := "infrastructure"
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

	fmt.Println("\nWorkspace is up and running.")
}
