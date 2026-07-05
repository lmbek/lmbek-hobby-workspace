package system

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const CLIDescription = "<cli> refers to the command you are using (e.g., ./workspace-controller.exe, go run main.go, or make)"

func PrintCLINote() {
	fmt.Printf("\nNote: %s\n", CLIDescription)
}

func GetCategoryDir(catName string) string {
	envKey := strings.ToUpper(catName) + "_DIR"
	switch catName {
	case "applications":
		envKey = "SERVICES_DIR"
	case "infrastructure":
		envKey = "INFRA_DIR"
	case "orchestrator":
		envKey = "ORCHESTRATOR_DIR"
	}

	if val, ok := os.LookupEnv(envKey); ok {
		return val
	}

	root := os.Getenv("WORKSPACE_ROOT")
	if root == "" {
		root = "."
	}

	return filepath.Join(root, catName)
}

func (sys *SystemDefinition) GetOrchestrationDir() string {
	if len(sys.Orchestrator) > 0 {
		for name := range sys.Orchestrator {
			return filepath.Join(GetCategoryDir("orchestrator"), name)
		}
	}
	return GetCategoryDir("infrastructure")
}

type SystemDefinition struct {
	SystemVersion  string   `yaml:"system-version"`
	Hooks          Hooks    `yaml:"hooks,omitempty"` // Added hooks
	Proxy          Category `yaml:"proxy,omitempty"` // Added proxy category
	Applications   Category `yaml:"applications"`
	Infrastructure Category `yaml:"infrastructure"`
	Orchestrator   Category `yaml:"orchestrator"`
	Platform       Category `yaml:"platform"`
	Tools          Category `yaml:"tools"`
	Docs           Category `yaml:"docs"`
}

type Category map[string]Component

func (c *Category) UnmarshalYAML(value *yaml.Node) error {
	// 1. Try to unmarshal as a map of components (the traditional way)
	type mapType map[string]Component
	var m mapType
	if err := value.Decode(&m); err == nil {
		// Heuristic: if the map successfully decoded but it contains "repository",
		// it might be a singleton Component that was mistakenly partially decoded as a map.
		// However, a struct Decode is more restrictive than map Decode.
		// In yaml.v3, Decode into map[string]Component will fail if values aren't maps/structs.
		*c = Category(m)
		return nil
	}

	// 2. Try as a single Component (the flat/singleton way)
	var comp Component
	if err := value.Decode(&comp); err == nil {
		// If it has at least one of the known component fields, it's a singleton.
		// We check Repository or Version as indicators.
		if comp.Repository != "" || comp.Version != "" {
			*c = Category{"": comp}
			return nil
		}
	}

	// 3. If it's an empty map or null, just return an empty category
	if value.Tag == "!!null" || (value.Kind == yaml.MappingNode && len(value.Content) == 0) {
		*c = make(Category)
		return nil
	}

	return fmt.Errorf("failed to unmarshal Category at line %d", value.Line)
}

type Hooks struct {
	PostSync []string `yaml:"post-sync,omitempty"` // Commands to run after sync
	PostUp   []string `yaml:"post-up,omitempty"`   // Commands to run after up
}

type Component struct {
	Repository  string            `yaml:"repository"`
	Version     string            `yaml:"version"`
	Environment map[string]string `yaml:"environment,omitempty"`
	HealthCheck string            `yaml:"health-check,omitempty"`
	DependsOn   []string          `yaml:"depends-on,omitempty"` // Added depends-on field
}
