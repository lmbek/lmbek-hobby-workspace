<?php

declare(strict_types=1);

namespace GitController\SshUtil;

final class SshUtil
{
    public static function startAgent(): bool
    {
        if (PHP_OS_FAMILY === 'Windows') {
            $result = self::exec('powershell -Command "Start-Service ssh-agent"');
            if ($result['code'] === 0) {
                return true;
            }
            $result = self::exec('powershell -Command "Set-Service ssh-agent -StartupType Automatic; Start-Service ssh-agent"');
            return $result['code'] === 0;
        }

        if (getenv('SSH_AUTH_SOCK') !== false && getenv('SSH_AUTH_SOCK') !== '') {
            return true;
        }

        if (PHP_OS_FAMILY === 'Linux') {
            $matches = glob('/tmp/ssh-*/agent.*');
            if ($matches !== false && count($matches) > 0) {
                putenv('SSH_AUTH_SOCK=' . $matches[0]);
                return true;
            }
        }

        $result = self::exec('ssh-agent -s');
        if ($result['code'] !== 0) {
            return false;
        }

        foreach (explode(';', $result['output']) as $part) {
            $part = trim($part);
            if (str_starts_with($part, 'SSH_AUTH_SOCK=')) {
                $val = explode('=', $part, 2)[1] ?? '';
                if (($sp = strpos($val, ' ')) !== false) {
                    $val = substr($val, 0, $sp);
                }
                putenv('SSH_AUTH_SOCK=' . $val);
            }
            if (str_starts_with($part, 'SSH_AGENT_PID=')) {
                $val = explode('=', $part, 2)[1] ?? '';
                putenv('SSH_AGENT_PID=' . $val);
            }
        }

        return true;
    }

    public static function checkGitHubConnectivity(): array
    {
        return self::checkSSH('git@github.com', false);
    }

    public static function checkGitHubConnectivityNonInteractive(): array
    {
        self::startAgent();
        return self::checkSSH('git@github.com', true);
    }

    public static function getGitHubIdentityFile(): string
    {
        $home = self::homeDir();
        if ($home === '') {
            return '';
        }

        $configPath = $home . DIRECTORY_SEPARATOR . '.ssh' . DIRECTORY_SEPARATOR . 'config';
        $content = @file_get_contents($configPath);
        if ($content === false) {
            return '';
        }

        $content = str_replace("\r\n", "\n", $content);
        $lines = explode("\n", $content);

        $foundIdentity = '';
        $bestSpecificity = -1;
        $inGitHubBlock = false;
        $currentBlockSpecificity = 0;

        foreach ($lines as $line) {
            $trimmed = trim($line);
            if ($trimmed === '' || str_starts_with($trimmed, '#')) {
                continue;
            }

            $lower = strtolower($trimmed);

            if (str_starts_with($lower, 'host ')) {
                $hosts = preg_split('/\s+/', substr($lower, 5));
                $inGitHubBlock = false;
                $currentBlockSpecificity = 0;
                foreach ($hosts as $h) {
                    if ($h === 'github.com') {
                        $inGitHubBlock = true;
                        $currentBlockSpecificity = 100;
                        break;
                    }
                    if ($h === '*') {
                        $inGitHubBlock = true;
                        $currentBlockSpecificity = 1;
                    } elseif (str_contains($h, '*') || str_contains($h, '?')) {
                        $inGitHubBlock = true;
                        $currentBlockSpecificity = 10;
                    }
                }
                continue;
            }

            if ($inGitHubBlock && str_starts_with($lower, 'identityfile ')) {
                $identityFile = trim(substr($trimmed, 13));
                $identityFile = trim($identityFile, '"');

                if (str_starts_with($identityFile, '~')) {
                    $identityFile = $home . substr($identityFile, 1);
                }

                if ($currentBlockSpecificity >= $bestSpecificity) {
                    $foundIdentity = $identityFile;
                    $bestSpecificity = $currentBlockSpecificity;
                }
            }
        }

        return $foundIdentity;
    }

    public static function detectSSHIssues(): array
    {
        $issues = [];
        $home = self::homeDir();
        if ($home === '') {
            return $issues;
        }

        $sshDir = $home . DIRECTORY_SEPARATOR . '.ssh';
        $defaultKeys = ['id_rsa.private', 'id_ed25519.private', 'id_ecdsa.private', 'id_dsa.private'];
        foreach ($defaultKeys as $key) {
            $path = $sshDir . DIRECTORY_SEPARATOR . $key;
            if (is_dir($path)) {
                $issues[] = "Conflict: '~/.ssh/$key' is a directory, but SSH expects it to be a key file. This can cause authentication failures.";
            }
        }

        return $issues;
    }

    public static function getConfiguredSSHCommand(): string
    {
        $result = self::exec('git config --get core.sshcommand');
        $output = trim($result['output']);
        if ($result['code'] !== 0 || $output === '') {
            if (PHP_OS_FAMILY === 'Windows') {
                $systemSSH = 'C:/Windows/System32/OpenSSH/ssh.exe';
                if (file_exists($systemSSH)) {
                    return $systemSSH;
                }
            }
            return 'ssh';
        }
        return $output;
    }

    public static function getLoadedKeys(): array
    {
        $result = self::exec('ssh-add -l');
        return ['output' => trim($result['output']), 'code' => $result['code']];
    }

    public static function getPublicKeyContent(string $privateKeyPath): array
    {
        // Try .public extension
        $pubPath = preg_replace('/\.private$/', '.public', $privateKeyPath);
        if ($pubPath !== $privateKeyPath) {
            $content = @file_get_contents($pubPath);
            if ($content !== false) {
                return ['content' => $content, 'path' => $pubPath];
            }
        }

        $pubPath = $privateKeyPath . '.public';
        $content = @file_get_contents($pubPath);
        if ($content !== false) {
            return ['content' => $content, 'path' => $pubPath];
        }

        // Try generating from private key
        $result = self::exec("ssh-keygen -y -f " . escapeshellarg($privateKeyPath));
        if ($result['code'] === 0) {
            return ['content' => $result['output'], 'path' => 'generated from private key'];
        }

        return ['content' => '', 'path' => '', 'error' => "could not find or generate public key for $privateKeyPath"];
    }

    private static function checkSSH(string $host, bool $nonInteractive): array
    {
        $sshCmd = self::getConfiguredSSHCommand();
        $args = '-T';
        if ($nonInteractive) {
            $args .= ' -o BatchMode=yes';
        }
        $args .= ' -o StrictHostKeyChecking=accept-new';
        $args .= ' ' . escapeshellarg($host);

        $result = self::exec("$sshCmd $args");
        $output = $result['output'];
        $success = str_contains($output, 'successfully authenticated');

        return ['success' => $success, 'output' => $output];
    }

    public static function homeDir(): string
    {
        return getenv('HOME') ?: (getenv('USERPROFILE') ?: '');
    }

    private static function exec(string $command): array
    {
        $output = '';
        $code = 1;
        exec($command . ' 2>&1', $lines, $code);
        $output = implode("\n", $lines);
        return ['output' => $output, 'code' => $code];
    }
}
