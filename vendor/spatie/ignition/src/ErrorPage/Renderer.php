<?php

namespace Spatie\Ignition\ErrorPage;

class Renderer
{
    /**
     * @param array<string, mixed> $data
     *
     * @return void
     */
    public function render(array $data): void
    {
        $viewFile = __DIR__ . '/../../resources/views/errorPage.php';

        extract($data, EXTR_OVERWRITE);

        include $viewFile;
    }
}
