<?php

declare(strict_types=1);

namespace Zend\Diactoros;

use Laminas\Diactoros\UploadedFile;

use function func_get_args;
use function Laminas\Diactoros\createUploadedFile as laminas_createUploadedFile;

/**
 * @deprecated Use Laminas\Diactoros\createUploadedFile instead
 */
function createUploadedFile(array $spec): UploadedFile
{
    return laminas_createUploadedFile(...func_get_args());
}
