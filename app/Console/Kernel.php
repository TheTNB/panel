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
            $schedule->exec('bash /www/server/cron/'.$cron->shell)->withoutOverlapping()->cron($cron->time)->appendOutputTo('/www/server/cron/logs/'.$cron->id.'.log')->when(function (
            ) use ($cron) {
                return (boolean) $cron->status;
            })->after(function () use ($cron) {
                $cron->updated_at = now();
                $cron->save();
            })->onSuccess(function () use ($cron) {
                shell_exec('echo "'.Carbon::now()->toDateTimeString().' 任务执行成功" >> /www/server/cron/logs/'.$cron->id.'.log');
            })->onFailure(function () use ($cron) {
                shell_exec('echo "'.Carbon::now()->toDateTimeString().' 任务执行失败" >> /www/server/cron/logs/'.$cron->id.'.log');
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
