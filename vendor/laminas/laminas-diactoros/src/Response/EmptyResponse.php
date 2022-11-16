<?php

declare(strict_types=1);

namespace Laminas\Diactoros\Response;

use Laminas\Diactoros\Response;
use Laminas\Diactoros\Stream;

/**
 * A class representing empty HTTP responses.
 */
class EmptyResponse extends Response
{
    /**
     * Create an empty response with the given status code.
     *
     * @param int $status Status code for the response, if any.
     * @param array $headers Headers for the response, if any.
     */
    public function __construct(int $status = 204, array $headers = [])
    {
        $body = new Stream('php://temp', 'r');
        parent::__construct($body, $status, $headers);
    }

    /**
     * Create an empty response with the given headers.
     *
     * @param array $headers Headers for the response.
     */
    public static function withHeaders(array $headers): EmptyResponse
    {
        return new static(204, $headers);
    }
}
