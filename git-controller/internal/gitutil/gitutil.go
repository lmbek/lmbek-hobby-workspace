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

func ProcessGitComponent(baseDir, name, repo string) error {
	targetPath := filepath.Join(baseDir, name)

	if repo == "" || strings.Contains(repo, "@company") {
		return nil
	}

	// Enforce SSH for remote repositories
	if strings.HasPrefix(repo, "http://") || strings.HasPrefix(repo, "https://") {
		return fmt.Errorf("insecure repository URL detected: %s. This project strictly enforces SSH for Git operations. Please update your git-repositories/system-definition.yaml to use SSH URLs (e.g., git@github.com:...) or local paths", repo)
	}

	gitDir := filepath.Join(targetPath, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		slog.Debug("Initializing/Cloning component", "name", name, "repo", repo)

		// Check if target directory exists and is not empty
		if fi, err := os.Stat(targetPath); err == nil && fi.IsDir() {
			entries, _ := os.ReadDir(targetPath)
			if len(entries) > 0 {
				slog.Debug("Target directory exists and is not empty, initializing git instead of cloning", "path", targetPath)
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
				} else {
					// Fallback: Check if origin/master exists
					if err := RunGitInDir(targetPath, "reset", "--soft", "origin/master"); err == nil {
						_ = RunGitInDir(targetPath, "checkout", "-b", "main")
						_ = RunGitInDir(targetPath, "branch", "--set-upstream-to=origin/master", "main")
					} else {
						_ = RunGitInDir(targetPath, "checkout", "-b", "main")
					}
				}
				return nil
			}
		}

		if err := RunGit("clone", repo, targetPath); err != nil {
			return EnhanceGitError(err)
		}
	} else {
		slog.Debug("Updating component", "name", name)
		if err := RunGitInDir(targetPath, "fetch", "--all"); err != nil {
			return EnhanceGitError(err)
		}
		if err := RunGitInDir(targetPath, "pull"); err != nil {
			return EnhanceGitError(err)
		}
	}
	return nil
}

func EnhanceGitError(err error) error {
	errStr := err.Error()
	var msg strings.Builder
	msg.WriteString(errStr)

	if strings.Contains(errStr, "Permission denied (publickey)") {
		msg.WriteString("\n\n" + ui.ColorRed + "SSH authentication failed." + ui.ColorReset)
		msg.WriteString("\nGitHub rejected your SSH key. Please ensure it's added to your account.")

		if identityFile := sshutil.GetGitHubIdentityFile(); identityFile != "" {
			msg.WriteString(fmt.Sprintf("\nIdentityFile: %s", identityFile))
			pubKey, _, err := sshutil.GetPublicKeyContent(identityFile)
			if err == nil {
				msg.WriteString(fmt.Sprintf("\nPublic Key:\n%s", strings.TrimSpace(pubKey)))
			}
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

func RunGit(args ...string) error {
	var stderr strings.Builder
	cmd := execCommand("git", args...)

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
	var filteredEnv []string
	for _, e := range env {
		upper := strings.ToUpper(e)
		if !strings.HasPrefix(upper, "GIT_SSH=") && !strings.HasPrefix(upper, "GIT_SSH_VARIANT=") {
			filteredEnv = append(filteredEnv, e)
		}
	}
	cmd.Env = filteredEnv

	slog.Debug("Running git", "args", args, "ssh", sshCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		slog.Debug("Git error output", "stderr", stderr.String())
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(stderr.String()))
	}
	return nil
}

func RunGitInDir(dir string, args ...string) error {
	var stderr strings.Builder
	cmd := execCommand("git", args...)
	cmd.Dir = dir

	sshCmd := sshutil.GetConfiguredSSHCommand()

	if runtime.GOOS == "windows" && filepath.IsAbs(sshCmd) {
		sshCmd = filepath.ToSlash(sshCmd)
		if strings.Contains(sshCmd, " ") {
			sshCmd = fmt.Sprintf("\"%s\"", sshCmd)
		}
	}

	env := os.Environ()
	env = append(env, "GIT_SSH_COMMAND="+sshCmd)
	var filteredEnv []string
	for _, e := range env {
		upper := strings.ToUpper(e)
		if !strings.HasPrefix(upper, "GIT_SSH=") && !strings.HasPrefix(upper, "GIT_SSH_VARIANT=") {
			filteredEnv = append(filteredEnv, e)
		}
	}
	cmd.Env = filteredEnv

	slog.Debug("Running git in dir", "dir", dir, "args", args, "ssh", sshCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		slog.Debug("Git error output", "stderr", stderr.String())
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(stderr.String()))
	}
	return nil
}

func EnsureDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0755)
	}
}

func RunHook(command string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = execCommand("powershell", "-Command", command)
	} else {
		cmd = execCommand("sh", "-c", command)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
