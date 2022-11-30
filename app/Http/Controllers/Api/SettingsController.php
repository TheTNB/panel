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
use Illuminate\Support\Facades\Hash;

class SettingsController extends Controller
{
    /**
     * 获取面板设置
     * @param  Request  $request
     * @return JsonResponse
     */
    public function get(Request $request): JsonResponse
    {
        $settings = Setting::all();
        // 隐藏字段
        $settings->makeHidden('id');
        $settings->makeHidden('created_at');
        $settings->makeHidden('updated_at');
        $settingArr = $settings->pluck('value', 'name');

        $data = [
            'name' => $settingArr['name'],
            'username' => $request->user()->username,
            'password' => '',
        ];
        $res['code'] = 0;
        $res['msg'] = 'success';
        $res['data'] = $data;
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
            if ($key == 'access_token' || $key == 'username' || $key == 'password') {
                continue;
            }
            Setting::query()->where('name', $key)->update(['value' => $value]);
        }
        // 单独处理用户名和密码
        if ($request->input('username') != $request->user()->username) {
            $request->user()->update(['username' => $request->input('username')]);
        }
        if ($request->input('password') != '') {
            $request->user()->update(['password' => Hash::make($request->input('password'))]);
        }
        $res['code'] = 0;
        $res['msg'] = 'success';
        return response()->json($res);
    }
}
