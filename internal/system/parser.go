package system

import (
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadDefinition(path string) (*SystemDefinition, error) {
	slog.Debug("Loading system definition", "path", path)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var system SystemDefinition
	err = yaml.Unmarshal(data, &system)
	if err != nil {
		return nil, err
	}

	return &system, nil
}
