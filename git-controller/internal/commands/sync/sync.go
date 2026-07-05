package sync

import (
	"fmt"
	"path/filepath"
	"workspace/git-controller/internal/gitutil"
	"workspace/git-controller/internal/system"
	"workspace/git-controller/internal/ui"
)

func Run() error {
	ui.Header("System Synchronization")

	sys, err := system.LoadDefinition("system-definition.yaml")
	if err != nil {
		return fmt.Errorf("error loading system definition: %w", err)
	}

	categories := []struct {
		name       string
		components map[string]system.Component
	}{
		{"proxy", sys.Proxy},
		{"applications", sys.Applications},
		{"infrastructure", sys.Infrastructure},
		{"orchestrator", sys.Orchestrator},
		{"platform", sys.Platform},
		{"tools", sys.Tools},
		{"docs", sys.Docs},
	}

	hasErrors := false

	for _, cat := range categories {
		if len(cat.components) == 0 {
			continue
		}

		ui.Step(2, fmt.Sprintf("Synchronizing %s", cat.name))
		catDir := system.GetCategoryDir(cat.name)

		if abs, err := filepath.Abs(catDir); err == nil {
			catDir = abs
		}

		gitutil.EnsureDir(catDir)

		for name, comp := range cat.components {
			displayName := name
			if displayName == "" {
				displayName = cat.name + " (singleton)"
			}
			ui.Info("→ %s", displayName)
			if err := gitutil.ProcessGitComponent(catDir, name, comp.Repository); err != nil {
				ui.Error("Failed to sync %s: %v", displayName, err)
				hasErrors = true
			}
		}
	}

	if hasErrors {
		return fmt.Errorf("synchronization failed with errors")
	}

	// Running Hooks
	if len(sys.Hooks.PostSync) > 0 {
		ui.Step(3, "Running Post-Sync Hooks")
		for _, hook := range sys.Hooks.PostSync {
			ui.Info("Running: %s", hook)
			if err := gitutil.RunHook(hook); err != nil {
				ui.Error("Hook failed: %s (%v)", hook, err)
			}
		}
	}

	ui.Success("Synchronization finished successfully!")
	return nil
}
