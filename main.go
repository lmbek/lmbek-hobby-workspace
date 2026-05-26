package main

import (
	"fmt"
	"log/slog"
	"os"
	"workspace-controller/internal/commands/doctor"
	"workspace-controller/internal/commands/down"
	initcmd "workspace-controller/internal/commands/init"
	"workspace-controller/internal/commands/sshsetup"
	"workspace-controller/internal/commands/sync"
	"workspace-controller/internal/commands/up"
	"workspace-controller/internal/commands/validate"
	"workspace-controller/internal/system"
)

func main() {
	setupLogger()

	if len(os.Args) < 2 {
		showHelp()
		os.Exit(1)
	}

	command := os.Args[1]

	commands := map[string]func(){
		"init":      initcmd.Run,
		"sync":      sync.Run,
		"validate":  validate.Run,
		"up":        up.Run,
		"down":      down.Run,
		"doctor":    doctor.RunFull,
		"ssh-setup": sshsetup.Run,
		"ssh":       sshsetup.Run,
		"help":      showHelp,
	}

	if run, ok := commands[command]; ok {
		run()
	} else {
		slog.Error("Unknown command", "command", command)
		showHelp()
		os.Exit(1)
	}
}

func setupLogger() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := slog.NewTextHandler(os.Stdout, opts)
	slog.SetDefault(slog.New(handler))
}

func showHelp() {
	fmt.Println("\nWorkspace Controller")
	fmt.Println("====================")
	fmt.Println("Usage: <cli> [command]")
	system.PrintCLINote()
	fmt.Println("\nAvailable Commands:")
	fmt.Println("  init       [1] Bootstrap workspace (Pre-flight checks + Planning)")
	fmt.Println("  sync       [2] Synchronize all repositories (Materialization + fetch/pull/hooks)")
	fmt.Println("  validate   [3] Validate system consistency and health")
	fmt.Println("  up         [4] Start the system (docker-compose up)")
	fmt.Println("  down           Stop the system (docker-compose down)")
	fmt.Println("  doctor     [D] Diagnose environmental issues (Git, SSH, Docker)")
	fmt.Println("  ssh-setup  [S] Interactive SSH key management tool (alias: ssh)")
	fmt.Println("  ssh        [S] Alias for ssh-setup")
	fmt.Println("  help           Show this help information")
}
