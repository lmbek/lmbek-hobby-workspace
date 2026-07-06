<?php

declare(strict_types=1);

namespace GitController\System;

use RuntimeException;

final class Parser
{
    /**
     * @return array{SystemDefinition, Workspace}
     */
    public static function loadDefinition(string $path): array
    {
        // resolveRoot returns the workspace root from a found definition file path.
        // The definition lives inside git-repositories/, so the workspace root is
        // two levels up from the file.
        $resolveRoot = static fn(string $defPath): string => dirname(dirname($defPath));

        // 1. Try path as-is
        if (file_exists($path)) {
            return self::loadFile($path, $resolveRoot($path));
        }

        // 2. Try searching parent directories for the workspace folder
        $currentDir = getcwd();
        if ($currentDir !== false) {
            while (true) {
                $altPath = $currentDir . DIRECTORY_SEPARATOR . $path;
                if (file_exists($altPath)) {
                    return self::loadFile($altPath, $resolveRoot($altPath));
                }
                $parent = dirname($currentDir);
                if ($parent === $currentDir) {
                    break;
                }
                $currentDir = $parent;
            }
        }

        throw new RuntimeException("could not find system definition: $path");
    }

    /**
     * @return array{SystemDefinition, Workspace}
     */
    private static function loadFile(string $path, string $root): array
    {
        $content = file_get_contents($path);
        if ($content === false) {
            throw new RuntimeException("failed to read file: $path");
        }

        // Expand environment variables (same as Go's os.ExpandEnv)
        $content = preg_replace_callback('/\$\{(\w+)\}|\$(\w+)/', static function (array $matches): string {
            $varName = $matches[1] !== '' ? $matches[1] : $matches[2];
            return getenv($varName) ?: '';
        }, $content);

        if (!function_exists('yaml_parse')) {
            // Fallback: simple YAML parser for the subset we need
            $data = self::parseSimpleYaml($content);
        } else {
            $data = yaml_parse($content);
        }

        if (!is_array($data)) {
            throw new RuntimeException("failed to parse YAML from: $path");
        }

        $root = realpath($root) ?: $root;

        return [self::buildDefinition($data), new Workspace(root: $root)];
    }

    private static function buildDefinition(array $data): SystemDefinition
    {
        $systemVersion = (string) ($data['system-version'] ?? '');
        $postCloneHooks = (array) ($data['hooks']['post-clone'] ?? []);

        $repos = [];
        foreach (($data['repos'] ?? []) as $catName => $catData) {
            if (!is_array($catData)) {
                continue;
            }

            // Check if this is a single component (flat/singleton)
            if (isset($catData['repository']) || isset($catData['version'])) {
                $repos[(string) $catName] = [
                    '' => new Component(
                        repository: (string) ($catData['repository'] ?? ''),
                        version: (string) ($catData['version'] ?? ''),
                    ),
                ];
                continue;
            }

            $components = [];
            foreach ($catData as $name => $compData) {
                if (!is_array($compData)) {
                    continue;
                }
                $components[(string) $name] = new Component(
                    repository: (string) ($compData['repository'] ?? ''),
                    version: (string) ($compData['version'] ?? ''),
                );
            }
            $repos[(string) $catName] = $components;
        }

        return new SystemDefinition(
            systemVersion: $systemVersion,
            postCloneHooks: $postCloneHooks,
            repos: $repos,
        );
    }

    /**
     * Minimal YAML parser for the system-definition.yaml subset.
     * Handles the specific structure used by the workspace controller.
     */
    private static function parseSimpleYaml(string $content): array
    {
        $result = [];
        $lines = explode("\n", str_replace("\r\n", "\n", $content));
        $stack = [&$result];
        $indentStack = [-1];

        foreach ($lines as $line) {
            $trimmed = trim($line);
            if ($trimmed === '' || str_starts_with($trimmed, '#')) {
                continue;
            }

            $indent = strlen($line) - strlen(ltrim($line));

            // Pop stack to find parent
            while (count($indentStack) > 1 && $indent <= end($indentStack)) {
                array_pop($stack);
                array_pop($indentStack);
            }

            if (preg_match('/^(\S.*?):\s*$/', $trimmed, $m)) {
                // Key with no value — start a new map
                $key = $m[1];
                $stack[count($stack) - 1][$key] = [];
                $stack[] = &$stack[count($stack) - 1][$key];
                $indentStack[] = $indent;
            } elseif (preg_match('/^(\S.*?):\s+(.+)$/', $trimmed, $m)) {
                // Key with inline value
                $key = $m[1];
                $value = trim($m[2]);
                // Remove surrounding quotes
                if ((str_starts_with($value, '"') && str_ends_with($value, '"'))
                    || (str_starts_with($value, "'") && str_ends_with($value, "'"))) {
                    $value = substr($value, 1, -1);
                }
                $stack[count($stack) - 1][$key] = $value;
            } elseif (preg_match('/^-\s+(.+)$/', $trimmed, $m)) {
                // List item
                $value = trim($m[1]);
                if ((str_starts_with($value, '"') && str_ends_with($value, '"'))
                    || (str_starts_with($value, "'") && str_ends_with($value, "'"))) {
                    $value = substr($value, 1, -1);
                }
                $stack[count($stack) - 1][] = $value;
            }
        }

        return $result;
    }
}
