package wsinit

import (
	"fmt"
	"os"
	"path/filepath"
	"workspace/git-controller/internal/ui"
)

const defaultDefinition = `# System Definition for Workspace Controller
system-version: main

hooks:
    post-clone:
        - echo "Clone complete! Run 'make status' to see the state of all repositories."

# Each key under "repos" is a category folder inside git-repositories/.
# Categories are dynamic — add whatever makes sense for your project.
repos:
    # applications:
    #     my-service:
    #         repository: git@github.com:org/my-service.git
    #         version: main
`

// workspaceRoot returns the workspace root directory. It uses WORKSPACE_ROOT
// if set, otherwise falls back to the parent of the current working directory
// (assuming the binary runs from git-controller/).
func workspaceRoot() string {
	if root := os.Getenv("WORKSPACE_ROOT"); root != "" {
		return root
	}
	if cwd, err := os.Getwd(); err == nil {
		return filepath.Dir(cwd)
	}
	return "."
}

// Run scaffolds a new workspace by creating the system-definition.yaml in the
// workspace root directory. It will not overwrite an existing definition file.
func Run() error {
	ui.Header("Initialise Workspace")

	root := workspaceRoot()
	defPath := filepath.Join(root, "system-definition.yaml")

	if _, err := os.Stat(defPath); err == nil {
		ui.Info("system-definition.yaml already exists at %s — skipping creation", defPath)
	} else {
		if err := os.WriteFile(defPath, []byte(defaultDefinition), 0644); err != nil {
			return fmt.Errorf("failed to create %s: %w", defPath, err)
		}
		ui.Success("Created %s", defPath)
	}

	gitReposDir := filepath.Join(root, "git-repositories")
	if err := os.MkdirAll(gitReposDir, 0755); err != nil {
		return fmt.Errorf("failed to create git-repositories/: %w", err)
	}

	gitignorePath := filepath.Join(root, ".gitignore")
	if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
		content := "git-repositories/\n"
		if err := os.WriteFile(gitignorePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create .gitignore: %w", err)
		}
		ui.Success("Created .gitignore")
	}

	makefile := filepath.Join(root, "Makefile")
	if _, err := os.Stat(makefile); os.IsNotExist(err) {
		content := `.PHONY: init clone pull push scaffold checkout status validate doctor ssh

export WORKSPACE_ROOT ?= $(abspath .)

init:
	cd git-controller && go run . init

clone:
	cd git-controller && go run . clone

pull:
	cd git-controller && go run . pull

push:
	cd git-controller && go run . push

scaffold:
	cd git-controller && go run . scaffold

checkout:
	cd git-controller && go run . checkout

status:
	cd git-controller && go run . status

validate:
	cd git-controller && go run . validate

doctor:
	cd git-controller && go run . doctor

ssh:
	cd git-controller && go run . ssh
`
		if err := os.WriteFile(makefile, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create Makefile: %w", err)
		}
		ui.Success("Created Makefile")
	}

	ui.Success("Workspace initialised! Edit system-definition.yaml to add your repositories, then run 'make clone'.")
	return nil
}
