<?php

namespace Alexusmai\LaravelFileManager\Controllers;

use Alexusmai\LaravelFileManager\Events\BeforeInitialization;
use Alexusmai\LaravelFileManager\Events\Deleting;
use Alexusmai\LaravelFileManager\Events\DirectoryCreated;
use Alexusmai\LaravelFileManager\Events\DirectoryCreating;
use Alexusmai\LaravelFileManager\Events\DiskSelected;
use Alexusmai\LaravelFileManager\Events\Download;
use Alexusmai\LaravelFileManager\Events\FileCreated;
use Alexusmai\LaravelFileManager\Events\FileCreating;
use Alexusmai\LaravelFileManager\Events\FilesUploaded;
use Alexusmai\LaravelFileManager\Events\FilesUploading;
use Alexusmai\LaravelFileManager\Events\FileUpdate;
use Alexusmai\LaravelFileManager\Events\Paste;
use Alexusmai\LaravelFileManager\Events\Rename;
use Alexusmai\LaravelFileManager\Events\Zip as ZipEvent;
use Alexusmai\LaravelFileManager\Events\Unzip as UnzipEvent;
use Alexusmai\LaravelFileManager\Requests\RequestValidator;
use Alexusmai\LaravelFileManager\FileManager;
use Alexusmai\LaravelFileManager\Services\Zip;
use Illuminate\Contracts\Container\BindingResolutionException;
use Illuminate\Contracts\View\Factory;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Response;
use Illuminate\Routing\Controller;
use Illuminate\View\View;
use League\Flysystem\FilesystemException;
use Symfony\Component\HttpFoundation\StreamedResponse;

class FileManagerController extends Controller
{
    /**
     * @var FileManager
     */
    public $fm;

    /**
     * FileManagerController constructor.
     *
     * @param  FileManager  $fm
     */
    public function __construct(FileManager $fm)
    {
        $this->fm = $fm;
    }

    /**
     * Initialize file manager
     *
     * @return JsonResponse
     */
    public function initialize(): JsonResponse
    {
        event(new BeforeInitialization());

        return response()->json(
            $this->fm->initialize()
        );
    }

    /**
     * Get files and directories for the selected path and disk
     *
     * @param  RequestValidator  $request
     *
     * @return JsonResponse
     * @throws FilesystemException
     */
    public function content(RequestValidator $request): JsonResponse
    {
        return response()->json(
            $this->fm->content(
                $request->input('disk'),
                $request->input('path')
            )
        );
    }

    /**
     * Directory tree
     *
     * @param  RequestValidator  $request
     *
     * @return JsonResponse
     * @throws FilesystemException
     */
    public function tree(RequestValidator $request): JsonResponse
    {
        return response()->json(
            $this->fm->tree(
                $request->input('disk'),
                $request->input('path')
            )
        );
    }

    /**
     * Check the selected disk
     *
     * @param  RequestValidator  $request
     *
     * @return JsonResponse
     */
    public function selectDisk(RequestValidator $request): JsonResponse
    {
        event(new DiskSelected($request->input('disk')));

        return response()->json([
            'result' => [
                'status'  => 'success',
                'message' => 'diskSelected',
            ],
        ]);
    }

    /**
     * Upload files
     *
     * @param  RequestValidator  $request
     *
     * @return JsonResponse
     */
    public function upload(RequestValidator $request): JsonResponse
    {
        event(new FilesUploading($request));

        $uploadResponse = $this->fm->upload(
            $request->input('disk'),
            $request->input('path'),
            $request->file('files'),
            $request->input('overwrite')
        );

        event(new FilesUploaded($request));

        return response()->json($uploadResponse);
    }

    /**
     * Delete files and folders
     *
     * @param  RequestValidator  $request
     *
     * @return JsonResponse
     */
    public function delete(RequestValidator $request): JsonResponse
    {
        event(new Deleting($request));

        $deleteResponse = $this->fm->delete(
            $request->input('disk'),
            $request->input('items')
        );

        return response()->json($deleteResponse);
    }

    /**
     * Copy / Cut files and folders
     *
     * @param  RequestValidator  $request
     *
     * @return JsonResponse
     */
    public function paste(RequestValidator $request): JsonResponse
    {
        event(new Paste($request));

        return response()->json(
            $this->fm->paste(
                $request->input('disk'),
                $request->input('path'),
                $request->input('clipboard')
            )
        );
    }

    /**
     * Rename
     *
     * @param  RequestValidator  $request
     *
     * @return JsonResponse
     */
    public function rename(RequestValidator $request): JsonResponse
    {
        event(new Rename($request));

        return response()->json(
            $this->fm->rename(
                $request->input('disk'),
                $request->input('newName'),
                $request->input('oldName')
            )
        );
    }

    /**
     * Download file
     *
     * @param  RequestValidator  $request
     *
     * @return StreamedResponse
     */
    public function download(RequestValidator $request): StreamedResponse
    {
        event(new Download($request));

        return $this->fm->download(
            $request->input('disk'),
            $request->input('path')
        );
    }

    /**
     * Create thumbnails
     *
     * @param  RequestValidator  $request
     *
     * @return Response|mixed
     * @throws BindingResolutionException
     */
    public function thumbnails(RequestValidator $request): mixed
    {
        return $this->fm->thumbnails(
            $request->input('disk'),
            $request->input('path')
        );
    }

    /**
     * Image preview
     *
     * @param  RequestValidator  $request
     *
     * @return mixed
     */
    public function preview(RequestValidator $request): mixed
    {
        return $this->fm->preview(
            $request->input('disk'),
            $request->input('path')
        );
    }

    /**
     * File url
     *
     * @param  RequestValidator  $request
     *
     * @return JsonResponse
     */
    public function url(RequestValidator $request): JsonResponse
    {
        return response()->json(
            $this->fm->url(
                $request->input('disk'),
                $request->input('path')
            )
        );
    }

    /**
     * Create new directory
     *
     * @param  RequestValidator  $request
     *
     * @return JsonResponse
     */
    public function createDirectory(RequestValidator $request): JsonResponse
    {
        event(new DirectoryCreating($request));

        $createDirectoryResponse = $this->fm->createDirectory(
            $request->input('disk'),
            $request->input('path'),
            $request->input('name')
        );

        if ($createDirectoryResponse['result']['status'] === 'success') {
            event(new DirectoryCreated($request));
        }

        return response()->json($createDirectoryResponse);
    }

    /**
     * Create new file
     *
     * @param  RequestValidator  $request
     *
     * @return JsonResponse
     */
    public function createFile(RequestValidator $request): JsonResponse
    {
        event(new FileCreating($request));

        $createFileResponse = $this->fm->createFile(
            $request->input('disk'),
            $request->input('path'),
            $request->input('name')
        );

        if ($createFileResponse['result']['status'] === 'success') {
            event(new FileCreated($request));
        }

        return response()->json($createFileResponse);
    }

    /**
     * Update file
     *
     * @param  RequestValidator  $request
     *
     * @return JsonResponse
     */
    public function updateFile(RequestValidator $request): JsonResponse
    {
        event(new FileUpdate($request));

        return response()->json(
            $this->fm->updateFile(
                $request->input('disk'),
                $request->input('path'),
                $request->file('file')
            )
        );
    }

    /**
     * Stream file
     *
     * @param  RequestValidator  $request
     *
     * @return mixed
     */
    public function streamFile(RequestValidator $request): mixed
    {
        return $this->fm->streamFile(
            $request->input('disk'),
            $request->input('path')
        );
    }

    /**
     * Create zip archive
     *
     * @param  RequestValidator  $request
     * @param  Zip  $zip
     *
     * @return array
     */
    public function zip(RequestValidator $request, Zip $zip)
    {
        event(new ZipEvent($request));

        return $zip->create();
    }

    /**
     * Extract zip archive
     *
     * @param  RequestValidator  $request
     * @param  Zip  $zip
     *
     * @return array
     */
    public function unzip(RequestValidator $request, Zip $zip)
    {
        event(new UnzipEvent($request));

        return $zip->extract();
    }

    /**
     * Integration with ckeditor 4
     *
     * @return Factory|View
     */
    public function ckeditor(): Factory|View
    {
        return view('file-manager::ckeditor');
    }

    /**
     * Integration with TinyMCE v4
     *
     * @return Factory|View
     */
    public function tinymce(): Factory|View
    {
        return view('file-manager::tinymce');
    }

    /**
     * Integration with TinyMCE v5
     *
     * @return Factory|View
     */
    public function tinymce5(): Factory|View
    {
        return view('file-manager::tinymce5');
    }

    /**
     * Integration with SummerNote
     *
     * @return Factory|View
     */
    public function summernote(): Factory|View
    {
        return view('file-manager::summernote');
    }

    /**
     * Simple integration with input field
     *
     * @return Factory|View
     */
    public function fmButton(): Factory|View
    {
        return view('file-manager::fmButton');
    }
}
