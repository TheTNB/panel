<?php

namespace App\Providers;

use Exception;
use Illuminate\Contracts\Container\BindingResolutionException;
use Illuminate\Contracts\Filesystem\FileNotFoundException;
use App\Services\Plugin;
use Illuminate\Support\ServiceProvider;

class PluginServiceProvider extends ServiceProvider
{
    /**
     * Bootstrap any application services.
     *
     * @param  Plugin  $plugins
     * @return void
     * @throws BindingResolutionException
     * @throws FileNotFoundException
     */
    public function boot(Plugin $plugins): void
    {

        $loader = $this->app->make('translation.loader');
        // Make view instead of view.finder since the finder is defined as not a singleton
        $finder = $this->app->make('view');

        foreach ($plugins->getPlugins() as $plugin) {

            // 加载视图路径
            $finder->addNamespace($plugin['name'], $plugin['path']."/views");

            // 加载语言包
            $loader->addNamespace($plugin['name'], $plugin['path']."/lang");
        }

        // 加载插件Composer装载文件
        try {
            foreach ($plugins->getComposerLoaders() as $autoloader) {
                require $autoloader;
            }
        } catch (Exception $e) {
            throw new BindingResolutionException($e->getMessage());
        }

        // 加载插件主文件
        try {
            foreach ($plugins->getPluginMainFiles() as $file) {
                require_once $file;
            }
        } catch (Exception $e) {
            throw new BindingResolutionException($e->getMessage());
        }
    }

    /**
     * 注册插件服务
     *
     * @return void
     */
    public function register(): void
    {
        $this->app->singleton('plugins', Plugin::class);
    }
}