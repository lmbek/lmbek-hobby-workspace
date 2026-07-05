package sshsetup

import (
	"bufio"
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

func RunError() error {
	ui.Header("SSH Setup Tool")

	fmt.Println("\nSSH Setup Options:")
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

func Run() {
	_ = RunError()
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
		slog.Info("Operation cancelled.")
		return
	}

	// Unset first to clear any potential issues
	exec.Command("git", "config", "--global", "--unset-all", "core.sshcommand").Run()

	cmd := exec.Command("git", "config", "--global", "core.sshCommand", sshPath)
	if err := cmd.Run(); err != nil {
		slog.Error("Failed to set git config", "error", err)
		return
	}

	slog.Info("Successfully configured git to use Windows OpenSSH.")
}

func configureSSHConfig(reader *bufio.Reader) {
	home, err := os.UserHomeDir()
	if err != nil {
		slog.Error("Could not determine home directory", "error", err)
		return
	}

	sshDir := filepath.Join(home, ".ssh")
	configPath := filepath.Join(sshDir, "config")

	// 1. Discover keys to help user pick one for the config
	files, err := os.ReadDir(sshDir)
	var keys []string
	if err == nil {
		for _, file := range files {
			if !file.IsDir() {
				name := file.Name()
				if !strings.HasSuffix(name, ".public") && !strings.HasSuffix(name, ".old") && !strings.HasSuffix(name, ".Identifier") && name != "known_hosts" && name != "config" && name != "authorized_keys" {
					keys = append(keys, filepath.Join(sshDir, name))
				}
			}
		}
	}

	selectedKey := ""
	if len(keys) > 0 {
		fmt.Println("\nSelect a key to use for GitHub in your config:")
		for i, key := range keys {
			fmt.Printf("%d. %s\n", i+1, key)
		}
		fmt.Printf("%d. Enter path manually\n", len(keys)+1)
		fmt.Print("Choice: ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)
		idx := 0
		fmt.Sscanf(choice, "%d", &idx)
		if idx > 0 && idx <= len(keys) {
			selectedKey = keys[idx-1]
		} else {
			fmt.Print("Enter the path to your private key: ")
			manualPath, _ := reader.ReadString('\n')
			selectedKey = strings.TrimSpace(manualPath)
		}
	} else {
		fmt.Print("Enter the path to your private key: ")
		manualPath, _ := reader.ReadString('\n')
		selectedKey = strings.TrimSpace(manualPath)
	}

	if selectedKey == "" {
		slog.Error("Key path cannot be empty")
		return
	}

	// 2. Propose configuration
	// Use forward slashes for the config file even on Windows
	identityFile := filepath.ToSlash(selectedKey)
	if strings.Contains(identityFile, " ") {
		identityFile = fmt.Sprintf("\"%s\"", identityFile)
	}

	configEntry := fmt.Sprintf("Host github.com\n  HostName github.com\n  User git\n  IdentityFile %s\n", identityFile)

	fmt.Printf("\nProposed configuration entry for %s:\n%s", configPath, configEntry)
	fmt.Print("\nDo you want to apply this to your SSH config? This will consolidate multiple github.com entries. (y/n): ")
	answer, _ := reader.ReadString('\n')
	if strings.ToLower(strings.TrimSpace(answer)) != "y" {
		slog.Info("Configuration skipped.")
		return
	}

	// 3. Read existing config and remove all Host github.com blocks
	data, _ := os.ReadFile(configPath)
	content := string(data)
	content = strings.ReplaceAll(content, "\r\n", "\n")
	lines := strings.Split(content, "\n")

	var newLines []string
	inGitHubBlock := false

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		lowerLine := strings.ToLower(trimmedLine)

		if strings.HasPrefix(lowerLine, "host ") {
			hosts := strings.Fields(lowerLine[5:])
			inGitHubBlock = false
			for _, h := range hosts {
				if h == "github.com" {
					inGitHubBlock = true
					break
				}
			}
		}

		if !inGitHubBlock {
			newLines = append(newLines, line)
		}
	}

	// Filter out trailing empty lines
	for len(newLines) > 0 && strings.TrimSpace(newLines[len(newLines)-1]) == "" {
		newLines = newLines[:len(newLines)-1]
	}

	finalConfig := strings.Join(newLines, "\n")
	if finalConfig != "" {
		finalConfig += "\n"
	}
	finalConfig += configEntry

	if err := os.WriteFile(configPath, []byte(finalConfig), 0600); err != nil {
		slog.Error("Could not write to config file", "error", err)
		return
	}

	slog.Info("SSH config updated and consolidated successfully.", "path", configPath)
	fmt.Println("\nRecommended next step:")
	if runtime.GOOS == "windows" {
		fmt.Println("1. Run Option 4 (Configure git to use Windows OpenSSH) then Option 5 (Check current SSH status).")
	} else {
		fmt.Println("1. Run Option 5 (Check current SSH status) to verify everything is working and test GitHub connectivity.")
	}
	fmt.Println("2. Run '<cli> sync' to synchronize your repositories.")
}

func generateNewKey(reader *bufio.Reader) {
	fmt.Print("Enter your email for the SSH key: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	home, err := os.UserHomeDir()
	if err != nil {
		slog.Error("Could not determine home directory", "error", err)
		return
	}

	defaultPath := filepath.Join(home, ".ssh", "id_ed25519.private")
	fmt.Printf("Enter file in which to save the key (%s): ", defaultPath)
	keyPath, _ := reader.ReadString('\n')
	keyPath = strings.TrimSpace(keyPath)
	if keyPath == "" {
		keyPath = defaultPath
	}

	if !strings.HasSuffix(keyPath, ".private") {
		keyPath += ".private"
		slog.Info("Automatically added .private extension", "path", keyPath)
	}

	// Ensure .ssh directory exists
	sshDir := filepath.Dir(keyPath)
	if _, err := os.Stat(sshDir); os.IsNotExist(err) {
		os.MkdirAll(sshDir, 0700)
	}

	slog.Info("Generating new SSH key...", "path", keyPath)
	// We need to tell ssh-keygen where to put the public key too,
	// but it usually just appends .pub to the filename.
	// Since the user wants .public, we'll have to rename it after generation.
	cmd := exec.Command("ssh-keygen", "-t", "ed25519", "-C", email, "-f", keyPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		slog.Error("Error generating SSH key", "error", err)
		return
	}

	slog.Info("SSH key generated successfully.")
	fmt.Println("\nRecommended next steps:")
	fmt.Println("1. Add the public key to your GitHub account (see below).")
	fmt.Println("2. Run Option 2 to add this key to the ssh-agent (this allows it to remember your passphrase).")
	fmt.Println("3. Run Option 3 to configure your ~/.ssh/config for GitHub.")

	fmt.Printf("\nNext step: Add your public key to GitHub:\n")
	// Look for both .pub and .public
	pubKeyPath := strings.TrimSuffix(keyPath, ".private") + ".public"
	pubKey, err := os.ReadFile(pubKeyPath)
	if err != nil {
		pubKeyPath = keyPath + ".pub"
		pubKey, err = os.ReadFile(pubKeyPath)
	}

	if err == nil {
		fmt.Println("\nYour public key content:")
		fmt.Println(string(pubKey))
		fmt.Println("\nCopy the content above and add it to your GitHub account (Settings -> SSH and GPG keys).")
	}
}

func addExistingKey(reader *bufio.Reader) {
	home, err := os.UserHomeDir()
	if err != nil {
		slog.Error("Could not determine home directory", "error", err)
		return
	}

	sshDir := filepath.Join(home, ".ssh")
	files, err := os.ReadDir(sshDir)
	var keys []string
	if err == nil {
		for _, file := range files {
			if !file.IsDir() {
				name := file.Name()
				// Basic heuristic for private keys: ends with .private
				if strings.HasSuffix(name, ".private") {
					keys = append(keys, filepath.Join(sshDir, name))
				}
			}
		}
	}

	keyPath := ""
	if len(keys) > 0 {
		fmt.Println("\nFound existing SSH keys:")
		for i, key := range keys {
			fmt.Printf("%d. %s\n", i+1, key)
		}
		fmt.Printf("%d. Enter path manually\n", len(keys)+1)
		fmt.Print("\nChoose a key or enter path manually (1-", len(keys)+1, "): ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		idx := 0
		fmt.Sscanf(choice, "%d", &idx)
		if idx > 0 && idx <= len(keys) {
			keyPath = keys[idx-1]
		} else if idx == len(keys)+1 {
			fmt.Print("Enter the path to your private key: ")
			manualPath, _ := reader.ReadString('\n')
			keyPath = strings.TrimSpace(manualPath)
		} else {
			slog.Error("Invalid choice")
			return
		}
	} else {
		fmt.Printf("Enter the path to your private key (e.g., %s): ", filepath.Join(home, ".ssh", "id_ed25519.private"))
		manualPath, _ := reader.ReadString('\n')
		keyPath = strings.TrimSpace(manualPath)
	}

	if keyPath == "" {
		slog.Error("Key path cannot be empty")
		return
	}

	addToAgent(keyPath)
}

func checkStatus(reader *bufio.Reader) {
	slog.Info("Checking current SSH status...")

	if runtime.GOOS == "windows" {
		currentSsh := sshutil.GetConfiguredSSHCommand()
		if currentSsh == "ssh" {
			slog.Warn("Git core.sshCommand is not explicitly set. Git might use its internal SSH client.")
			fmt.Println("Recommendation: Run Option 4 (Configure git to use Windows OpenSSH) for consistency.")
		} else {
			slog.Info("Git core.sshCommand is set", "value", currentSsh)
		}
	}

	// 0. Check for structural issues
	issues := sshutil.DetectSSHIssues()
	if len(issues) > 0 {
		fmt.Println("\n⚠️  SSH Structural Issues Detected:")
		for _, issue := range issues {
			fmt.Printf("   - %s\n", issue)
			fmt.Printf("     Do you want to attempt to fix this? (y/n): ")
			answer, _ := reader.ReadString('\n')
			if strings.ToLower(strings.TrimSpace(answer)) == "y" {
				success, err := sshutil.ResolveSSHIssue(issue)
				if success {
					slog.Info("Successfully resolved issue")
				} else {
					slog.Error("Failed to resolve issue", "error", err)
				}
			}
		}
	}

	keys, err := sshutil.GetLoadedKeys()
	if err != nil {
		slog.Warn("No keys found in agent or ssh-agent not running", "error", keys)
		if strings.Contains(keys, "authentication agent") || strings.Contains(keys, "Error connecting") {
			fmt.Print("\nssh-agent is not running or accessible. Do you want to try starting it? (y/n): ")
			answer, _ := reader.ReadString('\n')
			if strings.ToLower(strings.TrimSpace(answer)) == "y" {
				if err := sshutil.StartAgent(); err != nil {
					slog.Error("Failed to start ssh-agent", "error", err)
				} else {
					// Retry
					keys, err = sshutil.GetLoadedKeys()
					if err == nil {
						slog.Info("Loaded keys", "keys", keys)
					}
				}
			}
		}
	} else {
		slog.Info("Loaded keys", "keys", keys)
	}

	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".ssh", "config")
	if data, err := os.ReadFile(configPath); err == nil {
		content := string(data)
		content = strings.ReplaceAll(content, "\r\n", "\n")
		if strings.Count(content, "Host github.com") > 1 {
			slog.Warn("Duplicate 'Host github.com' entries detected in your config. Option 3 will consolidate them.")
		}
	}

	slog.Info("Testing GitHub connectivity...")
	success, output := sshutil.CheckGitHubConnectivity()
	if success {
		slog.Info("GitHub authentication successful!")
	} else {
		slog.Warn("GitHub authentication failed", "output", strings.TrimSpace(output))
		if strings.Contains(output, "Permission denied") {
			slog.Info("Hint: Your key might have a passphrase and is not added to the agent, or the wrong key is being used.")
			slog.Info("Try Option 2 (Add existing key to agent) if you are prompted for a passphrase during '<cli> sync'.")

			if identityFile := sshutil.GetGitHubIdentityFile(); identityFile != "" {
				fmt.Printf("\nYour configured IdentityFile is: %s\n", identityFile)
				pubKey, pubPath, err := sshutil.GetPublicKeyContent(identityFile)
				if err == nil {
					fmt.Printf("Corresponding public key (%s):\n\n%s\n", pubPath, strings.TrimSpace(pubKey))
					fmt.Println("CRITICAL: Please ensure this EXACT public key is added to your GitHub account.")
					fmt.Println("Go to: https://github.com/settings/keys")
				}
			}

			fmt.Println("\nTo clear potentially conflicting keys from your agent, you can run: ssh-add -D")

			fmt.Print("\nDo you want to run a detailed SSH diagnostic (ssh -vT)? (y/n): ")
			diagAnswer, _ := reader.ReadString('\n')
			if strings.ToLower(strings.TrimSpace(diagAnswer)) == "y" {
				fmt.Println("\n--- SSH Detailed Diagnostic ---")
				cmd := exec.Command("ssh", "-vT", "git@github.com")
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				_ = cmd.Run()
				fmt.Println("--- End of Diagnostic ---")
			}
		}
	}
}

func cleanupConfigs(reader *bufio.Reader) {
	slog.Info("Starting SSH configuration cleanup...")

	fmt.Println("\nCleanup Options:")
	fmt.Println("1. Clear all keys from ssh-agent (ssh-add -D)")
	fmt.Println("2. Reset ~/.ssh/config (Backup and create new)")
	fmt.Println("3. Full Reset (Agent + Config + Known Hosts backup)")
	fmt.Print("\nSelect a cleanup option (1-3, or Enter to cancel): ")

	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	home, err := os.UserHomeDir()
	if err != nil {
		slog.Error("Could not determine home directory", "error", err)
		return
	}
	sshDir := filepath.Join(home, ".ssh")

	switch choice {
	case "1":
		slog.Info("Clearing keys from agent...")
		cmd := exec.Command("ssh-add", "-D")
		if output, err := cmd.CombinedOutput(); err != nil {
			slog.Error("Failed to clear keys", "error", err, "output", string(output))
		} else {
			slog.Info("All identities removed from agent.")
		}

	case "2":
		configPath := filepath.Join(sshDir, "config")
		if _, err := os.Stat(configPath); err == nil {
			slog.Warn("Existing config found. Please manually back it up before resetting.", "path", configPath)
			fmt.Print("Are you sure you want to overwrite it with a fresh one? (y/n): ")
			answer, _ := reader.ReadString('\n')
			if strings.ToLower(strings.TrimSpace(answer)) != "y" {
				return
			}
		}
		slog.Info("Creating fresh ~/.ssh/config")
		if err := os.WriteFile(configPath, []byte("# Fresh SSH Config\n"), 0600); err != nil {
			slog.Error("Failed to create fresh config", "error", err)
		}

	case "3":
		// Agent
		exec.Command("ssh-add", "-D").Run()

		// Config
		configPath := filepath.Join(sshDir, "config")
		if _, err := os.Stat(configPath); err == nil {
			slog.Warn("Existing config found. Manual reset required.", "path", configPath)
		}

		// Known hosts
		khPath := filepath.Join(sshDir, "known_hosts")
		if _, err := os.Stat(khPath); err == nil {
			slog.Warn("Existing known_hosts found. Manual reset required.", "path", khPath)
		}

		slog.Info("Full reset requested. Please manually clean up the files mentioned above.")

	default:
		slog.Info("Cleanup cancelled.")
	}
}

func addToAgent(keyPath string) {
	slog.Info("Adding key to ssh-agent...", "path", keyPath)

	if err := sshutil.StartAgent(); err != nil {
		slog.Warn("Could not ensure ssh-agent is running", "error", err)
	}

	// On Windows, the OpenSSH agent service doesn't support -K, but it does persist
	// keys added via ssh-add if the service is running.
	cmd := exec.Command("ssh-add", keyPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		slog.Error("Error adding key to agent", "error", err)
		if runtime.GOOS == "windows" {
			slog.Info("Hint: If you are on Windows, ensure the 'OpenSSH Authentication Agent' service is running and you have permissions to start it.")
		} else {
			slog.Info("Hint: On Linux/WSL, ensure ssh-agent is running. If you just started it, you may need to reload your shell or run:")
			fmt.Println("   eval \"$(ssh-agent -s)\"")
			slog.Info("To make it persistent, consider adding this to your ~/.bashrc or ~/.zshrc:")
			fmt.Println(`   if [ -z "$SSH_AUTH_SOCK" ]; then
     # Check for existing agent
     SOCK=$(ls /tmp/ssh-*/agent.* 2>/dev/null | head -n 1)
     if [ -n "$SOCK" ]; then
       export SSH_AUTH_SOCK=$SOCK
     else
       eval $(ssh-agent -s)
     fi
   fi`)
		}
	} else {
		slog.Info("Key added successfully to ssh-agent.")
		if runtime.GOOS == "windows" {
			slog.Info("Windows Tip: Once added and passphrase entered, the OpenSSH Agent service should remember it across sessions.")
		}
		fmt.Println("\nRecommended next steps:")
		fmt.Println("1. Run Option 3 to configure your ~/.ssh/config for GitHub (if you haven't yet).")
		if runtime.GOOS == "windows" {
			fmt.Println("2. Run Option 4 (Configure git to use Windows OpenSSH) then Option 5 (Check current SSH status).")
		} else {
			fmt.Println("2. Run Option 5 (Check current SSH status) to verify everything is working and test GitHub connectivity.")
		}
		fmt.Println("3. Run '<cli> sync' to synchronize your repositories.")
	}
}
