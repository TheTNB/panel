<?php

declare(strict_types=1);

namespace Zend\Diactoros;

use function func_get_args;
use function Laminas\Diactoros\normalizeUploadedFiles as laminas_normalizeUploadedFiles;

/**
 * @deprecated Use Laminas\Diactoros\normalizeUploadedFiles instead
 */
function normalizeUploadedFiles(array $files): array
{
    return laminas_normalizeUploadedFiles(...func_get_args());
}
