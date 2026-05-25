/*
SYSTEM CONTROLLER - START

Purpose:
This module reads the system definition and produces a deterministic execution plan
for the local development environment.

It does NOT modify the system.
It does NOT clone repositories.
It does NOT start infrastructure.

It is a read-only planning step that answers:
"What would be started if the system was executed?"

Execution flow:
1. Load system-definition.yaml
2. Parse services and infrastructure
3. Validate structure (lightweight)
4. Print execution plan

Next step:
Run "system-controller sync" to materialize the system locally.
*/

package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type SystemDefinition struct {
	SystemVersion  string             `yaml:"system-version"`
	Services       map[string]Service `yaml:"services"`
	Infrastructure map[string]string  `yaml:"infrastructure"`
}

type Service struct {
	Repository string `yaml:"repository"`
	Version    string `yaml:"version"`
}

func main() {
	data, err := os.ReadFile("../system-definition.yaml")
	if err != nil {
		log.Fatal(err)
	}

	var system SystemDefinition
	err = yaml.Unmarshal(data, &system)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("SYSTEM START PLAN")
	fmt.Println("=================")

	fmt.Println("\nServices:")
	for name, svc := range system.Services {
		fmt.Printf("- %s @ %s (%s)\n", name, svc.Version, svc.Repository)
	}

	fmt.Println("\nInfrastructure:")
	for name, version := range system.Infrastructure {
		fmt.Printf("- %s (version %s)\n", name, version)
	}

	fmt.Println("\nNEXT STEP:")
	fmt.Println("Run: system-controller sync")
}
