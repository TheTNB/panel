<?php

declare(strict_types=1);

namespace Laminas\Diactoros\ServerRequestFilter;

use function explode;
use function inet_pton;
use function intval;
use function ip2long;
use function pack;
use function sprintf;
use function str_contains;
use function str_pad;
use function str_repeat;
use function substr_compare;
use function unpack;

/** @internal */
final class IPRange
{
    /**
     * Disable instantiation
     */
    private function __construct()
    {
    }

    /** @psalm-pure */
    public static function matches(string $ip, string $cidr): bool
    {
        if (str_contains($ip, ':')) {
            return self::matchesIPv6($ip, $cidr);
        }

        return self::matchesIPv4($ip, $cidr);
    }

    /** @psalm-pure */
    public static function matchesIPv4(string $ip, string $cidr): bool
    {
        $mask   = 32;
        $subnet = $cidr;

        if (str_contains($cidr, '/')) {
            [$subnet, $mask] = explode('/', $cidr, 2);
            $mask            = (int) $mask;
        }

        if ($mask < 0 || $mask > 32) {
            return false;
        }

        $ip     = ip2long($ip);
        $subnet = ip2long($subnet);
        if (false === $ip || false === $subnet) {
            // Invalid data
            return false;
        }

        return 0 === substr_compare(
            sprintf("%032b", $ip),
            sprintf("%032b", $subnet),
            0,
            $mask
        );
    }

    /** @psalm-pure */
    public static function matchesIPv6(string $ip, string $cidr): bool
    {
        $mask   = 128;
        $subnet = $cidr;

        if (str_contains($cidr, '/')) {
            [$subnet, $mask] = explode('/', $cidr, 2);
            $mask            = (int) $mask;
        }

        if ($mask < 0 || $mask > 128) {
            return false;
        }

        $ip     = inet_pton($ip);
        $subnet = inet_pton($subnet);

        if (false === $ip || false === $subnet) {
            // Invalid data
            return false;
        }

        // mask 0: if it's a valid IP, it's valid
        if ($mask === 0) {
            return (bool) unpack('n*', $ip);
        }

        // @see http://stackoverflow.com/questions/7951061/matching-ipv6-address-to-a-cidr-subnet, MW answer
        $binMask = str_repeat("f", intval($mask / 4));
        switch ($mask % 4) {
            case 0:
                break;
            case 1:
                $binMask .= "8";
                break;
            case 2:
                $binMask .= "c";
                break;
            case 3:
                $binMask .= "e";
                break;
        }

        $binMask = str_pad($binMask, 32, '0');
        $binMask = pack("H*", $binMask);

        return ($ip & $binMask) === $subnet;
    }
}
