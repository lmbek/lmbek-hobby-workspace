package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"workspace/git-controller/internal/commands/checkout"
	"workspace/git-controller/internal/commands/clone"
	"workspace/git-controller/internal/commands/doctor"
	"workspace/git-controller/internal/commands/fetch"
	"workspace/git-controller/internal/commands/initrepoenvs"
	"workspace/git-controller/internal/commands/sshhelper"
	"workspace/git-controller/internal/commands/status"
	"workspace/git-controller/internal/commands/sync"
	"workspace/git-controller/internal/system"
	"workspace/git-controller/internal/ui"
)

func main() {
	setupLogger()
	system.CLIName = filepath.Base(os.Args[0])

	if len(os.Args) < 2 {
		showHelp()
		os.Exit(0)
	}

	command := os.Args[1]

	commands := map[string]func() error{
		"clone":          clone.Run,
		"fetch":          fetch.Run,
		"sync":           sync.Run,
		"checkout":       checkout.Run,
		"init-repo-envs": initrepoenvs.Run,
		"doctor":         doctor.Run,
		"ssh-helper":     sshhelper.Run,
		"status":         status.Run,

		// Aliases
		"envs": initrepoenvs.Run,
		"ssh":  sshhelper.Run,
	}

	if command == "help" || command == "--help" || command == "-h" {
		showHelp()
		return
	}

	if run, ok := commands[command]; ok && run != nil {
		if err := run(); err != nil {
			ui.Error("%v", err)
			os.Exit(1)
		}
	} else {
		ui.Error("Unknown command: %s", command)
		showHelp()
		os.Exit(1)
	}
}

func setupLogger() {
	level := slog.LevelInfo
	if os.Getenv("DEBUG") == "true" {
		level = slog.LevelDebug
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}
	handler := slog.NewTextHandler(os.Stderr, opts)
	slog.SetDefault(slog.New(handler))
}

func showHelp() {
	fmt.Printf("\n%sWorkspace Controller%s\n", ui.ColorBold, ui.ColorReset)
	fmt.Println(strings.Repeat("=", 25))
	fmt.Printf("Usage: %s [command]\n", ui.ColorBold+system.CLIName+ui.ColorReset)

	fmt.Println("\nCommands:")
	fmt.Println("  clone           Clone all repositories defined in repo-definition.yaml")
	fmt.Println("  fetch           Fetch all remotes across all repositories")
	fmt.Println("  sync            Safely sync all repositories via readonly pulls (fast-forward)")
	fmt.Println("  checkout        Switch all repositories to their defined branch")
	fmt.Println("  init-repo-envs  Initialise .env files from .env.example across repositories (alias: envs)")
	fmt.Println("  doctor          Diagnose environment (Git, Go, SSH, Docker)")
	fmt.Println("  ssh-helper      Interactive SSH key management and configuration (alias: ssh)")
	fmt.Println("  status          Show dashboard overview of all repository states")
	fmt.Println("  help            Show this help")
}
