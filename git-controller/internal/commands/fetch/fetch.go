package fetch

import (
	"fmt"
	"path/filepath"
	"workspace/git-controller/internal/gitutil"
	"workspace/git-controller/internal/system"
	"workspace/git-controller/internal/ui"
)

// Run fetches all remotes for every repository defined in the system definition.
func Run() error {
	ui.Header("Fetch Repositories")

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
				ui.Info("→ Skipping %s (not cloned)", displayName)
				continue
			}

			ui.Info("→ Fetching %s", displayName)
			if err := gitutil.Fetch(targetPath); err != nil {
				ui.Error("Failed to fetch %s: %v", displayName, err)
				hasErrors = true
			}
		}
	}

	if hasErrors {
		return fmt.Errorf("fetch completed with errors")
	}

	ui.Success("All repositories fetched!")
	return nil
}
