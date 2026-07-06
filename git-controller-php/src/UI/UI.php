<?php

declare(strict_types=1);

namespace GitController\UI;

final class UI
{
    public const COLOR_RESET  = "\033[0m";
    public const COLOR_RED    = "\033[31m";
    public const COLOR_GREEN  = "\033[32m";
    public const COLOR_YELLOW = "\033[33m";
    public const COLOR_BLUE   = "\033[34m";
    public const COLOR_CYAN   = "\033[36m";
    public const COLOR_BOLD   = "\033[1m";

    public static function header(string $text): void
    {
        printf("\n%s%s%s\n%s\n", self::COLOR_BOLD, self::COLOR_CYAN, strtoupper($text), str_repeat('=', strlen($text)));
    }

    public static function info(string $format, mixed ...$args): void
    {
        printf("%s%s%s\n", self::COLOR_BLUE, sprintf($format, ...$args), self::COLOR_RESET);
    }

    public static function success(string $format, mixed ...$args): void
    {
        printf("%s✔ %s%s\n", self::COLOR_GREEN, sprintf($format, ...$args), self::COLOR_RESET);
    }

    public static function warn(string $format, mixed ...$args): void
    {
        printf("%s⚠ %s%s\n", self::COLOR_YELLOW, sprintf($format, ...$args), self::COLOR_RESET);
    }

    public static function error(string $format, mixed ...$args): void
    {
        printf("%s✖ %s%s\n", self::COLOR_RED, sprintf($format, ...$args), self::COLOR_RESET);
    }

    public static function step(int $num, string $text): void
    {
        printf("\n%s[%d] %s%s\n", self::COLOR_BOLD, $num, $text, self::COLOR_RESET);
    }

    public static function note(string $text): void
    {
        printf("\n%sNote: %s%s\n", self::COLOR_CYAN, $text, self::COLOR_RESET);
    }
}
