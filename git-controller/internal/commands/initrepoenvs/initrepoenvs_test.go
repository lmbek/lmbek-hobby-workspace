package initrepoenvs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyFile(t *testing.T) {
	tempDir := t.TempDir()
	src := filepath.Join(tempDir, ".env.example")
	dst := filepath.Join(tempDir, ".env")

	content := "PORT=8080\nDEBUG=true\n"
	if err := os.WriteFile(src, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	if err := copyFile(src, dst); err != nil {
		t.Fatalf("copyFile failed: %v", err)
	}

	dstContent, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("failed to read destination file: %v", err)
	}

	if string(dstContent) != content {
		t.Errorf("expected %q, got %q", content, string(dstContent))
	}
}

func TestCopyFile_SourceMissing(t *testing.T) {
	tempDir := t.TempDir()
	src := filepath.Join(tempDir, "nonexistent.example")
	dst := filepath.Join(tempDir, ".env")

	if err := copyFile(src, dst); err == nil {
		t.Error("expected error for missing source file, got nil")
	}
}
