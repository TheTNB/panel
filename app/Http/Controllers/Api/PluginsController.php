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
use Illuminate\Validation\ValidationException;


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
        $data['data'] = $this->pluginList();
        // 获取首页显示状态
        $shows = Plugin::query()->pluck('show', 'slug');
        foreach ($data['data'] as $k => $v) {
            // 如果本地已安装，则显示本地名称
            $data['data'][$k]['name'] = PLUGINS[$v['slug']]['name'] ?? $data['data'][$k]['name'] ?? '无名';
            // 作者名称
            $data['data'][$k]['author'] = PLUGINS[$v['slug']]['author'] ?? $data['data'][$k]['author'] ?? '耗子';
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
        // 插入本地插件数据
        $slugs = array_column($data['data'], 'slug');
        foreach (PLUGINS as $v) {
            // 如果本地已安装，则不显示
            if (in_array($v['slug'], $slugs)) {
                continue;
            }
            // 插入插件数据
            $data['data'][] = [
                'name' => $v['name'],
                'author' => $v['author'],
                'slug' => $v['slug'],
                'version' => $v['version'],
                'install_version' => $v['version'],
                'describe' => $v['describe'],
                'show' => $shows[$v['slug']] ?? 0,
                'control' => [
                    'installed' => true,
                    'allow_uninstall' => false,
                ],
            ];

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
        } catch (ValidationException $e) {
            return response()->json(['code' => 1, 'msg' => $e->getMessage()]);
        }
        // 通过slug获取插件信息
        $pluginList = $this->pluginList();
        // 循环插件列表查找对应的插件信息
        $pluginCheck = false;
        $pluginData = [];
        foreach ($pluginList as $v) {
            if ($v['slug'] == $slug) {
                $pluginData = $v;
                $pluginCheck = true;
            }
        }

        // 判断插件是否存在
        if (!$pluginCheck) {
            return response()->json(['code' => 1, 'msg' => '插件不存在']);
        }

        // 判断插件是否已经安装
        $installed = isset(PLUGINS[$slug]);
        if ($installed) {
            $data['code'] = 1;
            $data['msg'] = '请不要重复安装！';
            return response()->json($data);
        }

        // 入库等待安装
        $task = new Task();
        $task->name = '安装 '.$pluginData['name'];
        $task->shell = $pluginData['install'];
        $task->status = 'waiting';
        $task->log = '/tmp/'.$pluginData['slug'].'.log';
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
        } catch (ValidationException $e) {
            return response()->json(['code' => 1, 'msg' => $e->getMessage()]);
        }
        // 通过slug获取插件信息
        $pluginList = $this->pluginList();

        // 循环插件列表查找对应的插件信息
        $pluginCheck = false;
        $pluginData = [];
        foreach ($pluginList as $v) {
            if ($v['slug'] == $slug) {
                $pluginData = $v;
                $pluginCheck = true;
            }
        }

        // 判断插件是否存在
        if (!$pluginCheck) {
            return response()->json(['code' => 1, 'msg' => '插件不存在']);
        }

        // 判断是否是操作openresty
        if ($slug == 'openresty') {
            $data['code'] = 1;
            $data['msg'] = '请不要花样作死！';
            return response()->json($data);
        }

        // 入库等待卸载
        $task = new Task();
        $task->name = '卸载 '.$pluginData['name'];
        $task->shell = $pluginData['uninstall'];
        $task->status = 'waiting';
        $task->log = '/tmp/'.$pluginData['slug'].'.log';
        $task->save();
        // 塞入队列
        ProcessShell::dispatch($task->id)->delay(1);
        $res['code'] = 0;
        $res['msg'] = 'success';
        $res['data'] = '任务添加成功';

        return response()->json($res);
    }

    /**
     * 更新插件
     */
    public function update(Request $request): JsonResponse
    {
        // 消毒
        try {
            $credentials = $this->validate($request, [
                'slug' => 'required|max:255',
            ]);
            $slug = $credentials['slug'];
        } catch (ValidationException $e) {
            return response()->json(['code' => 1, 'msg' => $e->getMessage()]);
        }
        // 通过slug获取插件信息
        $pluginList = $this->pluginList();

        // 循环插件列表查找对应的插件信息
        $pluginCheck = false;
        $pluginData = [];
        foreach ($pluginList as $v) {
            if ($v['slug'] == $slug) {
                $pluginData = $v;
                $pluginCheck = true;
                break;
            }
        }

        // 判断插件是否存在
        if (!$pluginCheck) {
            return response()->json(['code' => 1, 'msg' => '插件不存在']);
        }

        // 判断插件是否已经安装
        $installed = isset(PLUGINS[$slug]);
        if (!$installed) {
            $data['code'] = 1;
            $data['msg'] = '插件未安装！';
            return response()->json($data);
        }

        // 入库等待更新
        $task = new Task();
        $task->name = '更新 '.$pluginData['name'];
        $task->shell = $pluginData['update'];
        $task->status = 'waiting';
        $task->log = '/tmp/'.$pluginData['slug'].'.log';
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
        } catch (ValidationException $e) {
            return response()->json(['code' => 1, 'msg' => $e->getMessage()]);
        }

        // 判断插件是否已经安装
        $installed = isset(PLUGINS[$slug]);
        if (!$installed) {
            $data['code'] = 1;
            $data['msg'] = '插件未安装！';
            return response()->json($data);
        }

        Plugin::query()->where('slug', $slug)->updateOrInsert(
            ['slug' => $slug],
            ['show' => $show]
        );
        $res['code'] = 0;
        $res['msg'] = 'success';
        $res['data'] = '设置成功';
        return response()->json($res);
    }
}
