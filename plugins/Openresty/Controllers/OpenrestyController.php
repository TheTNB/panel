<?php
/**
 * Name: OpenResty插件控制器
 * Author:耗子
 * Date: 2022-11-02
 */

namespace Plugins\Openresty\Controllers;


use App\Http\Controllers\Controller;

// HTTP
use Illuminate\Support\Facades\Http;
use Illuminate\Support\Facades\Request;
// Filesystem
use Illuminate\Filesystem\Filesystem;

class OpenrestyController extends Controller
{

    public function status()
    {
        $command = 'systemctl status nginx';
        $result = shell_exec($command);

        $res['code'] = 0;
        $res['msg'] = 'success';
        if (str_contains($result, 'inactive')) {
            $res['data'] = 'stopped';
        } else {
            $res['data'] = 'running';
        }

        return response()->json($res);
    }

    public function restart()
    {
        $command = 'nginx -t 2>&1';
        $result = shell_exec($command);

        $res['code'] = 0;
        $res['msg'] = 'success';
        if (str_contains($result, 'test failed')) {
            // 测试失败，则不允许重启
            $res['msg'] = 'error';
            $res['data'] = 'OpenResty配置有误，请修正后再重启：'.$result;
            return response()->json($res);
        }

        $command2 = 'systemctl restart nginx';
        $result2 = shell_exec($command2);
        if (str_contains($result2, 'done')) {
            $res['data'] = 'OpenResty已重启';
            return response()->json($res);
        }
        return response()->json($res);
    }

    public function reload()
    {
        $command = 'nginx -t 2>&1';
        $result = shell_exec($command);
        $res['code'] = 0;
        $res['msg'] = 'success';
        if (str_contains($result, 'test failed')) {
            // 测试失败，则不允许重载
            $res['msg'] = 'error';
            $res['data'] = 'OpenResty配置有误，请修正后再重载：'.$result;
            return response()->json($res);
        }

        $command2 = 'systemctl reload nginx';
        $result2 = shell_exec($command2);
        if (str_contains($result2, 'done')) {
            $res['data'] = 'OpenResty已重载';
        } else {
            $res['msg'] = 'error';
            $res['data'] = 'OpenResty重载失败';
        }
        return response()->json($res);
    }

    public function getConfig()
    {
        $res['code'] = 0;
        $res['msg'] = 'success';
        $res['data'] = file_get_contents('/www/server/nginx/conf/nginx.conf');
        return response()->json($res);
    }

    public function saveConfig()
    {
        $res['code'] = 0;
        $res['msg'] = 'success';
        // 获取配置内容
        $config = Request::post('config');
        // 备份一份旧配置
        $old_config = file_get_contents('/www/server/nginx/conf/nginx.conf');
        // 写入配置
        $result = file_put_contents('/www/server/nginx/conf/nginx.conf', $config);
        // 测试配置是否正确
        $test = shell_exec('nginx -t 2>&1');
        // 判断结果
        if (!str_contains($test, 'test is successful')) {
            // 测试失败，则不允许保存
            $res['msg'] = 'error';
            $res['data'] = 'OpenResty配置有误，请修正后再保存：'.$test;
            // 恢复旧配置
            file_put_contents('/www/server/nginx/conf/nginx.conf', $old_config);
            return response()->json($res);
        } else {
            // 测试成功，则重载OpenResty
            shell_exec('systemctl reload nginx');
            $res['data'] = 'OpenResty主配置已保存';
            return response()->json($res);
        }
    }

    public function load()
    {
        $raw_status = HTTP::get('http://127.0.0.1/nginx_status')->body();

        $res['code'] = 0;
        $res['msg'] = 'success';
        $res['data'][0]['name'] = '工作进程';
        $res['data'][0]['value'] = (int) shell_exec("ps aux|grep nginx|grep 'worker process'|wc -l");
        $res['data'][1]['name'] = '内存占用';
        $res['data'][1]['value'] = round(shell_exec("ps aux|grep nginx|grep 'worker process'|awk '{memsum+=$6};END {print memsum}'") / 1024,
                2).'MB';

        // 使用正则匹配Active connections: 的值
        preg_match('/Active connections:\s+(\d+)/', $raw_status, $matches);
        $res['data'][2]['name'] = '活跃连接数';
        $res['data'][2]['value'] = $matches[1] ?? 0;
        // 使用正则分别匹配server accepts handled requests的三个值
        preg_match('/server accepts handled requests\s+(\d+)\s+(\d+)\s+(\d+)/', $raw_status, $matches2);
        $res['data'][3]['name'] = '总连接次数';
        $res['data'][3]['value'] = $matches2[1] ?? 0;
        $res['data'][4]['name'] = '总握手次数';
        $res['data'][4]['value'] = $matches2[2] ?? 0;
        $res['data'][5]['name'] = '总请求次数';
        $res['data'][5]['value'] = $matches2[3] ?? 0;
        // 使用正则匹配Reading: 的值
        preg_match('/Reading:\s+(\d+)/', $raw_status, $matches3);
        $res['data'][6]['name'] = '请求数';
        $res['data'][6]['value'] = $matches3[1] ?? 0;
        // 使用正则匹配Writing: 的值
        preg_match('/Writing:\s+(\d+)/', $raw_status, $matches4);
        $res['data'][7]['name'] = '响应数';
        $res['data'][7]['value'] = $matches4[1] ?? 0;
        // 使用正则匹配Waiting: 的值
        preg_match('/Waiting:\s+(\d+)/', $raw_status, $matches5);
        $res['data'][8]['name'] = '驻留进程';
        $res['data'][8]['value'] = $matches5[1] ?? 0;

        return response()->json($res);
    }

    public function errorLog()
    {
        $res['code'] = 0;
        $res['msg'] = 'success';
        $res['data'] = file_get_contents('/www/wwwlogs/nginx_error.log');
        //如果data为换行符，则令返回空
        if ($res['data'] == "\n") {
            $res['data'] = '';
        }
        return response()->json($res);
    }

    public function cleanErrorLog()
    {
        $res['code'] = 0;
        $res['msg'] = 'success';
        shell_exec('echo "" > /www/wwwlogs/nginx_error.log');
        return response()->json($res);
    }

}
