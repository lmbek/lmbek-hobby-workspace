<?php

declare(strict_types=1);

namespace GitController\Commands;

use GitController\SshUtil\SshUtil;
use GitController\UI\UI;

final class SshSetupCommand
{
    public static function run(): void
    {
        UI::header('SSH Setup Tool');

        echo "\nSSH Setup Options:\n";
        UI::info('1. Generate a new SSH key');
        UI::info('2. Add an existing key to the agent');
        UI::info('3. Configure ~/.ssh/config for GitHub');
        if (PHP_OS_FAMILY === 'Windows') {
            UI::info('4. Configure git to use Windows OpenSSH (Windows only)');
        }
        UI::info('5. Check current SSH status');
        UI::info('6. Cleanup broken SSH configurations');
        echo "\nSelect an option (1-6): ";

        $option = trim(fgets(STDIN) ?: '');

        match ($option) {
            '1' => self::generateNewKey(),
            '2' => self::addExistingKey(),
            '3' => self::configureSSHConfig(),
            '4' => PHP_OS_FAMILY === 'Windows' ? self::configureGitSSH() : UI::error('Invalid option selected'),
            '5' => self::checkStatus(),
            '6' => self::cleanupConfigs(),
            default => UI::error('Invalid option selected'),
        };
    }

    private static function generateNewKey(): void
    {
        echo 'Enter your email for the SSH key: ';
        $email = trim(fgets(STDIN) ?: '');

        $home = SshUtil::homeDir();
        if ($home === '') {
            UI::error('Could not determine home directory');
            return;
        }

        $defaultPath = $home . DIRECTORY_SEPARATOR . '.ssh' . DIRECTORY_SEPARATOR . 'id_ed25519.private';
        printf('Enter file in which to save the key (%s): ', $defaultPath);
        $keyPath = trim(fgets(STDIN) ?: '');
        if ($keyPath === '') {
            $keyPath = $defaultPath;
        }

        if (!str_ends_with($keyPath, '.private')) {
            $keyPath .= '.private';
            UI::info('Automatically added .private extension: %s', $keyPath);
        }

        $sshDir = dirname($keyPath);
        if (!is_dir($sshDir)) {
            mkdir($sshDir, 0700, true);
        }

        UI::info('Generating new SSH key: %s', $keyPath);
        $cmd = sprintf('ssh-keygen -t ed25519 -C %s -f %s', escapeshellarg($email), escapeshellarg($keyPath));
        passthru($cmd, $code);

        if ($code !== 0) {
            UI::error('Error generating SSH key');
            return;
        }

        UI::success('SSH key generated successfully.');
        echo "\nRecommended next steps:\n";
        echo "1. Add the public key to your GitHub account (see below).\n";
        echo "2. Run Option 2 to add this key to the ssh-agent.\n";
        echo "3. Run Option 3 to configure your ~/.ssh/config for GitHub.\n";

        $pubKeyPath = preg_replace('/\.private$/', '.public', $keyPath);
        $pubKey = @file_get_contents($pubKeyPath);
        if ($pubKey === false) {
            $pubKeyPath = $keyPath . '.pub';
            $pubKey = @file_get_contents($pubKeyPath);
        }

        if ($pubKey !== false) {
            echo "\nYour public key content:\n";
            echo $pubKey . "\n";
            echo "Copy the content above and add it to your GitHub account (Settings -> SSH and GPG keys).\n";
        }
    }

    private static function addExistingKey(): void
    {
        $home = SshUtil::homeDir();
        if ($home === '') {
            UI::error('Could not determine home directory');
            return;
        }

        $sshDir = $home . DIRECTORY_SEPARATOR . '.ssh';
        $keys = [];
        if (is_dir($sshDir)) {
            foreach (scandir($sshDir) ?: [] as $file) {
                if ($file === '.' || $file === '..') {
                    continue;
                }
                if (str_ends_with($file, '.private')) {
                    $keys[] = $sshDir . DIRECTORY_SEPARATOR . $file;
                }
            }
        }

        $keyPath = '';
        if (count($keys) > 0) {
            echo "\nFound existing SSH keys:\n";
            foreach ($keys as $i => $key) {
                printf("%d. %s\n", $i + 1, $key);
            }
            printf("%d. Enter path manually\n", count($keys) + 1);
            printf("\nChoose a key or enter path manually (1-%d): ", count($keys) + 1);
            $choice = (int) trim(fgets(STDIN) ?: '');
            if ($choice > 0 && $choice <= count($keys)) {
                $keyPath = $keys[$choice - 1];
            } else {
                echo 'Enter the path to your private key: ';
                $keyPath = trim(fgets(STDIN) ?: '');
            }
        } else {
            printf('Enter the path to your private key (e.g., %s): ', $sshDir . DIRECTORY_SEPARATOR . 'id_ed25519.private');
            $keyPath = trim(fgets(STDIN) ?: '');
        }

        if ($keyPath === '') {
            UI::error('Key path cannot be empty.');
            return;
        }

        self::addToAgent($keyPath);
    }

    private static function configureSSHConfig(): void
    {
        $home = SshUtil::homeDir();
        if ($home === '') {
            UI::error('Could not determine home directory');
            return;
        }

        $sshDir = $home . DIRECTORY_SEPARATOR . '.ssh';
        $configPath = $sshDir . DIRECTORY_SEPARATOR . 'config';

        if (!is_dir($sshDir)) {
            mkdir($sshDir, 0700, true);
        }

        // Discover keys
        $keys = [];
        if (is_dir($sshDir)) {
            foreach (scandir($sshDir) ?: [] as $file) {
                if ($file === '.' || $file === '..' || is_dir($sshDir . DIRECTORY_SEPARATOR . $file)) {
                    continue;
                }
                if (!str_ends_with($file, '.public') && !str_ends_with($file, '.pub') && !str_ends_with($file, '.old')
                    && !str_ends_with($file, '.Identifier') && $file !== 'known_hosts' && $file !== 'known_hosts.old'
                    && $file !== 'config' && $file !== 'authorized_keys') {
                    $keys[] = $sshDir . DIRECTORY_SEPARATOR . $file;
                }
            }
        }

        $selectedKey = '';
        if (count($keys) > 0) {
            echo "\nSelect a key to use for GitHub in your config:\n";
            foreach ($keys as $i => $key) {
                printf("%d. %s\n", $i + 1, $key);
            }
            printf("%d. Enter path manually\n", count($keys) + 1);
            echo 'Choice: ';
            $choice = (int) trim(fgets(STDIN) ?: '');
            if ($choice > 0 && $choice <= count($keys)) {
                $selectedKey = $keys[$choice - 1];
            } else {
                echo 'Enter the path to your private key: ';
                $selectedKey = trim(fgets(STDIN) ?: '');
            }
        } else {
            echo 'Enter the path to your private key: ';
            $selectedKey = trim(fgets(STDIN) ?: '');
        }

        if ($selectedKey === '') {
            UI::error('Key path cannot be empty.');
            return;
        }

        $identityFile = str_replace('\\', '/', $selectedKey);
        if (str_contains($identityFile, ' ')) {
            $identityFile = "\"$identityFile\"";
        }

        $configEntry = "Host github.com\n  HostName github.com\n  User git\n  IdentityFile $identityFile\n  AddKeysToAgent yes\n  IdentitiesOnly yes\n";

        printf("\nProposed configuration entry for %s:\n%s", $configPath, $configEntry);
        echo "\nDo you want to apply this to your SSH config? This will consolidate multiple github.com entries. (y/n): ";
        $answer = strtolower(trim(fgets(STDIN) ?: ''));
        if ($answer !== 'y') {
            UI::info('Configuration skipped.');
            return;
        }

        // Read existing config and remove all Host github.com blocks
        $content = @file_get_contents($configPath) ?: '';
        $content = str_replace("\r\n", "\n", $content);
        $lines = explode("\n", $content);

        $newLines = [];
        $inGitHubBlock = false;

        foreach ($lines as $line) {
            $lower = strtolower(trim($line));
            if (str_starts_with($lower, 'host ')) {
                $hosts = preg_split('/\s+/', substr($lower, 5));
                $inGitHubBlock = in_array('github.com', $hosts, true);
            }
            if (!$inGitHubBlock) {
                $newLines[] = $line;
            }
        }

        // Filter trailing empty lines
        while (count($newLines) > 0 && trim(end($newLines)) === '') {
            array_pop($newLines);
        }

        $finalConfig = implode("\n", $newLines);
        if ($finalConfig !== '') {
            $finalConfig .= "\n";
        }
        $finalConfig .= $configEntry;

        file_put_contents($configPath, $finalConfig);
        chmod($configPath, 0600);

        UI::success('SSH config updated and consolidated successfully: %s', $configPath);
        echo "\nRecommended next step:\n";
        if (PHP_OS_FAMILY === 'Windows') {
            echo "1. Run Option 4 (Configure git to use Windows OpenSSH) then Option 5 (Check current SSH status).\n";
        } else {
            echo "1. Run Option 5 (Check current SSH status) to verify everything is working and test GitHub connectivity.\n";
        }
        printf("2. Run '%s clone' to clone your repositories.\n", DoctorCommand::$cliName);
    }

    private static function configureGitSSH(): void
    {
        if (PHP_OS_FAMILY !== 'Windows') {
            echo "This option is only relevant for Windows systems.\n";
            return;
        }

        $sshPath = 'C:/Windows/System32/OpenSSH/ssh.exe';
        printf("\nThis will set 'git config --global core.sshCommand' to: %s\n", $sshPath);
        echo 'Do you want to proceed? (y/n): ';
        $answer = strtolower(trim(fgets(STDIN) ?: ''));
        if ($answer !== 'y') {
            UI::info('Operation cancelled.');
            return;
        }

        exec('git config --global --unset-all core.sshcommand 2>&1');
        exec('git config --global core.sshCommand ' . escapeshellarg($sshPath) . ' 2>&1', $lines, $code);
        if ($code !== 0) {
            UI::error('Failed to set git config');
            return;
        }

        UI::success('Successfully configured git to use Windows OpenSSH.');
    }

    private static function checkStatus(): void
    {
        UI::info('Checking current SSH status...');

        if (PHP_OS_FAMILY === 'Windows') {
            $currentSsh = SshUtil::getConfiguredSSHCommand();
            if ($currentSsh === 'ssh') {
                UI::warn("Git core.sshCommand is not explicitly set. Git might use its internal SSH client.");
                echo "Recommendation: Run Option 4 (Configure git to use Windows OpenSSH) for consistency.\n";
            } else {
                UI::info('Git core.sshCommand is set to: %s', $currentSsh);
            }
        }

        // Check for structural issues
        $issues = SshUtil::detectSSHIssues();
        if (count($issues) > 0) {
            echo "\n⚠️  SSH Structural Issues Detected:\n";
            foreach ($issues as $issue) {
                printf("   - %s\n", $issue);
                echo "     Do you want to attempt to fix this? (y/n): ";
                $answer = strtolower(trim(fgets(STDIN) ?: ''));
                if ($answer === 'y') {
                    UI::warn('Manual intervention required for this issue.');
                }
            }
        }

        $keys = SshUtil::getLoadedKeys();
        if ($keys['code'] !== 0) {
            UI::warn('No keys found in agent or ssh-agent not running: %s', $keys['output']);
            if (str_contains($keys['output'], 'authentication agent') || str_contains($keys['output'], 'Error connecting')) {
                echo "\nssh-agent is not running or accessible. Do you want to try starting it? (y/n): ";
                $answer = strtolower(trim(fgets(STDIN) ?: ''));
                if ($answer === 'y') {
                    if (!SshUtil::startAgent()) {
                        UI::error('Failed to start ssh-agent');
                    } else {
                        $keys = SshUtil::getLoadedKeys();
                        if ($keys['code'] === 0) {
                            UI::success("Loaded keys:\n%s", $keys['output']);
                        }
                    }
                }
            }
        } else {
            UI::success("Loaded keys:\n%s", $keys['output']);
        }

        UI::info('Testing GitHub connectivity...');
        $result = SshUtil::checkGitHubConnectivity();
        if ($result['success']) {
            UI::success('GitHub authentication successful!');
        } else {
            UI::warn('GitHub authentication failed: %s', trim($result['output']));
            if (str_contains($result['output'], 'Permission denied')) {
                UI::info('Hint: Your key might have a passphrase and is not added to the agent, or the wrong key is being used.');
                UI::info("Try Option 2 (Add existing key to agent) if you are prompted for a passphrase during '%s clone'.", DoctorCommand::$cliName);

                $identityFile = SshUtil::getGitHubIdentityFile();
                if ($identityFile !== '') {
                    printf("\nYour configured IdentityFile is: %s\n", $identityFile);
                    $pubKey = SshUtil::getPublicKeyContent($identityFile);
                    if (!isset($pubKey['error'])) {
                        printf("Corresponding public key (%s):\n\n%s\n", $pubKey['path'], trim($pubKey['content']));
                        echo "CRITICAL: Please ensure this EXACT public key is added to your GitHub account.\n";
                        echo "Go to: https://github.com/settings/keys\n";
                    }
                }

                echo "\nTo clear potentially conflicting keys from your agent, you can run: ssh-add -D\n";

                echo "\nDo you want to run a detailed SSH diagnostic (ssh -vT)? (y/n): ";
                $answer = strtolower(trim(fgets(STDIN) ?: ''));
                if ($answer === 'y') {
                    echo "\n--- SSH Detailed Diagnostic ---\n";
                    $sshCmd = SshUtil::getConfiguredSSHCommand();
                    passthru("$sshCmd -vT git@github.com 2>&1");
                    echo "--- End of Diagnostic ---\n";
                }
            }
        }
    }

    private static function cleanupConfigs(): void
    {
        UI::info('Starting SSH configuration cleanup...');

        echo "\nCleanup Options:\n";
        echo "1. Clear all keys from ssh-agent (ssh-add -D)\n";
        echo "2. Reset ~/.ssh/config (Backup and create new)\n";
        echo "3. Full Reset (Agent + Config + Known Hosts backup)\n";
        echo "\nSelect a cleanup option (1-3, or Enter to cancel): ";

        $choice = trim(fgets(STDIN) ?: '');

        $home = SshUtil::homeDir();
        if ($home === '') {
            UI::error('Could not determine home directory');
            return;
        }
        $sshDir = $home . DIRECTORY_SEPARATOR . '.ssh';

        match ($choice) {
            '1' => self::cleanupAgent(),
            '2' => self::cleanupConfig($sshDir),
            '3' => self::cleanupFull($sshDir),
            default => UI::info('Cleanup cancelled.'),
        };
    }

    private static function cleanupAgent(): void
    {
        UI::info('Clearing keys from agent...');
        exec('ssh-add -D 2>&1', $lines, $code);
        if ($code !== 0) {
            UI::error('Failed to clear keys: %s', implode("\n", $lines));
        } else {
            UI::success('All identities removed from agent.');
        }
    }

    private static function cleanupConfig(string $sshDir): void
    {
        $configPath = $sshDir . DIRECTORY_SEPARATOR . 'config';
        if (file_exists($configPath)) {
            UI::warn('Existing config found at %s. Please manually back it up before resetting.', $configPath);
            echo "Are you sure you want to overwrite it with a fresh one? (y/n): ";
            $answer = strtolower(trim(fgets(STDIN) ?: ''));
            if ($answer !== 'y') {
                return;
            }
        }
        UI::info('Creating fresh ~/.ssh/config...');
        file_put_contents($configPath, "# Fresh SSH Config\n");
        chmod($configPath, 0600);
    }

    private static function cleanupFull(string $sshDir): void
    {
        exec('ssh-add -D 2>&1');

        $configPath = $sshDir . DIRECTORY_SEPARATOR . 'config';
        if (file_exists($configPath)) {
            UI::warn('Existing config found at %s. Manual reset required.', $configPath);
        }

        $khPath = $sshDir . DIRECTORY_SEPARATOR . 'known_hosts';
        if (file_exists($khPath)) {
            UI::warn('Existing known_hosts found at %s. Manual reset required.', $khPath);
        }

        UI::info('Full reset requested. Please manually clean up the files mentioned above.');
    }

    private static function addToAgent(string $keyPath): void
    {
        UI::info('Adding key to ssh-agent: %s', $keyPath);

        if (!SshUtil::startAgent()) {
            UI::warn('Could not ensure ssh-agent is running');
        }

        passthru('ssh-add ' . escapeshellarg($keyPath), $code);

        if ($code !== 0) {
            UI::error('Error adding key to agent');
            if (PHP_OS_FAMILY === 'Windows') {
                UI::info("Hint: If you are on Windows, ensure the 'OpenSSH Authentication Agent' service is running and you have permissions to start it.");
            } else {
                UI::info('Hint: On Linux/WSL, ensure ssh-agent is running. If you just started it, you may need to reload your shell or run:');
                echo "   eval \"\$(ssh-agent -s)\"\n";
            }
        } else {
            UI::success('Key added successfully to ssh-agent.');
            if (PHP_OS_FAMILY === 'Windows') {
                UI::info('Windows Tip: Once added and passphrase entered, the OpenSSH Agent service should remember it across sessions.');
            }
            echo "\nRecommended next steps:\n";
            echo "1. Run Option 3 to configure your ~/.ssh/config for GitHub (if you haven't yet).\n";
            if (PHP_OS_FAMILY === 'Windows') {
                echo "2. Run Option 4 (Configure git to use Windows OpenSSH) then Option 5 (Check current SSH status).\n";
            } else {
                echo "2. Run Option 5 (Check current SSH status) to verify everything is working and test GitHub connectivity.\n";
            }
            printf("3. Run '%s clone' to clone your repositories.\n", DoctorCommand::$cliName);
        }
    }
}
