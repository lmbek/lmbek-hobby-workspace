package sync

import (
	"fmt"
	"path/filepath"
	"workspace/git-controller/internal/gitutil"
	"workspace/git-controller/internal/system"
	"workspace/git-controller/internal/ui"
)

// Run synchronises all repositories defined in the system definition using read-only pulls (fast-forward only).
func Run() error {
	ui.Header("Sync Repositories (Readonly Pull)")

	sys, workspace, err := system.LoadDefinition(system.DefaultDefinitionFile)
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

		for name, comp := range repos {
			if !comp.IsGit() {
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

			ui.Info("→ Syncing %s", displayName)
			if err := gitutil.Sync(targetPath); err != nil {
				ui.Error("Failed to sync %s: %v", displayName, err)
				hasErrors = true
			}
		}
	}

	if hasErrors {
		return fmt.Errorf("sync completed with errors")
	}

	ui.Success("All repositories synchronised successfully!")
	return nil
}
