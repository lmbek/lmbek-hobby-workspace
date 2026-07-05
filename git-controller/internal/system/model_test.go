package system

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func normalize(p string) string {
	if runtime.GOOS == "windows" {
		return strings.ReplaceAll(filepath.Clean(p), "\\", "/")
	}
	return filepath.Clean(p)
}

func TestWorkspaceGetCategoryDir_DefaultsToGitRepositories(t *testing.T) {
	t.Setenv("APPLICATIONS_DIR", "")
	w := &Workspace{Root: filepath.FromSlash("/tmp/root")}
	dir := w.GetCategoryDir("applications")
	got := normalize(dir)
	want := normalize("/tmp/root/git-repositories/applications")
	if got != want {
		t.Fatalf("expected %s, got %s", want, got)
	}
}

func TestWorkspaceGetCategoryDir_EnvOverride(t *testing.T) {
	// ORCHESTRATOR_DIR should override orchestrator
	t.Setenv("ORCHESTRATOR_DIR", filepath.FromSlash("/orch"))
	w := &Workspace{Root: filepath.FromSlash("/tmp/root")}
	dir := w.GetCategoryDir("orchestrator")
	if normalize(dir) != normalize("/orch") {
		t.Fatalf("env override not respected for orchestrator: %s", dir)
	}

	// INFRASTRUCTURE_DIR should override infrastructure
	t.Setenv("INFRASTRUCTURE_DIR", filepath.FromSlash("/infra"))
	dir = w.GetCategoryDir("infrastructure")
	if normalize(dir) != normalize("/infra") {
		t.Fatalf("env override not respected for infrastructure: %s", dir)
	}

	// Unknown category falls back under git-repositories
	os.Unsetenv("MYCUSTOM_DIR")
	dir = w.GetCategoryDir("mycustom")
	if !strings.HasSuffix(normalize(dir), normalize("git-repositories/mycustom")) {
		t.Fatalf("unexpected default for mycustom: %s", dir)
	}
}
