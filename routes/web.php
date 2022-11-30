<?php

use Illuminate\Support\Facades\Route;

/*
|--------------------------------------------------------------------------
| Web Routes
|--------------------------------------------------------------------------
|
| Here is where you can register web routes for your application. These
| routes are loaded by the RouteServiceProvider within a group which
| contains the "web" middleware group. Now create something great!
|
*/

Route::view('/', 'index');
Route::prefix('panel/views')->group(function () {
    // 框架
    Route::prefix('ui')->group(function () {
        Route::view('layout', 'ui.layout');
        Route::view('theme', 'ui.theme');
        Route::view('404', 'ui.404');
        Route::view('error', 'ui.error');
    });


    // 主页
    Route::view('index', 'home');
    Route::view('setting', 'setting');
    // 网站
    Route::prefix('website')->group(function () {
        //全局设置
        Route::view('default_settings', 'website.default_settings');
        // 列表
        Route::view('list', 'website.list');
        // 添加
        Route::view('add', 'website.add');
        // 编辑
        Route::view('edit', 'website.edit');
    });

    // 数据库-MySQL
    Route::prefix('database')->group(function () {
        Route::view('mysql', 'database.mysql');
        Route::view('postgresql', 'database.postgresql');
    });
    // 监控
    Route::view('monitor', 'monitor');
    // 安全
    Route::view('safe', 'safe');
    // 文件
    Route::view('file', 'file');
    // 插件
    Route::view('plugin', 'plugin');

    // 其他独立页面
    // 登录
    Route::view('login', 'login');
    // 注销
    Route::view('logout', 'logout');
    // 主题设置
    Route::view('theme', 'theme');
    // 任务
    Route::view('task', 'task');
});
