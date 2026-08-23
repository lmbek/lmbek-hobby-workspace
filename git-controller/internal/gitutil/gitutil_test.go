package gitutil

import (
	"errors"
	"runtime"
	"strings"
	"testing"
)

func TestEnhanceGitError_PublicKey(t *testing.T) {
	base := errors.New("fatal: Could not read from remote repository.\nPermission denied (publickey).")
	out := EnhanceGitError(base).Error()
	if !strings.Contains(out, "SSH authentication failed") {
		t.Fatalf("expected SSH hint, got: %s", out)
	}
	if runtime.GOOS == "windows" && !strings.Contains(out, "OpenSSH Authentication Agent") {
		t.Fatalf("expected Windows agent hint in: %s", out)
	}
}

func TestEnhanceGitError_HostKey(t *testing.T) {
	base := errors.New("Host key verification failed")
	out := EnhanceGitError(base).Error()
	if !strings.Contains(out, "Host key verification failed") {
		t.Fatalf("expected host key hint, got: %s", out)
	}
}

func TestEnhanceGitError_Nil(t *testing.T) {
	if EnhanceGitError(nil) != nil {
		t.Fatal("expected nil for nil input")
	}
}

func TestClone_RejectsHTTP(t *testing.T) {
	err := Clone("https://github.com/user/repo.git", t.TempDir())
	if err == nil || !strings.Contains(err.Error(), "enforces SSH") {
		t.Fatalf("expected SSH enforcement error, got: %v", err)
	}
}

func TestValidateRemoteURL(t *testing.T) {
	if err := validateRemoteURL("git@github.com:user/repo.git"); err != nil {
		t.Fatalf("SSH URL should be allowed: %v", err)
	}
	if err := validateRemoteURL("http://github.com/user/repo.git"); err == nil {
		t.Fatal("HTTP URL should be rejected")
	}
	if err := validateRemoteURL("https://github.com/user/repo.git"); err == nil {
		t.Fatal("HTTPS URL should be rejected")
	}
}

func TestIsCloned(t *testing.T) {
	if IsCloned(t.TempDir()) {
		t.Fatal("empty dir should not be considered cloned")
	}
}

func TestIsNonEmptyDir(t *testing.T) {
	if IsNonEmptyDir(t.TempDir()) {
		t.Fatal("empty dir should return false")
	}
	if IsNonEmptyDir("/nonexistent/path") {
		t.Fatal("nonexistent path should return false")
	}
}

func TestSync_NonExistentDir(t *testing.T) {
	err := Sync("/nonexistent/repo/dir")
	if err == nil {
		t.Fatal("expected error when syncing non-existent repo, got nil")
	}
}
