<?php

declare(strict_types=1);

namespace Laminas\Diactoros\Response;

use Laminas\Diactoros\Exception;
use Laminas\Diactoros\Response;
use Laminas\Diactoros\Stream;
use Psr\Http\Message\StreamInterface;

use function gettype;
use function is_object;
use function is_string;
use function sprintf;

/**
 * Plain text response.
 *
 * Allows creating a response by passing a string to the constructor;
 * by default, sets a status code of 200 and sets the Content-Type header to
 * text/plain.
 */
class TextResponse extends Response
{
    use InjectContentTypeTrait;

    /**
     * Create a plain text response.
     *
     * Produces a text response with a Content-Type of text/plain and a default
     * status of 200.
     *
     * @param string|StreamInterface $text String or stream for the message body.
     * @param int $status Integer status code for the response; 200 by default.
     * @param array $headers Array of headers to use at initialization.
     * @throws Exception\InvalidArgumentException If $text is neither a string or stream.
     */
    public function __construct($text, int $status = 200, array $headers = [])
    {
        parent::__construct(
            $this->createBody($text),
            $status,
            $this->injectContentType('text/plain; charset=utf-8', $headers)
        );
    }

    /**
     * Create the message body.
     *
     * @param string|StreamInterface $text
     * @throws Exception\InvalidArgumentException If $text is neither a string or stream.
     */
    private function createBody($text): StreamInterface
    {
        if ($text instanceof StreamInterface) {
            return $text;
        }

        if (! is_string($text)) {
            throw new Exception\InvalidArgumentException(sprintf(
                'Invalid content (%s) provided to %s',
                is_object($text) ? $text::class : gettype($text),
                self::class
            ));
        }

        $body = new Stream('php://temp', 'wb+');
        $body->write($text);
        $body->rewind();
        return $body;
    }
}
