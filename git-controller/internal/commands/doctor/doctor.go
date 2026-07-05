package doctor

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"workspace/git-controller/internal/sshutil"
	"workspace/git-controller/internal/system"
	"workspace/git-controller/internal/ui"
)

func Run() error {
	if runCheck() {
		PrintSSHSetupInstructions()
	}
	return nil
}

func runCheck() bool {
	ui.Header("Workspace Doctor")
	slog.Debug("Environment Info", "os", runtime.GOOS, "arch", runtime.GOARCH)
	if isWSL() {
		slog.Debug("Environment Info", "wsl", true)
	}

	checkGit()
	checkGo()
	sshIssues := checkSSH()
	checkDocker()

	// 4. Check for conflicting GIT_SSH variables
	if val := os.Getenv("GIT_SSH"); val != "" {
		ui.Warn("Conflicting environment variable detected: GIT_SSH=%s", val)
		ui.Info("Hint: This can override our GIT_SSH_COMMAND. Consider unsetting it.")
	}
	if val := os.Getenv("GIT_SSH_COMMAND"); val != "" {
		slog.Debug("Current GIT_SSH_COMMAND", "value", val)
	}

	ui.Success("Doctor report complete.")
	return sshIssues
}

// isWSL returns true if the current environment appears to be Windows Subsystem for Linux.
// It checks for the presence of "microsoft" in /proc/version on linux systems.
func isWSL() bool {
	if runtime.GOOS != "linux" {
		return false
	}
	data, err := os.ReadFile("/proc/version")
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(data)), "microsoft")
}

func checkGit() {
	ui.Info("Checking Git installation...")
	cmd := exec.Command("git", "--version")
	output, err := cmd.Output()
	if err != nil {
		ui.Error("Git is not installed or not in PATH")
	} else {
		slog.Debug("Git version", "version", strings.TrimSpace(string(output)))
		ui.Success("Git installed")
	}
}

func checkGo() {
	ui.Info("Checking Go installation...")
	cmd := exec.Command("go", "version")
	output, err := cmd.Output()
	if err != nil {
		ui.Error("Go is not installed or not in PATH")
		return
	}
	versionLine := strings.TrimSpace(string(output))
	slog.Debug("Go version", "version", versionLine)

	// Check if the version meets the requirement (e.g., 1.26)
	parts := strings.Fields(versionLine)
	if len(parts) >= 3 && strings.HasPrefix(parts[2], "go") {
		version := strings.TrimPrefix(parts[2], "go")
		if isVersionOlder(version, "1.26") {
			ui.Warn("Go version might be too old: found %s, required 1.26+", version)
		} else {
			ui.Success("Go version OK (%s)", version)
		}
	}
}

func isVersionOlder(current, required string) bool {
	// Simple string comparison for go versions like "1.26.1" vs "1.26"
	// This is a naive implementation but works for major/minor versions.
	cParts := strings.Split(current, ".")
	rParts := strings.Split(required, ".")
	for i := 0; i < len(rParts) && i < len(cParts); i++ {
		if cParts[i] < rParts[i] {
			return true
		}
		if cParts[i] > rParts[i] {
			return false
		}
	}
	return len(cParts) < len(rParts)
}

func checkSSH() bool {
	ui.Info("Checking SSH configuration...")
	hasIssues := false

	// 0. Check for structural issues
	issues := sshutil.DetectSSHIssues()
	for _, issue := range issues {
		ui.Error("SSH Structural Issue: %s", issue)
		hasIssues = true
		if strings.Contains(issue, "is a directory") {
			ui.Info("Fix: Run '%s ssh' and select Option 6 (Cleanup broken SSH configurations) to resolve this automatically.", system.CLIName)
		}
	}

	// 1. Check if ssh-agent is running
	ui.Info("Checking if ssh-agent is running...")
	keys, err := sshutil.GetLoadedKeys()
	if err != nil {
		hasIssues = true
		ui.Warn("ssh-agent might not be running or accessible")
		if runtime.GOOS == "windows" {
			ui.Info("Hint: On Windows, the 'OpenSSH Authentication Agent' service is often disabled by default.")
			ui.Info("Fix (Admin PowerShell): Set-Service -Name ssh-agent -StartupType Manual; Start-Service ssh-agent")
		} else {
			ui.Info("Hint: On Linux/WSL, ensure ssh-agent is running. Start it with: eval \"$(ssh-agent -s)\"")
		}
	} else {
		ui.Success("ssh-agent is responsive")
	}

	// 2. Check loaded keys
	if err == nil {
		if strings.Contains(keys, "The agent has no identities") {
			ui.Warn("No SSH keys found in agent")
			ui.Info("Hint: Run '%s ssh' -> Option 2 to add a key to the agent.", system.CLIName)
		} else {
			ui.Success("SSH keys loaded in agent")
		}
	}

	// 3. Test GitHub connectivity
	ui.Info("Testing connectivity to GitHub...")
	success, output := sshutil.CheckGitHubConnectivity()
	if success {
		ui.Success("GitHub authentication successful!")
	} else {
		hasIssues = true
		ui.Error("GitHub authentication failed")
		slog.Debug("GitHub output", "output", strings.TrimSpace(output))
		ui.Info("Hint: Ensure your public key is added to your GitHub account settings.")
		if strings.Contains(output, "Permission denied") {
			ui.Info("Hint: If your key has a passphrase, it MUST be added to the ssh-agent.")
		}
	}

	return hasIssues
}

func checkDocker() {
	ui.Info("Checking Docker installation...")
	cmd := exec.Command("docker", "--version")
	output, err := cmd.Output()
	if err != nil {
		ui.Error("Docker is not installed or not in PATH")
	} else {
		ui.Success("Docker installed (%s)", strings.TrimSpace(string(output)))
	}

	cmd = exec.Command("docker-compose", "--version")
	output, err = cmd.Output()
	if err != nil {
		ui.Warn("docker-compose is not installed or not in PATH")
	} else {
		ui.Success("Docker Compose installed (%s)", strings.TrimSpace(string(output)))
	}
}

func PrintSSHSetupInstructions() {
	keyName := "id_ed25519.private"
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

	// Handle .private to .public conversion for instructions
	pubName := strings.TrimSuffix(keyName, ".private") + ".public"
	if runtime.GOOS == "windows" {
		fmt.Printf("   Get-Content ~/.ssh/%s | clip\n", pubName)
	} else {
		fmt.Printf("   cat ~/.ssh/%s\n", pubName)
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
	fmt.Printf("   %s ssh\n", system.CLIName)
}
