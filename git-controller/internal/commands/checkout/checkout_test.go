package checkout

import (
	"os"
	"testing"
)

func TestCheckout_RunInNonWorkspace(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "checkout-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}

	err = Run()
	if err == nil {
		t.Fatal("expected error when running checkout outside workspace, got nil")
	}
}
