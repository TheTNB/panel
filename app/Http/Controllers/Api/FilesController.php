<?php
/**
 * 耗子Linux面板 - 文件控制器
 * @author 耗子
 */

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Support\Carbon;
use Illuminate\Validation\ValidationException;
use Symfony\Component\HttpFoundation\BinaryFileResponse;

class FilesController extends Controller
{
    /**
     * 获取某个目录下的文件（夹）列表
     * @param  Request  $request
     * @return JsonResponse
     */
    public function getList(Request $request): JsonResponse
    {
        try {
            $credentials = $this->validate($request, [
                'path' => ['required', 'regex:/^\/.*$/'],
                'limit' => 'required|integer',
                'page' => 'required|integer',
            ]);
            $path = $credentials['path'];
            $limit = $credentials['limit'];
            $page = $credentials['page'];
        } catch (ValidationException $e) {
            return response()->json(['code' => 1, 'msg' => $e->getMessage()]);
        }

        if (!is_dir($path)) {
            return response()->json(['code' => 1, 'msg' => '目录不存在']);
        }

        $files = scandir($path);
        $fileData = [];
        foreach ($files as $k => $v) {
            if ($v == '.' || $v == '..') {
                continue;
            }
            $fileData[$k]['name'] = $v;
            $fileData[$k]['path'] = $path.'/'.$v;
            $fileData[$k]['type'] = filetype($path.'/'.$v);
            $fileData[$k]['size'] = filesize($path.'/'.$v);
            $fileData[$k]['mtime'] = Carbon::createFromTimestamp(filemtime($path.'/'.$v))->toDateTimeString();
        }

        // 分页

        $total = count($fileData);
        $filesArr = array_slice($fileData, ($page - 1) * $limit, $limit);
        $data['code'] = 0;
        $data['msg'] = 'success';
        $data['count'] = $total;
        $data['data'] = $filesArr;
        return response()->json($data);
    }

    /**
     * 获取文件内容
     * @param  Request  $request
     * @return JsonResponse
     */
    public function getFileContent(Request $request): JsonResponse
    {
        try {
            $credentials = $this->validate($request, [
                'file' => ['required', 'regex:/^\/.*$/'],
            ]);
            $file = $credentials['file'];
        } catch (ValidationException $e) {
            return response()->json(['code' => 1, 'msg' => $e->getMessage()]);
        }

        if (!is_file($file)) {
            return response()->json(['code' => 1, 'msg' => '文件不存在']);
        }

        $content = @file_get_contents($file);
        $data['code'] = 0;
        $data['msg'] = 'success';
        $data['data'] = $content;
        return response()->json($data);
    }

    /**
     * 保存文件内容
     * @param  Request  $request
     * @return JsonResponse
     */
    public function saveFileContent(Request $request): JsonResponse
    {
        try {
            $credentials = $this->validate($request, [
                'file' => ['required', 'regex:/^\/.*$/'],
                'content' => 'required',
            ]);
            $file = $credentials['file'];
            $content = $credentials['content'];
        } catch (ValidationException $e) {
            return response()->json(['code' => 1, 'msg' => $e->getMessage()]);
        }

        if (!is_file($file)) {
            return response()->json(['code' => 1, 'msg' => '文件不存在']);
        }

        $res = @file_put_contents($file, $content);
        if ($res === false) {
            return response()->json(['code' => 1, 'msg' => '保存失败']);
        }

        return response()->json(['code' => 0, 'msg' => 'success']);
    }

    /**
     * 创建目录
     * @param  Request  $request
     * @return JsonResponse
     */
    public function createDir(Request $request): JsonResponse
    {
        try {
            $credentials = $this->validate($request, [
                'path' => ['required', 'regex:/^\/.*$/'],
                'name' => 'required',
            ]);
            $path = $credentials['path'];
            $name = $credentials['name'];
        } catch (ValidationException $e) {
            return response()->json(['code' => 1, 'msg' => $e->getMessage()]);
        }

        if (!is_dir($path)) {
            return response()->json(['code' => 1, 'msg' => '目录不存在']);
        }

        $res = @mkdir($path.'/'.$name);
        if ($res === false) {
            return response()->json(['code' => 1, 'msg' => '创建失败']);
        }

        return response()->json(['code' => 0, 'msg' => 'success']);
    }

    /**
     * 重命名文件或目录
     * @param  Request  $request
     * @return JsonResponse
     */
    public function rename(Request $request): JsonResponse
    {
        try {
            $credentials = $this->validate($request, [
                'path' => ['required', 'regex:/^\/.*$/'],
                'name' => 'required',
            ]);
            $path = $credentials['path'];
            $name = $credentials['name'];
        } catch (ValidationException $e) {
            return response()->json(['code' => 1, 'msg' => $e->getMessage()]);
        }

        if (!is_file($path) && !is_dir($path)) {
            return response()->json(['code' => 1, 'msg' => '文件或目录不存在']);
        }

        $res = @rename($path, dirname($path).'/'.$name);
        if ($res === false) {
            return response()->json(['code' => 1, 'msg' => '重命名失败']);
        }

        return response()->json(['code' => 0, 'msg' => 'success']);
    }

    /**
     * 上传文件
     * @param  Request  $request
     * @return JsonResponse
     */
    public function uploadFile(Request $request): JsonResponse
    {
        try {
            $credentials = $this->validate($request, [
                'path' => ['required', 'regex:/^\/.*$/'],
                'file' => 'required|file',
            ]);
            $path = $credentials['path'];
            $file = $credentials['file'];
        } catch (ValidationException $e) {
            return response()->json(['code' => 1, 'msg' => $e->getMessage()]);
        }

        if (!is_dir($path)) {
            return response()->json(['code' => 1, 'msg' => '目录不存在']);
        }

        $res = @move_uploaded_file($file->getRealPath(), $path.'/'.$file->getClientOriginalName());
        if ($res === false) {
            return response()->json(['code' => 1, 'msg' => '上传失败']);
        }

        return response()->json(['code' => 0, 'msg' => 'success']);
    }

    /**
     * 下载文件
     * @param  Request  $request
     * @return JsonResponse|BinaryFileResponse
     */
    public function downloadFile(Request $request): BinaryFileResponse|JsonResponse
    {
        try {
            $credentials = $this->validate($request, [
                'path' => ['required', 'regex:/^\/.*$/'],
            ]);
            $path = $credentials['path'];
        } catch (ValidationException $e) {
            return response()->json(['code' => 1, 'msg' => $e->getMessage()]);
        }

        if (!is_file($path)) {
            return response()->json(['code' => 1, 'msg' => '文件不存在']);
        }

        return response()->download($path);
    }

    /**
     * 创建文件
     */
    public function createFile(Request $request): JsonResponse
    {
        try {
            $credentials = $this->validate($request, [
                'path' => ['required', 'regex:/^\/.*$/'],
                'name' => 'required',
            ]);
            $path = $credentials['path'];
            $name = $credentials['name'];
        } catch (ValidationException $e) {
            return response()->json(['code' => 1, 'msg' => $e->getMessage()]);
        }

        if (!is_dir($path)) {
            return response()->json(['code' => 1, 'msg' => '目录不存在']);
        }

        $res = @file_put_contents($path.'/'.$name, '');
        if ($res === false) {
            return response()->json(['code' => 1, 'msg' => '创建失败']);
        }

        return response()->json(['code' => 0, 'msg' => 'success']);
    }

    /**
     * 设置权限
     * @param  Request  $request
     * @return JsonResponse
     */
    public function chmod(Request $request): JsonResponse
    {
        try {
            $credentials = $this->validate($request, [
                'path' => ['required', 'regex:/^\/.*$/'],
                'mode' => 'required',
            ]);
            $path = $credentials['path'];
            $mode = $credentials['mode'];
        } catch (ValidationException $e) {
            return response()->json(['code' => 1, 'msg' => $e->getMessage()]);
        }

        if (!is_file($path) && !is_dir($path)) {
            return response()->json(['code' => 1, 'msg' => '文件或目录不存在']);
        }

        $res = @chmod($path, $mode);
        if ($res === false) {
            return response()->json(['code' => 1, 'msg' => '设置失败']);
        }

        return response()->json(['code' => 0, 'msg' => 'success']);
    }

    /**
     * 设置所有者
     */
    public function chown(Request $request): JsonResponse
    {
        try {
            $credentials = $this->validate($request, [
                'path' => ['required', 'regex:/^\/.*$/'],
                'user' => 'required',
            ]);
            $path = $credentials['path'];
            $user = $credentials['user'];
        } catch (ValidationException $e) {
            return response()->json(['code' => 1, 'msg' => $e->getMessage()]);
        }

        if (!is_file($path) && !is_dir($path)) {
            return response()->json(['code' => 1, 'msg' => '文件或目录不存在']);
        }

        $res = @chown($path, $user);
        if ($res === false) {
            return response()->json(['code' => 1, 'msg' => '设置失败']);
        }

        return response()->json(['code' => 0, 'msg' => 'success']);
    }

    /**
     * 设置所有者组
     * @param  Request  $request
     * @return JsonResponse
     */
    public function chgrp(Request $request): JsonResponse
    {
        try {
            $credentials = $this->validate($request, [
                'path' => ['required', 'regex:/^\/.*$/'],
                'group' => 'required',
            ]);
            $path = $credentials['path'];
            $group = $credentials['group'];
        } catch (ValidationException $e) {
            return response()->json(['code' => 1, 'msg' => $e->getMessage()]);
        }

        if (!is_file($path) && !is_dir($path)) {
            return response()->json(['code' => 1, 'msg' => '文件或目录不存在']);
        }

        $res = @chgrp($path, $group);
        if ($res === false) {
            return response()->json(['code' => 1, 'msg' => '设置失败']);
        }

        return response()->json(['code' => 0, 'msg' => 'success']);
    }

    /**
     * 复制文件/目录
     * @param  Request  $request
     * @return JsonResponse
     */
    public function copy(Request $request): JsonResponse
    {
        try {
            $credentials = $this->validate($request, [
                'path' => ['required', 'regex:/^\/.*$/'],
                'dest' => ['required', 'regex:/^\/.*$/'],
            ]);
            $path = $credentials['path'];
            $dest = $credentials['dest'];
        } catch (ValidationException $e) {
            return response()->json(['code' => 1, 'msg' => $e->getMessage()]);
        }

        if (!is_file($path) && !is_dir($path)) {
            return response()->json(['code' => 1, 'msg' => '文件或目录不存在']);
        }

        if (!is_dir($dest)) {
            return response()->json(['code' => 1, 'msg' => '目标目录不存在']);
        }

        $res = @copy($path, $dest.'/'.basename($path));
        if ($res === false) {
            return response()->json(['code' => 1, 'msg' => '复制失败']);
        }

        return response()->json(['code' => 0, 'msg' => 'success']);
    }

    /**
     * 移动文件/目录
     * @param  Request  $request
     * @return JsonResponse
     */
    public function move(Request $request): JsonResponse
    {
        try {
            $credentials = $this->validate($request, [
                'path' => ['required', 'regex:/^\/.*$/'],
                'dest' => ['required', 'regex:/^\/.*$/'],
            ]);
            $path = $credentials['path'];
            $dest = $credentials['dest'];
        } catch (ValidationException $e) {
            return response()->json(['code' => 1, 'msg' => $e->getMessage()]);
        }

        if (!is_file($path) && !is_dir($path)) {
            return response()->json(['code' => 1, 'msg' => '文件或目录不存在']);
        }

        if (!is_dir($dest)) {
            return response()->json(['code' => 1, 'msg' => '目标目录不存在']);
        }

        $res = @rename($path, $dest.'/'.basename($path));
        if ($res === false) {
            return response()->json(['code' => 1, 'msg' => '移动失败']);
        }

        return response()->json(['code' => 0, 'msg' => 'success']);
    }

    /**
     * 删除文件/目录
     * @param  Request  $request
     * @return JsonResponse
     */
    public function delete(Request $request): JsonResponse
    {
        try {
            $credentials = $this->validate($request, [
                'path' => ['required', 'regex:/^\/.*$/'],
            ]);
            $path = $credentials['path'];
        } catch (ValidationException $e) {
            return response()->json(['code' => 1, 'msg' => $e->getMessage()]);
        }

        if (!is_file($path) && !is_dir($path)) {
            return response()->json(['code' => 1, 'msg' => '文件或目录不存在']);
        }

        $res = @shell_exec('rm -rf '.escapeshellarg($path).' 2>&1');
        if (!empty($res)) {
            return response()->json(['code' => 1, 'msg' => '删除失败：'.$res]);
        }

        return response()->json(['code' => 0, 'msg' => 'success']);
    }

    /**
     * 压缩文件/目录
     * @param  Request  $request
     * @return JsonResponse
     */
    public function zip(Request $request): JsonResponse
    {
        try {
            $credentials = $this->validate($request, [
                'path' => ['required', 'regex:/^\/.*$/'],
            ]);
            $path = $credentials['path'];
        } catch (ValidationException $e) {
            return response()->json(['code' => 1, 'msg' => $e->getMessage()]);
        }

        if (!is_file($path) && !is_dir($path)) {
            return response()->json(['code' => 1, 'msg' => '文件或目录不存在']);
        }

        $zipPath = dirname($path).'/'.basename($path).'.zip';
        $res = @shell_exec('zip -r '.escapeshellarg($zipPath).' '.escapeshellarg($path).' 2>&1');
        if (!empty($res)) {
            return response()->json(['code' => 1, 'msg' => '压缩失败：'.$res]);
        }

        return response()->json(['code' => 0, 'msg' => 'success']);
    }

    /**
     * 解压文件/目录
     * @param  Request  $request
     * @return JsonResponse
     */
    public function unzip(Request $request): JsonResponse
    {
        try {
            $credentials = $this->validate($request, [
                'path' => ['required', 'regex:/^\/.*$/'],
            ]);
            $path = $credentials['path'];
        } catch (ValidationException $e) {
            return response()->json(['code' => 1, 'msg' => $e->getMessage()]);
        }

        if (!is_file($path)) {
            return response()->json(['code' => 1, 'msg' => '文件不存在']);
        }

        $res = @shell_exec('unzip -o '.escapeshellarg($path).' -d '.escapeshellarg(dirname($path)).' 2>&1');
        if (!empty($res)) {
            return response()->json(['code' => 1, 'msg' => '解压失败：'.$res]);
        }

        return response()->json(['code' => 0, 'msg' => 'success']);
    }
}
