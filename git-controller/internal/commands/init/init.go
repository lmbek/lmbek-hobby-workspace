package initcmd

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"workspace/git-controller/internal/gitutil"
	"workspace/git-controller/internal/sshutil"
	"workspace/git-controller/internal/system"
	"workspace/git-controller/internal/ui"
)

func Run() error {
	ui.Header("System Initialization")

	ui.Step(0, "Pre-flight checks")

	_ = sshutil.StartAgent()

	// Ensure GitHub is in known_hosts if it's missing (helps avoid interactive prompts)
	if runtime.GOOS == "windows" {
		home, _ := os.UserHomeDir()
		sshDir := filepath.Join(home, ".ssh")
		knownHosts := filepath.Join(sshDir, "known_hosts")
		if _, err := os.Stat(knownHosts); os.IsNotExist(err) || !containsHost(knownHosts, "github.com") {
			ui.Info("Adding github.com to known_hosts...")
			gitutil.EnsureDir(sshDir)
			// Use more robust PowerShell command for appending
			exec.Command("powershell", "-Command", "ssh-keyscan github.com | Out-File -FilePath \""+knownHosts+"\" -Append -Encoding ascii").Run()
		}
	}

	if err := checkGitHubConnectivity(); err != nil {
		return fmt.Errorf("GitHub connectivity check failed: %w\n\nRun 'make ssh' (or 'make ssh-setup') to set up your environment", err)
	}
	slog.Debug("GitHub connectivity verified")

	// 1. Planning & Analysis
	ui.Step(1, "Planning & Analysis")
	sys, _, err := system.LoadDefinition("system-definition.yaml")
	if err != nil {
		return fmt.Errorf("error loading system definition: %w", err)
	}

	planCount := len(sys.Proxy) + len(sys.Applications) + len(sys.Infrastructure) + len(sys.Orchestrator) + len(sys.Platform) + len(sys.Tools) + len(sys.Docs)
	ui.Success("System plan created with %d components", planCount)

	ui.Success("Initialization finished successfully!")
	ui.Info("Next step: Run '<cli> sync' to materialize the workspace")
	ui.Note(system.CLIDescription)
	return nil
}

func containsHost(filepath, host string) bool {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return false
	}
	return strings.Contains(string(data), host)
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
