<?php

declare(strict_types=1);

namespace Laminas\Diactoros;

use Laminas\Diactoros\ServerRequestFilter\FilterServerRequestInterface;
use Laminas\Diactoros\ServerRequestFilter\FilterUsingXForwardedHeaders;
use Psr\Http\Message\ServerRequestFactoryInterface;
use Psr\Http\Message\ServerRequestInterface;

use function array_change_key_case;
use function array_key_exists;
use function explode;
use function gettype;
use function implode;
use function is_array;
use function is_bool;
use function is_callable;
use function is_string;
use function ltrim;
use function preg_match;
use function preg_replace;
use function sprintf;
use function str_contains;
use function strlen;
use function strrpos;
use function strtolower;
use function substr;

use const CASE_LOWER;

/**
 * Class for marshaling a request object from the current PHP environment.
 */
class ServerRequestFactory implements ServerRequestFactoryInterface
{
    /**
     * Function to use to get apache request headers; present only to simplify mocking.
     *
     * @var callable
     */
    private static $apacheRequestHeaders = 'apache_request_headers';

    /**
     * Create a request from the supplied superglobal values.
     *
     * If any argument is not supplied, the corresponding superglobal value will
     * be used.
     *
     * The ServerRequest created is then passed to the fromServer() method in
     * order to marshal the request URI and headers.
     *
     * @see fromServer()
     *
     * @param array $server $_SERVER superglobal
     * @param array $query $_GET superglobal
     * @param array $body $_POST superglobal
     * @param array $cookies $_COOKIE superglobal
     * @param array $files $_FILES superglobal
     * @param null|FilterServerRequestInterface $requestFilter If present, the
     *     generated request will be passed to this instance and the result
     *     returned by this method. When not present, a default instance of
     *     FilterUsingXForwardedHeaders is created, using the `trustReservedSubnets()`
     *     constructor.
     */
    public static function fromGlobals(
        ?array $server = null,
        ?array $query = null,
        ?array $body = null,
        ?array $cookies = null,
        ?array $files = null,
        ?FilterServerRequestInterface $requestFilter = null
    ): ServerRequest {
        $requestFilter = $requestFilter ?: FilterUsingXForwardedHeaders::trustReservedSubnets();

        $server  = normalizeServer(
            $server ?: $_SERVER,
            is_callable(self::$apacheRequestHeaders) ? self::$apacheRequestHeaders : null
        );
        $files   = normalizeUploadedFiles($files ?: $_FILES);
        $headers = marshalHeadersFromSapi($server);

        if (null === $cookies && array_key_exists('cookie', $headers)) {
            $cookies = parseCookieHeader($headers['cookie']);
        }

        return $requestFilter(new ServerRequest(
            $server,
            $files,
            self::marshalUriFromSapi($server, $headers),
            marshalMethodFromSapi($server),
            'php://input',
            $headers,
            $cookies ?: $_COOKIE,
            $query ?: $_GET,
            $body ?: $_POST,
            marshalProtocolVersionFromSapi($server)
        ));
    }

    /**
     * {@inheritDoc}
     */
    public function createServerRequest(string $method, $uri, array $serverParams = []): ServerRequestInterface
    {
        $uploadedFiles = [];

        return new ServerRequest(
            $serverParams,
            $uploadedFiles,
            $uri,
            $method,
            'php://temp'
        );
    }

    /**
     * Marshal a Uri instance based on the values present in the $_SERVER array and headers.
     *
     * @param array<string, string|list<string>> $headers
     * @param array $server SAPI parameters
     */
    private static function marshalUriFromSapi(array $server, array $headers): Uri
    {
        $uri = new Uri('');

        // URI scheme
        $https = false;
        if (array_key_exists('HTTPS', $server)) {
            $https = self::marshalHttpsValue($server['HTTPS']);
        } elseif (array_key_exists('https', $server)) {
            $https = self::marshalHttpsValue($server['https']);
        }

        $uri = $uri->withScheme($https ? 'https' : 'http');

        // Set the host
        [$host, $port] = self::marshalHostAndPort($server, $headers);
        if (! empty($host)) {
            $uri = $uri->withHost($host);
            if (! empty($port)) {
                $uri = $uri->withPort($port);
            }
        }

        // URI path
        $path = self::marshalRequestPath($server);

        // Strip query string
        $path = explode('?', $path, 2)[0];

        // URI query
        $query = '';
        if (isset($server['QUERY_STRING'])) {
            $query = ltrim((string) $server['QUERY_STRING'], '?');
        }

        // URI fragment
        $fragment = '';
        if (str_contains($path, '#')) {
            [$path, $fragment] = explode('#', $path, 2);
        }

        return $uri
            ->withPath($path)
            ->withFragment($fragment)
            ->withQuery($query);
    }

    /**
     * Marshal the host and port from the PHP environment.
     *
     * @param array<string, string|list<string>> $headers
     * @return array{string, int|null} Array of two items, host and port,
     *     in that order (can be passed to a list() operation).
     */
    private static function marshalHostAndPort(array $server, array $headers): array
    {
        static $defaults = ['', null];

        $host = self::getHeaderFromArray('host', $headers, false);
        if ($host !== false) {
            // Ignore obviously malformed host headers:
            // - Whitespace is invalid within a hostname and break the URI representation within HTTP.
            //   non-printable characters other than SPACE and TAB are already rejected by HeaderSecurity.
            // - A comma indicates that multiple host headers have been sent which is not legal
            //   and might be used in an attack where a load balancer sees a different host header
            //   than Diactoros.
            if (! preg_match('/[\\t ,]/', $host)) {
                return self::marshalHostAndPortFromHeader($host);
            }
        }

        if (! isset($server['SERVER_NAME'])) {
            return $defaults;
        }

        $host = (string) $server['SERVER_NAME'];
        $port = isset($server['SERVER_PORT']) ? (int) $server['SERVER_PORT'] : null;

        if (
            ! isset($server['SERVER_ADDR'])
            || ! preg_match('/^\[[0-9a-fA-F\:]+\]$/', $host)
        ) {
            return [$host, $port];
        }

        // Misinterpreted IPv6-Address
        // Reported for Safari on Windows
        return self::marshalIpv6HostAndPort($server, $port);
    }

    /**
     * @return array{string, int|null} Array of two items, host and port,
     *     in that order (can be passed to a list() operation).
     */
    private static function marshalIpv6HostAndPort(array $server, ?int $port): array
    {
        $host             = '[' . (string) $server['SERVER_ADDR'] . ']';
        $port             = $port ?: 80;
        $portSeparatorPos = strrpos($host, ':');

        if (false === $portSeparatorPos) {
            return [$host, $port];
        }

        if ($port . ']' === substr($host, $portSeparatorPos + 1)) {
            // The last digit of the IPv6-Address has been taken as port
            // Unset the port so the default port can be used
            $port = null;
        }
        return [$host, $port];
    }

    /**
     * Detect the path for the request
     *
     * Looks at a variety of criteria in order to attempt to autodetect the base
     * request path, including:
     *
     * - IIS7 UrlRewrite environment
     * - REQUEST_URI
     * - ORIG_PATH_INFO
     */
    private static function marshalRequestPath(array $server): string
    {
        // IIS7 with URL Rewrite: make sure we get the unencoded url
        // (double slash problem).
        $iisUrlRewritten = $server['IIS_WasUrlRewritten'] ?? null;
        $unencodedUrl    = $server['UNENCODED_URL'] ?? '';
        if ('1' === $iisUrlRewritten && is_string($unencodedUrl) && '' !== $unencodedUrl) {
            return $unencodedUrl;
        }

        $requestUri = $server['REQUEST_URI'] ?? null;

        if (is_string($requestUri)) {
            return preg_replace('#^[^/:]+://[^/]+#', '', $requestUri);
        }

        $origPathInfo = $server['ORIG_PATH_INFO'] ?? '';
        if (! is_string($origPathInfo) || '' === $origPathInfo) {
            return '/';
        }

        return $origPathInfo;
    }

    private static function marshalHttpsValue(mixed $https): bool
    {
        if (is_bool($https)) {
            return $https;
        }

        if (! is_string($https)) {
            throw new Exception\InvalidArgumentException(sprintf(
                'SAPI HTTPS value MUST be a string or boolean; received %s',
                gettype($https)
            ));
        }

        return 'on' === strtolower($https);
    }

    /**
     * @param string|list<string> $host
     * @return array Array of two items, host and port, in that order (can be
     *     passed to a list() operation).
     */
    private static function marshalHostAndPortFromHeader($host): array
    {
        if (is_array($host)) {
            $host = implode(', ', $host);
        }

        $port = null;

        // works for regname, IPv4 & IPv6
        if (preg_match('|\:(\d+)$|', $host, $matches)) {
            $host = substr($host, 0, -1 * (strlen($matches[1]) + 1));
            $port = (int) $matches[1];
        }

        return [$host, $port];
    }

    /**
     * Retrieve a header value from an array of headers using a case-insensitive lookup.
     *
     * @template T
     * @param array<string, string|list<string>> $headers Key/value header pairs
     * @param T $default Default value to return if header not found
     * @return string|T
     */
    private static function getHeaderFromArray(string $name, array $headers, $default = null)
    {
        $header  = strtolower($name);
        $headers = array_change_key_case($headers, CASE_LOWER);
        if (! array_key_exists($header, $headers)) {
            return $default;
        }

        if (is_string($headers[$header])) {
            return $headers[$header];
        }

        return implode(', ', $headers[$header]);
    }
}
