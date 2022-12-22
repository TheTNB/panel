<?php

namespace App\Console;

use Illuminate\Console\Scheduling\Schedule;
use Illuminate\Foundation\Console\Kernel as ConsoleKernel;
use App\Models\Cron;
use Illuminate\Support\Carbon;

class Kernel extends ConsoleKernel
{
    /**
     * Define the application's command schedule.
     *
     * @param  \Illuminate\Console\Scheduling\Schedule  $schedule
     * @return void
     */
    protected function schedule(Schedule $schedule)
    {
        $schedule->command('monitor')->everyMinute();
        // 查询所有计划任务
        $crons = Cron::all();
        foreach ($crons as $cron) {
            $file = '/www/server/cron/'.$cron->shell;
            // 检查文件是否存在，及所有者是否为root
            if (!file_exists($file) || fileowner($file) != 0) {
                file_put_contents('/www/server/cron/logs/'.$cron->id.'.log',
                    PHP_EOL.'耗子Linux面板：检测到脚本文件异常，为确保安全已终止运行，如果你不知道发生了什么，这通常意味着服务器已被入侵。'.PHP_EOL,
                    FILE_APPEND);
                continue;
            }
            $schedule->exec('bash '.escapeshellarg($file))->withoutOverlapping()->cron($cron->time)->appendOutputTo('/www/server/cron/logs/'.$cron->id.'.log')->when(function (
            ) use ($cron) {
                return (boolean) $cron->status;
            })->after(function () use ($cron) {
                $cron->updated_at = now();
                $cron->save();
            })->onSuccess(function () use ($cron) {
                file_put_contents('/www/server/cron/logs/'.$cron->id.'.log',
                    PHP_EOL.Carbon::now()->toDateTimeString().' 任务执行成功'.PHP_EOL, FILE_APPEND);
            })->onFailure(function () use ($cron) {
                file_put_contents('/www/server/cron/logs/'.$cron->id.'.log',
                    PHP_EOL.Carbon::now()->toDateTimeString().' 任务执行失败'.PHP_EOL, FILE_APPEND);
            });
        }
    }

    /**
     * Register the commands for the application.
     *
     * @return void
     */
    protected function commands()
    {
        $this->load(__DIR__.'/Commands');

        require base_path('routes/console.php');
    }
}
