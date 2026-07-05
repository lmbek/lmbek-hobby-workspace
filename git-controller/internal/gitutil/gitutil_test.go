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

func TestProcessGitComponent_RejectsHTTP(t *testing.T) {
	err := ProcessGitComponent(t.TempDir(), "repo", "https://github.com/user/repo.git")
	if err == nil || !strings.Contains(err.Error(), "enforces SSH") {
		t.Fatalf("expected SSH enforcement error, got: %v", err)
	}
}
