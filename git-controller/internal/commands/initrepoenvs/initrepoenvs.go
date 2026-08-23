package initrepoenvs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"workspace/git-controller/internal/system"
	"workspace/git-controller/internal/ui"
)

// Run initialises .env files across all managed repositories from .env.example (or template) files.
func Run() error {
	ui.Header("Initialise Repository Environments")

	sys, workspace, err := system.LoadDefinition(system.DefaultDefinitionFile)
	if err != nil {
		return fmt.Errorf("error loading system definition: %w", err)
	}

	hasErrors := false
	var initialisedCount, existingCount, skippedCount int

	envCandidates := []string{".env.example", ".env.sample", ".env.template", "example.env"}

	for catName, repos := range sys.Repos {
		if len(repos) == 0 {
			continue
		}

		ui.Step(2, fmt.Sprintf("Category: %s", catName))
		catDir := workspace.GetCategoryDir(catName)
		if abs, err := filepath.Abs(catDir); err == nil {
			catDir = abs
		}

		for name := range repos {
			displayName := name
			if displayName == "" {
				displayName = catName
			}

			targetPath := filepath.Join(catDir, name)

			fi, err := os.Stat(targetPath)
			if err != nil || !fi.IsDir() {
				ui.Warn("⊘ %s (not cloned / directory missing, skipping)", displayName)
				skippedCount++
				continue
			}

			envPath := filepath.Join(targetPath, ".env")
			if _, err := os.Stat(envPath); err == nil {
				ui.Info("✓ %s: .env already exists", displayName)
				existingCount++
				continue
			}

			var templatePath string
			for _, candidate := range envCandidates {
				candidatePath := filepath.Join(targetPath, candidate)
				if _, err := os.Stat(candidatePath); err == nil {
					templatePath = candidatePath
					break
				}
			}

			if templatePath == "" {
				ui.Info("  %s: No .env template found (skipped)", displayName)
				skippedCount++
				continue
			}

			if err := copyFile(templatePath, envPath); err != nil {
				ui.Error("Failed to create .env for %s from %s: %v", displayName, filepath.Base(templatePath), err)
				hasErrors = true
			} else {
				ui.Success("✓ %s: Created .env from %s", displayName, filepath.Base(templatePath))
				initialisedCount++
			}
		}
	}

	// Also check workspace root .env if .env.example exists and .env does not
	rootEnv := ".env"
	rootExample := ".env.example"
	if _, err := os.Stat(rootExample); err == nil {
		if _, err := os.Stat(rootEnv); os.IsNotExist(err) {
			if err := copyFile(rootExample, rootEnv); err == nil {
				ui.Success("✓ root: Created .env from %s", rootExample)
				initialisedCount++
			}
		}
	}

	ui.Info("\nSummary: %d initialised, %d already existed, %d skipped", initialisedCount, existingCount, skippedCount)

	if hasErrors {
		return fmt.Errorf("environment initialisation completed with errors")
	}

	ui.Success("Repository environment setup complete!")
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}
