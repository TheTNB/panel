<?php
/**
 * 耗子Linux面板 - 设置控制器
 * @author 耗子
 */
namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use Illuminate\Http\Request;

class SettingsController extends Controller
{
    /**
     * 获取面板设置
     * @return
     */
    public function get_settings(Request $request)
    {
        $settings = Db::table('setting')->select()->toArray();
        foreach ($settings as $setting) {
            $res['data'][$setting['name']] = $setting['value'];
        }
        $user_password = Db::table('user')->where('username', $request->username)->value('password');
        $res['data']['username'] = $request->username;
        $res['data']['password'] = $user_password;

        if (!empty($settings)) {
            $res['code'] = 0;
            $res['msg'] = 'success';
            return response()->json($res);
        } else {
            $res['code'] = 1;
            $res['msg'] = '面板设置获取失败';
            $res['data'] = null;
            return response()->json($res);
        }
    }

    /**
     * 保存面板设置
     * @return
     */
    public function save_settings(Request $request)
    {
        // 获取前端传递过来的数据
        $settings = Request::post();
        // 将数据入库
        foreach ($settings as $key => $value) {
            if ($key == 'access_token' || $key == 'username' || $key == 'password') {
                continue;
            }
            if ($key == 'mysql_root_password') {
                $old_mysql_password = Db::table('setting')->where('name', 'mysql_root_password')->value('value');
                if ($old_mysql_password != $value) {
                    shell_exec('/www/server/mysql/bin/mysqladmin -uroot -p' . $old_mysql_password . ' password ' . $value);
                }
            }
            Db::table('setting')->where('name', $key)->update(['value' => $value]);
        }
        $res['code'] = 0;
        $res['msg'] = 'success';
        $old_user_info = Db::table('user')->where('username', $request->username)->select()->toArray();

        if ($old_user_info[0]['username'] != $settings['username'] || $old_user_info[0]['password'] != $settings['password']) {
            $res['msg'] = 'change';
            Db::table('user')->where('username', $request->username)->update(['username' => $settings['username']]);
            Db::table('user')->where('username', $settings['username'])->update(['password' => $settings['password']]);
        }
        return response()->json($res);
    }
}
