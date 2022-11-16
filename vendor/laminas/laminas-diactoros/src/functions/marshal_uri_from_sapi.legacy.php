<?php

declare(strict_types=1);

namespace Zend\Diactoros;

use Laminas\Diactoros\Uri;

use function func_get_args;
use function Laminas\Diactoros\marshalUriFromSapi as laminas_marshalUriFromSapi;

/**
 * @deprecated Use Laminas\Diactoros\marshalUriFromSapi instead
 */
function marshalUriFromSapi(array $server, array $headers): Uri
{
    return laminas_marshalUriFromSapi(...func_get_args());
}
