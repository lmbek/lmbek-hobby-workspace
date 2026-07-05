package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"workspace/git-controller/internal/commands/doctor"
	"workspace/git-controller/internal/commands/down"
	initcmd "workspace/git-controller/internal/commands/init"
	"workspace/git-controller/internal/commands/sshsetup"
	"workspace/git-controller/internal/commands/sync"
	"workspace/git-controller/internal/commands/up"
	"workspace/git-controller/internal/commands/validate"
	"workspace/git-controller/internal/system"
	"workspace/git-controller/internal/ui"
)

const version = "1.1.0"

func main() {
	setupLogger()

	if len(os.Args) < 2 {
		showHelp()
		os.Exit(0)
	}

	command := os.Args[1]

	commands := map[string]func() error{
		"init":      initcmd.Run,
		"sync":      sync.Run,
		"validate":  validate.Run,
		"up":        up.Run,
		"down":      down.Run,
		"doctor":    doctor.RunFullError,
		"ssh-setup": sshsetup.RunError,
		"ssh":       sshsetup.RunError,
	}

	if command == "help" {
		showHelp()
		return
	}

	if command == "version" || command == "v" {
		showVersion()
		return
	}

	if run, ok := commands[command]; ok {
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
	handler := slog.NewTextHandler(os.Stderr, opts) // Logs to Stderr
	slog.SetDefault(slog.New(handler))
}

func showHelp() {
	fmt.Printf("\n%sWorkspace Controller v%s%s\n", ui.ColorBold, version, ui.ColorReset)
	fmt.Println(strings.Repeat("=", 25))
	fmt.Printf("Usage: %s [command]\n", ui.ColorBold+"<cli>"+ui.ColorReset)
	ui.Note(system.CLIDescription)

	fmt.Println("\nAvailable Commands:")
	fmt.Println("  init       [1] Bootstrap workspace (Pre-flight checks + Planning)")
	fmt.Println("  sync       [2] Synchronize all repositories (Materialization + fetch/pull/hooks)")
	fmt.Println("  validate   [3] Validate system consistency and health")
	fmt.Println("  up         [4] Start the system (docker-compose up)")
	fmt.Println("  down           Stop the system (docker-compose down)")
	fmt.Println("  doctor     [D] Diagnose environmental issues (Git, SSH, Docker)")
	fmt.Println("  ssh-setup  [S] Interactive SSH key management tool (alias: ssh)")
	fmt.Println("  ssh        [S] Alias for ssh-setup")
	fmt.Println("  version    [V] Show version information (alias: v)")
	fmt.Println("  help           Show this help information")
}

func showVersion() {
	fmt.Printf("Workspace Controller version %s\n", version)
}
