<?php

namespace Alexusmai\LaravelFileManager\Services\TransferService;

use Alexusmai\LaravelFileManager\Traits\PathTrait;
use Illuminate\Support\Facades\Storage;
use League\Flysystem\FilesystemException;
use League\Flysystem\MountManager;

class ExternalTransfer extends Transfer
{
    use PathTrait;

    /**
     * @var MountManager
     */
    public $manager;

    /**
     * ExternalTransfer constructor.
     *
     * @param $disk
     * @param $path
     * @param $clipboard
     */
    public function __construct($disk, $path, $clipboard)
    {
        parent::__construct($disk, $path, $clipboard);

        $this->manager = new MountManager([
            'from' => Storage::drive($clipboard['disk'])->getDriver(),
            'to'   => Storage::drive($disk)->getDriver(),
        ]);
    }

    /**
     * Copy files and folders
     *
     * @return void
     * @throws FilesystemException
     */
    protected function copy()
    {
        // files
        foreach ($this->clipboard['files'] as $file) {
            $this->copyToDisk(
                $file,
                $this->renamePath($file, $this->path)
            );
        }

        // directories
        foreach ($this->clipboard['directories'] as $directory) {
            $this->copyDirectoryToDisk($directory);
        }
    }

    /**
     * Cut files and folders
     *
     * @return void
     * @throws FilesystemException
     */
    protected function cut()
    {
        // files
        foreach ($this->clipboard['files'] as $file) {
            $this->moveToDisk(
                $file,
                $this->renamePath($file, $this->path)
            );
        }

        // directories
        foreach ($this->clipboard['directories'] as $directory) {
            $this->copyDirectoryToDisk($directory);

            // remove directory
            Storage::disk($this->clipboard['disk'])
                ->deleteDirectory($directory);
        }
    }

    /**
     * Copy directory to another disk
     *
     * @param $directory
     *
     * @return void
     * @throws FilesystemException
     */
    protected function copyDirectoryToDisk($directory)
    {
        // get all directories in this directory
        $allDirectories = Storage::disk($this->clipboard['disk'])
            ->allDirectories($directory);

        $partsForRemove = count(explode('/', $directory)) - 1;

        // create this directories
        foreach ($allDirectories as $dir) {
            Storage::disk($this->disk)->makeDirectory(
                $this->transformPath($dir, $this->path, $partsForRemove)
            );
        }

        // get all files
        $allFiles = Storage::disk($this->clipboard['disk'])
            ->allFiles($directory);

        // copy files
        foreach ($allFiles as $file) {
            $this->copyToDisk($file,
                $this->transformPath($file, $this->path, $partsForRemove));
        }
    }

    /**
     * Copy files to disk
     *
     * @param $filePath
     * @param $newPath
     *
     * @return void
     * @throws FilesystemException
     */
    protected function copyToDisk($filePath, $newPath)
    {
        $this->manager->copy(
            'from://'.$filePath,
            'to://'.$newPath
        );
    }

    /**
     * Move files to disk
     *
     * @param $filePath
     * @param $newPath
     *
     * @return void
     * @throws FilesystemException
     */
    protected function moveToDisk($filePath, $newPath)
    {
        $this->manager->move(
            'from://'.$filePath,
            'to://'.$newPath
        );
    }
}
