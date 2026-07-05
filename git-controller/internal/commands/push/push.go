package push

import (
	"fmt"
	"path/filepath"
	"strings"
	"workspace/git-controller/internal/gitutil"
	"workspace/git-controller/internal/system"
	"workspace/git-controller/internal/ui"
)

// Run pushes local commits for all cloned repositories defined in the system definition.
func Run() error {
	ui.Header("Push Repositories")

	sys, workspace, err := system.LoadDefinition("system-definition.yaml")
	if err != nil {
		return fmt.Errorf("error loading system definition: %w", err)
	}

	hasErrors := false

	for catName, repos := range sys.Repos {
		if len(repos) == 0 {
			continue
		}

		catDir := workspace.GetCategoryDir(catName)
		if abs, err := filepath.Abs(catDir); err == nil {
			catDir = abs
		}

		for name, comp := range repos {
			if comp.Repository == "" || strings.Contains(comp.Repository, "@company") {
				continue
			}

			displayName := name
			if displayName == "" {
				displayName = catName
			}

			targetPath := filepath.Join(catDir, name)

			if !gitutil.IsCloned(targetPath) {
				ui.Warn("⊘ %s (not cloned, skipping)", displayName)
				continue
			}

			ui.Info("→ Pushing %s", displayName)
			if err := gitutil.Push(targetPath); err != nil {
				ui.Error("Failed to push %s: %v", displayName, err)
				hasErrors = true
			}
		}
	}

	if hasErrors {
		return fmt.Errorf("push completed with errors")
	}

	ui.Success("All repositories pushed successfully!")
	return nil
}
