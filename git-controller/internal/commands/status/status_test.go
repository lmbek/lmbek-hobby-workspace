package status

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStatusHelpers_NonGitDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "status-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	branch := getBranch(tmpDir)
	if branch != "unknown" {
		t.Errorf("expected branch 'unknown' for non-git dir, got %q", branch)
	}

	dirty := hasDirtyFiles(tmpDir)
	if dirty {
		t.Errorf("expected dirty=false for non-git dir, got true")
	}

	ahead, behind := getAheadBehind(tmpDir)
	if ahead != 0 || behind != 0 {
		t.Errorf("expected ahead=0, behind=0 for non-git dir, got ahead=%d, behind=%d", ahead, behind)
	}
}

func TestStatusHelpers_EmptyPath(t *testing.T) {
	nonExistent := filepath.Join(os.TempDir(), "definitely-does-not-exist-12345")
	branch := getBranch(nonExistent)
	if branch != "unknown" {
		t.Errorf("expected branch 'unknown', got %q", branch)
	}
}
