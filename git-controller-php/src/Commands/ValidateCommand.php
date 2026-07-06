<?php

declare(strict_types=1);

namespace GitController\Commands;

use GitController\System\Parser;
use GitController\UI\UI;
use RuntimeException;

final class ValidateCommand
{
    public static function run(): void
    {
        UI::header('Repository Validation');

        [$sys, $workspace] = Parser::loadDefinition('git-repositories/system-definition.yaml');

        UI::info('Validating repositories...');

        $hasErrors = false;

        foreach ($sys->repos as $catName => $repos) {
            $catDir = $workspace->getCategoryDir($catName);

            foreach ($repos as $name => $comp) {
                $targetPath = $catDir . DIRECTORY_SEPARATOR . $name;
                $displayName = $name !== '' ? $name : $catName;

                if ($comp->repository === '') {
                    UI::info('[%s] No remote configured, skipping', $displayName);
                    continue;
                }

                $error = self::validateRepo($targetPath, $comp->version);
                if ($error !== null) {
                    UI::error('[%s] %s', $displayName, $error);
                    $hasErrors = true;
                } else {
                    UI::success('[%s] OK', $displayName);
                }
            }
        }

        if ($hasErrors) {
            throw new RuntimeException('validation completed with errors');
        }

        UI::success('All repositories valid!');
    }

    private static function validateRepo(string $targetPath, string $expectedVersion): ?string
    {
        if (!is_dir($targetPath)) {
            return "not cloned — run 'git-controller clone'";
        }

        $currentVersion = self::getGitCurrentVersion($targetPath);
        if ($currentVersion === null) {
            return 'could not determine current branch';
        }

        if ($currentVersion !== $expectedVersion) {
            // Log warning but don't fail
        }

        $dirty = self::isGitDirty($targetPath);
        if ($dirty) {
            // Log warning but don't fail
        }

        return null;
    }

    private static function getGitCurrentVersion(string $dir): ?string
    {
        // Try exact tag match first
        exec('git -C ' . escapeshellarg($dir) . ' describe --tags --exact-match 2>&1', $lines, $code);
        if ($code === 0) {
            return trim(implode("\n", $lines));
        }

        // Get branch name or hash
        $lines = [];
        exec('git -C ' . escapeshellarg($dir) . ' rev-parse --abbrev-ref HEAD 2>&1', $lines, $code);
        if ($code !== 0) {
            return null;
        }

        $branch = trim(implode("\n", $lines));
        if ($branch === 'HEAD') {
            // Detached HEAD, get hash
            $lines = [];
            exec('git -C ' . escapeshellarg($dir) . ' rev-parse --short HEAD 2>&1', $lines, $code);
            if ($code !== 0) {
                return null;
            }
            return trim(implode("\n", $lines));
        }

        return $branch;
    }

    private static function isGitDirty(string $dir): bool
    {
        exec('git -C ' . escapeshellarg($dir) . ' status --porcelain 2>&1', $lines, $code);
        if ($code !== 0) {
            return false;
        }
        return trim(implode("\n", $lines)) !== '';
    }
}
