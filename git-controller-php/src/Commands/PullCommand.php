<?php

declare(strict_types=1);

namespace GitController\Commands;

use GitController\GitUtil\GitUtil;
use GitController\System\Parser;
use GitController\UI\UI;
use RuntimeException;

final class PullCommand
{
    public static function run(): void
    {
        UI::header('Pull Repositories');

        [$sys, $workspace] = Parser::loadDefinition('git-repositories/system-definition.yaml');

        $hasErrors = false;

        $catNum = 0;
        foreach ($sys->repos as $catName => $repos) {
            if (count($repos) === 0) {
                continue;
            }
            $catNum++;

            UI::step($catNum, "Category: $catName");
            $catDir = realpath($workspace->getCategoryDir($catName)) ?: $workspace->getCategoryDir($catName);
            GitUtil::ensureDir($catDir);

            foreach ($repos as $name => $comp) {
                $repo = $comp->repository;
                if ($repo === '' || str_contains($repo, '@company')) {
                    continue;
                }

                $displayName = $name !== '' ? $name : $catName;
                $targetPath = $catDir . DIRECTORY_SEPARATOR . $name;

                if (!GitUtil::isCloned($targetPath)) {
                    UI::info('→ Cloning %s (not yet cloned)', $displayName);
                    try {
                        if (GitUtil::isNonEmptyDir($targetPath)) {
                            GitUtil::initAndLink($targetPath, $repo);
                        } else {
                            GitUtil::clone($repo, $targetPath);
                        }
                    } catch (RuntimeException $e) {
                        UI::error('Failed to clone %s: %s', $displayName, $e->getMessage());
                        $hasErrors = true;
                    }
                } else {
                    UI::info('→ Pulling %s', $displayName);
                    try {
                        GitUtil::pull($targetPath);
                    } catch (RuntimeException $e) {
                        UI::error('Failed to pull %s: %s', $displayName, $e->getMessage());
                        $hasErrors = true;
                    }
                }
            }
        }

        if ($hasErrors) {
            throw new RuntimeException('pull completed with errors');
        }

        UI::success('All repositories up to date!');
    }
}
