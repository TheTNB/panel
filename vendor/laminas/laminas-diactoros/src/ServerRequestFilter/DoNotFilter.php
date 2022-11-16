<?php

declare(strict_types=1);

namespace Laminas\Diactoros\ServerRequestFilter;

use Psr\Http\Message\ServerRequestInterface;

final class DoNotFilter implements FilterServerRequestInterface
{
    public function __invoke(ServerRequestInterface $request): ServerRequestInterface
    {
        return $request;
    }
}
