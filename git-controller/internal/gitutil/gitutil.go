package gitutil

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"workspace/git-controller/internal/sshutil"
	"workspace/git-controller/internal/ui"
)

// execCommand allows injecting a stub in tests.
var execCommand = exec.Command

// Clone clones a remote repository into the given target path.
func Clone(repo, targetPath string) error {
	if err := validateRemoteURL(repo); err != nil {
		return err
	}
	slog.Debug("Cloning repository", "repo", repo, "target", targetPath)
	return EnhanceGitError(RunGit("clone", repo, targetPath))
}

// Pull fetches and pulls the latest changes in the given repository directory.
func Pull(repoDir string) error {
	slog.Debug("Pulling repository", "dir", repoDir)
	if err := RunGitInDir(repoDir, "fetch", "--all"); err != nil {
		return EnhanceGitError(err)
	}
	return EnhanceGitError(RunGitInDir(repoDir, "pull"))
}

// Checkout switches the given repository to the specified branch.
func Checkout(repoDir, branch string) error {
	slog.Debug("Checking out branch", "dir", repoDir, "branch", branch)
	return EnhanceGitError(RunGitInDir(repoDir, "checkout", branch))
}

// Push pushes local commits to the remote in the given repository directory.
func Push(repoDir string) error {
	slog.Debug("Pushing repository", "dir", repoDir)
	return EnhanceGitError(RunGitInDir(repoDir, "push"))
}

// InitAndLink initialises a git repo in an existing non-empty directory and
// links it to the given remote origin. It attempts to align with the remote
// main/master branch.
func InitAndLink(targetPath, repo string) error {
	if err := validateRemoteURL(repo); err != nil {
		return err
	}
	slog.Debug("Initialising and linking repository", "path", targetPath, "repo", repo)

	if err := RunGitInDir(targetPath, "init"); err != nil {
		return EnhanceGitError(err)
	}
	if err := RunGitInDir(targetPath, "remote", "add", "origin", repo); err != nil {
		return EnhanceGitError(err)
	}
	if err := RunGitInDir(targetPath, "fetch", "origin"); err != nil {
		return EnhanceGitError(err)
	}

	// Attempt to align with the remote main branch
	if err := RunGitInDir(targetPath, "reset", "--soft", "origin/main"); err == nil {
		_ = RunGitInDir(targetPath, "branch", "--set-upstream-to=origin/main", "main")
	} else if err := RunGitInDir(targetPath, "reset", "--soft", "origin/master"); err == nil {
		_ = RunGitInDir(targetPath, "checkout", "-b", "main")
		_ = RunGitInDir(targetPath, "branch", "--set-upstream-to=origin/master", "main")
	} else {
		_ = RunGitInDir(targetPath, "checkout", "-b", "main")
	}
	return nil
}

// Scaffold initialises a git repo and sets the remote origin without fetching
// or requiring access to the remote. Useful for bootstrapping repos you cannot
// yet clone.
func Scaffold(targetPath, repo string) error {
	if err := validateRemoteURL(repo); err != nil {
		return err
	}
	slog.Debug("Scaffolding repository", "path", targetPath, "repo", repo)

	if err := RunGitInDir(targetPath, "init"); err != nil {
		return EnhanceGitError(err)
	}
	if err := RunGitInDir(targetPath, "remote", "add", "origin", repo); err != nil {
		return EnhanceGitError(err)
	}
	return nil
}

// IsCloned returns true if the given path contains a .git directory.
func IsCloned(targetPath string) bool {
	gitDir := filepath.Join(targetPath, ".git")
	_, err := os.Stat(gitDir)
	return err == nil
}

// IsNonEmptyDir returns true if the path is an existing non-empty directory.
func IsNonEmptyDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil || !fi.IsDir() {
		return false
	}
	entries, _ := os.ReadDir(path)
	return len(entries) > 0
}

// validateRemoteURL rejects HTTP(S) URLs — this project enforces SSH.
func validateRemoteURL(repo string) error {
	if strings.HasPrefix(repo, "http://") || strings.HasPrefix(repo, "https://") {
		return fmt.Errorf("insecure repository URL detected: %s. This project strictly enforces SSH for Git operations. Please update your system-definition.yaml to use SSH URLs (e.g., git@github.com:...) or local paths", repo)
	}
	return nil
}

// EnhanceGitError adds user-friendly hints to common git errors.
func EnhanceGitError(err error) error {
	if err == nil {
		return nil
	}
	errStr := err.Error()
	var msg strings.Builder
	msg.WriteString(errStr)

	if strings.Contains(errStr, "Permission denied (publickey)") {
		msg.WriteString("\n\n" + ui.ColorRed + "SSH authentication failed." + ui.ColorReset)
		msg.WriteString("\nGitHub rejected your SSH key. Please ensure it's added to your account.")

		if identityFile := sshutil.GetGitHubIdentityFile(); identityFile != "" {
			msg.WriteString(fmt.Sprintf("\nIdentityFile: %s", identityFile))
		}

		if runtime.GOOS == "windows" {
			msg.WriteString("\n\nWindows Tip: Ensure the 'OpenSSH Authentication Agent' service is running and your key is loaded (run 'make ssh' Option 2).")
			msg.WriteString("\nAlso verify that your 'core.sshCommand' is set to the System OpenSSH (run 'make ssh' Option 4).")
		}

		msg.WriteString("\n\nRun 'make doctor' for more diagnostics.")
	} else if strings.Contains(errStr, "Host key verification failed") {
		msg.WriteString("\n\n" + ui.ColorRed + "Host key verification failed." + ui.ColorReset)
		msg.WriteString("\nGitHub's host key is missing from your known_hosts file.")
		msg.WriteString("\nRun 'make ssh' Option 5 to fix this automatically.")
	}

	return errors.New(msg.String())
}

// RunGit executes a git command with the configured SSH settings.
func RunGit(args ...string) error {
	var stderr strings.Builder
	cmd := execCommand("git", args...)
	cmd.Env = gitEnv()

	slog.Debug("Running git", "args", args)
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		slog.Debug("Git error output", "stderr", stderr.String())
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(stderr.String()))
	}
	return nil
}

// RunGitInDir executes a git command inside the given directory.
func RunGitInDir(dir string, args ...string) error {
	var stderr strings.Builder
	cmd := execCommand("git", args...)
	cmd.Dir = dir
	cmd.Env = gitEnv()

	slog.Debug("Running git in dir", "dir", dir, "args", args)
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		slog.Debug("Git error output", "stderr", stderr.String())
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(stderr.String()))
	}
	return nil
}

// EnsureDir creates the directory (and parents) if it does not exist.
func EnsureDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0755)
	}
}

// gitEnv returns the current environment with GIT_SSH_COMMAND configured.
func gitEnv() []string {
	sshCmd := sshutil.GetConfiguredSSHCommand()

	if runtime.GOOS == "windows" {
		if filepath.IsAbs(sshCmd) {
			sshCmd = filepath.ToSlash(sshCmd)
			if strings.Contains(sshCmd, " ") {
				sshCmd = fmt.Sprintf("\"%s\"", sshCmd)
			}
		}
	} else {
		if strings.Contains(sshCmd, " ") && !strings.HasPrefix(sshCmd, "\"") {
			sshCmd = fmt.Sprintf("\"%s\"", sshCmd)
		}
	}

	env := os.Environ()
	env = append(env, "GIT_SSH_COMMAND="+sshCmd)
	var filtered []string
	for _, e := range env {
		upper := strings.ToUpper(e)
		if !strings.HasPrefix(upper, "GIT_SSH=") && !strings.HasPrefix(upper, "GIT_SSH_VARIANT=") {
			filtered = append(filtered, e)
		}
	}
	return filtered
}
