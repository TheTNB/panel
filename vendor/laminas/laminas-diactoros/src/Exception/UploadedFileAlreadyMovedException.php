<?php

declare(strict_types=1);

namespace Laminas\Diactoros\Exception;

use RuntimeException;
use Throwable;

class UploadedFileAlreadyMovedException extends RuntimeException implements ExceptionInterface
{
    /** {@inheritDoc} */
    public function __construct(
        string $message = 'Cannot retrieve stream after it has already moved',
        int $code = 0,
        ?Throwable $previous = null
    ) {
        parent::__construct($message, $code, $previous);
    }
}
