<?php

declare(strict_types=1);

namespace Zend\Diactoros;

use function func_get_args;
use function Laminas\Diactoros\parseCookieHeader as laminas_parseCookieHeader;

/**
 * @deprecated Use {@see \Laminas\Diactoros\parseCookieHeader} instead
 *
 * @param string $cookieHeader A string cookie header value.
 * @return array<non-empty-string, string> key/value cookie pairs.
 */
function parseCookieHeader($cookieHeader): array
{
    return laminas_parseCookieHeader(...func_get_args());
}
