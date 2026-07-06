<?php

declare(strict_types=1);

namespace GitController\Commands;

use GitController\GitUtil\GitUtil;
use GitController\System\Parser;
use GitController\UI\UI;
use RuntimeException;

final class PushCommand
{
    public static function run(): void
    {
        UI::header('Push Repositories');

        [$sys, $workspace] = Parser::loadDefinition('git-repositories/system-definition.yaml');

        $hasErrors = false;

        // Push the workspace repository itself
        $wsRoot = $workspace->root !== '' ? $workspace->root : '.';
        $wsRoot = realpath($wsRoot) ?: $wsRoot;
        if (GitUtil::isCloned($wsRoot)) {
            if (!GitUtil::hasOutgoingCommits($wsRoot)) {
                UI::info('✓ workspace (nothing to push)');
            } else {
                UI::info('→ Pushing workspace');
                try {
                    GitUtil::push($wsRoot);
                } catch (RuntimeException $e) {
                    UI::error('Failed to push workspace: %s', $e->getMessage());
                    $hasErrors = true;
                }
            }
        }

        foreach ($sys->repos as $catName => $repos) {
            if (count($repos) === 0) {
                continue;
            }

            $catDir = realpath($workspace->getCategoryDir($catName)) ?: $workspace->getCategoryDir($catName);

            foreach ($repos as $name => $comp) {
                if ($comp->repository === '' || str_contains($comp->repository, '@company')) {
                    continue;
                }

                $displayName = $name !== '' ? $name : $catName;
                $targetPath = $catDir . DIRECTORY_SEPARATOR . $name;

                if (!GitUtil::isCloned($targetPath)) {
                    UI::warn('⊘ %s (not cloned, skipping)', $displayName);
                    continue;
                }

                if (!GitUtil::hasOutgoingCommits($targetPath)) {
                    UI::info('✓ %s (nothing to push)', $displayName);
                    continue;
                }

                UI::info('→ Pushing %s', $displayName);
                try {
                    GitUtil::push($targetPath);
                } catch (RuntimeException $e) {
                    UI::error('Failed to push %s: %s', $displayName, $e->getMessage());
                    $hasErrors = true;
                }
            }
        }

        if ($hasErrors) {
            throw new RuntimeException('push completed with errors');
        }

        UI::success('All repositories pushed successfully!');
    }
}
