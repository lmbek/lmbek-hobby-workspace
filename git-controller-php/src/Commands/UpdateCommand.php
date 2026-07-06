<?php

declare(strict_types=1);

namespace GitController\Commands;

use GitController\UI\UI;

final class UpdateCommand
{
    public static function run(): void
    {
        UI::header('Update Repositories');

        FetchCommand::run();
        PullCommand::run();
        StatusCommand::run();

        UI::success('Update complete!');
    }
}
