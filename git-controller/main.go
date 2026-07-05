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
	"workspace/git-controller/internal/commands/pull"
	"workspace/git-controller/internal/commands/push"
	"workspace/git-controller/internal/commands/sshsetup"
	"workspace/git-controller/internal/commands/status"
	"workspace/git-controller/internal/commands/validate"
	"workspace/git-controller/internal/commands/wsinit"
	"workspace/git-controller/internal/system"
	"workspace/git-controller/internal/ui"
)

const version = "5.0.0"

func main() {
	setupLogger()
	system.CLIName = filepath.Base(os.Args[0])

	if len(os.Args) < 2 {
		showHelp()
		os.Exit(0)
	}

	command := os.Args[1]

	commands := map[string]func() error{
		"init":     wsinit.Run,
		"clone":    clone.Run,
		"pull":     pull.Run,
		"push":     push.Run,
		"checkout": checkout.Run,
		"status":   status.Run,
		"validate": validate.Run,
		// Utilities
		"doctor":    doctor.Run,
		"ssh-setup": sshsetup.Run,
		"ssh":       sshsetup.Run,
	}

	if command == "help" {
		showHelp()
		return
	}

	if command == "version" || command == "v" {
		showVersion()
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
	fmt.Printf("\n%sWorkspace Controller v%s%s\n", ui.ColorBold, version, ui.ColorReset)
	fmt.Println(strings.Repeat("=", 25))
	fmt.Printf("Usage: %s [command]\n", ui.ColorBold+system.CLIName+ui.ColorReset)

	fmt.Println("\nWorkflow Commands:")
	fmt.Println("  init       Scaffold a new workspace (system-definition.yaml, Makefile, .gitignore)")
	fmt.Println("  clone      Clone all repositories defined in system-definition.yaml")
	fmt.Println("  pull       Pull latest changes across all repositories (clone if missing)")
	fmt.Println("  push       Push local commits across all repositories")
	fmt.Println("  checkout   Switch all repositories to their defined branch")
	fmt.Println("  status     Show dashboard overview of all repository states")
	fmt.Println("  validate   Validate repository consistency against the definition")

	fmt.Println("\nSetup Commands:")
	fmt.Println("  doctor     Diagnose environment (Git, Go, SSH, Docker)")
	fmt.Println("  ssh-setup  Interactive SSH key management (alias: ssh)")
	fmt.Println("  version    Show version (alias: v)")
	fmt.Println("  help       Show this help")
}

func showVersion() {
	fmt.Printf("Workspace Controller version %s\n", version)
}
