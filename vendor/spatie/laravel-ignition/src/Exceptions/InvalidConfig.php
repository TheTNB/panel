<?php

namespace Spatie\LaravelIgnition\Exceptions;

use Exception;
use Monolog\Logger;
use Spatie\Ignition\Contracts\BaseSolution;
use Spatie\Ignition\Contracts\ProvidesSolution;
use Spatie\Ignition\Contracts\Solution;

class InvalidConfig extends Exception implements ProvidesSolution
{
    public static function invalidLogLevel(string $logLevel): self
    {
        return new self("Invalid log level `{$logLevel}` specified.");
    }

    public function getSolution(): Solution
    {
        $validLogLevels = array_map(
            fn (string $level) => strtolower($level),
            array_keys(Logger::getLevels())
        );

        $validLogLevelsString = implode(',', $validLogLevels);

        return BaseSolution::create()
            ->setSolutionTitle('You provided an invalid log level')
            ->setSolutionDescription("Please change the log level in your `config/logging.php` file. Valid log levels are {$validLogLevelsString}.");
    }
}
