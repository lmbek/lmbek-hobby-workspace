package initcmd

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"workspace-controller/internal/sshutil"
	"workspace-controller/internal/system"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func Run() {
	fmt.Println("SYSTEM INITIALIZATION")
	fmt.Println("=====================")

	fmt.Println("\n[0] Pre-flight checks:")

	_ = sshutil.StartAgent()

	// Ensure GitHub is in known_hosts if it's missing (helps avoid interactive prompts)
	if runtime.GOOS == "windows" {
		home, _ := os.UserHomeDir()
		sshDir := filepath.Join(home, ".ssh")
		knownHosts := filepath.Join(sshDir, "known_hosts")
		if _, err := os.Stat(knownHosts); os.IsNotExist(err) || !containsHost(knownHosts, "github.com") {
			slog.Info("Adding github.com to known_hosts...")
			os.MkdirAll(sshDir, 0700)
			// Use more robust PowerShell command for appending
			exec.Command("powershell", "-Command", "ssh-keyscan github.com | Out-File -FilePath \""+knownHosts+"\" -Append -Encoding ascii").Run()
		}
	}

	if err := checkGitHubConnectivity(); err != nil {
		fmt.Printf("\nERROR: GitHub connectivity check failed.\n%v\n", err)
		fmt.Println("\nRun 'make ssh' to set up your environment.")
		os.Exit(1)
	}
	slog.Info("GitHub connectivity verified")

	sys, err := system.LoadDefinition("system/system-definition.yaml")
	if err != nil {
		slog.Error("Error loading system definition", "error", err)
		os.Exit(1)
	}

	// 1. Plan & Analyze
	fmt.Println("\n[1] Planning & Analysis:")
	planCount := len(sys.Services)
	if sys.Infrastructure != nil {
		planCount++
	}
	if sys.Tools != nil {
		planCount++
	}
	slog.Info("System plan created", "components", planCount)

	// 2. Materialize
	fmt.Println("\n[2] Materializing Workspace:")
	workspaceDir := getEnv("SERVICES_DIR", "../workspace/services")
	infraDir := getEnv("INFRA_DIR", "../workspace/infrastructure")
	toolsDir := getEnv("TOOLS_DIR", "../workspace/tools")

	// Ensure absolute paths for base directories to avoid confusion with git -C or Dir
	if abs, err := filepath.Abs(workspaceDir); err == nil {
		workspaceDir = abs
	}
	if abs, err := filepath.Abs(infraDir); err == nil {
		infraDir = abs
	}
	if abs, err := filepath.Abs(toolsDir); err == nil {
		toolsDir = abs
	}

	ensureDir(workspaceDir)
	ensureDir(infraDir)
	ensureDir(toolsDir)

	for name, svc := range sys.Services {
		if err := processGitComponent(workspaceDir, name, svc.Repository, svc.Version); err != nil {
			handleGitError(err)
			os.Exit(1)
		}
	}
	if sys.Infrastructure != nil {
		if err := processGitComponent(infraDir, "infrastructure", sys.Infrastructure.Repository, sys.Infrastructure.Version); err != nil {
			handleGitError(err)
			os.Exit(1)
		}
	}
	if sys.Tools != nil {
		if err := processGitComponent(toolsDir, "tools", sys.Tools.Repository, sys.Tools.Version); err != nil {
			handleGitError(err)
			os.Exit(1)
		}
	}

	// 3. Running Hooks
	if len(sys.Hooks.PostSync) > 0 {
		fmt.Println("\n[3] Running Hooks:")
		for _, hook := range sys.Hooks.PostSync {
			if err := runHook(hook); err != nil {
				slog.Error("Hook failed", "command", hook, "error", err)
			}
		}
	}

	fmt.Println("\nInitialization finished successfully!")
	fmt.Println("Next step: Run '<cli> up' to start the system.")
	system.PrintCLINote()
}

func containsHost(filepath, host string) bool {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return false
	}
	return strings.Contains(string(data), host)
}

func runHook(command string) error {
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

func processGitComponent(baseDir, name, repo, version string) error {
	targetPath := filepath.Join(baseDir, name)
	if filepath.Base(baseDir) == name {
		targetPath = baseDir
	}

	if repo == "" || strings.Contains(repo, "@company") {
		return nil
	}

	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		slog.Info("Sync component", "name", name, "action", "clone")
		if err := runGit("clone", repo, targetPath); err != nil {
			return fmt.Errorf("error cloning %s: %w", repo, err)
		}
	} else {
		slog.Info("Sync component", "name", name, "action", "update")
		if err := runGitInDir(targetPath, "fetch", "--all"); err != nil {
			return fmt.Errorf("error fetching %s: %w", name, err)
		}
		if err := runGitInDir(targetPath, "pull"); err != nil {
			return fmt.Errorf("error pulling %s: %w", name, err)
		}
	}
	return nil
}

func handleGitError(err error) {
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
			fmt.Println("Also verify that your 'core.sshCommand' is set to the System OpenSSH (run 'make ssh' Option 4).")
		}

		fmt.Println("\nRun 'make doctor' for more diagnostics.")
		os.Exit(1)
	} else if strings.Contains(errStr, "Host key verification failed") {
		fmt.Println("\nERROR: Host key verification failed.")
		fmt.Println("GitHub's host key is missing from your known_hosts file.")
		if runtime.GOOS == "windows" {
			fmt.Println("Run 'make ssh' Option 4 (Configure git to use Windows OpenSSH) to fix this.")
		} else {
			fmt.Println("Run 'make ssh' Option 5 (Check current SSH status) to fix this.")
		}
		os.Exit(1)
	} else {
		slog.Error("Git operation failed", "error", err)
	}
}

func ensureDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0755)
	}
}

func runGit(args ...string) error {
	var stderr strings.Builder
	cmd := exec.Command("git", args...)

	// Explicitly set GIT_SSH_COMMAND to ensure Git uses the same SSH as our pre-flight check
	sshCmd := sshutil.GetConfiguredSSHCommand()

	// Ensure the path uses forward slashes for Git's GIT_SSH_COMMAND on Windows
	// Git for Windows handles forward slashes correctly in shell environment variables.
	if runtime.GOOS == "windows" {
		if filepath.IsAbs(sshCmd) {
			sshCmd = filepath.ToSlash(sshCmd)
			// If there are spaces, quote it.
			if strings.Contains(sshCmd, " ") {
				sshCmd = fmt.Sprintf("\"%s\"", sshCmd)
			}
		}
	} else {
		// On Linux/WSL, we just want to ensure it's a valid command or path.
		// If it's a path with spaces, it might still need quoting for the shell.
		if strings.Contains(sshCmd, " ") && !strings.HasPrefix(sshCmd, "\"") {
			sshCmd = fmt.Sprintf("\"%s\"", sshCmd)
		}
	}

	env := os.Environ()
	env = append(env, "GIT_SSH_COMMAND="+sshCmd)
	// Clear potentially conflicting variables (GIT_SSH, GIT_SSH_VARIANT)
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
		// Print stderr so user can see it
		fmt.Fprint(os.Stderr, stderr.String())
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(stderr.String()))
	}
	return nil
}

func runGitInDir(dir string, args ...string) error {
	var stderr strings.Builder
	cmd := exec.Command("git", args...)
	cmd.Dir = dir

	// Explicitly set GIT_SSH_COMMAND to ensure Git uses the same SSH as our pre-flight check
	sshCmd := sshutil.GetConfiguredSSHCommand()

	// Ensure the path uses forward slashes for Git's GIT_SSH_COMMAND on Windows
	// Git for Windows handles forward slashes correctly in shell environment variables.
	if runtime.GOOS == "windows" && filepath.IsAbs(sshCmd) {
		sshCmd = filepath.ToSlash(sshCmd)
		// If there are spaces, quote it.
		if strings.Contains(sshCmd, " ") {
			sshCmd = fmt.Sprintf("\"%s\"", sshCmd)
		}
	}

	env := os.Environ()
	env = append(env, "GIT_SSH_COMMAND="+sshCmd)
	// Clear potentially conflicting variables (GIT_SSH, GIT_SSH_VARIANT)
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
		// Print stderr so user can see it
		fmt.Fprint(os.Stderr, stderr.String())
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(stderr.String()))
	}
	return nil
}

func checkGitHubConnectivity() error {
	success, output := sshutil.CheckGitHubConnectivityNonInteractive()
	if success {
		return nil
	}

	// If it fails, maybe it's because of a passphrase.
	// In that case, we should check if the key is in the agent.
	if strings.Contains(output, "Permission denied") || strings.Contains(output, "passphrase") {
		// Provide a more helpful error if it looks like a passphrase issue
		return fmt.Errorf("GitHub authentication failed (Permission Denied).\nThis often means your SSH key requires a passphrase but isn't added to the agent.\nTry running 'make ssh' and choose Option 2 to add your key.")
	}

	return fmt.Errorf("GitHub authentication failed: %s", strings.TrimSpace(output))
}
