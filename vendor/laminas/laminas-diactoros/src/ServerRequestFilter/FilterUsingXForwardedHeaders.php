<?php

declare(strict_types=1);

namespace Laminas\Diactoros\ServerRequestFilter;

use Laminas\Diactoros\Exception\InvalidForwardedHeaderNameException;
use Laminas\Diactoros\Exception\InvalidProxyAddressException;
use Psr\Http\Message\ServerRequestInterface;

use function explode;
use function filter_var;
use function in_array;
use function is_string;
use function str_contains;
use function strtolower;

use const FILTER_FLAG_IPV4;
use const FILTER_FLAG_IPV6;
use const FILTER_VALIDATE_IP;

/**
 * Modify the URI to reflect the X-Forwarded-* headers.
 *
 * If the request comes from a trusted proxy, this filter will analyze the
 * various X-Forwarded-* headers, if any, and if they are marked as trusted,
 * in order to return a new request that composes a URI instance that reflects
 * those headers.
 *
 * @psalm-immutable
*/
final class FilterUsingXForwardedHeaders implements FilterServerRequestInterface
{
    public const HEADER_HOST  = 'X-FORWARDED-HOST';
    public const HEADER_PORT  = 'X-FORWARDED-PORT';
    public const HEADER_PROTO = 'X-FORWARDED-PROTO';

    private const X_FORWARDED_HEADERS = [
        self::HEADER_HOST,
        self::HEADER_PORT,
        self::HEADER_PROTO,
    ];

    /**
     * Only allow construction via named constructors
     *
     * @param list<non-empty-string> $trustedProxies
     * @param list<FilterUsingXForwardedHeaders::HEADER_*> $trustedHeaders
     */
    private function __construct(private array $trustedProxies = [], private array $trustedHeaders = [])
    {
    }

    public function __invoke(ServerRequestInterface $request): ServerRequestInterface
    {
        $remoteAddress = $request->getServerParams()['REMOTE_ADDR'] ?? '';

        if ('' === $remoteAddress || ! is_string($remoteAddress)) {
            // Should we trigger a warning here?
            return $request;
        }

        if (! $this->isFromTrustedProxy($remoteAddress)) {
            // Do nothing
            return $request;
        }

        // Update the URI based on the trusted headers
        $uri = $originalUri = $request->getUri();
        foreach ($this->trustedHeaders as $headerName) {
            $header = $request->getHeaderLine($headerName);
            if ('' === $header || str_contains($header, ',')) {
                // Reject empty headers and/or headers with multiple values
                continue;
            }

            switch ($headerName) {
                case self::HEADER_HOST:
                    $uri = $uri->withHost($header);
                    break;
                case self::HEADER_PORT:
                    $uri = $uri->withPort((int) $header);
                    break;
                case self::HEADER_PROTO:
                    $scheme = strtolower($header) === 'https' ? 'https' : 'http';
                    $uri    = $uri->withScheme($scheme);
                    break;
            }
        }

        if ($uri !== $originalUri) {
            return $request->withUri($uri);
        }

        return $request;
    }

    /**
     * Indicate which proxies and which X-Forwarded headers to trust.
     *
     * @param list<non-empty-string> $proxyCIDRList Each element may
     *     be an IP address or a subnet specified using CIDR notation; both IPv4
     *     and IPv6 are supported. The special string "*" will be translated to
     *     two entries, "0.0.0.0/0" and "::/0". An empty list indicates no
     *     proxies are trusted.
     * @param list<FilterUsingXForwardedHeaders::HEADER_*> $trustedHeaders If
     *     the list is empty, all X-Forwarded headers are trusted.
     * @throws InvalidProxyAddressException
     * @throws InvalidForwardedHeaderNameException
     */
    public static function trustProxies(
        array $proxyCIDRList,
        array $trustedHeaders = self::X_FORWARDED_HEADERS
    ): self {
        $proxyCIDRList = self::normalizeProxiesList($proxyCIDRList);
        self::validateTrustedHeaders($trustedHeaders);

        return new self($proxyCIDRList, $trustedHeaders);
    }

    /**
     * Trust any X-FORWARDED-* headers from any address.
     *
     * This is functionally equivalent to calling `trustProxies(['*'])`.
     *
     * WARNING: Only do this if you know for certain that your application
     * sits behind a trusted proxy that cannot be spoofed. This should only
     * be the case if your server is not publicly addressable, and all requests
     * are routed via a reverse proxy (e.g., a load balancer, a server such as
     * Caddy, when using Traefik, etc.).
     */
    public static function trustAny(): self
    {
        return self::trustProxies(['*']);
    }

    /**
     * Trust X-Forwarded headers from reserved subnetworks.
     *
     * This is functionally equivalent to calling `trustProxies()` where the
     * `$proxcyCIDRList` argument is a list with the following:
     *
     * - 10.0.0.0/8
     * - 127.0.0.0/8
     * - 172.16.0.0/12
     * - 192.168.0.0/16
     * - ::1/128 (IPv6 localhost)
     * - fc00::/7 (IPv6 private networks)
     * - fe80::/10 (IPv6 local-link addresses)
     *
     * @param list<FilterUsingXForwardedHeaders::HEADER_*> $trustedHeaders If
     *     the list is empty, all X-Forwarded headers are trusted.
     * @throws InvalidForwardedHeaderNameException
     */
    public static function trustReservedSubnets(array $trustedHeaders = self::X_FORWARDED_HEADERS): self
    {
        return self::trustProxies([
            '10.0.0.0/8',
            '127.0.0.0/8',
            '172.16.0.0/12',
            '192.168.0.0/16',
            '::1/128', // ipv6 localhost
            'fc00::/7', // ipv6 private networks
            'fe80::/10', // ipv6 local-link addresses
        ], $trustedHeaders);
    }

    private function isFromTrustedProxy(string $remoteAddress): bool
    {
        foreach ($this->trustedProxies as $proxy) {
            if (IPRange::matches($remoteAddress, $proxy)) {
                return true;
            }
        }

        return false;
    }

    /** @throws InvalidForwardedHeaderNameException */
    private static function validateTrustedHeaders(array $headers): void
    {
        foreach ($headers as $header) {
            if (! in_array($header, self::X_FORWARDED_HEADERS, true)) {
                throw InvalidForwardedHeaderNameException::forHeader($header);
            }
        }
    }

    /**
     * @param list<non-empty-string> $proxyCIDRList
     * @return list<non-empty-string>
     * @throws InvalidProxyAddressException
     */
    private static function normalizeProxiesList(array $proxyCIDRList): array
    {
        $foundWildcard = false;

        foreach ($proxyCIDRList as $index => $cidr) {
            if ($cidr === '*') {
                unset($proxyCIDRList[$index]);
                $foundWildcard = true;
                continue;
            }

            if (! self::validateProxyCIDR($cidr)) {
                throw InvalidProxyAddressException::forAddress($cidr);
            }
        }

        if ($foundWildcard) {
            $proxyCIDRList[] = '0.0.0.0/0';
            $proxyCIDRList[] = '::/0';
        }

        return $proxyCIDRList;
    }

    private static function validateProxyCIDR(mixed $cidr): bool
    {
        if (! is_string($cidr) || '' === $cidr) {
            return false;
        }

        $address = $cidr;
        $mask    = null;
        if (str_contains($cidr, '/')) {
            [$address, $mask] = explode('/', $cidr, 2);
            $mask             = (int) $mask;
        }

        if (str_contains($address, ':')) {
            // is IPV6
            return filter_var($address, FILTER_VALIDATE_IP, FILTER_FLAG_IPV6)
                && (
                    $mask === null
                    || (
                        $mask <= 128
                        && $mask >= 0
                    )
                );
        }

        // is IPV4
        return filter_var($address, FILTER_VALIDATE_IP, FILTER_FLAG_IPV4)
            && (
                $mask === null
                || (
                    $mask <= 32
                    && $mask >= 0
                )
            );
    }
}
