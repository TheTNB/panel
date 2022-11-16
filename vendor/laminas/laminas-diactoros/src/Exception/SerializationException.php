<?php

declare(strict_types=1);

namespace Laminas\Diactoros\Exception;

use UnexpectedValueException;

class SerializationException extends UnexpectedValueException implements ExceptionInterface
{
    public static function forInvalidRequestLine(): self
    {
        return new self('Invalid request line detected');
    }

    public static function forInvalidStatusLine(): self
    {
        return new self('No status line detected');
    }
}
