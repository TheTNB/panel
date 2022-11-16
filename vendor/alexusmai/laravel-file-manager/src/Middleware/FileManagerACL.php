<?php

namespace Alexusmai\LaravelFileManager\Middleware;

use Alexusmai\LaravelFileManager\Services\ACLService\ACL;
use Alexusmai\LaravelFileManager\Services\ConfigService\ConfigRepository;
use Alexusmai\LaravelFileManager\Traits\PathTrait;
use Illuminate\Http\Request;
use Illuminate\Support\Arr;
use Closure;

class FileManagerACL
{
    use PathTrait;

    /**
     * Check method names
     */
    const CHECKERS = [
        'fm.tree'             => 'checkContent',
        'fm.content'          => 'checkContent',
        'fm.preview'          => 'checkContent',
        'fm.thumbnails'       => 'checkContent',
        'fm.url'              => 'checkContent',
        'fm.stream-file'      => 'checkContent',
        'fm.download'         => 'checkDownload',
        'fm.create-file'      => 'checkCreate',
        'fm.create-directory' => 'checkCreate',
        'fm.update-file'      => 'checkUpdate',
        'fm.upload'           => 'checkUpload',
        'fm.delete'           => 'checkDelete',
        'fm.paste'            => 'checkPaste',
        'fm.rename'           => 'checkRename',
        'fm.zip'              => 'checkZip',
        'fm.unzip'            => 'checkUnzip',
    ];

    /**
     * @var string|null
     */
    protected $disk;

    /**
     * @var string|null
     */
    protected $path;

    /**
     * @var ACL|mixed
     */
    protected $acl;

    /**
     * @var Request
     */
    protected $request;

    /**
     * FileManagerACL constructor.
     *
     * @param  Request  $request
     * @param  ACL  $acl
     */
    public function __construct(Request $request, ACL $acl)
    {
        $this->disk = $request->has('disk') ? $request->input('disk') : null;
        $this->path = $request->has('path') ? $request->input('path') : '/';

        $this->acl = $acl;

        $this->request = $request;
    }

    /**
     * Handle an incoming request.
     *
     * @param  \Illuminate\Http\Request  $request
     * @param  \Closure  $next
     *
     * @return mixed
     */
    public function handle($request, Closure $next)
    {
        $routeName = $request->route()->getName();

        // if ACL is OFF or route name wasn't found
        if ( ! resolve(ConfigRepository::class)->getAcl()
            || ! array_key_exists($routeName, self::CHECKERS)
        ) {
            return $next($request);
        }

        if ( ! call_user_func([$this, self::CHECKERS[$routeName]])) {
            return $this->errorMessage();
        }

        // return request
        return $next($request);
    }

    /**
     * ACL Error message
     *
     * @return \Illuminate\Http\JsonResponse
     */
    protected function errorMessage()
    {
        return response()->json([
            'result' => [
                'status'  => 'error',
                'message' => 'aclError',
            ],
        ]);
    }

    /**
     * Check content actions
     *
     * @return bool
     */
    protected function checkContent()
    {
        // need r access
        return ! ($this->acl->getAccessLevel($this->disk, $this->path) === 0);
    }

    /**
     * Check download actions
     */
    protected function checkDownload()
    {
        // need r access
        abort_if(
            $this->acl->getAccessLevel($this->disk, $this->path) === 0,
            403
        );

        return true;
    }

    /**
     * Check create actions
     *
     * @return bool
     */
    protected function checkCreate()
    {
        $name = $this->request->input('name');
        $pathToWrite = $this->request->input('path')
            ? $this->request->input('path').'/' : '';

        // need r/w access
        return ! ($this->acl->getAccessLevel($this->disk, $pathToWrite.$name) !== 2);
    }

    /**
     * Check update actions
     *
     * @return bool
     */
    protected function checkUpdate()
    {
        $pathToWrite = $this->request->input('path')
            ? $this->request->input('path').'/' : '';

        $name = $this->request->file('file')->getClientOriginalName();

        // need r/w access
        return ! ($this->acl->getAccessLevel($this->disk, $pathToWrite.$name) !== 2);
    }

    /**
     * Check upload actions
     *
     * @return bool
     */
    protected function checkUpload()
    {
        $pathToWrite = $this->request->input('path')
            ? $this->request->input('path').'/' : '';

        // filter
        $firstFall = Arr::first($this->request->file('files'),
            function ($value) use ($pathToWrite) {
                // need r/w access
                return $this->acl->getAccessLevel(
                        $this->disk,
                        $pathToWrite.$value->getClientOriginalName()
                    ) !== 2;
            }, null);

        // if founded one access error
        if ($firstFall) {
            return false;
        }

        return true;
    }

    /**
     * Check delete actions
     *
     * @return bool
     */
    protected function checkDelete()
    {
        $firstFall = Arr::first($this->request->input('items'),
            function ($value) {
                // need r/w access
                return $this->acl->getAccessLevel($this->disk, $value['path']) !== 2;
            }, null);

        if ($firstFall) {
            return false;
        }

        return true;
    }

    /**
     * Check paste action
     *
     * @return bool
     */
    protected function checkPaste()
    {
        // get clipboard data
        $clipboard = $this->request->input('clipboard');

        // copy - r, cut - rw
        $getLevel = $clipboard['type'] === 'copy' ? 1 : 2;

        // can user copy or cut selected files and folders
        $checkDirs = Arr::first($clipboard['directories'],
            function ($value) use ($clipboard, $getLevel) {
                return $this->acl->getAccessLevel($clipboard['disk'], $value) < $getLevel;
            }, null);

        $checkFiles = Arr::first($clipboard['files'],
            function ($value) use ($clipboard, $getLevel) {
                return $this->acl->getAccessLevel($clipboard['disk'], $value) < $getLevel;
            }, null);

        // can user write to selected folder?
        $writeToFolder = $this->acl->getAccessLevel($this->disk, $this->path);

        return ! ($checkDirs || $checkFiles || $writeToFolder !== 2);
    }

    /**
     * Check rename actions
     *
     * @return bool
     */
    protected function checkRename()
    {
        // old path
        $oldPath = $this->request->input('oldName');

        // new path
        $newPath = $this->request->input('newName');

        // need r/w access
        return ! ($this->acl->getAccessLevel($this->disk, $oldPath) !== 2
            || $this->acl->getAccessLevel($this->disk, $newPath) !== 2);
    }

    /**
     * Check zip actions
     *
     * @return bool
     */
    protected function checkZip()
    {
        // can user write to selected folder?
        $writeToFolder = $this->acl->getAccessLevel(
            $this->disk,
            $this->newPath(
                $this->request->input('path'),
                $this->request->input('name')
            )
        );

        // need r/w access
        if ($writeToFolder !== 2) {
            return false;
        }

        // data to zip
        $elements = $this->request->input('elements');

        // can user read selected files and folders?
        $checkDirs = Arr::first($elements['directories'],
            function ($value) {
                // need r access
                return $this->acl->getAccessLevel($this->disk, $value) === 0;
            }, null);


        $checkFiles = Arr::first($elements['files'],
            function ($value) {
                // need r access
                return $this->acl->getAccessLevel($this->disk, $value) === 0;
            }, null);

        return ! ($checkDirs || $checkFiles);
    }

    /**
     * Check unzip actions
     *
     * @return bool
     */
    protected function checkUnzip()
    {
        if ($this->request->input('folder')) {
            $dirname = dirname($this->path) === '.' ? ''
                : dirname($this->path).'/';
            $pathToWrite = $dirname.$this->request->input('folder');
        } else {
            $pathToWrite = dirname($this->path) === '.' ? '/'
                : dirname($this->path);
        }

        return ! ($this->acl->getAccessLevel($this->disk, $pathToWrite) !== 2
            || $this->acl->getAccessLevel($this->disk, $this->path) === 0);
    }
}
