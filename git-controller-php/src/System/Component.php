<?php

declare(strict_types=1);

namespace GitController\System;

final class Component
{
    public function __construct(
        public readonly string $repository = '',
        public readonly string $version = '',
    ) {}
}
