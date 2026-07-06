package system

import (
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
	w := &Workspace{Root: filepath.FromSlash("/tmp/root")}
	dir := w.GetCategoryDir("applications")
	got := normalize(dir)
	want := normalize("/tmp/root/git-repositories/applications")
	if got != want {
		t.Fatalf("expected %s, got %s", want, got)
	}
}

func TestWorkspaceGetCategoryDir_AlwaysUsesGitRepositories(t *testing.T) {
	w := &Workspace{Root: filepath.FromSlash("/tmp/root")}
	for _, cat := range []string{"orchestrator", "infrastructure", "mycustom"} {
		dir := w.GetCategoryDir(cat)
		want := normalize("/tmp/root/git-repositories/" + cat)
		if normalize(dir) != want {
			t.Fatalf("expected %s, got %s", want, normalize(dir))
		}
	}
}
