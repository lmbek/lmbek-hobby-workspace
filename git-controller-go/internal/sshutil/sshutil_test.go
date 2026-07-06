package sshutil

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetPublicKeyContent(t *testing.T) {
	tmpDir := t.TempDir()
	privKeyPath := filepath.Join(tmpDir, "id_ed25519.private")
	pubKeyPath := filepath.Join(tmpDir, "id_ed25519.public")
	pubContent := "ssh-ed25519 ABC...fake...key...content"

	err := os.WriteFile(pubKeyPath, []byte(pubContent), 0644)
	if err != nil {
		t.Fatalf("failed to write pub key: %v", err)
	}

	content, path, err := GetPublicKeyContent(privKeyPath)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !strings.Contains(content, pubContent) {
		t.Errorf("expected content to contain %s, got %s", pubContent, content)
	}
	if path != pubKeyPath {
		t.Errorf("expected path %s, got %s", pubKeyPath, path)
	}
}
