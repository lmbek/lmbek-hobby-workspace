package system

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func LoadDefinition(path string) (*SystemDefinition, error) {
	// 1. Try path as is
	if _, err := os.Stat(path); err == nil {
		setWorkspaceRoot(filepath.Dir(path))
		return loadFile(path)
	}

	// 2. Try relative to WORKSPACE_ROOT if set
	workspaceRoot := os.Getenv("WORKSPACE_ROOT")
	if workspaceRoot != "" {
		altPath := filepath.Join(workspaceRoot, path)
		if _, err := os.Stat(altPath); err == nil {
			setWorkspaceRoot(workspaceRoot)
			return loadFile(altPath)
		}
	}

	// 3. Try searching parent directories for the workspace folder
	currentDir, err := os.Getwd()
	if err == nil {
		for {
			altPath := filepath.Join(currentDir, path)
			if _, err := os.Stat(altPath); err == nil {
				setWorkspaceRoot(currentDir)
				return loadFile(altPath)
			}
			parent := filepath.Dir(currentDir)
			if parent == currentDir {
				break
			}
			currentDir = parent
		}
	}

	return nil, fmt.Errorf("could not find system definition: %s", path)
}

func setWorkspaceRoot(root string) {
	if abs, err := filepath.Abs(root); err == nil {
		root = abs
	}
	os.Setenv("WORKSPACE_ROOT", root)
	slog.Debug("Setting effective WORKSPACE_ROOT", "root", root)
}

func loadFile(path string) (*SystemDefinition, error) {
	slog.Debug("Loading system definition", "path", path)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var system SystemDefinition
	err = yaml.Unmarshal([]byte(os.ExpandEnv(string(data))), &system)
	if err != nil {
		return nil, err
	}

	return &system, nil
}
