<?php

namespace Alexusmai\LaravelFileManager;

use Alexusmai\LaravelFileManager\Middleware\FileManagerACL;
use Alexusmai\LaravelFileManager\Services\ACLService\ACLRepository;
use Alexusmai\LaravelFileManager\Services\ConfigService\ConfigRepository;
use Illuminate\Support\ServiceProvider;

class FileManagerServiceProvider extends ServiceProvider
{
    /**
     * Bootstrap the application services.
     *
     * @return void
     */
    public function boot()
    {
        // routes
        $this->loadRoutesFrom(__DIR__.'/routes.php');

        // views
        $this->loadViewsFrom(__DIR__.'/../resources/views', 'file-manager');

        // publish config
        $this->publishes([
            __DIR__
            .'/../config/file-manager.php' => config_path('file-manager.php'),
        ], 'fm-config');

        // publish views
        $this->publishes([
            __DIR__
            .'/../resources/views' => resource_path('views/vendor/file-manager'),
        ], 'fm-views');

        // publish js and css files - vue-file-manager module
        $this->publishes([
            __DIR__
            .'/../resources/assets' => public_path('vendor/file-manager'),
        ], 'fm-assets');

        // publish migrations
        $this->publishes([
            __DIR__
            .'/../migrations' => database_path('migrations'),
        ], 'fm-migrations');
    }

    /**
     * Register the application services.
     *
     * @return void
     */
    public function register()
    {
        $this->mergeConfigFrom(
            __DIR__.'/../config/file-manager.php',
            'file-manager'
        );

        // Config Repository
        $this->app->bind(
            ConfigRepository::class,
            $this->app['config']['file-manager.configRepository']
        );

        // ACL Repository
        $this->app->bind(
            ACLRepository::class,
            $this->app->make(ConfigRepository::class)->getAclRepository()
        );

        // register ACL middleware
        $this->app['router']->aliasMiddleware('fm-acl', FileManagerACL::class);
    }
}
