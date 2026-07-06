<?php

declare(strict_types=1);

namespace GitController\Commands;

use GitController\SshUtil\SshUtil;
use GitController\UI\UI;

final class DoctorCommand
{
    public static string $cliName = 'git-controller-php';

    public static function run(): void
    {
        if (self::runCheck()) {
            self::printSSHSetupInstructions();
        }
    }

    private static function runCheck(): bool
    {
        UI::header('Workspace Doctor');
        UI::info('OS: %s/%s', PHP_OS_FAMILY, php_uname('m'));

        if (self::isWSL()) {
            UI::info('Environment: WSL (Windows Subsystem for Linux)');
        }

        self::checkGit();
        self::checkPHP();
        $sshIssues = self::checkSSH();
        self::checkDocker();

        // Check for conflicting GIT_SSH variables
        $gitSsh = getenv('GIT_SSH');
        if ($gitSsh !== false && $gitSsh !== '') {
            UI::warn('Conflicting environment variable detected: GIT_SSH=%s', $gitSsh);
            UI::info('Hint: This can override our GIT_SSH_COMMAND. Consider unsetting it.');
        }

        UI::success('Doctor report complete.');
        return $sshIssues;
    }

    private static function isWSL(): bool
    {
        if (PHP_OS_FAMILY !== 'Linux') {
            return false;
        }
        $version = @file_get_contents('/proc/version');
        if ($version === false) {
            return false;
        }
        return str_contains(strtolower($version), 'microsoft');
    }

    private static function checkGit(): void
    {
        UI::info('Checking Git installation...');
        exec('git --version 2>&1', $lines, $code);
        if ($code !== 0) {
            UI::error('Git is not installed or not in PATH');
        } else {
            UI::success('Git installed');
        }
    }

    private static function checkPHP(): void
    {
        UI::info('Checking PHP installation...');
        $version = PHP_VERSION;
        if (version_compare($version, '8.1.0', '<')) {
            UI::warn('PHP version might be too old: found %s, required 8.1+', $version);
        } else {
            UI::success('PHP version OK (%s)', $version);
        }
    }

    private static function checkSSH(): bool
    {
        UI::info('Checking SSH configuration...');
        $hasIssues = false;

        // Check for structural issues
        $issues = SshUtil::detectSSHIssues();
        foreach ($issues as $issue) {
            UI::error('SSH Structural Issue: %s', $issue);
            $hasIssues = true;
            if (str_contains($issue, 'is a directory')) {
                UI::info("Fix: Run '%s ssh' and select Option 6 (Cleanup broken SSH configurations) to resolve this automatically.", self::$cliName);
            }
        }

        // Check if ssh-agent is running
        UI::info('Checking if ssh-agent is running...');
        $keys = SshUtil::getLoadedKeys();
        if ($keys['code'] !== 0) {
            $hasIssues = true;
            UI::warn('ssh-agent might not be running or accessible');
            if (PHP_OS_FAMILY === 'Windows') {
                UI::info("Hint: On Windows, the 'OpenSSH Authentication Agent' service is often disabled by default.");
                UI::info('Fix (Admin PowerShell): Set-Service -Name ssh-agent -StartupType Manual; Start-Service ssh-agent');
            } elseif (self::isWSL()) {
                UI::info('Hint: In WSL, the ssh-agent does not persist across sessions by default.');
                UI::info('Fix: Add \'eval "$(ssh-agent -s)"\' to your ~/.bashrc or ~/.zshrc');
            } else {
                UI::info('Hint: On Linux, ensure ssh-agent is running. Start it with: eval "$(ssh-agent -s)"');
            }
        } else {
            UI::success('ssh-agent is responsive');
        }

        // Check loaded keys
        if ($keys['code'] === 0) {
            if (str_contains($keys['output'], 'The agent has no identities')) {
                UI::warn('No SSH keys found in agent');
                UI::info("Hint: Run '%s ssh' -> Option 2 to add a key to the agent.", self::$cliName);
            } else {
                UI::success('SSH keys loaded in agent');
            }
        }

        // Test GitHub connectivity
        UI::info('Testing connectivity to GitHub...');
        $result = SshUtil::checkGitHubConnectivity();
        if ($result['success']) {
            UI::success('GitHub authentication successful!');
        } else {
            $hasIssues = true;
            UI::error('GitHub authentication failed');
            UI::info('Hint: Ensure your public key is added to your GitHub account settings.');
            if (str_contains($result['output'], 'Permission denied')) {
                UI::info('Hint: If your key has a passphrase, it MUST be added to the ssh-agent.');
            }
        }

        return $hasIssues;
    }

    private static function checkDocker(): void
    {
        UI::info('Checking Docker installation...');
        exec('docker --version 2>&1', $lines, $code);
        if ($code !== 0) {
            UI::error('Docker is not installed or not in PATH');
        } else {
            UI::success('Docker installed (%s)', trim(implode("\n", $lines)));
        }

        $lines = [];
        exec('docker-compose --version 2>&1', $lines, $code);
        if ($code !== 0) {
            UI::warn('docker-compose is not installed or not in PATH');
        } else {
            UI::success('Docker Compose installed (%s)', trim(implode("\n", $lines)));
        }
    }

    private static function printSSHSetupInstructions(): void
    {
        $keyName = 'id_ed25519.private';
        $identityFile = SshUtil::getGitHubIdentityFile();
        if ($identityFile !== '') {
            $keyName = basename($identityFile);
        }

        echo "\nSSH Setup Guide:\n";
        echo "\n--- Manual Configuration ---\n";
        echo "1. Generate a new SSH key (if missing):\n";
        printf("   ssh-keygen -t ed25519 -C \"your_email@example.com\" -f ~/.ssh/%s\n", $keyName);
        echo "2. Start the ssh-agent:\n";
        if (PHP_OS_FAMILY === 'Windows') {
            echo "   PowerShell (Admin): Set-Service -Name ssh-agent -StartupType Manual; Start-Service ssh-agent\n";
        } else {
            echo "   eval \"\$(ssh-agent -s)\"\n";
        }
        echo "3. Add your SSH key to the agent:\n";
        printf("   ssh-add ~/.ssh/%s\n", $keyName);
        echo "4. Add the public key to your GitHub account:\n";

        $pubName = preg_replace('/\.private$/', '.public', $keyName);
        if (PHP_OS_FAMILY === 'Windows') {
            printf("   Get-Content ~/.ssh/%s | clip\n", $pubName);
        } else {
            printf("   cat ~/.ssh/%s\n", $pubName);
        }
        echo "   GitHub -> Settings -> SSH and GPG keys -> New SSH key\n";
        echo "5. (Optional) Configure ~/.ssh/config for host-specific settings:\n";
        echo "   Host github.com\n";
        echo "     HostName github.com\n";
        echo "     User git\n";
        printf("     IdentityFile ~/.ssh/%s\n", $keyName);
        echo "     AddKeysToAgent yes\n";
        echo "     IdentitiesOnly yes\n";

        echo "\n--- Automated Configuration ---\n";
        echo "   You can use the built-in automated tool for these steps:\n";
        printf("   %s ssh\n", self::$cliName);
    }
}
