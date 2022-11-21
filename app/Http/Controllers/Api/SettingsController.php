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
     * @param  Request  $request
     * @return JsonResponse
     */
    public function get(Request $request)
    {
        $settings = Setting::query()->get()->toArray();
        foreach ($settings as $setting) {
            $res['data'][$setting['name']] = $setting['value'];
        }

        if (!empty($settings)) {
            $res['code'] = 0;
            $res['msg'] = 'success';
        } else {
            $res['code'] = 1;
            $res['msg'] = '面板设置获取失败';
            $res['data'] = null;
        }
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
            if ($key == 'mysql_root_password') {
                $old_mysql_password = Setting::query()->where('name', 'mysql_root_password')->value('value');
                if ($old_mysql_password != $value) {
                    shell_exec('mysql -uroot -p'.$old_mysql_password.' -e "ALTER USER \'root\'@\'localhost\' IDENTIFIED BY \''.$value.'\';"');
                    shell_exec('mysql -uroot -p'.$old_mysql_password.' -e "flush privileges;"');
                }
            }
            Setting::query()->where('name', $key)->update(['value' => $value]);
        }
        $res['code'] = 0;
        $res['msg'] = 'success';
        return response()->json($res);
    }
}
