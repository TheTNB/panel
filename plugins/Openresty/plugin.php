<?php
/**
 * Name: OpenResty插件
 * Author: 耗子
 * Date: 2022-11-21
 */

use Illuminate\Support\Facades\Route;
use Plugins\Openresty\Controllers\OpenrestyController;

// 视图
app('router')->group([
    'prefix' => 'panel/views/plugin/openresty',
    //'middleware' => ['auth:sanctum'],
], function () {
    Route::view('/', 'openresty::index');
});
// 控制器
app('router')->group([
    'prefix' => 'api/plugin/openresty',
    'middleware' => ['auth:sanctum'],
], function () {
    Route::get('status', [OpenrestyController::class, 'status']);
    Route::get('load', [OpenrestyController::class, 'load']);
    Route::get('errorLog', [OpenrestyController::class, 'errorLog']);
    Route::get('config', [OpenrestyController::class, 'getConfig']);
    Route::post('config', [OpenrestyController::class, 'saveConfig']);
    Route::get('cleanErrorLog', [OpenrestyController::class, 'cleanErrorLog']);
    Route::get('restart', [OpenrestyController::class, 'restart']);
    Route::get('reload', [OpenrestyController::class, 'reload']);
});

