<?php

declare(strict_types=1);

namespace Zend\Diactoros;

use function func_get_args;
use function Laminas\Diactoros\marshalHeadersFromSapi as laminas_marshalHeadersFromSapi;

/**
 * @deprecated Use Laminas\Diactoros\marshalHeadersFromSapi instead
 */
function marshalHeadersFromSapi(array $server): array
{
    return laminas_marshalHeadersFromSapi(...func_get_args());
}
