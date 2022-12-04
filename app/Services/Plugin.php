<?php

namespace App\Services;

use Exception;
use Illuminate\Contracts\Filesystem\FileNotFoundException;
use Illuminate\Support\Arr;
use Illuminate\Support\Collection;
use Illuminate\Filesystem\Filesystem;
use Illuminate\Contracts\Events\Dispatcher;
use Illuminate\Contracts\Foundation\Application;

class Plugin
{
    /**
     * @var Application
     */
    protected Application $app;

    /**
     * @var Dispatcher
     */
    protected Dispatcher $dispatcher;

    /**
     * @var Filesystem
     */
    protected Filesystem $filesystem;

    /**
     * @var Collection|null
     */
    protected ?Collection $plugins;

    public function __construct(
        Application $app,
        Dispatcher $dispatcher,
        Filesystem $filesystem,
        Collection $plugins = new Collection()
    ) {
        $this->app = $app;
        $this->dispatcher = $dispatcher;
        $this->filesystem = $filesystem;
        $this->plugins = $plugins;
    }

    /**
     * 读取所有插件
     *
     * @return Collection|null
     * @throws FileNotFoundException
     * @throws Exception
     */
    public function getPlugins(): ?Collection
    {
        if ($this->plugins->isEmpty()) {
            $plugins = new Collection();

            $installed = [];

            try {
                $resource = opendir($this->getPluginsDir());
            } catch (Exception $e) {
                // 输出错误信息
                throw new Exception($e->getMessage());
            }

            // 读取插件目录
            while ($filename = @readdir($resource)) {
                if ($filename == '.' || $filename == '..') {
                    continue;
                }

                $path = $this->getPluginsDir().DIRECTORY_SEPARATOR.$filename;

                if (is_dir($path)) {
                    $pluginJsonPath = $path.DIRECTORY_SEPARATOR.'plugin.json';

                    if (file_exists($pluginJsonPath)) {
                        // 读取插件配置文件
                        $installed[$filename] = json_decode($this->filesystem->get($pluginJsonPath), true);
                    }
                }

            }
            closedir($resource);

            foreach ($installed as $dirname => $package) {

                // 初始化插件信息
                $plugin = [];
                $plugin['name'] = (Arr::get($package, 'name'));
                $plugin['author'] = (Arr::get($package, 'author'));
                $plugin['describe'] = (Arr::get($package, 'describe'));
                $plugin['slug'] = (Arr::get($package, 'slug'));
                $plugin['version'] = (Arr::get($package, 'version'));
                $plugin['path'] = $this->getPluginsDir().DIRECTORY_SEPARATOR.$dirname;

                if ($plugins->has($plugin['slug'])) {
                    continue;
                }

                $plugins->put($plugin['slug'], $plugin);
            }

            define('PLUGINS', $plugins->toArray());
            $this->plugins = $plugins;
        }

        return $this->plugins;
    }

    /**
     * 获取全部插件的主文件
     *
     * @return Collection
     * @throws Exception
     */
    public function getPluginMainFiles(): Collection
    {
        $pluginMainFiles = new Collection;

        try {
            foreach ($this->getPlugins() as $plugin) {
                if ($this->filesystem->exists($file = $plugin['path'].'/plugin.php')) {
                    $pluginMainFiles->push($file);
                }
            }
        } catch (Exception $e) {
            // 输出错误信息
            throw new Exception($e->getMessage());
        }

        return $pluginMainFiles;
    }

    /**
     * 获取全部插件的Composer装载文件
     *
     * @return Collection
     * @throws Exception
     */
    public function getComposerLoaders(): Collection
    {
        $composerLoaders = new Collection;

        foreach ($this->getPlugins() as $plugin) {
            if ($this->filesystem->exists($file = $plugin['path'].'/vendor/autoload.php')) {
                $composerLoaders->push($file);
            }
        }

        return $composerLoaders;
    }

    /**
     * 插件目录
     *
     * @return string
     */
    public function getPluginsDir(): string
    {
        return config('plugins.directory') ?: base_path('plugins');
    }
}
