package doctor

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

func Run() {
	slog.Info("Starting Workspace Doctor...")
	slog.Info("Environment Info", "os", runtime.GOOS, "arch", runtime.GOARCH)
	if isWSL() {
		slog.Info("Environment Info", "wsl", true)
	}

	checkGit()
	checkSSH()
	checkDocker()

	// 4. Check for conflicting GIT_SSH variables
	if val := os.Getenv("GIT_SSH"); val != "" {
		slog.Warn("Conflicting environment variable detected", "GIT_SSH", val)
		slog.Info("Hint: This can override our GIT_SSH_COMMAND. Consider unsetting it.")
	}
	if val := os.Getenv("GIT_SSH_COMMAND"); val != "" {
		slog.Info("Current GIT_SSH_COMMAND in environment", "value", val)
	}

	slog.Info("Doctor report complete.")
}

func isWSL() bool {
	if runtime.GOOS != "linux" {
		return false
	}
	// WSL typically has "microsoft" in /proc/version
	data, err := os.ReadFile("/proc/version")
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(data)), "microsoft")
}

func RunFull() {
	Run()
	PrintSSHSetupInstructions()
}

func checkGit() {
	slog.Info("Checking Git installation...")
	cmd := exec.Command("git", "--version")
	output, err := cmd.Output()
	if err != nil {
		slog.Error("Git is not installed or not in PATH", "error", err)
	} else {
		slog.Info("Git version", "version", strings.TrimSpace(string(output)))
	}
}

func checkSSH() {
	slog.Info("Checking SSH configuration...")

	// 0. Check for structural issues
	issues := sshutil.DetectSSHIssues()
	for _, issue := range issues {
		slog.Error("SSH Structural Issue", "detail", issue)
		if strings.Contains(issue, "is a directory") {
			slog.Info("Fix: Run '<cli> ssh-setup' and select Option 4 to resolve this automatically.")
		}
	}

	// 1. Check if ssh-agent is running
	slog.Info("Checking if ssh-agent is running...")
	keys, err := sshutil.GetLoadedKeys()
	if err != nil {
		outStr := keys
		slog.Warn("ssh-agent might not be running or accessible", "error", outStr)
		if runtime.GOOS == "windows" {
			slog.Info("Hint: On Windows, the 'OpenSSH Authentication Agent' service is often disabled by default.")
			slog.Info("Fix (Admin PowerShell): Set-Service -Name ssh-agent -StartupType Manual; Start-Service ssh-agent")
			slog.Info("Or use '<cli> ssh-setup' -> Option 4 to try starting it automatically.")
		} else {
			slog.Info("Hint: On Linux/WSL, ensure ssh-agent is running. Start it with: eval \"$(ssh-agent -s)\"")
			slog.Info("To persist, add the following snippet to your ~/.bashrc or ~/.zshrc:")
			fmt.Println(`   if [ -z "$SSH_AUTH_SOCK" ]; then
     SOCK=$(ls /tmp/ssh-*/agent.* 2>/dev/null | head -n 1)
     if [ -n "$SOCK" ]; then
       export SSH_AUTH_SOCK=$SOCK
     else
       eval $(ssh-agent -s)
     fi
   fi`)
			slog.Info("Or use '<cli> ssh-setup' -> Option 4 to try starting it automatically.")
		}
	} else {
		slog.Info("ssh-agent is responsive.")
	}

	// 2. Check loaded keys
	slog.Info("Checking loaded SSH keys...")
	if err != nil && !strings.Contains(keys, "The agent has no identities") {
		slog.Warn("Could not check loaded keys", "output", keys)
	} else if strings.Contains(keys, "The agent has no identities") || err != nil {
		slog.Warn("No SSH keys found in agent", "output", keys)
		slog.Info("Hint: Run '<cli> ssh-setup' -> Option 2 to add a key to the agent.")
	} else {
		slog.Info("Loaded SSH keys", "keys", keys)
	}

	// 3. Test GitHub connectivity
	slog.Info("Testing connectivity to GitHub...")
	success, output := sshutil.CheckGitHubConnectivity()
	if success {
		slog.Info("GitHub authentication successful!")
	} else {
		slog.Error("GitHub authentication failed", "output", strings.TrimSpace(output))
		slog.Info("Hint: Ensure your public key is added to your GitHub account settings and check your ~/.ssh/config if you use custom keys.")
		if strings.Contains(output, "Permission denied") {
			slog.Info("Hint: If your key has a passphrase, it MUST be added to the ssh-agent.")
		}
	}
}

func checkDocker() {
	slog.Info("Checking Docker installation...")
	cmd := exec.Command("docker", "--version")
	output, err := cmd.Output()
	if err != nil {
		slog.Error("Docker is not installed or not in PATH", "error", err)
	} else {
		slog.Info("Docker version", "version", strings.TrimSpace(string(output)))
	}

	cmd = exec.Command("docker-compose", "--version")
	output, err = cmd.Output()
	if err != nil {
		slog.Warn("docker-compose is not installed or not in PATH", "error", err)
	} else {
		slog.Info("Docker Compose version", "version", strings.TrimSpace(string(output)))
	}
}

func PrintSSHSetupInstructions() {
	keyName := "id_ed25519"
	identityFile := sshutil.GetGitHubIdentityFile()
	if identityFile != "" {
		keyName = filepath.Base(identityFile)
	}

	fmt.Println("\nSSH Setup Guide:")
	fmt.Println("\n--- Manual Configuration ---")
	fmt.Println("1. Generate a new SSH key (if missing):")
	fmt.Printf("   ssh-keygen -t ed25519 -C \"your_email@example.com\" -f ~/.ssh/%s\n", keyName)
	fmt.Println("2. Start the ssh-agent:")
	if runtime.GOOS == "windows" {
		fmt.Println("   PowerShell (Admin): Set-Service -Name ssh-agent -StartupType Manual; Start-Service ssh-agent")
	} else {
		fmt.Println("   eval \"$(ssh-agent -s)\"")
	}
	fmt.Println("3. Add your SSH key to the agent:")
	fmt.Printf("   ssh-add ~/.ssh/%s\n", keyName)
	fmt.Println("4. Add the public key to your GitHub account:")
	if runtime.GOOS == "windows" {
		fmt.Printf("   Get-Content ~/.ssh/%s.pub | clip\n", keyName)
	} else {
		fmt.Printf("   cat ~/.ssh/%s.pub\n", keyName)
	}
	fmt.Println("   GitHub -> Settings -> SSH and GPG keys -> New SSH key")
	fmt.Println("5. (Optional) Configure ~/.ssh/config for host-specific settings:")
	fmt.Println("   Host github.com")
	fmt.Println("     HostName github.com")
	fmt.Println("     User git")
	fmt.Printf("     IdentityFile ~/.ssh/%s\n", keyName)
	fmt.Println("     AddKeysToAgent yes")
	fmt.Println("     IdentitiesOnly yes")

	fmt.Println("\n--- Automated Configuration ---")
	fmt.Println("   You can use the built-in automated tool for these steps:")
	fmt.Println("   <cli> ssh-setup")
	system.PrintCLINote()
}
