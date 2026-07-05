package status

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"workspace/git-controller/internal/gitutil"
	"workspace/git-controller/internal/system"
	"workspace/git-controller/internal/ui"
)

// Run displays a dashboard-style overview of all repositories.
func Run() error {
	ui.Header("Repository Status")

	sys, workspace, err := system.LoadDefinition("system-definition.yaml")
	if err != nil {
		return fmt.Errorf("error loading system definition: %w", err)
	}

	var totalRepos, cloned, dirty, ahead, behind int

	for catName, repos := range sys.Repos {
		if len(repos) == 0 {
			continue
		}

		catDir := workspace.GetCategoryDir(catName)
		if abs, err := filepath.Abs(catDir); err == nil {
			catDir = abs
		}

		ui.Step(0, fmt.Sprintf("Category: %s", catName))

		for name, comp := range repos {
			totalRepos++

			displayName := name
			if displayName == "" {
				displayName = catName
			}

			if comp.Repository == "" || strings.Contains(comp.Repository, "@company") {
				ui.Info("  %-30s  %s", displayName, "⊘ no remote configured")
				continue
			}

			targetPath := filepath.Join(catDir, name)

			if !gitutil.IsCloned(targetPath) {
				ui.Warn("  %-30s  %s", displayName, "not cloned")
				continue
			}
			cloned++

			branch := getBranch(targetPath)
			isDirty := hasDirtyFiles(targetPath)
			aheadN, behindN := getAheadBehind(targetPath)

			if isDirty {
				dirty++
			}
			if aheadN > 0 {
				ahead++
			}
			if behindN > 0 {
				behind++
			}

			var parts []string
			parts = append(parts, fmt.Sprintf("branch:%s", branch))
			if isDirty {
				parts = append(parts, ui.ColorYellow+"dirty"+ui.ColorReset)
			} else {
				parts = append(parts, ui.ColorGreen+"clean"+ui.ColorReset)
			}
			if aheadN > 0 {
				parts = append(parts, fmt.Sprintf("↑%d", aheadN))
			}
			if behindN > 0 {
				parts = append(parts, fmt.Sprintf("↓%d", behindN))
			}

			ui.Info("  %-30s  %s", displayName, strings.Join(parts, "  "))
		}
	}

	fmt.Println()
	ui.Info("Total: %d repos | Cloned: %d | Dirty: %d | Ahead: %d | Behind: %d",
		totalRepos, cloned, dirty, ahead, behind)

	if dirty > 0 {
		ui.Warn("Some repositories have uncommitted changes.")
	}
	if behind > 0 {
		ui.Warn("Some repositories are behind remote. Run 'make pull' to update.")
	}
	if cloned < totalRepos {
		ui.Warn("Some repositories are not cloned. Run 'make clone' to set up.")
	}

	return nil
}

func getBranch(dir string) string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(out))
}

func hasDirtyFiles(dir string) bool {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return len(strings.TrimSpace(string(out))) > 0
}

func getAheadBehind(dir string) (aheadN, behindN int) {
	cmd := exec.Command("git", "rev-list", "--left-right", "--count", "HEAD...@{upstream}")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return 0, 0
	}
	parts := strings.Fields(strings.TrimSpace(string(out)))
	if len(parts) == 2 {
		fmt.Sscanf(parts[0], "%d", &aheadN)
		fmt.Sscanf(parts[1], "%d", &behindN)
	}
	return
}
