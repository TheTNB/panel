<?php

namespace App\Console\Commands;

use App\Models\Plugin;
use App\Models\Setting;
use App\Models\Task;
use App\Models\User;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Hash;
use Illuminate\Support\Str;

class Panel extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'panel {action?} {a1?} {a2?} {a3?} {a4?} {a5?} {a6?} {a7?} {a8?} {a9?} {a10?}';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '耗子Linux面板命令行工具';

    /**
     * Execute the console command.
     *
     * @return int
     */
    public function handle()
    {
        $action = $this->argument('action');
        switch ($action) {
            case 'init':
                $this->init();
                break;
            case 'update':
                $this->update();
                break;
            case 'getInfo':
                $this->getInfo();
                break;
            case 'getPort':
                $this->getPort();
                break;
            case 'writePluginInstall':
                $this->writePluginInstall();
                break;
            case 'writePluginUnInstall':
                $this->writePluginUnInstall();
                break;
            case 'writeMysqlPassword':
                $this->writeMysqlPassword();
                break;
            case 'cleanRunningTask':
                $this->cleanRunningTask();
                break;
            default:
                $this->error('错误的操作');
                break;
        }
        return Command::SUCCESS;
    }

    /**
     * 初始化
     * @return void
     */
    private function init(): void
    {
        Setting::query()->updateOrCreate(['name' => 'name'], ['value' => '耗子Linux面板']);
        Setting::query()->updateOrCreate(['name' => 'monitor'], ['value' => '1']);
        Setting::query()->updateOrCreate(['name' => 'monitor_days'], ['value' => '30']);
        Setting::query()->updateOrCreate(['name' => 'mysql_root_password'], ['value' => '']);
        Setting::query()->updateOrCreate(['name' => 'postgresql_root_password'], ['value' => '']);
        User::query()->create([
            'id' => 1,
            'username' => 'admin',
            'email' => 'panel@haozi.net',
            'password' => Hash::make(Str::random()),
        ]);
    }

    /**
     * 更新面板
     * @return void
     */
    private function update(): void
    {
        $this->info('正在下载面板...');
        $this->info(shell_exec('wget -O /tmp/panel.zip https://api.panel.haozi.xyz/api/version/latest'));
        $this->info('正在备份数据库...');
        $this->info(shell_exec('\cp /www/panel/database/database.sqlite /tmp/database.sqlite'));
        $this->info('正在备份插件...');
        $this->info(shell_exec('rm -rf /tmp/plugins'));
        $this->info(shell_exec('mkdir /tmp/plugins'));
        $this->info(shell_exec('\cp -r /www/panel/plugins/* /tmp/plugins'));
        $this->info('正在删除旧版本...');
        $this->info(shell_exec('rm -rf /www/panel'));
        $this->info(shell_exec('mkdir /www/panel'));
        $this->info('正在解压新版本...');
        $this->info(shell_exec('unzip /tmp/panel.zip -d /www/panel'));
        $this->info('正在恢复数据库...');
        $this->info(shell_exec('\cp /tmp/database.sqlite /www/panel/database/database.sqlite'));
        $this->info('正在恢复插件...');
        $this->info(shell_exec('\cp -r /tmp/plugins/* /www/panel/plugins'));
        $this->info('正在清理临时文件...');
        $this->info(shell_exec('rm -rf /tmp/panel.zip'));
        $this->info(shell_exec('rm -rf /tmp/database.sqlite'));
        $this->info(shell_exec('rm -rf /tmp/plugins'));
        $this->info('正在更新面板数据库...');
        $this->info(shell_exec('cd /www/panel && php-panel artisan migrate'));
        $this->info('正在重启面板服务...');
        $this->info(shell_exec('systemctl restart panel.service'));
        $this->info('更新完成');
    }

    /**
     * 获取面板信息
     * @return void
     */
    private function getInfo(): void
    {
        $user = User::query()->where('id', 1);
        // 生成唯一信息
        $username = Str::random(6);
        $password = Str::random(12);
        // 入库
        $user->update([
            'username' => $username,
            'password' => Hash::make($password),
        ]);

        // 从nginx配置文件中获取面板端口
        $nginxConf = file_get_contents('/www/server/nginx/conf/nginx.conf');
        preg_match('/listen\s+(\d+)/', $nginxConf, $matches);

        if (!isset($matches[1])) {
            $this->info('获取面板端口失败，请检查nginx主配置文件');
        }

        $this->info('面板用户名：'.$username);
        $this->info('面板密码：'.$password);
        $this->info('访问地址：http://IP:'.$matches[1]);
    }

    /**
     * 获取端口
     * @return void
     */
    private function getPort(): void
    {
        // 从nginx配置文件中获取面板端口
        $nginxConf = file_get_contents('/www/server/nginx/conf/nginx.conf');
        preg_match('/listen\s+(\d+)/', $nginxConf, $matches);

        if (!isset($matches[1])) {
            $this->info($matches[1]);
        }
    }

    /**
     * 写入插件安装状态
     * @return void
     */
    private function writePluginInstall(): void
    {
        $pluginSlug = $this->argument('a1');

        // 判空
        if (empty($pluginSlug)) {
            $this->error('参数错误');
            return;
        }
        // 入库
        Plugin::query()->create([
            'slug' => $pluginSlug,
            'show' => 0,
        ]);
        $this->info('成功');
    }

    /**
     * 写入插件卸载状态
     * @return void
     */
    private function writePluginUnInstall(): void
    {
        $pluginSlug = $this->argument('a1');

        // 判空
        if (empty($pluginSlug)) {
            $this->error('参数错误');
            return;
        }
        if ($pluginSlug == 'openresty') {
            $this->error('耗子Linux面板：请不要花样作死！');
            return;
        }
        // 入库
        Plugin::query()->where('slug', $pluginSlug)->delete();
        $this->info('成功');
    }

    /**
     * 写入MySQL密码
     */
    private function writeMysqlPassword(): void
    {
        $password = $this->argument('a1');

        // 判空
        if (empty($password)) {
            $this->error('参数错误');
            return;
        }
        // 入库
        Setting::query()->where('name', 'mysql_root_password')->update(['value' => $password]);
        $this->info('成功');
    }

    /**
     * 清理所有运行中和等待中的任务
     */
    private function cleanRunningTask(): void
    {
        Task::query()->update(['status' => 'finished']);
        $this->info('成功');
    }
}
