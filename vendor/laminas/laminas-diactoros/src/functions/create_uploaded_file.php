<?php

declare(strict_types=1);

namespace Laminas\Diactoros;

use function sprintf;

/**
 * Create an uploaded file instance from an array of values.
 *
 * @param array $spec A single $_FILES entry.
 * @throws Exception\InvalidArgumentException If one or more of the tmp_name,
 *     size, or error keys are missing from $spec.
 */
function createUploadedFile(array $spec): UploadedFile
{
    if (
        ! isset($spec['tmp_name'])
        || ! isset($spec['size'])
        || ! isset($spec['error'])
    ) {
        throw new Exception\InvalidArgumentException(sprintf(
            '$spec provided to %s MUST contain each of the keys "tmp_name",'
            . ' "size", and "error"; one or more were missing',
            __FUNCTION__
        ));
    }

    return new UploadedFile(
        $spec['tmp_name'],
        (int) $spec['size'],
        $spec['error'],
        $spec['name'] ?? null,
        $spec['type'] ?? null
    );
}
