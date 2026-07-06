<?php

declare(strict_types=1);

namespace GitController\Commands;

use GitController\GitUtil\GitUtil;
use GitController\System\Parser;
use GitController\UI\UI;

final class StatusCommand
{
    public static function run(): void
    {
        UI::header('Repository Status');

        [$sys, $workspace] = Parser::loadDefinition('git-repositories/system-definition.yaml');

        $totalRepos = 0;
        $cloned = 0;
        $dirty = 0;
        $ahead = 0;
        $behind = 0;

        foreach ($sys->repos as $catName => $repos) {
            if (count($repos) === 0) {
                continue;
            }

            $catDir = realpath($workspace->getCategoryDir($catName)) ?: $workspace->getCategoryDir($catName);

            UI::step(0, "Category: $catName");

            foreach ($repos as $name => $comp) {
                $totalRepos++;

                $displayName = $name !== '' ? $name : $catName;

                if ($comp->repository === '' || str_contains($comp->repository, '@company')) {
                    UI::info('  %-30s  %s', $displayName, '⊘ no remote configured');
                    continue;
                }

                $targetPath = $catDir . DIRECTORY_SEPARATOR . $name;

                if (!GitUtil::isCloned($targetPath)) {
                    UI::warn('  %-30s  %s', $displayName, 'not cloned');
                    continue;
                }
                $cloned++;

                $branch = self::getBranch($targetPath);
                $isDirty = self::hasDirtyFiles($targetPath);
                [$aheadN, $behindN] = self::getAheadBehind($targetPath);

                if ($isDirty) {
                    $dirty++;
                }
                if ($aheadN > 0) {
                    $ahead++;
                }
                if ($behindN > 0) {
                    $behind++;
                }

                $parts = ["branch:$branch"];
                if ($isDirty) {
                    $parts[] = UI::COLOR_YELLOW . 'dirty' . UI::COLOR_RESET;
                } else {
                    $parts[] = UI::COLOR_GREEN . 'clean' . UI::COLOR_RESET;
                }
                if ($aheadN > 0) {
                    $parts[] = "↑$aheadN";
                }
                if ($behindN > 0) {
                    $parts[] = "↓$behindN";
                }

                UI::info('  %-30s  %s', $displayName, implode('  ', $parts));
            }
        }

        echo "\n";
        UI::info('Total: %d repos | Cloned: %d | Dirty: %d | Ahead: %d | Behind: %d',
            $totalRepos, $cloned, $dirty, $ahead, $behind);

        if ($dirty > 0) {
            UI::warn('Some repositories have uncommitted changes.');
        }
        if ($behind > 0) {
            UI::warn("Some repositories are behind remote. Run 'make pull' to update.");
        }
        if ($cloned < $totalRepos) {
            UI::warn("Some repositories are not cloned. Run 'make clone' to set up.");
        }
    }

    private static function getBranch(string $dir): string
    {
        exec('git -C ' . escapeshellarg($dir) . ' rev-parse --abbrev-ref HEAD 2>&1', $lines, $code);
        if ($code !== 0) {
            return 'unknown';
        }
        return trim(implode("\n", $lines));
    }

    private static function hasDirtyFiles(string $dir): bool
    {
        exec('git -C ' . escapeshellarg($dir) . ' status --porcelain 2>&1', $lines, $code);
        if ($code !== 0) {
            return false;
        }
        return trim(implode("\n", $lines)) !== '';
    }

    /**
     * @return array{int, int}
     */
    private static function getAheadBehind(string $dir): array
    {
        exec('git -C ' . escapeshellarg($dir) . ' rev-list --left-right --count HEAD...@{upstream} 2>&1', $lines, $code);
        if ($code !== 0) {
            return [0, 0];
        }
        $parts = preg_split('/\s+/', trim(implode("\n", $lines)));
        if (count($parts) === 2) {
            return [(int) $parts[0], (int) $parts[1]];
        }
        return [0, 0];
    }
}
