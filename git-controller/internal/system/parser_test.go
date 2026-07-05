package system

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefinition(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "system-definition.yaml")

	content := `
system-version: 1.0.0
repos:
  applications:
    test-service:
      repository: git@github.com:test/repo.git
      version: main
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	sys, workspace, err := LoadDefinition(configPath)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if sys == nil || workspace == nil {
		t.Fatal("expected system definition and workspace, got nil")
	}

	if sys.SystemVersion != "1.0.0" {
		t.Errorf("expected version 1.0.0, got %s", sys.SystemVersion)
	}

	apps, ok := sys.Repos["applications"]
	if !ok {
		t.Fatal("expected applications category to be present")
	}
	if _, ok := apps["test-service"]; !ok {
		t.Error("expected test-service to be present")
	}
}

func TestLoadDefinitionWithEnv(t *testing.T) {
	t.Setenv("TEST_VAR", "my-repo")

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "system-definition-env.yaml")

	content := `
system-version: 1.0.0
repos:
  applications:
    test-service:
      repository: git@github.com:test/${TEST_VAR}.git
      version: main
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	sys, _, err := LoadDefinition(configPath)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	repo := sys.Repos["applications"]["test-service"].Repository
	expected := "git@github.com:test/my-repo.git"
	if repo != expected {
		t.Errorf("expected repository %s, got %s", expected, repo)
	}
}

func TestLoadDefinitionDynamicCategories(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "system-definition.yaml")

	content := `
system-version: 1.0.0
repos:
  custom-category:
    my-repo:
      repository: git@github.com:test/custom.git
      version: develop
  another-one:
    repository: git@github.com:test/singleton.git
    version: main
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	sys, _, err := LoadDefinition(configPath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(sys.Repos) != 2 {
		t.Fatalf("expected 2 categories, got %d", len(sys.Repos))
	}

	custom, ok := sys.Repos["custom-category"]
	if !ok {
		t.Fatal("expected custom-category to be present")
	}
	if custom["my-repo"].Version != "develop" {
		t.Errorf("expected version develop, got %s", custom["my-repo"].Version)
	}

	another, ok := sys.Repos["another-one"]
	if !ok {
		t.Fatal("expected another-one to be present")
	}
	// Singleton should have empty-string key
	if another[""].Repository != "git@github.com:test/singleton.git" {
		t.Errorf("expected singleton repo, got %s", another[""].Repository)
	}
}
