package start

import (
	"fmt"
	"log/slog"
	"os"
	"workspace-controller/internal/system"
)

func Run() {
	sys, err := system.LoadDefinition("system/system-definition.yaml")
	if err != nil {
		slog.Error("Error loading system definition", "error", err)
		os.Exit(1)
	}

	fmt.Println("SYSTEM START PLAN")
	fmt.Println("=================")

	fmt.Println("\nServices:")
	for name, svc := range sys.Services {
		slog.Info("Service planned", "name", name, "version", svc.Version, "repository", svc.Repository)
	}

	fmt.Println("\nInfrastructure:")
	if sys.Infrastructure != nil {
		slog.Info("Infrastructure planned", "version", sys.Infrastructure.Version, "repository", sys.Infrastructure.Repository)
	} else {
		slog.Warn("No infrastructure defined")
	}

	fmt.Println("\n[OK] Infrastructure plan is managed within the infrastructure repository.")

	fmt.Println("\nNEXT STEPS:")
	fmt.Println("1. Run: go run main.go sync")
	fmt.Println("2. Run: go run main.go up")
}
