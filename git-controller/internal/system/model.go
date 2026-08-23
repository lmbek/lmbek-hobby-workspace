package system

import (
	"fmt"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// CLIName is set by main() at startup to the actual binary name (e.g. "git-controller.exe").
var CLIName = "git-controller"

// DefaultDefinitionFile is the default path to the repository definition manifest.
const DefaultDefinitionFile = "git-repositories/repo-definition.yaml"

type Workspace struct {
	Root string
}

func (w *Workspace) GetCategoryDir(catName string) string {
	root := w.Root
	if root == "" {
		root = "."
	}

	// Place all managed categories under git-repositories/<category>
	return filepath.Join(root, "git-repositories", catName)
}

type SystemDefinition struct {
	SystemVersion string              `yaml:"system-version,omitempty"`
	Hooks         Hooks               `yaml:"hooks,omitempty"`
	Repos         map[string]Category `yaml:"repos"`
}

type Category map[string]Component

func (c *Category) UnmarshalYAML(value *yaml.Node) error {
	// 1. Try to unmarshal as a map of components (the traditional multi-repo category way)
	type mapType map[string]Component
	var m mapType
	if err := value.Decode(&m); err == nil && len(m) > 0 {
		*c = Category(m)
		return nil
	}

	// 2. Try as a single Component (the flat/singleton way or local entry)
	var comp Component
	if err := value.Decode(&comp); err == nil {
		*c = Category{"": comp}
		return nil
	}

	// 3. If it's an empty map, null, or empty scalar, return singleton with empty Component
	if value.Tag == "!!null" || (value.Kind == yaml.MappingNode && len(value.Content) == 0) || (value.Kind == yaml.ScalarNode && value.Value == "") {
		*c = Category{"": Component{}}
		return nil
	}

	return fmt.Errorf("failed to unmarshal Category at line %d", value.Line)
}

type Hooks struct {
	PostClone []string `yaml:"post-clone,omitempty"` // Commands to run after clone
}

type Component struct {
	Repository string `yaml:"repository,omitempty"`
	Version    string `yaml:"version,omitempty"`
}

// IsGit returns true if the component is configured with a valid remote Git repository.
func (c Component) IsGit() bool {
	repo := c.Repository
	if repo == "" || repo == "none" || repo == "local" || repo == "false" || repo == "null" || repo == "off" {
		return false
	}
	if repo == "local" || strings.Contains(repo, "@company") {
		return false
	}
	return true
}
