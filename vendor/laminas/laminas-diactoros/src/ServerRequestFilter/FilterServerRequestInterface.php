<?php

declare(strict_types=1);

namespace Laminas\Diactoros\ServerRequestFilter;

use Psr\Http\Message\ServerRequestInterface;

/**
 * Filter/initialize a server request.
 *
 * Implementations of this interface will take an incoming request, and
 * decide if additional modifications are necessary. As examples:
 *
 * - Injecting a unique request identifier header.
 * - Using the X-Forwarded-* headers to rewrite the URI to reflect the original request.
 * - Using the Forwarded header to rewrite the URI to reflect the original request.
 *
 * This functionality is consumed by the ServerRequestFactory using the request
 * instance it generates, just prior to returning a request.
 */
interface FilterServerRequestInterface
{
    /**
     * Determine if a request needs further modification, and if so, return a
     * new instance reflecting those modifications.
     */
    public function __invoke(ServerRequestInterface $request): ServerRequestInterface;
}
