<?php
/**
 * 耗子Linux面板 - 信息控制器
 * @author 耗子
 */

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Models\Plugin;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class InfosController extends Controller
{
    public function getMenu(): JsonResponse
    {
        $menu = array(
            "code" => 0,
            "msg" => "",
            "data" => array(
                array(
                    "name" => "home",
                    "title" => "主页",
                    "icon" => "layui-icon-home",
                    "jump" => "/"
                ),
                array(
                    "name" => "website",
                    "title" => "网站",
                    "icon" => "layui-icon-website",
                    "jump" => "website/list"
                ),
                array(
                    "name" => "database",
                    "title" => "数据库",
                    "icon" => "layui-icon-template-1",
                    "list" => array(
                        array(
                            "name" => "mysql",
                            "title" => "MySQL",
                            "jump" => "database/mysql"
                        ),
                        array(
                            "name" => "postgresql",
                            "title" => "PostgreSQL",
                            "jump" => "database/postgresql"
                        )
                    )
                ),
                array(
                    "name" => "monitor",
                    "title" => "监控",
                    "icon" => "layui-icon-chart-screen",
                    "jump" => "monitor"
                ),
                array(
                    "name" => "safe",
                    "title" => "安全",
                    "icon" => "layui-icon-auz",
                    "jump" => "safe"
                ),
                array(
                    "name" => "file",
                    "title" => "文件",
                    "icon" => "layui-icon-file",
                    "jump" => "file"
                ),
                array(
                    "name" => "plugin",
                    "title" => "插件",
                    "icon" => "layui-icon-app",
                    "jump" => "plugin"
                ),
                array(
                    "name" => "setting",
                    "title" => "设置",
                    "icon" => "layui-icon-set",
                    "jump" => "setting"
                ),
                array(
                    "name" => "logout",
                    "title" => "退出",
                    "icon" => "layui-icon-logout",
                    "jump" => "logout"
                )
            )
        );

        return response()->json($menu);
    }

    /**
     * 系统资源统计
     * @return JsonResponse
     */
    public function getNowMonitor(): JsonResponse
    {
        // 第一次获取网络信息
        $net_info1 = $this->getNetInfo();
        // 卡它一秒钟
        sleep(1);
        // 第二次获取网络信息
        $net_info2 = $this->getNetInfo();

        // CPU统计信息及负载
        $cpu_info = file_get_contents('/proc/cpuinfo');
        $physical_list = array();
        $physical_sum = 0;
        $cores_sum = 0;
        $siblings_sum = 0;
        preg_match("/model name\s*:(.*)/", $cpu_info, $name);
        preg_match("/vendor_id\s*:(.*)/", $cpu_info, $vendor);
        preg_match("/cpu family\s*:(.*)/", $cpu_info, $family);
        preg_match("/cpu MHz\s*:(.*)/", $cpu_info, $MHz);
        preg_match("/cache size\s*:(.*)/", $cpu_info, $cache);
        preg_match("/(\d+\.\d+), (\d+\.\d+), (\d+\.\d+)/", exec('uptime'), $uptime);
        $name = $name[1] ?? 'No';
        $vendor = $vendor[1] ?? 'No';
        $family = $family[1] ?? 'No';
        $MHz = isset($MHz[1]) ? number_format($MHz[1] / 1000, 2) : 'No';
        $cache = $cache[1] ?? 'No';
        $uptime_1 = $uptime[1] ?? 'No';
        $uptime_5 = $uptime[2] ?? 'No';
        $uptime_15 = $uptime[3] ?? 'No';

        $p_list = explode("\nprocessor", $cpu_info);
        foreach ($p_list as $key => $val) {
            preg_match("/physical id\s*:(.*)/", $val, $physical);
            preg_match("/cpu cores\s*:(.*)/", $val, $cores);
            preg_match("/siblings\s*:(.*)/", $val, $siblings);
            if (isset($physical[1])) {
                if (!in_array($physical[1], $physical_list)) {
                    $physical_sum += 1;
                    if (isset($cores[1])) {
                        $cores_sum += $cores[1];
                    }

                    if (isset($siblings[1])) {
                        $siblings_sum += $siblings[1];
                    }
                }
                $physical_list[] = $physical[1];
            }
        }

        // CPU使用率
        $cpu_use = 0.1;

        $result = explode("\n", shell_exec('ps aux'));
        foreach ($result as $key => $val) {
            $val = preg_replace("/\s+/", " ", $val);
            $val = (explode(' ', $val));
            $cpu_use += isset($val[2]) ? (float) $val[2] : 0;
        }
        $cpu_use = $siblings_sum > 0 ? ($cpu_use / $siblings_sum) : $cpu_use;
        $cpu_use = round($cpu_use, 2);
        $cpu_use = $cpu_use > 100 ? 100 .'%' : $cpu_use.'%';


        // 内存使用率
        $result = explode("\n", shell_exec('free -m'));
        foreach ($result as $key => $val) {
            if (str_contains($val, 'Mem')) {
                $mem_list = preg_replace("/\s+/", " ", $val);
            }
        }
        $mem_arr = explode(' ', $mem_list);
        // 内存大小MB
        $mem_total = $mem_arr[1];
        // 使用中MB
        $mem_use = (str_contains($result[0], 'buff/cache')) ? $mem_arr[2] : ($mem_arr[2] - $mem_arr[5] - $mem_arr[6]);
        // 使用中%
        $mem_use_p = round($mem_use / $mem_total, 2) * 100 .'%';
        // 1分钟负载%
        $uptime_1_p = $uptime_1 * 10;
        $uptime_1_p = $uptime_1_p > 100 ? 100 .'%' : $uptime_1_p.'%';
        // 5分钟负载%
        $uptime_5_p = $uptime_5 * 10;
        $uptime_5_p = $uptime_5_p > 100 ? 100 .'%' : $uptime_5_p.'%';
        // 15分钟负载%
        $uptime_15_p = $uptime_15 * 10;
        $uptime_15_p = $uptime_15_p > 100 ? 100 .'%' : $uptime_15_p.'%';

        // 构建返回数组
        $res['code'] = 0;
        $res['msg'] = 'success';
        $res['data'] = [
            'cpu_use' => $cpu_use,
            'uptime_1' => $uptime_1,
            'uptime_1_p' => $uptime_1_p,
            'uptime_5' => $uptime_5,
            'uptime_5_p' => $uptime_5_p,
            'uptime_15' => $uptime_15,
            'uptime_15_p' => $uptime_15_p,
            'mem_total' => $mem_total,
            'mem_use' => $mem_use,
            'mem_use_p' => $mem_use_p,
            'tx_total' => $this->formatBytes($net_info1['tx']),
            'rx_total' => $this->formatBytes($net_info1['rx']),
            'tx_now' => $this->formatBytes($net_info2['tx'] - $net_info1['tx']),
            'rx_now' => $this->formatBytes($net_info2['rx'] - $net_info1['rx']),
            'cpu_info' => [
                'name' => $name,
                'cores' => $cores_sum,
                'physical' => $physical_sum,
                'siblings' => $siblings_sum,
                'vendor' => $vendor,
                'family' => $family,
                'MHz' => $MHz,
                'cache' => $cache,
            ],
        ];
        return response()->json($res);
    }

    /**
     * 获取系统信息
     * @return JsonResponse
     */
    public function getSystemInfo(): JsonResponse
    {
        $os_name = file_get_contents('/etc/redhat-release');
        $uptime = file_get_contents('/proc/uptime');
        $uptime = explode(' ', $uptime)[0];
        $uptime = round($uptime / 86400, 1);
        $panel_version = config('panel.version');
        // 构建返回数组
        $res['code'] = 0;
        $res['msg'] = 'success';
        $res['data'] = [
            'os_name' => $os_name,
            'uptime' => $uptime,
            'panel_version' => $panel_version,
        ];
        return response()->json($res);
    }

    /**
     * 获取首页插件列表
     */
    public function getHomePlugins(): JsonResponse
    {
        $plugins = Plugin::where('show', 1)->get();
        // 判空
        if ($plugins->isEmpty()) {
            $res['code'] = 0;
            $res['msg'] = 'success';
            $res['data'] = [];
        } else {
            $plugins = $plugins->toArray();
            $res['code'] = 0;
            $res['msg'] = 'success';
            $res['data'] = $plugins;
        }
        return response()->json($res);
    }

    /**
     * 获取网络统计信息
     * @return array
     */
    private function getNetInfo(): array
    {
        $net_result = file_get_contents('/proc/net/dev');
        $net_result = explode("\n", $net_result);
        foreach ($net_result as $key => $val) {
            if ($key < 2) {
                continue;
            }
            $val = str_replace(':', ' ', trim($val));
            $val = preg_replace("/[ ]+/", " ", $val);
            $arr = explode(' ', $val);
            if (!empty($arr[0])) {
                $arr = array($arr[0], $arr[1], $arr[9]);
                $all_rs[$arr[0].$key] = $arr;
            }
        }
        ksort($all_rs);
        $tx = 0;
        $rx = 0;
        foreach ($all_rs as $key => $val) {
            // 排除本地lo
            if (str_contains($key, 'lo')) {
                continue;
            }
            $tx += $val[2];
            $rx += $val[1];
        }
        $res['tx'] = $tx;
        $res['rx'] = $rx;
        return $res;
    }

    /**
     * 格式化bytes
     * @param $size
     * @return string
     */
    private function formatBytes($size): string
    {
        $size = is_numeric($size) ? $size : 0;
        $units = array(' B', ' KB', ' MB', ' GB', ' TB');
        for ($i = 0; $size >= 1024 && $i < 4; $i++) {
            $size /= 1024;
        }
        return round($size, 2).$units[$i];
    }

    /**
     * 获取已安装的数据库和PHP版本
     */
    public function getInstalledDbAndPhp()
    {
        // 判断mysql插件目录是否存在
        if (is_dir('/www/panel/plugins/mysql')) {
            $mysql_version = 80;
        } else {
            $mysql_version = false;
        }
        /**
         * TODO: PostgreSQL版本
         */
        $php_version = shell_exec('ls /www/server/php');
        $php_version = trim($php_version);

        if (!empty($php_version)) {
            $php_versions = explode("\n", $php_version);
        }
        $php_versions[] = '00';
        unset($php_versions[array_search('panel', $php_versions)]);

        $res['code'] = 0;
        $res['msg'] = 'success';
        $res['data'] = array(
            'db_version' => [
                'mysql' => $mysql_version,
                'postgresql' => false
            ],
            'php_version' => $php_versions
        );
        return response()->json($res);
    }
}
