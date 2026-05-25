package system

type SystemDefinition struct {
	SystemVersion  string             `yaml:"system-version"`
	Services       map[string]Service `yaml:"services"`
	Infrastructure map[string]Infra   `yaml:"infrastructure"`
}

type Service struct {
	Repository  string            `yaml:"repository"`
	Version     string            `yaml:"version"`
	Environment map[string]string `yaml:"environment,omitempty"`
}

type Infra struct {
	Version     string            `yaml:"version"`
	Environment map[string]string `yaml:"environment,omitempty"`
}
