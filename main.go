package main

import (
	"fmt"
	"log/slog"
	"os"
	"workspace-controller/internal/commands/down"
	"workspace-controller/internal/commands/start"
	"workspace-controller/internal/commands/sync"
	"workspace-controller/internal/commands/up"
	"workspace-controller/internal/commands/validate"
)

func main() {
	// Initialize structured logger
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)

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
		slog.Error("Unknown command", "command", command)
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
