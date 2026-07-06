<?php

declare(strict_types=1);

namespace GitController\Commands;

use GitController\GitUtil\GitUtil;
use GitController\System\Parser;
use GitController\UI\UI;
use RuntimeException;

final class FetchCommand
{
    public static function run(): void
    {
        UI::header('Fetch Repositories');

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

            foreach ($repos as $name => $comp) {
                if ($comp->repository === '' || str_contains($comp->repository, '@company')) {
                    continue;
                }

                $displayName = $name !== '' ? $name : $catName;
                $targetPath = $catDir . DIRECTORY_SEPARATOR . $name;

                if (!GitUtil::isCloned($targetPath)) {
                    UI::info('→ Skipping %s (not cloned)', $displayName);
                    continue;
                }

                UI::info('→ Fetching %s', $displayName);
                try {
                    GitUtil::fetch($targetPath);
                } catch (RuntimeException $e) {
                    UI::error('Failed to fetch %s: %s', $displayName, $e->getMessage());
                    $hasErrors = true;
                }
            }
        }

        if ($hasErrors) {
            throw new RuntimeException('fetch completed with errors');
        }

        UI::success('All repositories fetched!');
    }
}
