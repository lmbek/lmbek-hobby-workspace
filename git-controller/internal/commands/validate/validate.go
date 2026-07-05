package validate

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"workspace/git-controller/internal/system"
	"workspace/git-controller/internal/ui"
)

func Run() error {
	ui.Header("Repository Validation")

	sys, workspace, err := system.LoadDefinition("system-definition.yaml")
	if err != nil {
		return fmt.Errorf("could not read system definition: %w", err)
	}

	ui.Info("Validating repositories...")

	hasErrors := false

	for catName, repos := range sys.Repos {
		catDir := workspace.GetCategoryDir(catName)

		for name, comp := range repos {
			slog.Debug("Checking repository", "category", catName, "name", name)
			targetPath := filepath.Join(catDir, name)

			displayName := name
			if displayName == "" {
				displayName = catName
			}

			if comp.Repository == "" {
				ui.Info("[%s] No remote configured, skipping", displayName)
				continue
			}

			if err := validateRepo(targetPath, comp.Version); err != nil {
				ui.Error("[%s] %v", displayName, err)
				hasErrors = true
			} else {
				ui.Success("[%s] OK", displayName)
			}
		}
	}

	if hasErrors {
		return fmt.Errorf("validation completed with errors")
	}

	ui.Success("All repositories valid!")
	return nil
}

func validateRepo(targetPath string, expectedVersion string) error {
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		return fmt.Errorf("not cloned — run '%s clone'", system.CLIName)
	}

	currentVersion, err := getGitCurrentVersion(targetPath)
	if err != nil {
		return fmt.Errorf("could not determine current branch: %v", err)
	}

	if currentVersion != expectedVersion {
		slog.Warn("Branch mismatch", "path", targetPath, "expected", expectedVersion, "found", currentVersion)
	}

	isDirty, err := isGitDirty(targetPath)
	if err != nil {
		return fmt.Errorf("could not check git status: %v", err)
	}
	if isDirty {
		slog.Warn("Uncommitted changes", "path", targetPath)
	}

	return nil
}

func getGitCurrentVersion(dir string) (string, error) {
	// Try to get exact tag match first
	cmd := exec.Command("git", "describe", "--tags", "--exact-match")
	cmd.Dir = dir
	output, err := cmd.Output()
	if err == nil {
		return strings.TrimSpace(string(output)), nil
	}

	// If no tag, get branch name or hash
	cmd = exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = dir
	output, err = cmd.Output()
	if err != nil {
		return "", err
	}

	branch := strings.TrimSpace(string(output))
	if branch == "HEAD" {
		// Detached HEAD, get hash
		cmd = exec.Command("git", "rev-parse", "--short", "HEAD")
		cmd.Dir = dir
		output, err = cmd.Output()
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(output)), nil
	}

	return branch, nil
}

func isGitDirty(dir string) (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return len(strings.TrimSpace(string(output))) > 0, nil
}
