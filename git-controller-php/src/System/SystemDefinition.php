<?php

declare(strict_types=1);

namespace GitController\System;

final class SystemDefinition
{
    /**
     * @param string $systemVersion
     * @param string[] $postCloneHooks
     * @param array<string, array<string, Component>> $repos category => [name => Component]
     */
    public function __construct(
        public readonly string $systemVersion = '',
        public readonly array $postCloneHooks = [],
        public readonly array $repos = [],
    ) {}
}
