<?php
/**
 * 耗子Linux面板 - 计划任务控制器
 * @author 耗子
 */

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Models\Cron;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Support\Carbon;
use Illuminate\Validation\ValidationException;

class CronsController extends Controller
{
    /**
     * 面板计划任务列表
     */
    public function getList(Request $request): JsonResponse
    {
        $limit = $request->input('limit', 10);

        $crons = Cron::query()->orderBy('id', 'desc')->paginate($limit);
        $cronData = [];

        foreach ($crons as $k => $v) {
            // 格式化时间
            $cronData[$k]['id'] = $v['id'];
            $cronData[$k]['name'] = $v['name'];
            $cronData[$k]['status'] = $v['status'];
            $cronData[$k]['type'] = $v['type'];
            $cronData[$k]['time'] = $v['time'];
            $cronData[$k]['shell'] = $v['shell'];
            $cronData[$k]['script'] = @file_get_contents('/www/server/cron/'.$v['shell']);
            $cronData[$k]['created_at'] = Carbon::create($v['created_at'])->toDateTimeString();
            $cronData[$k]['updated_at'] = Carbon::create($v['updated_at'])->toDateTimeString();
        }

        $data['code'] = 0;
        $data['msg'] = 'success';
        $data['count'] = $crons->total();
        $data['data'] = $cronData;
        return response()->json($data);
    }

    /**
     * 添加计划任务
     */
    public function add(Request $request): JsonResponse
    {
        // 消毒
        try {
            $credentials = $this->validate($request, [
                'name' => 'required|max:255',
                'time' => ['required', 'regex:/^((\*|\d+|\d+-\d+|\d+\/\d+|\d+-\d+\/\d+|\*\/\d+)(\,(\*|\d+|\d+-\d+|\d+\/\d+|\d+-\d+\/\d+|\*\/\d+))*\s?){5}$/'],
                'script' => 'required',
            ]);
        } catch (ValidationException $e) {
            return response()->json(['code' => 1, 'msg' => $e->getMessage()]);
        }

        // 将script写入shell文件
        $shellDir = '/www/server/cron/';
        $shellLogDir = '/www/server/cron/logs/';
        if (!is_dir($shellDir)) {
            mkdir($shellDir, 0755, true);
        }
        if (!is_dir($shellLogDir)) {
            mkdir($shellLogDir, 0755, true);
        }
        $shellFile = uniqid().'.sh';
        file_put_contents($shellDir.$shellFile, $credentials['script']);

        $cron = new Cron();
        $cron->name = $credentials['name'];
        $cron->status = 1;
        $cron->type = '脚本';
        $cron->time = $credentials['time'];
        $cron->shell = $shellFile;
        $cron->save();

        $data['code'] = 0;
        $data['msg'] = 'success';
        return response()->json($data);
    }

    /**
     * 修改计划任务
     */
    public function edit(Request $request): JsonResponse
    {
        // 消毒
        try {
            $credentials = $this->validate($request, [
                'id' => 'required|integer',
                'name' => 'required|max:255',
                'time' => ['required', 'regex:/^((\*|\d+|\d+-\d+|\d+\/\d+)(\,(\*|\d+|\d+-\d+|\d+\/\d+))*\s?){5}$/'],
                'script' => 'required',
            ]);
        } catch (ValidationException $e) {
            return response()->json(['code' => 1, 'msg' => $e->getMessage()]);
        }

        $cron = Cron::query()->find($credentials['id']);
        $cron->name = $credentials['name'];
        $cron->time = $credentials['time'];
        // 将script写入shell文件
        $shellDir = '/www/server/cron/';
        $shellLogDir = '/www/server/cron/logs/';
        if (!is_dir($shellDir)) {
            mkdir($shellDir, 0755, true);
        }
        if (!is_dir($shellLogDir)) {
            mkdir($shellLogDir, 0755, true);
        }
        $shellFile = $cron->shell;
        file_put_contents($shellDir.$shellFile, $credentials['script']);
        $cron->save();

        $data['code'] = 0;
        $data['msg'] = 'success';
        return response()->json($data);
    }

    /**
     * 删除计划任务
     */
    public function delete(Request $request): JsonResponse
    {
        // 消毒
        try {
            $credentials = $this->validate($request, [
                'id' => 'required|integer',
            ]);
        } catch (ValidationException $e) {
            return response()->json(['code' => 1, 'msg' => $e->getMessage()]);
        }

        $cron = Cron::query()->find($credentials['id']);
        // 删除shell文件
        $shellDir = '/www/server/cron/';
        $shellFile = $cron->shell;
        @unlink($shellDir.$shellFile);
        // 删除日志文件
        $shellLogDir = '/www/server/cron/logs/';
        $shellLogFile = $shellFile.'.log';
        @unlink($shellLogDir.$shellLogFile);
        $cron->delete();

        $data['code'] = 0;
        $data['msg'] = 'success';
        return response()->json($data);
    }

    /**
     * 修改计划任务状态
     */
    public function setStatus(Request $request): JsonResponse
    {
        // 消毒
        try {
            $credentials = $this->validate($request, [
                'id' => 'required|integer',
                'status' => 'required|integer',
            ]);
        } catch (ValidationException $e) {
            return response()->json(['code' => 1, 'msg' => $e->getMessage()]);
        }

        $cron = Cron::query()->find($credentials['id']);
        $cron->status = $credentials['status'];
        $cron->save();

        $data['code'] = 0;
        $data['msg'] = 'success';
        return response()->json($data);
    }

    /**
     * 获取计划任务日志
     */
    public function getLog(Request $request): JsonResponse
    {
        // 消毒
        try {
            $credentials = $this->validate($request, [
                'id' => 'required|integer',
            ]);
        } catch (ValidationException $e) {
            return response()->json(['code' => 1, 'msg' => $e->getMessage()]);
        }

        $log = @file_get_contents('/www/server/cron/logs/'.$credentials['id'].'.log');
        if ($log === false) {
            $log = '暂无日志';
        }

        $data['code'] = 0;
        $data['msg'] = 'success';
        $data['data'] = $log;
        return response()->json($data);
    }
}
