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

// Run scaffolds a new workspace in the current directory.
func Run() error {
	ui.Header("Initialise Workspace")

	defPath := "system-definition.yaml"
	if _, err := os.Stat(defPath); err == nil {
		return fmt.Errorf("%s already exists — this directory is already a workspace", defPath)
	}

	if err := os.WriteFile(defPath, []byte(defaultDefinition), 0644); err != nil {
		return fmt.Errorf("failed to create %s: %w", defPath, err)
	}
	ui.Success("Created %s", defPath)

	gitReposDir := filepath.Join(".", "git-repositories")
	if err := os.MkdirAll(gitReposDir, 0755); err != nil {
		return fmt.Errorf("failed to create git-repositories/: %w", err)
	}

	gitignorePath := ".gitignore"
	if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
		content := "git-repositories/\n"
		if err := os.WriteFile(gitignorePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create .gitignore: %w", err)
		}
		ui.Success("Created .gitignore")
	}

	makefile := "Makefile"
	if _, err := os.Stat(makefile); os.IsNotExist(err) {
		content := `.PHONY: clone pull push status validate doctor ssh version

clone:
	cd git-controller && go run . clone

pull:
	cd git-controller && go run . pull

push:
	cd git-controller && go run . push

status:
	cd git-controller && go run . status

validate:
	cd git-controller && go run . validate

doctor:
	cd git-controller && go run . doctor

ssh:
	cd git-controller && go run . ssh

version:
	cd git-controller && go run . version
`
		if err := os.WriteFile(makefile, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create Makefile: %w", err)
		}
		ui.Success("Created Makefile")
	}

	ui.Success("Workspace initialised! Edit %s to add your repositories, then run 'make clone'.", defPath)
	return nil
}
