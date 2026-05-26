package gitutil

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"workspace-controller/internal/sshutil"
)

func ProcessGitComponent(baseDir, name, repo, version string) error {
	targetPath := filepath.Join(baseDir, name)
	if filepath.Base(baseDir) == name {
		targetPath = baseDir
	}

	if repo == "" || strings.Contains(repo, "@company") {
		return nil
	}

	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		slog.Info("Sync component", "name", name, "action", "clone")
		if err := RunGit("clone", repo, targetPath); err != nil {
			return fmt.Errorf("error cloning %s: %w", repo, err)
		}
	} else {
		slog.Info("Sync component", "name", name, "action", "update")
		if err := RunGitInDir(targetPath, "fetch", "--all"); err != nil {
			return fmt.Errorf("error fetching %s: %w", name, err)
		}
		if err := RunGitInDir(targetPath, "pull"); err != nil {
			return fmt.Errorf("error pulling %s: %w", name, err)
		}
	}
	return nil
}

func HandleGitError(err error) {
	errStr := err.Error()
	if strings.Contains(errStr, "Permission denied (publickey)") {
		fmt.Println("\nERROR: SSH authentication failed.")
		fmt.Println("GitHub rejected your SSH key. Please ensure it's added to your account.")

		if identityFile := sshutil.GetGitHubIdentityFile(); identityFile != "" {
			fmt.Printf("\nIdentityFile: %s\n", identityFile)
			pubKey, _, err := sshutil.GetPublicKeyContent(identityFile)
			if err == nil {
				fmt.Printf("Public Key:\n%s\n", strings.TrimSpace(pubKey))
			}
		}

		if runtime.GOOS == "windows" {
			fmt.Println("\nWindows Tip: Ensure the 'OpenSSH Authentication Agent' service is running and your key is loaded (run 'make ssh' Option 2).")
			fmt.Println("Also verify that your 'core.sshCommand' is set to the System OpenSSH (run 'make ssh' Option 4 - Configure git to use Windows OpenSSH).")
		}

		fmt.Println("\nRun 'make doctor' for more diagnostics.")
		fmt.Println("You can also run 'make ssh' (alias for ssh-setup) to manage your keys.")
		os.Exit(1)
	} else if strings.Contains(errStr, "Host key verification failed") {
		fmt.Println("\nERROR: Host key verification failed.")
		fmt.Println("GitHub's host key is missing from your known_hosts file.")
		fmt.Println("Run 'make ssh' Option 5 (Check current SSH status) to fix this automatically.")
		os.Exit(1)
	} else {
		slog.Error("Git operation failed", "error", err)
	}
}

func RunGit(args ...string) error {
	var stderr strings.Builder
	cmd := exec.Command("git", args...)

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
	err := cmd.Run()
	if err != nil {
		fmt.Fprint(os.Stderr, stderr.String())
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(stderr.String()))
	}
	return nil
}

func RunGitInDir(dir string, args ...string) error {
	var stderr strings.Builder
	cmd := exec.Command("git", args...)
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
	err := cmd.Run()
	if err != nil {
		fmt.Fprint(os.Stderr, stderr.String())
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
		cmd = exec.Command("powershell", "-Command", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
