package sshhelper

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"workspace/git-controller/internal/sshutil"
	"workspace/git-controller/internal/ui"
)

// Run starts the interactive SSH helper wizard.
func Run() error {
	ui.Header("SSH Helper Tool")

	fmt.Println("\nSSH Setup & Diagnostic Options:")
	ui.Info("1. Generate a new SSH key")
	ui.Info("2. Add an existing key to the agent")
	ui.Info("3. Configure ~/.ssh/config for GitHub")
	if runtime.GOOS == "windows" {
		ui.Info("4. Configure git to use Windows OpenSSH (Windows only)")
	}
	ui.Info("5. Check current SSH status")
	ui.Info("6. Cleanup broken SSH configurations")
	fmt.Print("\nSelect an option (1-6): ")

	reader := bufio.NewReader(os.Stdin)
	option, _ := reader.ReadString('\n')
	option = strings.TrimSpace(option)

	switch option {
	case "1":
		generateNewKey(reader)
	case "2":
		addExistingKey(reader)
	case "3":
		configureSSHConfig(reader)
	case "4":
		if runtime.GOOS == "windows" {
			configureGitSSH(reader)
		} else {
			ui.Error("Invalid option selected")
		}
	case "5":
		checkStatus(reader)
	case "6":
		cleanupConfigs(reader)
	default:
		ui.Error("Invalid option selected")
	}
	return nil
}

func configureGitSSH(reader *bufio.Reader) {
	if runtime.GOOS != "windows" {
		fmt.Println("This option is only relevant for Windows systems.")
		return
	}

	sshPath := "C:/Windows/System32/OpenSSH/ssh.exe"
	fmt.Printf("\nThis will set 'git config --global core.sshCommand' to: %s\n", sshPath)
	fmt.Print("Do you want to proceed? (y/n): ")
	answer, _ := reader.ReadString('\n')
	if strings.ToLower(strings.TrimSpace(answer)) != "y" {
		ui.Info("Operation cancelled.")
		return
	}

	exec.Command("git", "config", "--global", "--unset-all", "core.sshcommand").Run()

	cmd := exec.Command("git", "config", "--global", "core.sshCommand", sshPath)
	if err := cmd.Run(); err != nil {
		ui.Error("Failed to set git config: %v", err)
		return
	}

	ui.Success("Successfully configured git to use Windows OpenSSH.")
}

func configureSSHConfig(reader *bufio.Reader) {
	home, err := os.UserHomeDir()
	if err != nil {
		ui.Error("Could not determine home directory: %v", err)
		return
	}

	sshDir := filepath.Join(home, ".ssh")
	configPath := filepath.Join(sshDir, "config")

	if _, err := os.Stat(sshDir); os.IsNotExist(err) {
		if err := os.MkdirAll(sshDir, 0700); err != nil {
			ui.Error("Could not create %s: %v", sshDir, err)
			return
		}
	}

	files, err := os.ReadDir(sshDir)
	var keys []string
	if err == nil {
		for _, file := range files {
			if !file.IsDir() {
				name := file.Name()
				if !strings.HasSuffix(name, ".public") && !strings.HasSuffix(name, ".pub") && !strings.HasSuffix(name, ".old") && !strings.HasSuffix(name, ".Identifier") && name != "known_hosts" && name != "known_hosts.old" && name != "config" && name != "authorized_keys" {
					keys = append(keys, name)
				}
			}
		}
	}

	var keyName string
	if len(keys) == 0 {
		ui.Warn("No private keys found in %s.", sshDir)
		fmt.Print("Enter the key filename you want to configure (e.g. id_ed25519): ")
		input, _ := reader.ReadString('\n')
		keyName = strings.TrimSpace(input)
		if keyName == "" {
			ui.Error("Key name cannot be empty.")
			return
		}
	} else {
		fmt.Println("\nAvailable keys:")
		for i, k := range keys {
			fmt.Printf("  %d) %s\n", i+1, k)
		}
		fmt.Printf("Select a key (1-%d) or type a filename: ", len(keys))
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		var selectedIndex int
		if _, err := fmt.Sscanf(input, "%d", &selectedIndex); err == nil && selectedIndex >= 1 && selectedIndex <= len(keys) {
			keyName = keys[selectedIndex-1]
		} else if input != "" {
			keyName = input
		} else {
			keyName = keys[0]
		}
	}

	var identityPath string
	if runtime.GOOS == "windows" {
		identityPath = "~/.ssh/" + keyName
	} else {
		identityPath = filepath.Join(home, ".ssh", keyName)
	}

	configEntry := fmt.Sprintf("\nHost github.com\n    HostName github.com\n    User git\n    IdentityFile %s\n    IdentitiesOnly yes\n", identityPath)

	fmt.Printf("\nGenerated configuration entry:\n%s\n", configEntry)
	fmt.Printf("This will be appended to: %s\n", configPath)
	fmt.Print("Proceed? (y/n): ")
	answer, _ := reader.ReadString('\n')
	if strings.ToLower(strings.TrimSpace(answer)) != "y" {
		ui.Info("Operation cancelled.")
		return
	}

	var existingContent string
	if data, err := os.ReadFile(configPath); err == nil {
		existingContent = string(data)
	}

	if strings.Contains(existingContent, "Host github.com") {
		ui.Warn("~/.ssh/config already contains a 'Host github.com' block.")
		fmt.Print("Do you want to append anyway? (y/n): ")
		ans, _ := reader.ReadString('\n')
		if strings.ToLower(strings.TrimSpace(ans)) != "y" {
			ui.Info("Operation cancelled. Please edit ~/.ssh/config manually.")
			return
		}
	}

	f, err := os.OpenFile(configPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		ui.Error("Failed to open %s: %v", configPath, err)
		return
	}
	defer f.Close()

	if _, err := f.WriteString(configEntry); err != nil {
		ui.Error("Failed to write to %s: %v", configPath, err)
		return
	}

	ui.Success("Successfully updated %s", configPath)
	ui.Info("Hint: Test your connection by choosing option 5 (Check current SSH status).")
}

func generateNewKey(reader *bufio.Reader) {
	fmt.Println("\nSelect key type:")
	fmt.Println("1. ED25519 (Recommended - modern, secure, compact)")
	fmt.Println("2. RSA 4096 (Legacy - for compatibility)")
	fmt.Print("Choice (1-2) [1]: ")

	keyTypeChoice, _ := reader.ReadString('\n')
	keyTypeChoice = strings.TrimSpace(keyTypeChoice)

	keyType := "ed25519"
	defaultFilename := "id_ed25519"
	bits := ""

	if keyTypeChoice == "2" {
		keyType = "rsa"
		defaultFilename = "id_rsa"
		bits = "4096"
	}

	fmt.Print("Enter email address (comment): ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	home, err := os.UserHomeDir()
	if err != nil {
		ui.Error("Could not determine home directory: %v", err)
		return
	}

	sshDir := filepath.Join(home, ".ssh")
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		ui.Error("Could not create ~/.ssh directory: %v", err)
		return
	}

	fmt.Printf("Enter filename for the key [%s]: ", defaultFilename)
	filename, _ := reader.ReadString('\n')
	filename = strings.TrimSpace(filename)
	if filename == "" {
		filename = defaultFilename
	}

	keyPath := filepath.Join(sshDir, filename)

	if _, err := os.Stat(keyPath); err == nil {
		ui.Warn("Key already exists at %s", keyPath)
		fmt.Print("Overwrite? (y/N): ")
		overwrite, _ := reader.ReadString('\n')
		if strings.ToLower(strings.TrimSpace(overwrite)) != "y" {
			ui.Info("Aborted.")
			return
		}
	}

	fmt.Print("Enter passphrase (press Enter for none): ")
	passphrase, _ := reader.ReadString('\n')
	passphrase = strings.TrimSpace(passphrase)

	args := []string{"-t", keyType, "-f", keyPath, "-N", passphrase}
	if email != "" {
		args = append(args, "-C", email)
	}
	if bits != "" {
		args = append(args, "-b", bits)
	}

	cmd := exec.Command("ssh-keygen", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		ui.Error("Failed to generate SSH key: %v", err)
		return
	}

	ui.Success("SSH key generated at %s", keyPath)

	pubKeyPath := keyPath + ".pub"
	if pubKey, err := os.ReadFile(pubKeyPath); err == nil {
		fmt.Println("\nPublic key:")
		fmt.Println(strings.TrimSpace(string(pubKey)))
		fmt.Println("\nCopy this public key and add it to your GitHub account:")
		fmt.Println("https://github.com/settings/keys")
	}

	fmt.Print("\nAdd this key to ssh-agent now? (y/n) [y]: ")
	addNow, _ := reader.ReadString('\n')
	addNow = strings.TrimSpace(addNow)
	if addNow == "" || strings.ToLower(addNow) == "y" {
		addKeyToAgent(keyPath)
	}
}

func addExistingKey(reader *bufio.Reader) {
	home, err := os.UserHomeDir()
	if err != nil {
		ui.Error("Could not determine home directory: %v", err)
		return
	}

	sshDir := filepath.Join(home, ".ssh")
	files, err := os.ReadDir(sshDir)
	if err != nil {
		ui.Error("Could not read ~/.ssh directory: %v", err)
		return
	}

	var keys []string
	for _, file := range files {
		if !file.IsDir() {
			name := file.Name()
			if !strings.HasSuffix(name, ".pub") && !strings.HasSuffix(name, ".old") && name != "known_hosts" && name != "known_hosts.old" && name != "config" && name != "authorized_keys" {
				keys = append(keys, name)
			}
		}
	}

	if len(keys) == 0 {
		ui.Warn("No private keys found in ~/.ssh/")
		return
	}

	fmt.Println("\nAvailable keys:")
	for i, key := range keys {
		fmt.Printf("%d. %s\n", i+1, key)
	}

	fmt.Print("Select key (1-", len(keys), "): ")
	choiceStr, _ := reader.ReadString('\n')
	var choice int
	fmt.Sscanf(strings.TrimSpace(choiceStr), "%d", &choice)

	if choice < 1 || choice > len(keys) {
		ui.Error("Invalid choice")
		return
	}

	keyPath := filepath.Join(sshDir, keys[choice-1])
	addKeyToAgent(keyPath)
}

func addKeyToAgent(keyPath string) {
	if runtime.GOOS == "windows" {
		out, err := exec.Command("powershell", "-NoProfile", "-Command", "Get-Service ssh-agent | Select-Object -ExpandProperty Status").Output()
		if err == nil {
			status := strings.TrimSpace(string(out))
			if status != "Running" {
				ui.Warn("Windows ssh-agent service is not running (status: %s)", status)
				ui.Info("To start it, open PowerShell as Administrator and run:")
				ui.Info("  Get-Service ssh-agent | Set-Service -StartupType Automatic")
				ui.Info("  Start-Service ssh-agent")
			}
		}
	}

	cmd := exec.Command("ssh-add", keyPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		ui.Error("Failed to add key to agent: %v", err)
		ui.Info("Make sure ssh-agent is running.")
		if runtime.GOOS != "windows" {
			ui.Info("Start it with: eval $(ssh-agent)")
		}
		return
	}

	ui.Success("Key added to ssh-agent successfully.")
}

func checkStatus(reader *bufio.Reader) {
	fmt.Println("\nRunning SSH diagnostics...")

	gitSSH := os.Getenv("GIT_SSH")
	gitSSHCommand := os.Getenv("GIT_SSH_COMMAND")
	if gitSSH != "" {
		ui.Warn("GIT_SSH environment variable is set: %s", gitSSH)
	}
	if gitSSHCommand != "" {
		ui.Info("GIT_SSH_COMMAND environment variable is set: %s", gitSSHCommand)
	}

	agentRunning := false
	if runtime.GOOS == "windows" {
		out, err := exec.Command("powershell", "-NoProfile", "-Command", "Get-Service ssh-agent | Select-Object -ExpandProperty Status").Output()
		if err == nil && strings.TrimSpace(string(out)) == "Running" {
			agentRunning = true
			ui.Success("Windows ssh-agent service is running")
		} else {
			ui.Warn("Windows ssh-agent service is not running")
		}
	} else {
		if os.Getenv("SSH_AUTH_SOCK") != "" {
			agentRunning = true
			ui.Success("SSH_AUTH_SOCK is set (%s)", os.Getenv("SSH_AUTH_SOCK"))
		} else {
			ui.Warn("SSH_AUTH_SOCK is not set - ssh-agent might not be running")
		}
	}

	if agentRunning {
		cmd := exec.Command("ssh-add", "-l")
		output, err := cmd.Output()
		if err != nil {
			ui.Warn("ssh-agent has no identities loaded (or ssh-add failed)")
		} else {
			ui.Success("ssh-agent loaded keys:")
			for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
				if line != "" {
					fmt.Printf("    %s\n", line)
				}
			}
		}
	}

	home, err := os.UserHomeDir()
	if err == nil {
		configPath := filepath.Join(home, ".ssh", "config")
		if _, err := os.Stat(configPath); err == nil {
			ui.Success("~/.ssh/config exists")
		} else {
			ui.Info("~/.ssh/config does not exist")
		}
	}

	fmt.Println("\nTesting connection to GitHub via SSH...")
	if ok, _ := sshutil.CheckGitHubConnectivity(); ok {
		ui.Success("GitHub SSH authentication successful!")
	} else {
		ui.Error("GitHub SSH authentication failed")
		ui.Info("Troubleshooting steps:")
		ui.Info("1. Ensure your public key is added to https://github.com/settings/keys")
		ui.Info("2. Ensure your private key is added to ssh-agent: ssh-add ~/.ssh/<your_key>")
		ui.Info("3. Test manually: ssh -Tv git@github.com")
	}
}

func cleanupConfigs(reader *bufio.Reader) {
	fmt.Println("\nSSH Configuration Cleanup")
	ui.Warn("This will inspect and optionally remove problematic settings.")

	cmd := exec.Command("git", "config", "--global", "--get", "core.sshcommand")
	out, err := cmd.Output()
	if err == nil && len(strings.TrimSpace(string(out))) > 0 {
		currentVal := strings.TrimSpace(string(out))
		ui.Info("Current git config core.sshCommand: %s", currentVal)

		broken := false
		if runtime.GOOS == "windows" {
			if strings.Contains(currentVal, "OpenSSH") {
				path := strings.Trim(currentVal, "\"")
				if _, err := os.Stat(path); os.IsNotExist(err) {
					ui.Warn("The configured path does not exist!")
					broken = true
				}
			}
		}

		if broken {
			ui.Warn("This configuration appears broken.")
		}

		fmt.Print("Unset git config core.sshCommand? (y/n): ")
		ans, _ := reader.ReadString('\n')
		if strings.ToLower(strings.TrimSpace(ans)) == "y" {
			exec.Command("git", "config", "--global", "--unset-all", "core.sshcommand").Run()
			ui.Success("Unset core.sshCommand.")
		}
	} else {
		ui.Info("No global core.sshCommand configured in git.")
	}

	home, err := os.UserHomeDir()
	if err == nil {
		configPath := filepath.Join(home, ".ssh", "config")
		if _, err := os.Stat(configPath); err == nil {
			fmt.Printf("\n~/.ssh/config exists (%s)\n", configPath)
			fmt.Print("Do you want to view its contents? (y/n): ")
			ans, _ := reader.ReadString('\n')
			if strings.ToLower(strings.TrimSpace(ans)) == "y" {
				if data, err := os.ReadFile(configPath); err == nil {
					fmt.Println("--- ~/.ssh/config ---")
					fmt.Println(string(data))
					fmt.Println("---------------------")
				}
			}
		}
	}
}
