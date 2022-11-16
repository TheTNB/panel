<?php

declare(strict_types=1);

namespace Zend\Diactoros;

use function func_get_args;
use function Laminas\Diactoros\marshalMethodFromSapi as laminas_marshalMethodFromSapi;

/**
 * @deprecated Use Laminas\Diactoros\marshalMethodFromSapi instead
 */
function marshalMethodFromSapi(array $server): string
{
    return laminas_marshalMethodFromSapi(...func_get_args());
}
