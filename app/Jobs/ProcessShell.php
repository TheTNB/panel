<?php

namespace App\Jobs;

use Illuminate\Bus\Queueable;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Foundation\Bus\Dispatchable;
use Illuminate\Queue\InteractsWithQueue;
use Illuminate\Queue\SerializesModels;
use App\Models\Task;

class ProcessShell implements ShouldQueue
{
    use Dispatchable, InteractsWithQueue, Queueable, SerializesModels;

    /**
     * 任务最大尝试次数
     * @var int
     */
    public int $tries = 1;
    /**
     * 任务运行的超时时间
     * @var int
     */
    public int $timeout = 7200;

    /**
     * 任务名
     */
    public string $task_id;

    /**
     * 任务脚本
     */
    public string $shell;

    /**
     * Create a new job instance.
     *
     * @return void
     */
    public function __construct($task_id)
    {
        $this->task_id = $task_id;
    }

    /**
     * Execute the job.
     *
     * @return void
     */
    public function handle(): void
    {
        // 检查当前是否有任务正在运行
        $taskCheck = Task::query()->where('status', 'running')->get();
        if ($taskCheck->isNotEmpty()) {
            $this->release(10);
            return;
        }
        // 查询任务
        $task = Task::query()->where('id', $this->task_id)->first();
        echo $task->name."开始执行".PHP_EOL;
        // 更新任务状态为running
        $task->job_id = $this->job->getJobId();
        $task->status = 'running';
        $task->save();
        shell_exec($task->shell);
        // 更新任务状态
        $task->status = 'finished';
        $task->save();
        echo $task->name."执行完毕".PHP_EOL;
    }
}
