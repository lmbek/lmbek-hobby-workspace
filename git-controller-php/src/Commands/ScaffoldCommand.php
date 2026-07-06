<?php

declare(strict_types=1);

namespace GitController\Commands;

use GitController\GitUtil\GitUtil;
use GitController\System\Parser;
use GitController\UI\UI;
use RuntimeException;

final class ScaffoldCommand
{
    public static function run(): void
    {
        UI::header('Scaffold Repositories');

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
                    UI::info('✓ %s (already initialised)', $displayName);
                    continue;
                }

                GitUtil::ensureDir($targetPath);
                UI::info('→ Scaffolding %s', $displayName);

                try {
                    GitUtil::scaffold($targetPath, $repo);
                } catch (RuntimeException $e) {
                    UI::error('Failed to scaffold %s: %s', $displayName, $e->getMessage());
                    $hasErrors = true;
                }
            }
        }

        if ($hasErrors) {
            throw new RuntimeException('scaffold completed with errors');
        }

        UI::success('All repositories scaffolded — git init, remote origin set, fetched, and default branch configured.');
    }
}
