package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
	"workspace/git-controller/internal/gitutil"
	"workspace/git-controller/internal/system"
	"workspace/git-controller/internal/ui"
)

// Run initialises .git and sets the remote origin for every repository in the
// system definition — without cloning or fetching. This is useful when you
// don't yet have SSH access to the remotes but want the directory structure
// and git configuration ready.
func Run() error {
	ui.Header("Scaffold Repositories")

	sys, workspace, err := system.LoadDefinition("git-repositories/system-definition.yaml")
	if err != nil {
		return fmt.Errorf("error loading system definition: %w", err)
	}

	hasErrors := false

	for catName, repos := range sys.Repos {
		if len(repos) == 0 {
			continue
		}

		ui.Step(2, fmt.Sprintf("Category: %s", catName))
		catDir := workspace.GetCategoryDir(catName)
		if abs, err := filepath.Abs(catDir); err == nil {
			catDir = abs
		}
		gitutil.EnsureDir(catDir)

		for name, comp := range repos {
			repo := comp.Repository
			if repo == "" || strings.Contains(repo, "@company") {
				continue
			}

			displayName := name
			if displayName == "" {
				displayName = catName
			}

			targetPath := filepath.Join(catDir, name)

			if gitutil.IsCloned(targetPath) {
				ui.Info("✓ %s (already initialised)", displayName)
				continue
			}

			gitutil.EnsureDir(targetPath)
			ui.Info("→ Scaffolding %s", displayName)

			if err := gitutil.Scaffold(targetPath, repo); err != nil {
				ui.Error("Failed to scaffold %s: %v", displayName, err)
				hasErrors = true
			}
		}
	}

	if hasErrors {
		return fmt.Errorf("scaffold completed with errors")
	}

	ui.Success("All repositories scaffolded — git init, remote origin set, fetched, and default branch configured.")
	return nil
}
