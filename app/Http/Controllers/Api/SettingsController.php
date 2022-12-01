<?php
/**
 * 耗子Linux面板 - 设置控制器
 * @author 耗子
 */

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Models\Setting;
use App\Models\User;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;
use Illuminate\Support\Facades\Hash;
use Illuminate\Support\Str;

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

        // 从nginx配置文件中获取面板端口
        $nginxConf = file_get_contents('/www/server/nginx/conf/nginx.conf');
        preg_match('/listen\s+(\d+)/', $nginxConf, $matches);

        $api = 0;
        $apiToken = '';
        if (isset($settingArr['api']) && $settingArr['api'] == 1) {
            $api = 1;
            $apiToken = $settingArr['api_token'] ?? '';
        }

        if (!isset($matches[1])) {
            $res['code'] = 1;
            $res['msg'] = '获取面板端口失败，请检查nginx主配置文件';
            return response()->json($res);
        }
        $data = [
            'name' => $settingArr['name'],
            'username' => $request->user()->username,
            'password' => '',
            'port' => $matches[1],
            'api' => $api,
            'api_token' => $apiToken,
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
            if ($key == 'access_token' || $key == 'username' || $key == 'password' || $key == 'api_token' || $key == 'api') {
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
        // 处理面板端口
        $port = $request->input('port');
        $nginxConf = file_get_contents('/www/server/nginx/conf/nginx.conf');
        preg_match('/listen\s+(\d+)/', $nginxConf, $matches);
        if (!isset($matches[1])) {
            $res['code'] = 1;
            $res['msg'] = '获取面板端口失败，请检查nginx主配置文件';
            return response()->json($res);
        }
        if ($port != $matches[1]) {
            $nginxConf = preg_replace('/listen\s+(\d+)/', 'listen '.$port, $nginxConf);
            file_put_contents('/www/server/nginx/conf/nginx.conf', $nginxConf);
            // 重载nginx
            shell_exec('systemctl reload nginx');
            // 防火墙放行端口
            shell_exec('firewall-cmd --permanent --zone=public --add-port='.$port.'/tcp >/dev/null 2>&1');
            shell_exec('firewall-cmd --reload');
        }
        // 处理api
        $api = $request->input('api', 0);
        $apiCheck = Setting::query()->where('name', 'api')->value('value');
        if (empty($apiCheck)) {
            $apiCheck = 0;
        }
        if ($api != $apiCheck) {
            if ($api) {
                Setting::query()->insert([
                    'name' => 'api',
                    'value' => 1,
                ]);
                // 生成api用户
                $username = 'api_'.Str::random();
                $apiUser = User::query()->create([
                    'username' => $username,
                    'password' => Hash::make(Str::random()),
                    'email' => 'panel_api@haozi.net',
                ]);
                // 生成api token
                $apiToken = $apiUser->createToken('api')->plainTextToken;
                Setting::query()->insert([
                    'name' => 'api_user',
                    'value' => $username,
                ]);
                Setting::query()->insert([
                    'name' => 'api_token',
                    'value' => $apiToken,
                ]);
            } else {
                Setting::query()->where('name', 'api')->delete();
                Setting::query()->where('name', 'api_token')->delete();
                $username = Setting::query()->where('name', 'api_user')->value('value');
                Setting::query()->where('name', 'api_user')->delete();
                Setting::query()->where('name', 'api')->delete();
                $apiUser = User::query()->where('username', $username)->first();
                // 删除api用户的所有token
                $apiUser->tokens()->delete();
                // 删除api用户
                $apiUser->delete();
            }
        }
        $res['code'] = 0;
        $res['msg'] = 'success';
        return response()->json($res);
    }
}
