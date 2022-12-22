<?php

use App\Http\Controllers\Api\CronsController;
use App\Http\Controllers\Api\FilesController;
use App\Http\Controllers\Api\MonitorsController;
use App\Http\Controllers\Api\PluginsController;
use App\Http\Controllers\Api\SafesController;
use App\Http\Controllers\Api\SettingsController;
use App\Http\Controllers\Api\TasksController;
use App\Http\Controllers\Api\UsersController;
use App\Http\Controllers\Api\WebsitesController;
use App\Http\Controllers\Api\InfosController;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Route;

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
    // 网站管理
    Route::middleware('auth:sanctum')->prefix('website')->group(function () {
        // 获取默认设置
        Route::get('getDefaultSettings', [WebsitesController::class, 'getDefaultSettings']);
        // 保存默认设置
        Route::post('saveDefaultSettings', [WebsitesController::class, 'saveDefaultSettings']);
        // 获取网站列表
        Route::get('getList', [WebsitesController::class, 'getList']);
        // 添加网站
        Route::post('add', [WebsitesController::class, 'add']);
        // 删除网站
        Route::post('delete', [WebsitesController::class, 'delete']);
        // 获取网站设置
        Route::get('getSiteSettings', [WebsitesController::class, 'getSiteSettings']);
        // 保存网站设置
        Route::post('saveSiteSettings', [WebsitesController::class, 'saveSiteSettings']);
        // 清空网站日志
        Route::post('clearSiteLog', [WebsitesController::class, 'clearSiteLog']);
        // 更新网站备注
        Route::post('updateSiteNote', [WebsitesController::class, 'updateSiteNote']);
        // 设置网站状态
        Route::post('setSiteStatus', [WebsitesController::class, 'setSiteStatus']);
        // 获取备份列表
        Route::get('getBackupList', [WebsitesController::class, 'getBackupList']);
        // 创建备份
        Route::post('createBackup', [WebsitesController::class, 'createBackup']);
        // 上传备份
        Route::post('uploadBackup', [WebsitesController::class, 'uploadBackup']);
        // 恢复备份
        Route::post('restoreBackup', [WebsitesController::class, 'restoreBackup']);
        // 删除备份
        Route::post('deleteBackup', [WebsitesController::class, 'deleteBackup']);
        // 重置网站配置
        Route::post('resetSiteConfig', [WebsitesController::class, 'resetSiteConfig']);
        // 签发SSL证书
        Route::post('issueSsl', [WebsitesController::class, 'issueSsl']);
    });
    // 资源监控
    Route::middleware('auth:sanctum')->prefix('monitor')->group(function () {
        // 获取监控数据
        Route::get('getMonitorData', [MonitorsController::class, 'getMonitorData']);
        // 获取监控开关和保存天数
        Route::get('getMonitorSwitchAndDays', [MonitorsController::class, 'getMonitorSwitchAndDays']);
        // 设置监控开关
        Route::post('setMonitorSwitch', [MonitorsController::class, 'setMonitorSwitch']);
        // 设置保存天数
        Route::post('setMonitorSaveDays', [MonitorsController::class, 'setMonitorSaveDays']);
        // 清空监控数据
        Route::post('clearMonitorData', [MonitorsController::class, 'clearMonitorData']);
    });
    // 系统安全
    Route::middleware('auth:sanctum')->prefix('safe')->group(function () {
        // 获取防火墙状态
        Route::get('getFirewallStatus', [SafesController::class, 'getFirewallStatus']);
        // 设置防火墙状态
        Route::post('setFirewallStatus', [SafesController::class, 'setFirewallStatus']);
        // 获取SSH状态
        Route::get('getSshStatus', [SafesController::class, 'getSshStatus']);
        // 设置SSH状态
        Route::post('setSshStatus', [SafesController::class, 'setSshStatus']);
        // 获取SSH端口
        Route::get('getSshPort', [SafesController::class, 'getSshPort']);
        // 设置SSH端口
        Route::post('setSshPort', [SafesController::class, 'setSshPort']);
        // 获取ping状态
        Route::get('getPingStatus', [SafesController::class, 'getPingStatus']);
        // 设置ping状态
        Route::post('setPingStatus', [SafesController::class, 'setPingStatus']);
        // 获取防火墙规则
        Route::get('getFirewallRules', [SafesController::class, 'getFirewallRules']);
        // 添加防火墙规则
        Route::post('addFirewallRule', [SafesController::class, 'addFirewallRule']);
        // 删除防火墙规则
        Route::post('deleteFirewallRule', [SafesController::class, 'deleteFirewallRule']);
    });
    // 文件管理
    Route::middleware('auth:sanctum')->prefix('file')->group(function () {
        // 获取文件（夹）列表
        Route::get('getList', [FilesController::class, 'getList']);
    });
    // 计划任务
    Route::middleware('auth:sanctum')->prefix('cron')->group(function () {
        // 获取计划任务列表
        Route::get('getList', [CronsController::class, 'getList']);
        // 添加计划任务
        Route::post('add', [CronsController::class, 'add']);
        // 编辑计划任务
        Route::post('edit', [CronsController::class, 'edit']);
        // 删除计划任务
        Route::post('delete', [CronsController::class, 'delete']);
        // 设置计划任务状态
        Route::post('setStatus', [CronsController::class, 'setStatus']);
        // 获取计划任务日志
        Route::get('getLog', [CronsController::class, 'getLog']);
    });
    // 插件中心
    Route::middleware('auth:sanctum')->prefix('plugin')->group(function () {
        // 获取插件列表
        Route::get('getList', [PluginsController::class, 'getList']);
        // 安装插件
        Route::post('install', [PluginsController::class, 'install']);
        // 卸载插件
        Route::post('uninstall', [PluginsController::class, 'uninstall']);
        // 更新插件
        Route::post('update', [PluginsController::class, 'update']);
        // 设置插件首页显示
        Route::post('setShowHome', [PluginsController::class, 'setShowHome']);
    });
    // 面板设置
    Route::middleware('auth:sanctum')->prefix('setting')->group(function () {
        // 获取设置
        Route::get('get', [SettingsController::class, 'get']);
        // 保存设置
        Route::post('save', [SettingsController::class, 'save']);
    });
});
