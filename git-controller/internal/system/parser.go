package system

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func LoadDefinition(path string) (*SystemDefinition, *Workspace, error) {
	// resolveRoot returns the workspace root from a found definition file path.
	// The definition lives inside git-repositories/, so the workspace root is
	// two levels up from the file (git-repositories/ → workspace root).
	resolveRoot := func(defPath string) string {
		return filepath.Dir(filepath.Dir(defPath))
	}

	// 1. Try path as is
	if _, err := os.Stat(path); err == nil {
		return loadFile(path, resolveRoot(path))
	}

	// 1b. If default repo-definition.yaml not found, check legacy system-definition.yaml
	if filepath.Base(path) == "repo-definition.yaml" {
		legacyPath := filepath.Join(filepath.Dir(path), "system-definition.yaml")
		if _, err := os.Stat(legacyPath); err == nil {
			return loadFile(legacyPath, resolveRoot(legacyPath))
		}
	}

	// 2. Try searching parent directories for the workspace folder
	currentDir, err := os.Getwd()
	if err == nil {
		for {
			altPath := filepath.Join(currentDir, path)
			if _, err := os.Stat(altPath); err == nil {
				return loadFile(altPath, resolveRoot(altPath))
			}
			parent := filepath.Dir(currentDir)
			if parent == currentDir {
				break
			}
			currentDir = parent
		}
	}

	return nil, nil, fmt.Errorf("could not find system definition: %s", path)
}

func loadFile(path string, root string) (*SystemDefinition, *Workspace, error) {
	slog.Debug("Loading system definition", "path", path)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	var system SystemDefinition
	err = yaml.Unmarshal([]byte(os.ExpandEnv(string(data))), &system)
	if err != nil {
		return nil, nil, err
	}

	// Normalize empty or null categories to singleton local component
	for catName, repos := range system.Repos {
		if repos == nil || len(repos) == 0 {
			system.Repos[catName] = Category{"": Component{}}
		}
	}

	if abs, err := filepath.Abs(root); err == nil {
		root = abs
	}
	return &system, &Workspace{Root: root}, nil
}
