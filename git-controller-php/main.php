<?php

declare(strict_types=1);

require_once __DIR__ . '/vendor/autoload.php';

use GitController\Commands\CheckoutCommand;
use GitController\Commands\CloneCommand;
use GitController\Commands\DoctorCommand;
use GitController\Commands\FetchCommand;
use GitController\Commands\PullCommand;
use GitController\Commands\PushCommand;
use GitController\Commands\ScaffoldCommand;
use GitController\Commands\SshSetupCommand;
use GitController\Commands\StatusCommand;
use GitController\Commands\UpdateCommand;
use GitController\Commands\ValidateCommand;
use GitController\Commands\WsInitCommand;
use GitController\UI\UI;

// Set CLI name for doctor/ssh messages
DoctorCommand::$cliName = basename($argv[0]);

if ($argc < 2) {
    showHelp();
    exit(0);
}

$command = $argv[1];

$commands = [
    'init'      => [WsInitCommand::class, 'run'],
    'clone'     => [CloneCommand::class, 'run'],
    'fetch'     => [FetchCommand::class, 'run'],
    'pull'      => [PullCommand::class, 'run'],
    'push'      => [PushCommand::class, 'run'],
    'scaffold'  => [ScaffoldCommand::class, 'run'],
    'checkout'  => [CheckoutCommand::class, 'run'],
    'status'    => [StatusCommand::class, 'run'],
    'update'    => [UpdateCommand::class, 'run'],
    'validate'  => [ValidateCommand::class, 'run'],
    'doctor'    => [DoctorCommand::class, 'run'],
    'ssh-setup' => [SshSetupCommand::class, 'run'],
    'ssh'       => [SshSetupCommand::class, 'run'],
];

if ($command === 'help') {
    showHelp();
    exit(0);
}

if (isset($commands[$command])) {
    try {
        call_user_func($commands[$command]);
    } catch (Throwable $e) {
        UI::error('%s', $e->getMessage());
        exit(1);
    }
} else {
    UI::error('Unknown command: %s', $command);
    showHelp();
    exit(1);
}

function showHelp(): void
{
    $cliName = basename($GLOBALS['argv'][0]);

    printf("\n%sWorkspace Controller (PHP)%s\n", UI::COLOR_BOLD, UI::COLOR_RESET);
    echo str_repeat('=', 30) . "\n";
    printf("Usage: %sphp %s%s [command]\n", UI::COLOR_BOLD, $cliName, UI::COLOR_RESET);

    echo "\nWorkflow Commands:\n";
    echo "  init       Scaffold a new workspace (system-definition.yaml, Makefile, .gitignore)\n";
    echo "  clone      Clone all repositories defined in system-definition.yaml\n";
    echo "  fetch      Fetch all remotes across all repositories\n";
    echo "  pull       Pull latest changes across all repositories (clone if missing)\n";
    echo "  push       Push local commits across all repositories\n";
    echo "  scaffold   Initialise .git and set remote origin (no clone/fetch needed)\n";
    echo "  checkout   Switch all repositories to their defined branch\n";
    echo "  status     Show dashboard overview of all repository states\n";
    echo "  update     Fetch, pull, and show status across all repositories\n";
    echo "  validate   Validate repository consistency against the definition\n";

    echo "\nSetup Commands:\n";
    echo "  doctor     Diagnose environment (Git, PHP, SSH, Docker)\n";
    echo "  ssh-setup  Interactive SSH key management (alias: ssh)\n";
    echo "  help       Show this help\n";
}
