<?php

declare(strict_types=1);

namespace GitController\Commands;

use GitController\GitUtil\GitUtil;
use GitController\System\Parser;
use GitController\UI\UI;
use RuntimeException;

final class CheckoutCommand
{
    public static function run(): void
    {
        UI::header('Checkout Branches');

        [$sys, $workspace] = Parser::loadDefinition('git-repositories/system-definition.yaml');

        $overrideBranch = getenv('BRANCH') ?: '';
        $hasErrors = false;

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

                $branch = $overrideBranch !== '' ? $overrideBranch : $comp->version;

                UI::info('→ %s → %s', $displayName, $branch);
                try {
                    GitUtil::checkout($targetPath, $branch);
                } catch (RuntimeException $e) {
                    UI::error('Failed to checkout %s: %s', $displayName, $e->getMessage());
                    $hasErrors = true;
                }
            }
        }

        if ($hasErrors) {
            throw new RuntimeException('checkout completed with errors');
        }

        UI::success('All repositories on correct branches!');
    }
}
