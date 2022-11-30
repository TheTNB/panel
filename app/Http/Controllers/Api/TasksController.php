<?php
/**
 * 耗子Linux面板 - 任务控制器
 * @author 耗子
 */
namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Models\Task;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class TasksController extends Controller
{
    /**
     * 获取是否有进行中/未完成的任务
     * 前台会根据这个返回渲染顶部工具栏的图标
     * @return JsonResponse
     */
    public function getStatus(): JsonResponse
    {
        // 获取任务表中的等待中/进行中的任务
        $task = Task::query()->where('status', 'running')->orWhere('status', 'waiting')->first();
        $res['code'] = 0;
        $res['msg'] = 'success';
        if (empty($task)) {
            $res['data'] = false;
        } else {
            $res['data'] = true;
        }
        return response()->json($res);
    }

    /**
     * 获取进行中的那个任务
     * @return JsonResponse
     */
    public function getListRunning(): JsonResponse
    {
        // 获取进行中的任务列表
        $task = Task::query()->where('status', 'running')->get()->toArray();

        // 判断任务是否存在
        if (empty($task)) {
            $task[0] = "";
        }
        // 构建返回数据
        $res['code'] = 0;
        $res['msg'] = 'success';
        $res['data'] = $task[0];
        return response()->json($res);
    }

    /**
     * 获取等待中的任务列表
     * @return JsonResponse
     */
    public function getListWaiting(): JsonResponse
    {
        $task_list = Task::query()->where('status', 'waiting')->get()->toArray();
        $res['code'] = 0;
        $res['msg'] = 'success';
        // 判断任务是否存在
        if (empty($task_list)) {
            $res['code'] = 1;
            $res['msg'] = '无任务！';
        }
        $res['data'] = $task_list;
        return response()->json($res);
    }

    /**
     * 获取已完成的任务列表
     * @return JsonResponse
     */
    public function getListFinished(): JsonResponse
    {
        $task_list = Task::query()->where('status', 'finished')->get()->toArray();
        $res['code'] = 0;
        $res['msg'] = 'success';
        if (empty($task_list)) {
            $res['code'] = 1;
            $res['msg'] = '无任务！';
        }
        $res['data'] = $task_list;
        return response()->json($res);
    }

    /**
     * 获取单个任务的log
     * @param  Request  $request
     * @return JsonResponse
     */
    public function getTaskLog(Request $request): JsonResponse
    {
        $name = $request->get('name');
        $log_file = Task::query()->where('name', $name)->value('log');

        $log = shell_exec('tail -n 30 '.$log_file);

        $res['code'] = 0;
        $res['msg'] = 'success';
        $res['data'] = $log;
        return response()->json($res);
    }

    /**
     * 删除任务
     * @param  Request  $request
     * @return JsonResponse
     */
    public function deleteTask(Request $request): JsonResponse
    {
        $name = $request->get('name');
        $res['code'] = 0;
        $res['msg'] = 'success';
        $res['data'] = Task::query()->where('name', $name)->delete();
        return response()->json($res);
    }
}
