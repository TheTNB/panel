<?php
/**
 * 耗子Linux面板 - 监控控制器
 * @author 耗子
 */

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Models\Monitor;
use App\Models\Setting;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Support\Carbon;

class MonitorsController extends Controller
{
    /**
     * 修改监控开关
     */
    public function setMonitorSwitch(Request $request): JsonResponse
    {
        $switch = $request->input('switch');
        if ($switch) {
            $status = true;
        } else {
            $status = false;
        }
        Setting::query()->where('name', 'monitor')->update(['value' => $status]);
        return response()->json(['code' => 0, 'msg' => '修改成功']);
    }

    /**
     * 修改保存天数
     */
    public function setMonitorSaveDays(Request $request): JsonResponse
    {
        $days = $request->input('days');
        Setting::query()->where('name', 'monitor_days')->update(['value' => $days]);
        return response()->json(['code' => 0, 'msg' => '修改成功']);
    }

    /**
     * 清空监控数据
     */
    public function clearMonitorData(): JsonResponse
    {
        Monitor::query()->truncate();
        return response()->json(['code' => 0, 'msg' => '清空成功']);
    }

    /**
     * 获取监控开关和保存天数
     */
    public function getMonitorSwitchAndDays(): JsonResponse
    {
        $monitor = Setting::query()->where('name', 'monitor')->first();
        $monitor_days = Setting::query()->where('name', 'monitor_days')->first();
        return response()->json([
            'code' => 0, 'msg' => '获取成功',
            'data' => ['monitor' => $monitor->value, 'monitor_days' => $monitor_days->value]
        ]);
    }

    /**
     * 获取监控数据
     */
    public function getMonitorData(Request $request): JsonResponse
    {
        $start = $request->input('start') ?? now();
        $end = $request->input('end') ?? now();
        $start = Carbon::create($start);
        $end = Carbon::create($end);
        $data = Monitor::query()->where('created_at', '>=', $start)->where('created_at', '<=', $end)->get()->toArray();
        $res['code'] = 0;
        $res['msg'] = 'success';
        if (empty($data)) {
            $res['code'] = 1;
            $res['msg'] = '暂无数据';
            return response()->json($res);
        }
        foreach ($data as $key => $value) {
            $info = json_decode($value['info'], true);
            $res['data']['times'][] = Carbon::create($value['created_at'])->tz(config('app.timezone',
                'PRC'))->isoFormat('MM-DD HH:mm');
            $res['data']['uptime']['uptime'][] = round($info['uptime'], 2);
            $res['data']['cpu']['use'][] = round($info['cpu_use'], 2);
            $res['data']['memory']['mem_use'][] = round($info['mem_use'], 2);
            $res['data']['memory']['mem_use_p'][] = round($info['mem_use_p'], 2);
            $res['data']['memory']['swap_use'][] = round($info['swap_use'], 2);
            $res['data']['memory']['swap_use_p'][] = round($info['swap_use_p'], 2);
            $res['data']['network']['tx_now'][] = round($info['tx_now'] / 1024, 2);
            $res['data']['network']['rx_now'][] = round($info['rx_now'] / 1024, 2);
        }
        // 插入总内存大小
        $result = explode("\n", shell_exec('free -m'));
        foreach ($result as $key => $val) {
            if (str_contains($val, 'Mem')) {
                $mem_list = preg_replace("/\s+/", " ", $val);
                break;
            }
        }
        $mem_arr = explode(' ', $mem_list);
        // 内存大小MB
        $mem_total = $mem_arr[1];
        $res['data']['mem_total'] = round($mem_total, 2);
        return response()->json($res);
    }
}
