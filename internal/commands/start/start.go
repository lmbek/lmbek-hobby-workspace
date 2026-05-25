package start

import (
	"fmt"
	"os"
	"workspace-controller/internal/system"
)

func Run() {
	sys, err := system.LoadDefinition("system/system-definition.yaml")
	if err != nil {
		fmt.Printf("Error loading system definition: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("SYSTEM START PLAN")
	fmt.Println("=================")

	fmt.Println("\nServices:")
	for name, svc := range sys.Services {
		fmt.Printf("- %s @ %s (%s)\n", name, svc.Version, svc.Repository)
	}

	fmt.Println("\nInfrastructure:")
	if sys.Infrastructure != nil {
		fmt.Printf("- infrastructure @ %s (%s)\n", sys.Infrastructure.Version, sys.Infrastructure.Repository)
	} else {
		fmt.Println("  [SKIP] No infrastructure defined.")
	}

	fmt.Println("\n[OK] Infrastructure plan is managed within the infrastructure repository.")

	fmt.Println("\nNEXT STEPS:")
	fmt.Println("1. Run: go run main.go sync")
	fmt.Println("2. Run: go run main.go up")
}
