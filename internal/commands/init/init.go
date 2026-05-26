package initcmd

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"workspace-controller/internal/gitutil"
	"workspace-controller/internal/sshutil"
	"workspace-controller/internal/system"
)

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
			gitutil.EnsureDir(sshDir)
			// Use more robust PowerShell command for appending
			exec.Command("powershell", "-Command", "ssh-keyscan github.com | Out-File -FilePath \""+knownHosts+"\" -Append -Encoding ascii").Run()
		}
	}

	if err := checkGitHubConnectivity(); err != nil {
		fmt.Printf("\nERROR: GitHub connectivity check failed.\n%v\n", err)
		fmt.Println("\nRun 'make ssh' (or 'make ssh-setup') to set up your environment.")
		os.Exit(1)
	}
	slog.Info("GitHub connectivity verified")

	// 1. Planning & Analysis
	fmt.Println("\n[1] Planning & Analysis:")
	sys, err := system.LoadDefinition("system/system-definition.yaml")
	if err != nil {
		slog.Error("Error loading system definition", "error", err)
		os.Exit(1)
	}

	planCount := len(sys.Services)
	if sys.Infrastructure != nil {
		planCount++
	}
	if sys.Tools != nil {
		planCount++
	}
	slog.Info("System plan created", "components", planCount)

	fmt.Println("\nInitialization finished successfully!")
	fmt.Println("Next step: Run '<cli> sync' to materialize the workspace and synchronize all repositories.")
	system.PrintCLINote()
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
