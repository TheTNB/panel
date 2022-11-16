<?php

namespace Alexusmai\LaravelFileManager\Traits;

use Alexusmai\LaravelFileManager\Services\ACLService\ACL;
use Illuminate\Support\Facades\Storage;
use League\Flysystem\AwsS3V3\AwsS3V3Adapter;
use League\Flysystem\FilesystemException;

trait ContentTrait
{

    /**
     * Get content for the selected disk and path
     *
     * @param $disk
     * @param  null  $path
     *
     * @return array
     * @throws FilesystemException
     */
    public function getContent($disk, $path = null): array
    {
        $content = Storage::disk($disk)->listContents($path ?: '')->toArray();

        $directories = $this->filterDir($disk, $content);
        $files = $this->filterFile($disk, $content);

        return compact('directories', 'files');
    }

    /**
     * Get directories with properties
     *
     * @param $disk
     * @param  null  $path
     *
     * @return array
     * @throws FilesystemException
     */
    public function directoriesWithProperties($disk, $path = null): array
    {
        $content = Storage::disk($disk)->listContents($path ?: '')->toArray();

        return $this->filterDir($disk, $content);
    }

    /**
     * Get files with properties
     *
     * @param       $disk
     * @param  null  $path
     *
     * @return array
     * @throws FilesystemException
     */
    public function filesWithProperties($disk, $path = null): array
    {
        $content = Storage::disk($disk)->listContents($path ?: '');

        return $this->filterFile($disk, $content);
    }

    /**
     * Get directories for tree module
     *
     * @param $disk
     * @param  null  $path
     *
     * @return array
     * @throws FilesystemException
     */
    public function getDirectoriesTree($disk, $path = null): array
    {
        $directories = $this->directoriesWithProperties($disk, $path);

        foreach ($directories as $index => $dir) {
            $directories[$index]['props'] = [
                'hasSubdirectories' => (bool) Storage::disk($disk)->directories($dir['path']),
            ];
        }

        return $directories;
    }

    /**
     * File properties
     *
     * @param $disk
     * @param $path
     *
     * @return mixed
     */
    public function fileProperties($disk, $path = null): mixed
    {
        $pathInfo = pathinfo($path);

        $properties = [
            'type'       => 'file',
            'path'       => $path,
            'basename'   => $pathInfo['basename'],
            'dirname'    => $pathInfo['dirname'] === '.' ? '' : $pathInfo['dirname'],
            'extension'  => $pathInfo['extension'] ?? '',
            'filename'   => $pathInfo['filename'],
            'size'       => Storage::disk($disk)->size($path),
            'timestamp'  => Storage::disk($disk)->lastModified($path),
            'visibility' => Storage::disk($disk)->getVisibility($path),
        ];

        // if ACL ON
        if ($this->configRepository->getAcl()) {
            return $this->aclFilter($disk, [$properties])[0];
        }

        return $properties;
    }

    /**
     * Get properties for the selected directory
     *
     * @param       $disk
     * @param  null  $path
     *
     * @return array|false
     */
    public function directoryProperties($disk, $path = null): bool|array
    {
        $adapter = Storage::drive($disk)->getAdapter();

        $pathInfo = pathinfo($path);

        $properties = [
            'type'       => 'dir',
            'path'       => $path,
            'basename'   => $pathInfo['basename'],
            'dirname'    => $pathInfo['dirname'] === '.' ? '' : $pathInfo['dirname'],
            'timestamp'  => $adapter instanceof AwsS3V3Adapter ? null : Storage::disk($disk)->lastModified($path),
            'visibility' => $adapter instanceof AwsS3V3Adapter ? null : Storage::disk($disk)->getVisibility($path),
        ];

        // if ACL ON
        if ($this->configRepository->getAcl()) {
            return $this->aclFilter($disk, [$properties])[0];
        }

        return $properties;
    }

    /**
     * Get only directories
     *
     * @param $disk
     * @param $content
     *
     * @return array
     */
    protected function filterDir($disk, $content): array
    {
        // select only dir
        $dirsList = array_filter($content, fn($item) => $item['type'] === 'dir');

        $dirs = array_map(function ($item) {
            $pathInfo = pathinfo($item['path']);

            return [
                'type'       => $item['type'],
                'path'       => $item['path'],
                'basename'   => $pathInfo['basename'],
                'dirname'    => $pathInfo['dirname'] === '.' ? '' : $pathInfo['dirname'],
                'timestamp'  => $item['lastModified'],
                'visibility' => $item['visibility'],
            ];
        }, $dirsList);

        // if ACL ON
        if ($this->configRepository->getAcl()) {
            return array_values($this->aclFilter($disk, $dirs));
        }

        return array_values($dirs);
    }

    /**
     * Get only files
     *
     * @param $disk
     * @param $content
     *
     * @return array
     */
    protected function filterFile($disk, $content): array
    {
        // select only dir
        $filesList = array_filter($content, fn($item) => $item['type'] === 'file');

        $files = array_map(function ($item) {
            $pathInfo = pathinfo($item['path']);

            return [
                'type'       => $item['type'],
                'path'       => $item['path'],
                'basename'   => $pathInfo['basename'],
                'dirname'    => $pathInfo['dirname'] === '.' ? '' : $pathInfo['dirname'],
                'extension'  => $pathInfo['extension'] ?? '',
                'filename'   => $pathInfo['filename'],
                'size'       => $item['fileSize'],
                'timestamp'  => $item['lastModified'],
                'visibility' => $item['visibility'],
            ];
        }, $filesList);

        // if ACL ON
        if ($this->configRepository->getAcl()) {
            return array_values($this->aclFilter($disk, $files));
        }

        return array_values($files);
    }

    /**
     * ACL filter
     *
     * @param $disk
     * @param $content
     *
     * @return mixed
     */
    protected function aclFilter($disk, $content): mixed
    {
        $acl = resolve(ACL::class);

        $withAccess = array_map(function ($item) use ($acl, $disk) {
            // add acl access level
            $item['acl'] = $acl->getAccessLevel($disk, $item['path']);

            return $item;
        }, $content);

        // filter files and folders
        if ($this->configRepository->getAclHideFromFM()) {
            return array_filter($withAccess, function ($item) {
                return $item['acl'] !== 0;
            });
        }

        return $withAccess;
    }
}
