<?php

declare(strict_types=1);

namespace GitController\Commands;

use GitController\UI\UI;
use RuntimeException;

final class WsInitCommand
{
    private const DEFAULT_DEFINITION = <<<'YAML'
# System Definition for Workspace Controller
system-version: main

hooks:
    post-clone:
        - echo "Clone complete! Run 'make status' to see the state of all repositories."

# Each key under "repos" is a category folder inside git-repositories/.
# Categories are dynamic — add whatever makes sense for your project.
repos:
    # applications:
    #     my-service:
    #         repository: git@github.com:org/my-service.git
    #         version: main
YAML;

    private static function workspaceRoot(): string
    {
        $cwd = getcwd();
        if ($cwd !== false) {
            return dirname($cwd);
        }
        return '.';
    }

    public static function run(): void
    {
        UI::header('Initialise Workspace');

        $root = self::workspaceRoot();

        $gitReposDir = $root . DIRECTORY_SEPARATOR . 'git-repositories';
        if (!is_dir($gitReposDir)) {
            if (!mkdir($gitReposDir, 0755, true)) {
                throw new RuntimeException('failed to create git-repositories/');
            }
        }

        $defPath = $gitReposDir . DIRECTORY_SEPARATOR . 'system-definition.yaml';

        if (file_exists($defPath)) {
            UI::info('system-definition.yaml already exists at %s — skipping creation', $defPath);
        } else {
            if (file_put_contents($defPath, self::DEFAULT_DEFINITION) === false) {
                throw new RuntimeException("failed to create $defPath");
            }
            UI::success('Created %s', $defPath);
        }

        $gitignorePath = $root . DIRECTORY_SEPARATOR . '.gitignore';
        if (!file_exists($gitignorePath)) {
            file_put_contents($gitignorePath, "git-repositories/\n");
            UI::success('Created .gitignore');
        }

        $makefile = $root . DIRECTORY_SEPARATOR . 'Makefile';
        if (!file_exists($makefile)) {
            $content = <<<'MAKEFILE'
.PHONY: init clone pull push scaffold checkout status validate doctor ssh

init:
	cd git-controller-php && php main.php init

clone:
	cd git-controller-php && php main.php clone

pull:
	cd git-controller-php && php main.php pull

push:
	cd git-controller-php && php main.php push

scaffold:
	cd git-controller-php && php main.php scaffold

checkout:
	cd git-controller-php && php main.php checkout

status:
	cd git-controller-php && php main.php status

validate:
	cd git-controller-php && php main.php validate

doctor:
	cd git-controller-php && php main.php doctor

ssh:
	cd git-controller-php && php main.php ssh
MAKEFILE;
            if (file_put_contents($makefile, $content) === false) {
                throw new RuntimeException('failed to create Makefile');
            }
            UI::success('Created Makefile');
        }

        UI::success("Workspace initialised! Edit git-repositories/system-definition.yaml to add your repositories, then run 'make clone'.");
    }
}
