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
        $netInfo1 = getNetInfo();
        // 卡它一秒钟
        sleep(1);
        // 第二次获取网络信息
        $netInfo2 = getNetInfo();

        // CPU统计信息及负载
        $cpuInfoRaw = file_get_contents('/proc/cpuinfo');
        $physicalArr = array();
        $siblingsSum = 0;
        preg_match("/(\d+\.\d+), (\d+\.\d+), (\d+\.\d+)/", exec('uptime'), $uptime);
        $uptime1 = $uptime[1] ?? 0;

        $processorArr = explode("\nprocessor", $cpuInfoRaw);
        foreach ($processorArr as $v) {
            preg_match("/physical id\s*:\s(.*)/", $v, $physical);
            preg_match("/siblings\s*:\s(.*)/", $v, $siblings);
            if (isset($physical[1])) {
                if (!in_array($physical[1], $physicalArr)) {
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
        $cpuUse = min($cpuUse, 100);


        // 内存使用率
        $memRaw = explode("\n", shell_exec('free -m'));
        foreach ($memRaw as $v) {
            if (str_contains($v, 'Mem')) {
                $memList = preg_replace("/\s+/", " ", $v);
            } elseif (str_contains($v, 'Swap')) {
                $swapList = preg_replace("/\s+/", " ", $v);
            }
        }
        $memArr = explode(' ', $memList);
        $swapArr = explode(' ', $swapList);
        // 内存大小MB
        $memTotal = $memArr[1];
        // Swap大小MB
        $swapTotal = $swapArr[1];
        // 使用中MB
        $memUse = (str_contains($memRaw[0], 'buff/cache')) ? $memArr[2] : ($memArr[2] - $memArr[5] - $memArr[6]);
        // Swap使用中MB
        $swapUse = $swapArr[2];
        // 使用中%
        $memUseP = round($memUse / $memTotal, 2) * 100;
        // Swap使用中%
        $swapUseP = round($swapUse / $swapTotal, 2) * 100;
        // 1分钟负载%
        $uptime1P = round(min($uptime1 * 10, 100), 2);

        // 构建返回数组
        $res['cpu_use'] = $cpuUse;
        $res['uptime'] = $uptime1;
        $res['uptime_p'] = $uptime1P;
        $res['mem_total'] = $memTotal;
        $res['mem_use'] = $memUse;
        $res['mem_use_p'] = $memUseP;
        $res['swap_total'] = $swapTotal;
        $res['swap_use'] = $swapUse;
        $res['swap_use_p'] = $swapUseP;
        $res['tx_now'] = $netInfo2['tx'] - $netInfo1['tx'];
        $res['rx_now'] = $netInfo2['rx'] - $netInfo1['rx'];
        return $res;
    }
}
