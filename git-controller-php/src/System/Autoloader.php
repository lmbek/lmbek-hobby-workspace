<?php

declare(strict_types=1);

namespace GitController\System;

/**
 * Minimal PSR-4 Autoloader to remove hard dependency on Composer.
 */
final class Autoloader
{
    /**
     * @param array<string, string> $prefixes Prefix => Directory mapping
     */
    public function __construct(
        private readonly array $prefixes = []
    ) {}

    public static function register(string $prefix, string $baseDir): void
    {
        $loader = new self([$prefix => $baseDir]);
        spl_autoload_register([$loader, 'loadClass']);
    }

    public function loadClass(string $class): bool
    {
        foreach ($this->prefixes as $prefix => $baseDir) {
            if (str_starts_with($class, $prefix)) {
                $relativeClass = substr($class, strlen($prefix));
                $file = $baseDir . DIRECTORY_SEPARATOR . str_replace('\\', DIRECTORY_SEPARATOR, $relativeClass) . '.php';

                if (file_exists($file)) {
                    require $file;
                    return true;
                }
            }
        }

        return false;
    }
}
