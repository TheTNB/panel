<?php

namespace App\Console\Commands;

use App\Models\Monitor as MonitorModel;
use App\Models\Setting;
use Illuminate\Console\Command;
use Illuminate\Support\Carbon;

class Monitor extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'monitor';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '耗子Linux面板 - 系统监控';

    /**
     * Execute the console command.
     *
     * @return int
     */
    public function handle()
    {
        if (Setting::query()->where('name', 'monitor')->value('value')) {
            $info = self::getNowMonitor();
            MonitorModel::query()->create(['info' => json_encode($info)]);
            // 删除过期的记录
            $days = Setting::query()->where('name', 'monitor_days')->value('value');
            MonitorModel::query()->where('created_at', '<', Carbon::now()->subDays($days))->delete();
            $this->info(time().' 监控完成');
        } else {
            $this->info('监控未开启');
        }
        return Command::SUCCESS;
    }

    /**
     * 系统资源统计
     * @return array
     */
    private function getNowMonitor(): array
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
        $siblings_sum = 0;
        preg_match("/(\d+\.\d+), (\d+\.\d+), (\d+\.\d+)/", exec('uptime'), $uptime);
        $uptime_1 = $uptime[1] ?? 'No';

        $p_list = explode("\nprocessor", $cpu_info);
        foreach ($p_list as $key => $val) {
            preg_match("/physical id\s*:(.*)/", $val, $physical);
            preg_match("/cpu cores\s*:(.*)/", $val, $cores);
            preg_match("/siblings\s*:(.*)/", $val, $siblings);
            if (isset($physical[1])) {
                if (!in_array($physical[1], $physical_list)) {
                    $physical_sum += 1;

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
        $cpu_use = min($cpu_use, 100);


        // 内存使用率
        $result = explode("\n", shell_exec('free -m'));
        foreach ($result as $key => $val) {
            if (str_contains($val, 'Mem')) {
                $mem_list = preg_replace("/\s+/", " ", $val);
            } elseif (str_contains($val, 'Swap')) {
                $swap_list = preg_replace("/\s+/", " ", $val);
            }
        }
        $mem_arr = explode(' ', $mem_list);
        $swap_arr = explode(' ', $swap_list);
        // 内存大小MB
        $mem_total = $mem_arr[1];
        // Swap大小MB
        $swap_total = $swap_arr[1];
        // 使用中MB
        $mem_use = (str_contains($result[0], 'buff/cache')) ? $mem_arr[2] : ($mem_arr[2] - $mem_arr[5] - $mem_arr[6]);
        // Swap使用中MB
        $swap_use = $swap_arr[2];
        // 使用中%
        $mem_use_p = round($mem_use / $mem_total, 2) * 100;
        // Swap使用中%
        $swap_use_p = round($swap_use / $swap_total, 2) * 100;
        // 1分钟负载%
        $uptime_1_p = $uptime_1 * 10;
        $uptime_1_p = min($uptime_1_p, 100);

        // 构建返回数组
        $res['cpu_use'] = $cpu_use;
        $res['uptime'] = $uptime_1;
        $res['uptime_p'] = $uptime_1_p;
        $res['mem_total'] = $mem_total;
        $res['mem_use'] = $mem_use;
        $res['mem_use_p'] = $mem_use_p;
        $res['swap_total'] = $swap_total;
        $res['swap_use'] = $swap_use;
        $res['swap_use_p'] = $swap_use_p;
        $res['tx_now'] = $net_info2['tx'] - $net_info1['tx'];
        $res['rx_now'] = $net_info2['rx'] - $net_info1['rx'];
        return $res;
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
}
