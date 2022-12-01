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
        $netInfo1 = getNetInfo();
        // 卡它一秒钟
        sleep(1);
        // 第二次获取网络信息
        $netInfo2 = getNetInfo();

        // CPU统计信息及负载
        $cpuInfoRaw = file_get_contents('/proc/cpuinfo');
        $physicalArr = array();
        $physicalSum = 0;
        $coresSum = 0;
        $siblingsSum = 0;
        preg_match("/model name\s*:\s(.*)/", $cpuInfoRaw, $cpuName);
        preg_match("/vendor_id\s*:\s(.*)/", $cpuInfoRaw, $cpuVendor);
        preg_match("/cpu family\s*:\s(.*)/", $cpuInfoRaw, $cpuFamily);
        preg_match("/cpu MHz\s*:\s(.*)/", $cpuInfoRaw, $cpuFreq);
        preg_match("/cache size\s*:\s(.*)/", $cpuInfoRaw, $cpuCache);
        preg_match("/(\d+\.\d+), (\d+\.\d+), (\d+\.\d+)/", exec('uptime'), $uptime);
        $cpuName = $cpuName[1] ?? 'No';
        $cpuVendor = $cpuVendor[1] ?? 'No';
        $cpuFamily = $cpuFamily[1] ?? 'No';
        $cpuFreq = isset($cpuFreq[1]) ? round($cpuFreq[1], 2) : 'No';
        $cpuCache = $cpuCache[1] ?? 'No';
        $uptime1 = $uptime[1] ?? 0;
        $uptime5 = $uptime[2] ?? 0;
        $uptime15 = $uptime[3] ?? 0;

        $processorArr = explode("\nprocessor", $cpuInfoRaw);
        foreach ($processorArr as $v) {
            preg_match("/physical id\s*:\s(.*)/", $v, $physical);
            preg_match("/cpu cores\s*:\s(.*)/", $v, $cores);
            preg_match("/siblings\s*:\s(.*)/", $v, $siblings);
            if (isset($physical[1])) {
                if (!in_array($physical[1], $physicalArr)) {
                    $physicalSum += 1;
                    if (isset($cores[1])) {
                        $coresSum += $cores[1];
                    }

                    if (isset($siblings[1])) {
                        $siblingsSum += $siblings[1];
                    }
                }
                $physicalArr[] = $physical[1];
            }
        }

        // CPU使用率
        $cpuUse = 0.1;
        $cpuRaw = explode("\n", shell_exec('ps aux'));
        // 弹出第一项和最后一项
        array_pop($cpuRaw);
        array_shift($cpuRaw);
        // 获取当前php进程的pid
        $pid = getmypid();
        foreach ($cpuRaw as $v) {
            $v = preg_replace("/\s+/", " ", $v);
            $v = (explode(' ', $v));
            // 排除当前进程
            if ($v[1] == $pid) {
                continue;
            }
            $cpuUse += isset($v[2]) ? (float) $v[2] : 0;
        }
        $cpuUse = $siblingsSum > 0 ? ($cpuUse / $siblingsSum) : $cpuUse;
        $cpuUse = round($cpuUse, 2);
        $cpuUse = $cpuUse > 100 ? 100 .'%' : $cpuUse.'%';

        // 内存使用率
        $memRaw = explode("\n", shell_exec('free -m'));
        foreach ($memRaw as $v) {
            if (str_contains($v, 'Mem')) {
                $memList = preg_replace("/\s+/", " ", $v);
            }
        }
        $memArr = explode(' ', $memList);
        // 内存大小MB
        $memTotal = $memArr[1];
        // 使用中MB
        $memUse = (str_contains($memRaw[0], 'buff/cache')) ? $memArr[2] : ($memArr[2] - $memArr[5] - $memArr[6]);
        // 使用中%
        $memUseP = round($memUse / $memTotal, 2) * 100 .'%';
        // 1分钟负载%
        $uptime1P = $uptime1 * 10;
        $uptime1P = $uptime1P > 100 ? 100 .'%' : $uptime1P.'%';
        // 5分钟负载%
        $uptime5P = $uptime5 * 10;
        $uptime5P = $uptime5P > 100 ? 100 .'%' : $uptime5P.'%';
        // 15分钟负载%
        $uptime15P = $uptime15 * 10;
        $uptime15P = $uptime15P > 100 ? 100 .'%' : $uptime15P.'%';

        // 构建返回数组
        $res['code'] = 0;
        $res['msg'] = 'success';
        $res['data'] = [
            'cpu_use' => $cpuUse,
            'uptime_1' => $uptime1,
            'uptime_1_p' => $uptime1P,
            'uptime_5' => $uptime5,
            'uptime_5_p' => $uptime5P,
            'uptime_15' => $uptime15,
            'uptime_15_p' => $uptime15P,
            'mem_total' => $memTotal,
            'mem_use' => $memUse,
            'mem_use_p' => $memUseP,
            'tx_total' => formatBytes($netInfo1['tx']),
            'rx_total' => formatBytes($netInfo1['rx']),
            'tx_now' => formatBytes($netInfo2['tx'] - $netInfo1['tx']),
            'rx_now' => formatBytes($netInfo2['rx'] - $netInfo1['rx']),
            'cpu_info' => [
                'name' => $cpuName,
                'cores' => $coresSum,
                'physical' => $physicalSum,
                'siblings' => $siblingsSum,
                'vendor' => $cpuVendor,
                'family' => $cpuFamily,
                'freq' => $cpuFreq,
                'cache' => $cpuCache,
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
        $plugins = Plugin::query()->where('show', 1)->get();
        // 判空
        if ($plugins->isEmpty()) {
            $res['code'] = 0;
            $res['msg'] = 'success';
            $res['data'] = [];
        } else {
            $plugins = $plugins->toArray();
            $plugins = array_map(function ($item) {
                $item['name'] = PLUGINS[$item['slug']]['name'];
                return $item;
            }, $plugins);
            $res['code'] = 0;
            $res['msg'] = 'success';
            $res['data'] = $plugins;
        }
        return response()->json($res);
    }

    /**
     * 获取已安装的数据库和PHP版本
     */
    public function getInstalledDbAndPhp(): JsonResponse
    {
        $dbVersions = [];
        // 判断mysql插件是否安装
        if (isset(PLUGINS['mysql'])) {
            $dbVersions['mysql'] = PLUGINS['mysql']['version'];
        } else {
            $dbVersions['mysql'] = false;
        }
        // 判断postgresql插件是否安装
        if (isset(PLUGINS['postgresql'])) {
            $dbVersions['postgresql'] = PLUGINS['postgresql']['version'];
        } else {
            $dbVersions['postgresql'] = false;
        }
        // 循环获取已安装的PHP版本
        $php_versions = Plugin::query()->where('slug', 'like', 'php%')->get();
        $php_versions = $php_versions->toArray();
        $php_versions = array_column($php_versions, 'slug');
        $php_versions = array_map(function ($item) {
            return str_replace('php', '', $item);
        }, $php_versions);

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
            'db_version' => $dbVersions,
            'php_version' => $php_versions
        );
        return response()->json($res);
    }
}
