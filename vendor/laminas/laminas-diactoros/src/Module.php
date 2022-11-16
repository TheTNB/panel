<?php

declare(strict_types=1);

namespace Laminas\Diactoros;

class Module
{
    public function getConfig(): array
    {
        return [
            'service_manager' => (new ConfigProvider())->getDependencies(),
        ];
    }
}
