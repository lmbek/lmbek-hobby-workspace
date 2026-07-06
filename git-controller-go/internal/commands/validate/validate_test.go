package validate

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestValidateRepo_NotCloned(t *testing.T) {
	err := validateRepo(filepath.Join(t.TempDir(), "nonexistent"), "main")
	if err == nil {
		t.Error("expected error for non-existent directory, got nil")
	}
}

func TestValidateRepo_ClonedOnCorrectBranch(t *testing.T) {
	dir := t.TempDir()

	run := func(args ...string) {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), "GIT_CONFIG_GLOBAL=/dev/null")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v failed: %s", args, out)
		}
	}

	run("init", "-b", "main")
	run("config", "user.email", "test@test.com")
	run("config", "user.name", "Test")

	// Create an initial commit so HEAD is valid
	f, _ := os.Create(filepath.Join(dir, "README.md"))
	f.Close()
	run("add", ".")
	run("commit", "-m", "init")

	err := validateRepo(dir, "main")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestGetGitCurrentVersion(t *testing.T) {
	dir := t.TempDir()

	run := func(args ...string) {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), "GIT_CONFIG_GLOBAL=/dev/null")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v failed: %s", args, out)
		}
	}

	run("init", "-b", "main")
	run("config", "user.email", "test@test.com")
	run("config", "user.name", "Test")
	f, _ := os.Create(filepath.Join(dir, "file.txt"))
	f.Close()
	run("add", ".")
	run("commit", "-m", "init")

	version, err := getGitCurrentVersion(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if version != "main" {
		t.Errorf("expected 'main', got %q", version)
	}
}
