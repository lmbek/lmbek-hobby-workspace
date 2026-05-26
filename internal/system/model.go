package system

import "fmt"

const CLIDescription = "<cli> refers to the command you are using (e.g., ./workspace-controller.exe, go run main.go, or make)"

func PrintCLINote() {
	fmt.Printf("\nNote: %s\n", CLIDescription)
}

type SystemDefinition struct {
	SystemVersion  string             `yaml:"system-version"`
	Hooks          Hooks              `yaml:"hooks,omitempty"` // Added hooks
	Services       map[string]Service `yaml:"services"`
	Infrastructure *Infra             `yaml:"infrastructure,omitempty"`
	Tools          *Tool              `yaml:"tools,omitempty"` // Updated tools to single object
}

type Hooks struct {
	PostSync []string `yaml:"post-sync,omitempty"` // Commands to run after sync
	PostUp   []string `yaml:"post-up,omitempty"`   // Commands to run after up
}

type Service struct {
	Repository  string            `yaml:"repository"`
	Version     string            `yaml:"version"`
	Environment map[string]string `yaml:"environment,omitempty"`
	HealthCheck string            `yaml:"health-check,omitempty"`
	DependsOn   []string          `yaml:"depends-on,omitempty"` // Added depends-on field
}

type Infra struct {
	Repository string `yaml:"repository"`
	Version    string `yaml:"version"`
}

type Tool struct {
	Repository string `yaml:"repository"`
	Version    string `yaml:"version"`
}
