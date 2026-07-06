<?php

declare(strict_types=1);

namespace GitController\System;

final class Workspace
{
    public function __construct(
        public readonly string $root = '',
    ) {}

    public function getCategoryDir(string $catName): string
    {
        $root = $this->root !== '' ? $this->root : '.';

        return $root . DIRECTORY_SEPARATOR . 'git-repositories' . DIRECTORY_SEPARATOR . $catName;
    }
}
