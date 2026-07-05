package system

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefinition(t *testing.T) {
	// Create a temporary system definition file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "system-definition.yaml")

	content := `
system-version: 1.0.0
applications:
  test-service:
    repository: git@github.com:test/repo.git
    version: main
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	// Test loading
	sys, err := LoadDefinition(configPath)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if sys == nil {
		t.Fatal("expected system definition, got nil")
	}

	if sys.SystemVersion != "1.0.0" {
		t.Errorf("expected version 1.0.0, got %s", sys.SystemVersion)
	}

	if _, ok := sys.Applications["test-service"]; !ok {
		t.Error("expected test-service to be present")
	}
}

func TestLoadDefinitionWithEnv(t *testing.T) {
	os.Setenv("TEST_VAR", "my-repo")
	defer os.Unsetenv("TEST_VAR")

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "system-definition-env.yaml")

	content := `
system-version: 1.0.0
applications:
  test-service:
    repository: git@github.com:test/${TEST_VAR}.git
    version: main
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	sys, err := LoadDefinition(configPath)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	repo := sys.Applications["test-service"].Repository
	expected := "git@github.com:test/my-repo.git"
	if repo != expected {
		t.Errorf("expected repository %s, got %s", expected, repo)
	}
}
