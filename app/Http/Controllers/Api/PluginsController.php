<?php
/**
 * 耗子Linux面板 - 插件控制器
 * @author 耗子
 */

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Jobs\ProcessShell;
use App\Models\Plugin;
use App\Models\Task;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Cache;
use Illuminate\Support\Facades\Http;


class PluginsController extends Controller
{

    /**
     * 面板插件列表
     * @return JsonResponse
     */
    public function getList(): JsonResponse
    {
        $data['code'] = 0;
        $data['msg'] = 'success';
        $data['data'] = $this->pluginList(false);
        foreach ($data['data'] as $k => $v) {
            // 获取首页显示状态
            $shows = Plugin::query()->pluck('show', 'slug');
            // 如果本地已安装，则显示本地名称
            $data['data'][$k]['name'] = PLUGINS[$v['slug']]['name'] ?? $data['data'][$k]['name'];
            // 已装版本
            $data['data'][$k]['install_version'] = PLUGINS[$v['slug']]['version'] ?? '';
            // 首页显示
            $data['data'][$k]['show'] = $shows[$v['slug']] ?? 0;
            // 去除不需要的字段
            unset($data['data'][$k]['url']);
            unset($data['data'][$k]['install']);
            unset($data['data'][$k]['uninstall']);
            unset($data['data'][$k]['update']);
            if (isset(PLUGINS[$v['slug']])) {
                $data['data'][$k]['control']['installed'] = true;
                $data['data'][$k]['control']['allow_uninstall'] = true;
                // 判断是否有更新
                $data['data'][$k]['control']['update'] = version_compare($v['version'],
                    $data['data'][$k]['install_version'], '>');
                if ($v['slug'] == 'openresty') {
                    $data['data'][$k]['control']['allow_uninstall'] = false;
                }
            } else {
                $data['data'][$k]['control']['installed'] = false;
                $data['data'][$k]['control']['allow_uninstall'] = false;
            }
        }
        return response()->json($data);
    }

    /**
     * 安装插件
     * @param  Request  $request
     * @return JsonResponse
     */
    public function install(Request $request): JsonResponse
    {
        // 消毒
        try {
            $credentials = $this->validate($request, [
                'slug' => 'required|max:255',
            ]);
            $slug = $credentials['slug'];
        } catch (Exception $e) {
            return response()->json(['code' => 1, 'msg' => $e->getMessage()]);
        }
        // 通过slug获取插件信息
        $pluginList = $this->pluginList();
        // 循环插件列表查找对应的插件信息
        $plugin_check = false;
        $plugin_data = [];
        foreach ($pluginList as $v) {
            if ($v['slug'] == $slug) {
                $plugin_data = $v;
                $plugin_check = true;
            }
        }
        // 判断插件是否存在
        if (!$plugin_check) {
            return response()->json(['code' => 1, 'msg' => '插件不存在']);
        }

        // 判断有无任务记录
        $task_check = Task::query()->where('name', '安装'.$plugin_data['name'])->first();
        if ($task_check) {
            $data['code'] = 1;
            $data['msg'] = '此插件已存在安装记录，请先删除！';
            return response()->json($data);
        }
        // 判断插件是否已经安装
        $installed = Task::query()->where('slug', $slug)->first();
        if ($installed) {
            $data['code'] = 1;
            $data['msg'] = '请不要重复安装！';
            return response()->json($data);
        }

        // 入库等待安装
        $task = new Task();
        $task->name = '安装'.$plugin_data['name'];
        $task->shell = $plugin_data['install'];
        $task->status = 'waiting';
        $task->log = '/tmp/'.$plugin_data['slug'].'.log';
        $task->save();
        // 塞入队列
        ProcessShell::dispatch($task->id)->delay(1);
        $res['code'] = 0;
        $res['msg'] = 'success';
        $res['data'] = '任务添加成功';

        return response()->json($res);
    }

    /**
     * 卸载插件
     * @param  Request  $request
     * @return JsonResponse
     */
    public function uninstall(Request $request): JsonResponse
    {
        // 消毒
        try {
            $credentials = $this->validate($request, [
                'slug' => 'required|max:255',
            ]);
            $slug = $credentials['slug'];
        } catch (Exception $e) {
            return response()->json(['code' => 1, 'msg' => $e->getMessage()]);
        }
        // 通过slug获取插件信息
        $pluginList = $this->pluginList();

        // 循环插件列表查找对应的插件信息
        $plugin_check = false;
        $plugin_data = [];
        foreach ($pluginList as $v) {
            if ($v['slug'] == $slug) {
                $plugin_data = $v;
                $plugin_check = true;
            }
        }
        // 判断插件是否存在
        if (!$plugin_check) {
            return response()->json(['code' => 1, 'msg' => '插件不存在']);
        }

        // 判断有无任务记录
        $task_check = Task::query()->where('name', '卸载'.$plugin_data['name'])->first();
        if ($task_check) {
            $data['code'] = 1;
            $data['msg'] = '此插件已存在卸载记录，请先删除！';
            return response()->json($data);
        }

        // 判断是否是操作openresty
        if ($slug == 'openresty') {
            $data['code'] = 1;
            $data['msg'] = '请不要花样作死！';
            return response()->json($data);
        }

        // 入库等待卸载
        $task = new Task();
        $task->name = '卸载'.$plugin_data['name'];
        $task->shell = $plugin_data['uninstall'];
        $task->status = 'waiting';
        $task->log = '/tmp/'.$plugin_data['slug'].'.log';
        $task->save();
        // 塞入队列
        ProcessShell::dispatch($task->id)->delay(1);
        $res['code'] = 0;
        $res['msg'] = 'success';
        $res['data'] = '任务添加成功';

        return response()->json($res);
    }

    /**
     * 读取插件列表
     */
    public function pluginList($cache = true)
    {
        // 判断刷新缓存
        if (!$cache) {
            Cache::forget('pluginList');
        }
        if (!Cache::has('pluginList')) {
            return Cache::remember('pluginList', 3600, function () {
                $response = Http::get('https://api.panel.haozi.xyz/api/plugin/list');
                // 判断请求是否成功，如果不成功则抛出异常
                if ($response->failed()) {
                    throw new Exception('获取插件列表失败，请求错误');
                }
                // 判断返回的JSON数据中code是否为0，如果不为0则抛出异常
                if (!$response->json('code') == 0) {
                    throw new Exception('获取插件列表失败，服务器未返回正确的状态码');
                }

                return $response->json('data');
            });
        } else {
            // 从缓存中获取access_token
            return Cache::get('pluginList');
        }
    }

    /**
     * 设置插件首页显示
     */
    public function setShowHome(Request $request): JsonResponse
    {
        // 消毒
        try {
            $credentials = $this->validate($request, [
                'slug' => 'required|max:255',
                'show' => 'required|boolean',
            ]);
            $slug = $credentials['slug'];
            $show = $credentials['show'];
        } catch (Exception $e) {
            return response()->json(['code' => 1, 'msg' => $e->getMessage()]);
        }

        Plugin::query()->where('slug', $slug)->update(['show' => $show]);
        $res['code'] = 0;
        $res['msg'] = 'success';
        $res['data'] = '设置成功';
        return response()->json($res);
    }
}
