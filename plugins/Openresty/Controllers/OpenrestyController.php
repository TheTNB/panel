<?php
/**
 * Name: OpenResty插件控制器
 * Author:耗子
 * Date: 2022-12-10
 */

namespace Plugins\Openresty\Controllers;


use App\Http\Controllers\Controller;

use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Http;

class OpenrestyController extends Controller
{

    /**
     * 获取运行状态
     * @return JsonResponse
     */
    public function status(): JsonResponse
    {
        $status = shell_exec('systemctl status nginx | grep Active | grep -v grep | awk \'{print $2}\'');
        // 格式化掉换行符
        $status = trim($status);
        if (empty($status)) {
            return response()->json(['code' => 1, 'msg' => '获取服务运行状态失败']);
        }
        if ($status == 'active') {
            $status = 1;
        } else {
            $status = 0;
        }

        // 返回结果
        $res['code'] = 0;
        $res['msg'] = 'success';
        $res['data'] = $status;
        return response()->json($res);
    }

    /**
     * 重启
     * @return JsonResponse
     */
    public function restart(): JsonResponse
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

        shell_exec('systemctl restart nginx');
        $status = shell_exec('systemctl status nginx | grep Active | grep -v grep | awk \'{print $2}\'');
        // 格式化掉换行符
        $status = trim($status);
        if (empty($status)) {
            return response()->json(['code' => 1, 'msg' => '获取服务运行状态失败']);
        }
        if ($status != 'active') {
            return response()->json(['code' => 1, 'msg' => '重启服务失败']);
        }

        // 返回结果
        $res['code'] = 0;
        $res['msg'] = 'success';
        return response()->json($res);
    }

    /**
     * 重载
     * @return JsonResponse
     */
    public function reload(): JsonResponse
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

        shell_exec('systemctl reload nginx');
        $status = shell_exec('systemctl status nginx | grep Active | grep -v grep | awk \'{print $2}\'');
        // 格式化掉换行符
        $status = trim($status);
        if (empty($status)) {
            return response()->json(['code' => 1, 'msg' => '获取服务运行状态失败']);
        }
        if ($status != 'active') {
            return response()->json(['code' => 1, 'msg' => '重载服务失败']);
        }

        // 返回结果
        $res['code'] = 0;
        $res['msg'] = 'success';
        return response()->json($res);
    }

    /**
     * 获取配置文件
     * @return JsonResponse
     */
    public function getConfig(): JsonResponse
    {
        $res['code'] = 0;
        $res['msg'] = 'success';
        $res['data'] = @file_get_contents('/www/server/nginx/conf/nginx.conf');
        return response()->json($res);
    }

    /**
     * 保存配置文件
     * @param  Request  $request
     * @return JsonResponse
     */
    public function saveConfig(Request $request): JsonResponse
    {
        // 获取配置内容
        $config = $request->input('config');
        // 备份一份旧配置
        $old_config = @file_get_contents('/www/server/nginx/conf/nginx.conf');
        // 写入配置
        @file_put_contents('/www/server/nginx/conf/nginx.conf', $config);
        // 测试配置是否正确
        $test = shell_exec('nginx -t 2>&1');
        // 判断结果
        if (!str_contains($test, 'test is successful')) {
            // 测试失败，则不允许保存
            @file_put_contents('/www/server/nginx/conf/nginx.conf', $old_config);
            return response()->json(['code' => 1, 'msg' => 'OpenResty配置有误，请修正后再保存：'.$test]);
        } else {
            // 测试成功，则重载OpenResty
            shell_exec('systemctl reload nginx');
            return response()->json(['code' => 0, 'msg' => 'success']);
        }
    }

    /**
     * 获取负载状态
     * @return JsonResponse
     */
    public function load(): JsonResponse
    {
        $status = HTTP::get('http://127.0.0.1/nginx_status');
        // 判断状态码
        if ($status->status() != 200) {
            return response()->json(['code' => 1, 'msg' => '获取状态失败']);
        }
        $statusRaw = $status->body();

        $res['code'] = 0;
        $res['msg'] = 'success';
        $res['data'][0]['name'] = '工作进程';
        $res['data'][0]['value'] = (int) shell_exec("ps aux|grep nginx|grep 'worker process'|wc -l");
        $res['data'][1]['name'] = '内存占用';
        $res['data'][1]['value'] = round(shell_exec("ps aux|grep nginx|grep 'worker process'|awk '{memsum+=$6};END {print memsum}'") / 1024,
                2).'MB';

        // 使用正则匹配Active connections: 的值
        preg_match('/Active connections:\s+(\d+)/', $statusRaw, $matches);
        $res['data'][2]['name'] = '活跃连接数';
        $res['data'][2]['value'] = $matches[1] ?? 0;
        // 使用正则分别匹配server accepts handled requests的三个值
        preg_match('/server accepts handled requests\s+(\d+)\s+(\d+)\s+(\d+)/', $statusRaw, $matches2);
        $res['data'][3]['name'] = '总连接次数';
        $res['data'][3]['value'] = $matches2[1] ?? 0;
        $res['data'][4]['name'] = '总握手次数';
        $res['data'][4]['value'] = $matches2[2] ?? 0;
        $res['data'][5]['name'] = '总请求次数';
        $res['data'][5]['value'] = $matches2[3] ?? 0;
        // 使用正则匹配Reading: 的值
        preg_match('/Reading:\s+(\d+)/', $statusRaw, $matches3);
        $res['data'][6]['name'] = '请求数';
        $res['data'][6]['value'] = $matches3[1] ?? 0;
        // 使用正则匹配Writing: 的值
        preg_match('/Writing:\s+(\d+)/', $statusRaw, $matches4);
        $res['data'][7]['name'] = '响应数';
        $res['data'][7]['value'] = $matches4[1] ?? 0;
        // 使用正则匹配Waiting: 的值
        preg_match('/Waiting:\s+(\d+)/', $statusRaw, $matches5);
        $res['data'][8]['name'] = '驻留进程';
        $res['data'][8]['value'] = $matches5[1] ?? 0;

        return response()->json($res);
    }

    /**
     * 获取错误日志
     * @return JsonResponse
     */
    public function errorLog(): JsonResponse
    {
        $res['code'] = 0;
        $res['msg'] = 'success';
        $res['data'] = @file_get_contents('/www/wwwlogs/nginx_error.log');
        //如果data为换行符，则令返回空
        if ($res['data'] == "\n") {
            $res['data'] = '';
        }
        return response()->json($res);
    }

    /**
     * 清空错误日志
     * @return JsonResponse
     */
    public function cleanErrorLog(): JsonResponse
    {
        $res['code'] = 0;
        $res['msg'] = 'success';
        shell_exec('echo "" > /www/wwwlogs/nginx_error.log');
        return response()->json($res);
    }

}
