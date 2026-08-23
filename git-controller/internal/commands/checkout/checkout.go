package checkout

import (
	"fmt"
	"os"
	"path/filepath"
	"workspace/git-controller/internal/gitutil"
	"workspace/git-controller/internal/system"
	"workspace/git-controller/internal/ui"
)

// Run switches all cloned repositories to the branch defined in repo-definition.yaml.
// If a branch argument is provided via CLI argument or BRANCH env var, it overrides the definition.
func Run() error {
	ui.Header("Checkout Branches")

	sys, workspace, err := system.LoadDefinition(system.DefaultDefinitionFile)
	if err != nil {
		return fmt.Errorf("error loading system definition: %w", err)
	}

	overrideBranch := os.Getenv("BRANCH")
	if overrideBranch == "" && len(os.Args) > 2 {
		overrideBranch = os.Args[2]
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

			branch := comp.Version
			if overrideBranch != "" {
				branch = overrideBranch
			}

			ui.Info("→ %s → %s", displayName, branch)
			if err := gitutil.Checkout(targetPath, branch); err != nil {
				ui.Error("Failed to checkout %s: %v", displayName, err)
				hasErrors = true
			}
		}
	}

	if hasErrors {
		return fmt.Errorf("checkout completed with errors")
	}

	ui.Success("All repositories on correct branches!")
	return nil
}
