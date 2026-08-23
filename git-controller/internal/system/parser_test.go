package system

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefinition(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "repo-definition.yaml")

	content := `
repos:
  services:
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

	services, ok := sys.Repos["services"]
	if !ok {
		t.Fatal("expected services category to be present")
	}
	if _, ok := services["test-service"]; !ok {
		t.Error("expected test-service to be present")
	}
}

func TestLoadDefinitionWithEnv(t *testing.T) {
	t.Setenv("TEST_VAR", "my-repo")

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "system-definition-env.yaml")

	content := `
repos:
  services:
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

	repo := sys.Repos["services"]["test-service"].Repository
	expected := "git@github.com:test/my-repo.git"
	if repo != expected {
		t.Errorf("expected repository %s, got %s", expected, repo)
	}
}

func TestLoadDefinitionDynamicCategories(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "repo-definition.yaml")

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

func TestLoadDefinitionNonGitAndLocalEntries(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "repo-definition.yaml")

	content := `
system-version: 1.0.0
repos:
  local-flat-category:
  local-nested-category:
    local-tool1:
    local-tool2:
      repository: local
    local-tool3:
      repository: none
    git-tool:
      repository: git@github.com:test/tool.git
      version: main
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	sys, _, err := LoadDefinition(configPath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	flat, ok := sys.Repos["local-flat-category"]
	if !ok {
		t.Fatal("expected local-flat-category to be present")
	}
	if comp, ok := flat[""]; !ok || comp.IsGit() {
		t.Errorf("expected local-flat-category to have non-git singleton, got %+v", comp)
	}

	nested, ok := sys.Repos["local-nested-category"]
	if !ok {
		t.Fatal("expected local-nested-category to be present")
	}

	if t1, ok := nested["local-tool1"]; !ok || t1.IsGit() {
		t.Errorf("expected local-tool1 to be non-git, got %+v", t1)
	}
	if t2, ok := nested["local-tool2"]; !ok || t2.IsGit() {
		t.Errorf("expected local-tool2 to be non-git, got %+v", t2)
	}
	if t3, ok := nested["local-tool3"]; !ok || t3.IsGit() {
		t.Errorf("expected local-tool3 to be non-git, got %+v", t3)
	}
	if gt, ok := nested["git-tool"]; !ok || !gt.IsGit() {
		t.Errorf("expected git-tool to be git, got %+v", gt)
	}
}
