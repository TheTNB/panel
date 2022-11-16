<?php

use App\Http\Controllers\Api\PluginsController;
use App\Http\Controllers\Api\TasksController;
use App\Http\Controllers\Api\UsersController;
use App\Http\Controllers\Api\WebsitesController;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Route;

use App\Http\Controllers\Api\InfosController;

/*
|--------------------------------------------------------------------------
| API Routes
|--------------------------------------------------------------------------
|
| Here is where you can register API routes for your application. These
| routes are loaded by the RouteServiceProvider within a group which
| is assigned the "api" middleware group. Enjoy building your API!
|
*/

Route::middleware('auth:sanctum')->get('/user', function (Request $request) {
    return $request->user();
});
Route::middleware('auth:sanctum')->get('/fm', function () {
    return view('fm');
});
Route::prefix('panel')->group(function () {
    Route::prefix('user')->group(function () {
        // 登录
        Route::post('login', [UsersController::class, 'login']);
        // 获取用户信息
        Route::middleware('auth:sanctum')->get('getInfo', [UsersController::class, 'getInfo']);
    });
    // 顶栏任务中心列表
    Route::middleware('auth:sanctum')->prefix('task')->group(function () {
        // 当前是否有任务在运行
        Route::get('getStatus', [TasksController::class, 'getStatus']);
        // 任务中心进行中的那个任务
        Route::get('getListRunning', [TasksController::class, 'getListRunning']);
        // 任务中心等待中的一堆任务
        Route::get('getListWaiting', [TasksController::class, 'getListWaiting']);
        // 任务中心已完成的一堆任务
        Route::get('getListFinished', [TasksController::class, 'getListFinished']);
        // 任务中心删除任务
        Route::post('deleteTask', [TasksController::class, 'deleteTask']);
        // 任务中心任务的实时运行日志
        Route::get('getTaskLog', [TasksController::class, 'getTaskLog']);
    });

    Route::middleware('auth:sanctum')->prefix('info')->group(function () {
        // 菜单
        Route::get('getMenu', [InfosController::class, 'getMenu']);
        // 系统资源统计
        Route::get('getNowMonitor', [InfosController::class, 'getNowMonitor']);
        // 系统信息
        Route::get('getSystemInfo', [InfosController::class, 'getSystemInfo']);
        // 首页应用
        Route::get('getHomePlugins', [InfosController::class, 'getHomePlugins']);
        // 已安装PHP和DB版本
        Route::get('getInstalledDbAndPhp', [InfosController::class, 'getInstalledDbAndPhp']);

    });
    Route::middleware('auth:sanctum')->prefix('settings')->group(function () {
        // 获取面板设置
        Route::get('get', [TasksController::class, 'get']);
        // 保存面板设置
        Route::post('save', [TasksController::class, 'save']);
    });
    // 网站
    Route::middleware('auth:sanctum')->prefix('website')->group(function () {
        // 获取网站列表
        Route::get('getList', [WebsitesController::class, 'getList']);
        Route::post('add', [WebsitesController::class, 'add']);
    });
    // 插件
    Route::middleware('auth:sanctum')->prefix('plugin')->group(function () {
        // 获取插件列表
        Route::get('getList', [PluginsController::class, 'getList']);
        Route::get('install', [PluginsController::class, 'install']);
        Route::get('uninstall', [PluginsController::class, 'uninstall']);
    });
});

