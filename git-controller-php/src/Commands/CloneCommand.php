<?php

declare(strict_types=1);

namespace GitController\Commands;

use GitController\GitUtil\GitUtil;
use GitController\System\Parser;
use GitController\UI\UI;
use RuntimeException;

final class CloneCommand
{
    public static function run(): void
    {
        UI::header('Clone Repositories');

        [$sys, $workspace] = Parser::loadDefinition('git-repositories/system-definition.yaml');

        $hasErrors = false;

        foreach ($sys->repos as $catName => $repos) {
            if (count($repos) === 0) {
                continue;
            }

            UI::step(2, "Category: $catName");
            $catDir = realpath($workspace->getCategoryDir($catName)) ?: $workspace->getCategoryDir($catName);
            GitUtil::ensureDir($catDir);

            foreach ($repos as $name => $comp) {
                $repo = $comp->repository;
                if ($repo === '' || str_contains($repo, '@company')) {
                    continue;
                }

                $displayName = $name !== '' ? $name : $catName;
                $targetPath = $catDir . DIRECTORY_SEPARATOR . $name;

                if (GitUtil::isCloned($targetPath)) {
                    UI::info('✓ %s (already cloned)', $displayName);
                    continue;
                }

                UI::info('→ Cloning %s', $displayName);

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
            }
        }

        if ($hasErrors) {
            throw new RuntimeException('clone completed with errors');
        }

        UI::success('All repositories cloned successfully!');
    }
}
