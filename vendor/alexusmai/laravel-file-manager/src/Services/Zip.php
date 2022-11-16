<?php

namespace Alexusmai\LaravelFileManager\Services;

use Alexusmai\LaravelFileManager\Events\UnzipCreated;
use Alexusmai\LaravelFileManager\Events\UnzipFailed;
use Alexusmai\LaravelFileManager\Events\ZipCreated;
use Alexusmai\LaravelFileManager\Events\ZipFailed;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Storage;
use RecursiveIteratorIterator;
use RecursiveDirectoryIterator;
use ZipArchive;

class Zip
{
    protected $zip;
    protected $request;
    //protected $pathPrefix;

    /**
     * Zip constructor.
     *
     * @param  ZipArchive  $zip
     * @param  Request  $request
     */
    public function __construct(ZipArchive $zip, Request $request)
    {
        $this->zip = $zip;
        $this->request = $request;
        //$this->pathPrefix = Storage::disk($request->input('disk'))->path();
            //->getDriver()
            //->getAdapter()
            //->getPathPrefix();
    }

    /**
     * Create new zip archive
     *
     * @return array
     */
    public function create(): array
    {

        if ($this->createArchive()) {
            return [
                'result' => [
                    'status'  => 'success',
                    'message' => null,
                ],
            ];
        }

        return [
            'result' => [
                'status'  => 'warning',
                'message' => 'zipError',
            ],
        ];
    }

    /**
     * Extract
     *
     * @return array
     */
    public function extract(): array
    {
        if ($this->extractArchive()) {
            return [
                'result' => [
                    'status'  => 'success',
                    'message' => null,
                ],
            ];
        }

        return [
            'result' => [
                'status'  => 'warning',
                'message' => 'zipError',
            ],
        ];
    }


    protected function prefixer($path): string
    {
        return Storage::disk($this->request->input('disk'))->path($path);
    }

    /**
     * Create zip archive
     *
     * @return bool
     */
    protected function createArchive(): bool
    {
        // elements list
        $elements = $this->request->input('elements');

        // create or overwrite archive
        if ($this->zip->open(
                $this->createName(),
                ZIPARCHIVE::OVERWRITE | ZIPARCHIVE::CREATE
            ) === true
        ) {
            // files processing
            if ($elements['files']) {
                foreach ($elements['files'] as $file) {
                    $this->zip->addFile(
                        $this->prefixer($file),
                        basename($file)
                    );
                }
            }

            // directories processing
            if ($elements['directories']) {
                $this->addDirs($elements['directories']);
            }

            $this->zip->close();

            event(new ZipCreated($this->request));

            return true;
        }

        event(new ZipFailed($this->request));

        return false;
    }

    /**
     * Archive extract
     *
     * @return bool
     */
    protected function extractArchive(): bool
    {
        $zipPath = $this->prefixer($this->request->input('path'));

        $rootPath = dirname($zipPath);

        // extract to new folder
        $folder = $this->request->input('folder');

        if ($this->zip->open($zipPath) === true) {
            $this->zip->extractTo($folder ? $rootPath.'/'.$folder : $rootPath);
            $this->zip->close();

            event(new UnzipCreated($this->request));

            return true;
        }

        event(new UnzipFailed($this->request));

        return false;
    }

    /**
     * Add directories - recursive
     *
     * @param  array  $directories
     */
    protected function addDirs(array $directories)
    {
        foreach ($directories as $directory) {

            // Create recursive directory iterator
            $files = new RecursiveIteratorIterator(
                new RecursiveDirectoryIterator($this->prefixer($directory)),
                RecursiveIteratorIterator::LEAVES_ONLY
            );

            foreach ($files as $name => $file) {
                // Get real and relative path for current item
                $filePath = $file->getRealPath();
                $relativePath = substr(
                    $filePath,
                    strlen($this->fullPath($this->request->input('path')))
                );

                if (!$file->isDir()) {
                    // Add current file to archive
                    $this->zip->addFile($filePath, $relativePath);
                } else {
                    // add empty folders
                    if (!glob($filePath.'/*')) {
                        $this->zip->addEmptyDir($relativePath);
                    }
                }
            }
        }
    }

    /**
     * Create archive name with full path
     *
     * @return string
     */
    protected function createName(): string
    {
        return $this->fullPath($this->request->input('path'))
            .$this->request->input('name');
    }

    /**
     * Generate full path
     *
     * @param $path
     *
     * @return string
     */
    protected function fullPath($path): string
    {
        return $path ? $this->prefixer($path).'/' : $this->prefixer('');
    }
}
