package main

import (
	"fmt"
	"log/slog"
	"os"
	"workspace-controller/internal/commands/doctor"
	"workspace-controller/internal/commands/down"
	initcmd "workspace-controller/internal/commands/init"
	"workspace-controller/internal/commands/sshsetup"
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
		"validate":  validate.Run,
		"up":        up.Run,
		"down":      down.Run,
		"doctor":    doctor.RunFull,
		"ssh-setup": sshsetup.Run,
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
	fmt.Println("  init       Initialize workspace (Pre-flight checks + Planning + Materialization)")
	fmt.Println("  validate   Validate system consistency and health")
	fmt.Println("  up         Start the system (docker-compose up)")
	fmt.Println("  down       Stop the system (docker-compose down)")
	fmt.Println("  doctor     Diagnose environmental issues (Git, SSH, Docker)")
	fmt.Println("  ssh-setup  Interactive SSH key management tool")
	fmt.Println("  help       Show this help information")
}
