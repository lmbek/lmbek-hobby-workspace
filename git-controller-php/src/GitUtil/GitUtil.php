<?php

declare(strict_types=1);

namespace GitController\GitUtil;

use GitController\SshUtil\SshUtil;
use GitController\UI\UI;
use RuntimeException;

final class GitUtil
{
    public static function clone(string $repo, string $targetPath): void
    {
        self::validateRemoteURL($repo);
        self::runGit('clone', $repo, $targetPath);
    }

    public static function fetch(string $repoDir): void
    {
        self::runGitInDir($repoDir, 'fetch', '--all');
    }

    public static function pull(string $repoDir): void
    {
        self::runGitInDir($repoDir, 'fetch', '--all');
        self::runGitInDir($repoDir, 'pull');
    }

    public static function checkout(string $repoDir, string $branch): void
    {
        self::runGitInDir($repoDir, 'checkout', $branch);
    }

    public static function push(string $repoDir): void
    {
        if (!self::hasCommits($repoDir)) {
            $branch = self::currentBranch($repoDir);
            if ($branch !== '') {
                self::ensureGitIdentity($repoDir);
                self::runGitInDir($repoDir, 'commit', '--allow-empty', '-m', 'Initial commit');
            }
        }

        if (!self::hasUpstream($repoDir)) {
            $branch = self::currentBranch($repoDir);
            if ($branch !== '') {
                self::runGitInDir($repoDir, 'push', '-u', 'origin', $branch);
                return;
            }
        }
        self::runGitInDir($repoDir, 'push');
    }

    public static function initAndLink(string $targetPath, string $repo): void
    {
        self::validateRemoteURL($repo);
        self::runGitInDir($targetPath, 'init');
        self::runGitInDir($targetPath, 'remote', 'add', 'origin', $repo);
        self::runGitInDir($targetPath, 'fetch', 'origin');

        // Attempt to align with the remote main branch
        if (self::tryRunGitInDir($targetPath, 'reset', '--soft', 'origin/main')) {
            self::tryRunGitInDir($targetPath, 'branch', '--set-upstream-to=origin/main', 'main');
        } elseif (self::tryRunGitInDir($targetPath, 'reset', '--soft', 'origin/master')) {
            self::tryRunGitInDir($targetPath, 'checkout', '-b', 'main');
            self::tryRunGitInDir($targetPath, 'branch', '--set-upstream-to=origin/master', 'main');
        } else {
            self::tryRunGitInDir($targetPath, 'checkout', '-b', 'main');
        }
    }

    public static function scaffold(string $targetPath, string $repo): void
    {
        self::validateRemoteURL($repo);
        self::runGitInDir($targetPath, 'init');
        self::runGitInDir($targetPath, 'remote', 'add', 'origin', $repo);
        self::runGitInDir($targetPath, 'fetch', 'origin');

        if (!self::tryRunGitInDir($targetPath, 'checkout', '-b', 'main', '--track', 'origin/main')) {
            if (!self::tryRunGitInDir($targetPath, 'checkout', '-b', 'master', '--track', 'origin/master')) {
                self::tryRunGitInDir($targetPath, 'checkout', '-b', 'main');
            }
        }
    }

    public static function hasOutgoingCommits(string $repoDir): bool
    {
        if (!self::hasCommits($repoDir)) {
            return true;
        }
        if (!self::hasUpstream($repoDir)) {
            return true;
        }

        $result = self::execGitInDir($repoDir, 'rev-list', '--count', '@{u}..HEAD');
        if ($result['code'] !== 0) {
            return true;
        }
        return trim($result['output']) !== '0';
    }

    public static function isCloned(string $targetPath): bool
    {
        return is_dir($targetPath . DIRECTORY_SEPARATOR . '.git');
    }

    public static function isNonEmptyDir(string $path): bool
    {
        if (!is_dir($path)) {
            return false;
        }
        $entries = scandir($path);
        // scandir returns . and .. at minimum
        return $entries !== false && count($entries) > 2;
    }

    public static function ensureDir(string $path): void
    {
        if (!is_dir($path)) {
            mkdir($path, 0755, true);
        }
    }

    // --- Private helpers ---

    private static function validateRemoteURL(string $repo): void
    {
        if (str_starts_with($repo, 'http://') || str_starts_with($repo, 'https://')) {
            throw new RuntimeException(
                "insecure repository URL detected: $repo. This project strictly enforces SSH for Git operations. "
                . "Please update your system-definition.yaml to use SSH URLs (e.g., git@github.com:...) or local paths"
            );
        }
    }

    private static function ensureGitIdentity(string $repoDir): void
    {
        $result = self::execGitInDir($repoDir, 'config', 'user.name');
        if ($result['code'] !== 0) {
            throw new RuntimeException(
                "git author identity not configured. Please run:\n\n"
                . "  git config --global user.name \"Your Name\"\n"
                . "  git config --global user.email \"you@example.com\""
            );
        }
    }

    private static function hasCommits(string $repoDir): bool
    {
        $result = self::execGitInDir($repoDir, 'rev-parse', 'HEAD');
        return $result['code'] === 0;
    }

    private static function hasUpstream(string $repoDir): bool
    {
        $result = self::execGitInDir($repoDir, 'rev-parse', '--abbrev-ref', '--symbolic-full-name', '@{u}');
        return $result['code'] === 0;
    }

    private static function currentBranch(string $repoDir): string
    {
        $result = self::execGitInDir($repoDir, 'symbolic-ref', '--short', 'HEAD');
        if ($result['code'] !== 0) {
            return '';
        }
        return trim($result['output']);
    }

    private static function configureGitEnv(): void
    {
        static $configured = false;
        if ($configured) {
            return;
        }
        $configured = true;

        $sshCmd = SshUtil::getConfiguredSSHCommand();

        if (PHP_OS_FAMILY === 'Windows') {
            if (DIRECTORY_SEPARATOR === '\\' && str_contains($sshCmd, ' ')) {
                $sshCmd = str_replace('\\', '/', $sshCmd);
                $sshCmd = "\"$sshCmd\"";
            }
        } else {
            if (str_contains($sshCmd, ' ') && !str_starts_with($sshCmd, '"')) {
                $sshCmd = "\"$sshCmd\"";
            }
        }

        putenv('GIT_SSH_COMMAND=' . $sshCmd);
        // Remove GIT_SSH / GIT_SSH_VARIANT to avoid conflicts.
        putenv('GIT_SSH');
        putenv('GIT_SSH_VARIANT');
    }

    public static function runGit(string ...$args): void
    {
        self::configureGitEnv();
        $cmd = 'git ' . implode(' ', array_map('escapeshellarg', $args));
        $result = self::execCmd($cmd);
        if ($result['code'] !== 0) {
            throw new RuntimeException(self::enhanceGitError($result['output']));
        }
    }

    public static function runGitInDir(string $dir, string ...$args): void
    {
        self::configureGitEnv();
        $cmd = 'git -C ' . escapeshellarg($dir) . ' ' . implode(' ', array_map('escapeshellarg', $args));
        $result = self::execCmd($cmd);
        if ($result['code'] !== 0) {
            throw new RuntimeException(self::enhanceGitError($result['output']));
        }
    }

    private static function tryRunGitInDir(string $dir, string ...$args): bool
    {
        try {
            self::runGitInDir($dir, ...$args);
            return true;
        } catch (RuntimeException) {
            return false;
        }
    }

    private static function execGitInDir(string $dir, string ...$args): array
    {
        $cmd = 'git -C ' . escapeshellarg($dir) . ' ' . implode(' ', array_map('escapeshellarg', $args));
        return self::execCmd($cmd);
    }

    private static function execCmd(string $command): array
    {
        $lines = [];
        $code = 1;
        exec($command . ' 2>&1', $lines, $code);
        return ['output' => implode("\n", $lines), 'code' => $code];
    }

    private static function enhanceGitError(string $errStr): string
    {
        $msg = $errStr;

        if (str_contains($errStr, 'Permission denied (publickey)')) {
            $msg .= "\n\n" . UI::COLOR_RED . "SSH authentication failed." . UI::COLOR_RESET;
            $msg .= "\nGitHub rejected your SSH key. Please ensure it's added to your account.";

            $identityFile = SshUtil::getGitHubIdentityFile();
            if ($identityFile !== '') {
                $msg .= "\nIdentityFile: $identityFile";
            }

            if (PHP_OS_FAMILY === 'Windows') {
                $msg .= "\n\nWindows Tip: Ensure the 'OpenSSH Authentication Agent' service is running and your key is loaded (run 'make ssh' Option 2).";
                $msg .= "\nAlso verify that your 'core.sshCommand' is set to the System OpenSSH (run 'make ssh' Option 4).";
            }

            $msg .= "\n\nRun 'make doctor' for more diagnostics.";
        } elseif (str_contains($errStr, 'Host key verification failed')) {
            $msg .= "\n\n" . UI::COLOR_RED . "Host key verification failed." . UI::COLOR_RESET;
            $msg .= "\nGitHub's host key is missing from your known_hosts file.";
            $msg .= "\nRun 'make ssh' Option 5 to fix this automatically.";
        }

        return $msg;
    }
}
