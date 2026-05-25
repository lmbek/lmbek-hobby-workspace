package system

type SystemDefinition struct {
	SystemVersion  string             `yaml:"system-version"`
	Hooks          Hooks              `yaml:"hooks,omitempty"` // Added hooks
	Services       map[string]Service `yaml:"services"`
	Infrastructure map[string]Infra   `yaml:"infrastructure"`
	Tools          map[string]Tool    `yaml:"tools,omitempty"` // Added tools
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
	Repository  string            `yaml:"repository,omitempty"` // Added repository field
	Version     string            `yaml:"version"`
	Environment map[string]string `yaml:"environment,omitempty"`
	HealthCheck string            `yaml:"health-check,omitempty"` // Added health-check field
}

type Tool struct {
	Repository string `yaml:"repository"`
	Version    string `yaml:"version"`
}
