package start

import (
	"fmt"
	"os"
	"path/filepath"
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
	for name, infra := range sys.Infrastructure {
		fmt.Printf("- %s (version %s)\n", name, infra.Version)
	}

	// Phase 4: Generate docker-compose.yaml
	err = generateDockerCompose(sys)
	if err != nil {
		fmt.Printf("\n[ERROR] Failed to generate docker-compose.yaml: %v\n", err)
	} else {
		fmt.Println("\n[OK] docker-compose.yaml generated in infrastructure/")
	}

	fmt.Println("\nNEXT STEPS:")
	fmt.Println("1. Run: workspace-controller sync")
	fmt.Println("2. Run: workspace-controller up")
}

func generateDockerCompose(sys *system.SystemDefinition) error {
	infraDir := "infrastructure"
	// Ensure infra dir exists
	if _, err := os.Stat(infraDir); os.IsNotExist(err) {
		err := os.MkdirAll(infraDir, 0755)
		if err != nil {
			return err
		}
	}

	filePath := filepath.Join(infraDir, "docker-compose.yaml")
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString("version: '3.8'\n\nservices:\n")
	if err != nil {
		return err
	}

	// 1. Add services from workspace
	for name, svc := range sys.Services {
		_, err = f.WriteString(fmt.Sprintf("  %s:\n", name))
		if err != nil {
			return err
		}

		// Point to the service directory in workspace.
		// Docker Compose will look for a Dockerfile there.
		_, err = f.WriteString(fmt.Sprintf("    build: ../workspace/%s\n", name))
		if err != nil {
			return err
		}

		if len(svc.Environment) > 0 {
			_, err = f.WriteString("    environment:\n")
			for k, v := range svc.Environment {
				_, err = f.WriteString(fmt.Sprintf("      %s: %s\n", k, v))
			}
		}

		if len(svc.DependsOn) > 0 {
			_, err = f.WriteString("    depends_on:\n")
			for _, dep := range svc.DependsOn {
				_, err = f.WriteString(fmt.Sprintf("      - %s\n", dep))
			}
		}

		_, err = f.WriteString("    networks:\n      - workspace-net\n")
		if err != nil {
			return err
		}
	}

	// 2. Add infrastructure
	for name, infra := range sys.Infrastructure {
		_, err = f.WriteString(fmt.Sprintf("  %s:\n", name))
		if err != nil {
			return err
		}

		switch name {
		case "postgres":
			_, err = f.WriteString(fmt.Sprintf("    image: postgres:%s\n", infra.Version))
			if err != nil {
				return err
			}

			// Add environment variables if present
			if len(infra.Environment) > 0 {
				_, err = f.WriteString("    environment:\n")
				for k, v := range infra.Environment {
					_, err = f.WriteString(fmt.Sprintf("      %s: %s\n", k, v))
				}
			} else {
				// Default if none provided
				_, err = f.WriteString("    environment:\n      POSTGRES_PASSWORD: password\n")
			}

			if err != nil {
				return err
			}

			_, err = f.WriteString("    ports:\n      - \"5432:5432\"\n    networks:\n      - workspace-net\n")

		default:
			// Generic fallback for other infra components
			_, err = f.WriteString(fmt.Sprintf("    image: %s:%s\n", name, infra.Version))
			if err != nil {
				return err
			}

			if len(infra.Environment) > 0 {
				_, err = f.WriteString("    environment:\n")
				for k, v := range infra.Environment {
					_, err = f.WriteString(fmt.Sprintf("      %s: %s\n", k, v))
				}
			}

			_, err = f.WriteString("    networks:\n      - workspace-net\n")
		}
		if err != nil {
			return err
		}
	}

	_, err = f.WriteString("\nnetworks:\n  workspace-net:\n    driver: bridge\n")
	return err
}
