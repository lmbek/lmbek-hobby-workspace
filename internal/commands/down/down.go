package down

import (
	"fmt"
	"os"
	"os/exec"
)

func Run() {
	fmt.Println("STOPPING WORKSPACE")
	fmt.Println("==================")

	infraDir := "../workspace/infrastructure"
	if _, err := os.Stat(infraDir); os.IsNotExist(err) {
		fmt.Println("Error: infrastructure directory not found.")
		os.Exit(1)
	}

	fmt.Println("Running: docker-compose down")
	cmd := exec.Command("docker-compose", "down")
	cmd.Dir = infraDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error stopping docker-compose: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nWorkspace has been stopped.")
}
