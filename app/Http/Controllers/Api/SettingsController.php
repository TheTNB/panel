<?php
/**
 * 耗子Linux面板 - 设置控制器
 * @author 耗子
 */

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Models\Setting;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class SettingsController extends Controller
{
    /**
     * 获取面板设置
     * @return JsonResponse
     */
    public function get(): JsonResponse
    {
        $settings = Setting::all();
        // 隐藏字段
        $settings->makeHidden('id');
        $settings->makeHidden('created_at');
        $settings->makeHidden('updated_at');
        $res['code'] = 0;
        $res['msg'] = 'success';
        $res['data'] = $settings->pluck('value', 'name');
        return response()->json($res);
    }

    /**
     * 保存面板设置
     * @param  Request  $request
     * @return JsonResponse
     */
    public function save(Request $request): JsonResponse
    {
        // 获取前端传递过来的数据
        $settings = $request->all();
        // 将数据入库
        foreach ($settings as $key => $value) {
            if ($key == 'access_token') {
                continue;
            }
            Setting::query()->where('name', $key)->update(['value' => $value]);
        }
        $res['code'] = 0;
        $res['msg'] = 'success';
        return response()->json($res);
    }
}
