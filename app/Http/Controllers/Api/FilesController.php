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

class FilesController extends Controller
{
    /**
     * 获取某个目录的文件列表
     */
    public function getDirList(Request $request): JsonResponse
    {
        $limit = $request->input('limit', 10);


        $data['code'] = 0;
        $data['msg'] = 'success';
        $data['count'] = '';
        $data['data'] = '';
        return response()->json($data);
    }
}
