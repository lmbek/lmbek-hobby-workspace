package system

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// CLIName is set by main() at startup to the actual binary name (e.g. "git-controller.exe").
var CLIName = "git-controller"

type Workspace struct {
	Root string
}

func (w *Workspace) GetCategoryDir(catName string) string {
	envKey := strings.ToUpper(catName) + "_DIR"
	if val, ok := os.LookupEnv(envKey); ok && val != "" {
		return val
	}

	root := w.Root
	if root == "" {
		root = "."
	}

	// Place all managed categories under git-repositories/<category>
	return filepath.Join(root, "git-repositories", catName)
}

type SystemDefinition struct {
	SystemVersion string              `yaml:"system-version"`
	Hooks         Hooks               `yaml:"hooks,omitempty"`
	Repos         map[string]Category `yaml:"repos"`
}

type Category map[string]Component

func (c *Category) UnmarshalYAML(value *yaml.Node) error {
	// 1. Try to unmarshal as a map of components (the traditional way)
	type mapType map[string]Component
	var m mapType
	if err := value.Decode(&m); err == nil {
		*c = Category(m)
		return nil
	}

	// 2. Try as a single Component (the flat/singleton way)
	var comp Component
	if err := value.Decode(&comp); err == nil {
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
	PostClone []string `yaml:"post-clone,omitempty"` // Commands to run after clone
}

type Component struct {
	Repository string `yaml:"repository"`
	Version    string `yaml:"version"`
}
