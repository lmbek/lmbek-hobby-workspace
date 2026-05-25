package main

import (
	"fmt"
	"os"
	"workspace-controller/internal/commands/down"
	"workspace-controller/internal/commands/start"
	"workspace-controller/internal/commands/sync"
	"workspace-controller/internal/commands/up"
	"workspace-controller/internal/commands/validate"
)

func main() {
	if len(os.Args) < 2 {
		showHelp()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "start":
		start.Run()
	case "sync":
		sync.Run()
	case "validate":
		validate.Run()
	case "up":
		up.Run()
	case "down":
		down.Run()
	case "help":
		showHelp()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		showHelp()
		os.Exit(1)
	}
}

func showHelp() {
	fmt.Println("\nWorkspace Controller - Help")
	fmt.Println("==========================")
	fmt.Println("Usage:")
	fmt.Println("  workspace-controller <command>")
	fmt.Println("\nAvailable Commands:")
	fmt.Println("  start     Generate execution plan and infrastructure files")
	fmt.Println("  sync      Materialize the system (clone repos and setup workspace)")
	fmt.Println("  validate  Check system consistency and local state")
	fmt.Println("  up        Start the workspace (docker-compose up)")
	fmt.Println("  down      Stop the workspace (docker-compose down)")
	fmt.Println("  help      Show this help message")
}
