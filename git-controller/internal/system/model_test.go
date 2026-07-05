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
	t.Setenv("SERVICES_DIR", "")
	w := &Workspace{Root: filepath.FromSlash("/tmp/root")}
	dir := w.GetCategoryDir("applications")
	got := normalize(dir)
	want := normalize("/tmp/root/git-repositories/applications")
	if got != want {
		t.Fatalf("expected %s, got %s", want, got)
	}
}

func TestWorkspaceGetCategoryDir_EnvOverride(t *testing.T) {
	// SERVICES_DIR should override applications
	t.Setenv("SERVICES_DIR", filepath.FromSlash("/custom/services"))
	w := &Workspace{Root: filepath.FromSlash("/tmp/root")}
	dir := w.GetCategoryDir("applications")
	if normalize(dir) != normalize("/custom/services") {
		t.Fatalf("env override not respected: %s", dir)
	}

	// ORCHESTRATOR_DIR should override orchestrator
	t.Setenv("ORCHESTRATOR_DIR", filepath.FromSlash("/orch"))
	dir = w.GetCategoryDir("orchestrator")
	if normalize(dir) != normalize("/orch") {
		t.Fatalf("env override not respected for orchestrator: %s", dir)
	}

	// INFRA_DIR should override infrastructure
	t.Setenv("INFRA_DIR", filepath.FromSlash("/infra"))
	dir = w.GetCategoryDir("infrastructure")
	if normalize(dir) != normalize("/infra") {
		t.Fatalf("env override not respected for infrastructure: %s", dir)
	}

	// Unknown category falls back under git-repositories
	os.Unsetenv("UNKNOWN_DIR")
	dir = w.GetCategoryDir("tools")
	if !strings.HasSuffix(normalize(dir), normalize("git-repositories/tools")) {
		t.Fatalf("unexpected default for tools: %s", dir)
	}
}

func TestGetOrchestrationDir(t *testing.T) {
	w := &Workspace{Root: filepath.FromSlash("/root")}
	sys := &SystemDefinition{
		Orchestrator:   Category{"orchestrator-stack": {Repository: "git@example.com:x/y.git"}},
		Infrastructure: Category{"infra": {Repository: "git@example.com:i/j.git"}},
	}
	got := normalize(sys.GetOrchestrationDir(w))
	// Should be <root>/git-repositories/orchestrator/orchestrator-stack
	want := normalize("/root/git-repositories/orchestrator/orchestrator-stack")
	if got != want {
		t.Fatalf("want %s, got %s", want, got)
	}

	// If orchestrator missing, fall back to infrastructure base dir
	sys2 := &SystemDefinition{}
	got2 := normalize(sys2.GetOrchestrationDir(w))
	want2 := normalize("/root/git-repositories/infrastructure")
	if got2 != want2 {
		t.Fatalf("fallback want %s, got %s", want2, got2)
	}
}
