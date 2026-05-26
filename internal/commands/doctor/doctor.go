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

func Run() bool {
	slog.Info("Starting Workspace Doctor...")
	slog.Info("Environment Info", "os", runtime.GOOS, "arch", runtime.GOARCH)
	if isWSL() {
		slog.Info("Environment Info", "wsl", true)
	}

	checkGit()
	checkGo()
	sshIssues := checkSSH()
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
	return sshIssues
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
	if Run() {
		PrintSSHSetupInstructions()
	}
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

func checkGo() {
	slog.Info("Checking Go installation...")
	cmd := exec.Command("go", "version")
	output, err := cmd.Output()
	if err != nil {
		slog.Error("Go is not installed or not in PATH", "error", err)
		return
	}
	versionLine := strings.TrimSpace(string(output))
	slog.Info("Go version", "version", versionLine)

	// Check if the version meets the requirement (e.g., 1.26 as per go.mod)
	// Output format is usually: go version go1.26.1 windows/amd64
	parts := strings.Fields(versionLine)
	if len(parts) >= 3 && strings.HasPrefix(parts[2], "go") {
		version := strings.TrimPrefix(parts[2], "go")
		if isVersionOlder(version, "1.26") {
			slog.Warn("Go version might be too old", "found", version, "required", "1.26")
			if isWSL() {
				slog.Info("Hint: On WSL, you can install the latest Go version using the following commands:")
				slog.Info("1. wget https://go.dev/dl/go1.26.1.linux-amd64.tar.gz")
				slog.Info("2. sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.26.1.linux-amd64.tar.gz")
				slog.Info("3. export PATH=$PATH:/usr/local/go/bin (add this to your ~/.bashrc)")
			}
		}
	}

	// Check GOTOOLCHAIN
	gtc := os.Getenv("GOTOOLCHAIN")
	if gtc == "" {
		gtc = "auto (default)"
	}
	slog.Info("GOTOOLCHAIN setting", "value", gtc)
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
	slog.Info("Checking SSH configuration...")
	hasIssues := false

	// 0. Check for structural issues
	issues := sshutil.DetectSSHIssues()
	for _, issue := range issues {
		slog.Error("SSH Structural Issue", "detail", issue)
		hasIssues = true
		if strings.Contains(issue, "is a directory") {
			slog.Info("Fix: Run '<cli> ssh' and select Option 6 (Cleanup broken SSH configurations) to resolve this automatically.")
		}
	}

	// 1. Check if ssh-agent is running
	slog.Info("Checking if ssh-agent is running...")
	keys, err := sshutil.GetLoadedKeys()
	if err != nil {
		hasIssues = true
		outStr := keys
		slog.Warn("ssh-agent might not be running or accessible", "error", outStr)
		if runtime.GOOS == "windows" {
			slog.Info("Hint: On Windows, the 'OpenSSH Authentication Agent' service is often disabled by default.")
			slog.Info("Fix (Admin PowerShell): Set-Service -Name ssh-agent -StartupType Manual; Start-Service ssh-agent")
			slog.Info("Or use '<cli> ssh' -> Option 5 (Check current SSH status) to try starting it automatically.")
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
			slog.Info("Or use '<cli> ssh' -> Option 5 (Check current SSH status) to try starting it automatically.")
		}
	} else {
		slog.Info("ssh-agent is responsive.")
	}

	// 2. Check loaded keys
	slog.Info("Checking loaded SSH keys...")
	if err != nil && !strings.Contains(keys, "The agent has no identities") {
		hasIssues = true
		slog.Warn("Could not check loaded keys", "output", keys)
	} else if strings.Contains(keys, "The agent has no identities") {
		// No keys is not necessarily an "issue" that requires the guide if everything else works,
		// but usually it is what users need help with if GitHub fails.
		// However, if GitHub check succeeds (step 3), then having no keys in agent might be fine
		// (e.g. using a key from config without agent).
		slog.Warn("No SSH keys found in agent", "output", keys)
		slog.Info("Hint: Run '<cli> ssh' -> Option 2 to add a key to the agent.")
	} else if err != nil {
		hasIssues = true
		slog.Warn("Error checking loaded keys", "error", err)
	} else {
		slog.Info("Loaded SSH keys", "keys", keys)
	}

	// 3. Test GitHub connectivity
	slog.Info("Testing connectivity to GitHub...")
	success, output := sshutil.CheckGitHubConnectivity()
	if success {
		slog.Info("GitHub authentication successful!")
	} else {
		hasIssues = true
		slog.Error("GitHub authentication failed", "output", strings.TrimSpace(output))
		slog.Info("Hint: Ensure your public key is added to your GitHub account settings and check your ~/.ssh/config if you use custom keys.")
		if strings.Contains(output, "Permission denied") {
			slog.Info("Hint: If your key has a passphrase, it MUST be added to the ssh-agent.")
		}
	}

	return hasIssues
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
	fmt.Println("   <cli> ssh")
	system.PrintCLINote()
}
