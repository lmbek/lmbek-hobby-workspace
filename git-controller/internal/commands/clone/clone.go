package clone

import (
	"fmt"
	"path/filepath"
	"strings"
	"workspace/git-controller/internal/gitutil"
	"workspace/git-controller/internal/system"
	"workspace/git-controller/internal/ui"
)

// Run clones all repositories defined in the system definition.
// Repositories that are already cloned are skipped.
func Run() error {
	ui.Header("Clone Repositories")

	sys, workspace, err := system.LoadDefinition("system-definition.yaml")
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
				ui.Info("✓ %s (already cloned)", displayName)
				continue
			}

			ui.Info("→ Cloning %s", displayName)

			if gitutil.IsNonEmptyDir(targetPath) {
				if err := gitutil.InitAndLink(targetPath, repo); err != nil {
					ui.Error("Failed to initialise %s: %v", displayName, err)
					hasErrors = true
				}
			} else {
				if err := gitutil.Clone(repo, targetPath); err != nil {
					ui.Error("Failed to clone %s: %v", displayName, err)
					hasErrors = true
				}
			}
		}
	}

	if hasErrors {
		return fmt.Errorf("clone completed with errors")
	}

	ui.Success("All repositories cloned successfully!")
	return nil
}
