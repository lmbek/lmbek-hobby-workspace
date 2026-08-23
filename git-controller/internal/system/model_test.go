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
	dir := w.GetCategoryDir("services")
	got := normalize(dir)
	want := normalize("/tmp/root/git-repositories/services")
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

func TestComponent_IsGit(t *testing.T) {
	tests := []struct {
		repo string
		want bool
	}{
		{"git@github.com:user/repo.git", true},
		{"ssh://git@github.com/user/repo.git", true},
		{"", false},
		{"none", false},
		{"local", false},
		{"false", false},
		{"null", false},
		{"git@company/repo.git", false},
	}

	for _, tt := range tests {
		comp := Component{Repository: tt.repo}
		if got := comp.IsGit(); got != tt.want {
			t.Errorf("Component{Repository: %q}.IsGit() = %v, want %v", tt.repo, got, tt.want)
		}
	}
}
