<?php

declare(strict_types=1);

namespace Laminas\Diactoros;

use Psr\Http\Message\RequestFactoryInterface;
use Psr\Http\Message\ResponseFactoryInterface;
use Psr\Http\Message\ServerRequestFactoryInterface;
use Psr\Http\Message\StreamFactoryInterface;
use Psr\Http\Message\UploadedFileFactoryInterface;
use Psr\Http\Message\UriFactoryInterface;

class ConfigProvider
{
    public const CONFIG_KEY                  = 'laminas-diactoros';
    public const X_FORWARDED                 = 'x-forwarded-request-filter';
    public const X_FORWARDED_TRUSTED_PROXIES = 'trusted-proxies';
    public const X_FORWARDED_TRUSTED_HEADERS = 'trusted-headers';

    /**
     * Retrieve configuration for laminas-diactoros.
     *
     * @return array
     */
    public function __invoke(): array
    {
        return [
            'dependencies'   => $this->getDependencies(),
            self::CONFIG_KEY => $this->getComponentConfig(),
        ];
    }

    /**
     * Returns the container dependencies.
     * Maps factory interfaces to factories.
     */
    public function getDependencies(): array
    {
        // @codingStandardsIgnoreStart
        return [
            'invokables' => [
                RequestFactoryInterface::class => RequestFactory::class,
                ResponseFactoryInterface::class => ResponseFactory::class,
                StreamFactoryInterface::class => StreamFactory::class,
                ServerRequestFactoryInterface::class => ServerRequestFactory::class,
                UploadedFileFactoryInterface::class => UploadedFileFactory::class,
                UriFactoryInterface::class => UriFactory::class
            ],
        ];
        // @codingStandardsIgnoreEnd
    }

    public function getComponentConfig(): array
    {
        return [
            self::X_FORWARDED => [
                self::X_FORWARDED_TRUSTED_PROXIES => '',
                self::X_FORWARDED_TRUSTED_HEADERS => [],
            ],
        ];
    }
}
