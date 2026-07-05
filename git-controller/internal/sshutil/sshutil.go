package sshutil

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// StartAgent ensures the ssh-agent is running, specifically handling Windows service requirements.
func StartAgent() error {
	if runtime.GOOS == "windows" {
		slog.Debug("Ensuring ssh-agent service is running on Windows...")

		// 1. Try to start the service directly first
		startCmd := exec.Command("powershell", "-Command", "Start-Service ssh-agent")
		if err := startCmd.Run(); err == nil {
			return nil
		}

		// 2. If it failed, try to check if it's disabled and set to Manual/Automatic
		setCmd := exec.Command("powershell", "-Command", "Set-Service ssh-agent -StartupType Automatic; Start-Service ssh-agent")
		if err := setCmd.Run(); err != nil {
			return fmt.Errorf("could not start ssh-agent service. Please run 'Set-Service -Name ssh-agent -StartupType Automatic' in an Administrative PowerShell window: %w", err)
		}
		return nil
	}

	// Linux / WSL logic
	if os.Getenv("SSH_AUTH_SOCK") != "" {
		// Agent already seems to be configured
		return nil
	}

	slog.Debug("Starting ssh-agent...")
	// Try to find an existing agent first
	if runtime.GOOS == "linux" {
		// Common patterns for WSL/Linux
		matches, _ := filepath.Glob("/tmp/ssh-*/agent.*")
		if len(matches) > 0 {
			os.Setenv("SSH_AUTH_SOCK", matches[0])
			slog.Debug("Found existing ssh-agent socket", "path", matches[0])
			return nil
		}
	}

	cmd := exec.Command("ssh-agent", "-s")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to start ssh-agent: %w", err)
	}

	// ssh-agent -s output looks like:
	// SSH_AUTH_SOCK=/tmp/ssh-XXXXXX/agent.XXXX; export SSH_AUTH_SOCK;
	// SSH_AGENT_PID=XXXXX; export SSH_AGENT_PID;
	// echo Agent pid XXXXX;
	lines := strings.Split(string(output), ";")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "SSH_AUTH_SOCK=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				val := parts[1]
				// Remove trailing export or other junk if any
				if idx := strings.Index(val, " "); idx != -1 {
					val = val[:idx]
				}
				os.Setenv("SSH_AUTH_SOCK", val)
				slog.Debug("ssh-agent started and environment set", "SSH_AUTH_SOCK", val)
			}
		}
		if strings.HasPrefix(line, "SSH_AGENT_PID=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				os.Setenv("SSH_AGENT_PID", parts[1])
			}
		}
	}

	return nil
}

// CheckGitHubConnectivity tests if the user can authenticate with GitHub via SSH.
func CheckGitHubConnectivity() (bool, string) {
	return checkSSH("git@github.com", false)
}

// CheckGitHubConnectivityNonInteractive tests if authentication works without prompting for passphrases.
func CheckGitHubConnectivityNonInteractive() (bool, string) {
	// Attempt to find or start agent before checking connectivity
	_ = StartAgent()
	return checkSSH("git@github.com", true)
}

// GetGitHubIdentityFile attempts to find the IdentityFile configured for github.com in ~/.ssh/config.
func GetGitHubIdentityFile() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	configPath := filepath.Join(home, ".ssh", "config")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return ""
	}

	content := string(data)
	content = strings.ReplaceAll(content, "\r\n", "\n")
	lines := strings.Split(content, "\n")

	var foundIdentity string
	bestSpecificity := -1
	inGitHubBlock := false
	currentBlockSpecificity := 0

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "#") {
			continue
		}

		lowerLine := strings.ToLower(trimmedLine)

		if strings.HasPrefix(lowerLine, "host ") {
			hosts := strings.Fields(lowerLine[5:])
			inGitHubBlock = false
			currentBlockSpecificity = 0
			for _, h := range hosts {
				if h == "github.com" {
					inGitHubBlock = true
					currentBlockSpecificity = 100 // Highest priority
					break
				} else if h == "*" {
					inGitHubBlock = true
					currentBlockSpecificity = 1
				} else if strings.Contains(h, "*") || strings.Contains(h, "?") {
					inGitHubBlock = true
					currentBlockSpecificity = 10
				}
			}
			continue
		}

		if inGitHubBlock && strings.HasPrefix(lowerLine, "identityfile ") {
			identityFile := strings.TrimSpace(trimmedLine[13:])
			identityFile = strings.Trim(identityFile, "\"")

			if strings.HasPrefix(identityFile, "~") {
				identityFile = filepath.Join(home, identityFile[1:])
			}

			newIdentity := filepath.Clean(identityFile)

			// Last match wins within the same specificity
			if currentBlockSpecificity >= bestSpecificity {
				foundIdentity = newIdentity
				bestSpecificity = currentBlockSpecificity
			}
		}
	}

	if foundIdentity != "" && !filepath.IsAbs(foundIdentity) {
		if abs, err := filepath.Abs(foundIdentity); err == nil {
			foundIdentity = abs
		}
	}
	return foundIdentity
}

// DetectSSHIssues checks for common SSH configuration pitfalls.
func DetectSSHIssues() []string {
	var issues []string
	home, err := os.UserHomeDir()
	if err != nil {
		return issues
	}

	sshDir := filepath.Join(home, ".ssh")
	// Check for common identity names that might be directories
	defaultKeys := []string{"id_rsa.private", "id_ed25519.private", "id_ecdsa.private", "id_dsa.private"}
	for _, key := range defaultKeys {
		path := filepath.Join(sshDir, key)
		if fi, err := os.Stat(path); err == nil && fi.IsDir() {
			issues = append(issues, fmt.Sprintf("Conflict: '~/.ssh/%s' is a directory, but SSH expects it to be a key file. This can cause authentication failures.", key))
		}
	}

	return issues
}

// ResolveSSHIssue attempts to fix a specific SSH issue.
func ResolveSSHIssue(issue string) (bool, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return false, err
	}

	if strings.Contains(issue, "is a directory, but SSH expects it to be a key file") {
		// Extract the path from the issue string or reconstruct it
		// Issue format: "Conflict: '~/.ssh/%s' is a directory..."
		start := strings.Index(issue, "~/.ssh/")
		if start == -1 {
			return false, fmt.Errorf("could not parse issue: %s", issue)
		}
		end := strings.Index(issue[start:], "'")
		if end == -1 {
			return false, fmt.Errorf("could not parse issue: %s", issue)
		}
		keyName := issue[start+7 : start+end]
		path := filepath.Join(home, ".ssh", keyName)

		slog.Warn("SSH Structural Issue", "path", path, "detail", "is a directory, but SSH expects it to be a key file.")
		slog.Debug("Manual intervention required", "path", path)
		return false, fmt.Errorf("manual intervention required: %s is a directory", path)
	}

	return false, fmt.Errorf("unsupported issue type")
}

// GetConfiguredSSHCommand returns the SSH command configured in git, or "ssh" if not set.
func GetConfiguredSSHCommand() string {
	cmd := exec.Command("git", "config", "--get", "core.sshcommand")
	output, err := cmd.Output()
	if err != nil || len(strings.TrimSpace(string(output))) == 0 {
		// Fallback to searching for system ssh on Windows if not configured
		if runtime.GOOS == "windows" {
			systemSSH := "C:/Windows/System32/OpenSSH/ssh.exe"
			if _, err := os.Stat(systemSSH); err == nil {
				return systemSSH
			}
		}
		return "ssh"
	}
	return strings.TrimSpace(string(output))
}

func checkSSH(host string, nonInteractive bool) (bool, string) {
	sshCmd := GetConfiguredSSHCommand()
	args := []string{"-T"}
	if nonInteractive {
		args = append(args, "-o", "BatchMode=yes")
	}
	// Add StrictHostKeyChecking=accept-new to handle new hosts automatically if they are valid
	args = append(args, "-o", "StrictHostKeyChecking=accept-new")

	args = append(args, host)

	cmd := exec.Command(sshCmd, args...)
	// Explicitly pass the environment so we don't lose SSH_AUTH_SOCK or other variables
	cmd.Env = os.Environ()
	output, _ := cmd.CombinedOutput()
	outStr := string(output)
	if strings.Contains(outStr, "successfully authenticated") {
		return true, outStr
	}

	// If it failed, try once more with more verbosity to help debugging
	if !nonInteractive {
		slog.Debug("SSH check failed, retrying with verbose output for diagnostics...", "command", sshCmd)
		vArgs := append([]string{"-v"}, args...)
		vCmd := exec.Command(sshCmd, vArgs...)
		vOutput, _ := vCmd.CombinedOutput()
		slog.Debug("Verbose SSH output", "output", string(vOutput))
		// We still return the original failure output but now it's in the logs if someone looks
	}

	return false, outStr
}

// GetLoadedKeys returns the output of ssh-add -l
func GetLoadedKeys() (string, error) {
	cmd := exec.Command("ssh-add", "-l")
	output, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

// GetPublicKeyContent attempts to find and return the public key content for a given private key path.
func GetPublicKeyContent(privateKeyPath string) (string, string, error) {
	// Try .public extension (replacing .private with .public if it exists)
	pubPath := strings.TrimSuffix(privateKeyPath, ".private") + ".public"
	if data, err := os.ReadFile(pubPath); err == nil {
		return string(data), pubPath, nil
	}

	// Try just appending .public
	pubPath = privateKeyPath + ".public"
	if data, err := os.ReadFile(pubPath); err == nil {
		return string(data), pubPath, nil
	}

	// Try generating it from private key
	cmd := exec.Command("ssh-keygen", "-y", "-f", privateKeyPath)
	output, err := cmd.CombinedOutput()
	if err == nil {
		return string(output), "generated from private key", nil
	}

	return "", "", fmt.Errorf("could not find or generate public key for %s", privateKeyPath)
}
