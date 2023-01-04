<?php

namespace App\Console\Commands;

use App\Models\Plugin;
use App\Models\Setting;
use App\Models\Task;
use App\Models\User;
use App\Models\Website;
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
        // 检测是否以root用户运行
        if (trim(shell_exec('whoami')) != 'root') {
            $this->error('耗子Linux面板：请以root用户运行');
            return 1;
        }
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
            case 'backup':
                $this->backup();
                break;
            case 'writeSite':
                $this->writeSite();
                break;
            case 'deleteSite':
                $this->deleteSite();
                break;
            case 'writeSetting':
                $this->writeSetting();
                break;
            case 'deleteSetting':
                $this->deleteSetting();
                break;
            default:
                $this->info('耗子Linux面板命令行工具');
                $this->info('请使用以下命令:');
                $this->info('panel update 更新/修复面板到最新版本');
                $this->info('panel getInfo 重新初始化面板账号信息');
                $this->info('panel getPort 获取面板访问端口');
                $this->info('panel cleanRunningTask 强制清理面板正在运行的任务');
                $this->info('panel backup {website/mysql/postgresql} {name} {path} 备份网站/MySQL数据库/PostgreSQL数据库到指定目录');
                $this->warn('以下命令请在开发者指导下使用：');
                $this->info('panel init 初始化面板');
                $this->info('panel writePluginInstall {slug} 写入插件安装状态');
                $this->info('panel writePluginUnInstall {slug} 移除插件安装状态');
                $this->info('panel writeMysqlPassword {password} 写入MySQL root密码');
                $this->info('panel writeSite {name} {status} {path} {php} {ssl} 写入网站数据到面板');
                $this->info('panel deleteSite {name} 删除面板网站数据');
                $this->info('panel writeSetting {name} {value} 写入/更新面板设置数据');
                $this->info('panel deleteSetting {name} 删除面板设置数据');
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
        /**
         * 检查当前是否有任务正在运行
         */
        if (Task::query()->where('status', '!=', 'finished')->count()) {
            $this->error('当前有任务正在运行，请稍后再试');
            $this->info('如需强制更新，请先执行：panel cleanRunningTask');
            return;
        }
        $this->info('正在下载面板...');
        shell_exec('wget -O /tmp/panel.zip https://api.panel.haozi.xyz/api/version/latest');
        $this->info('正在备份数据库...');
        shell_exec('\cp /www/panel/database/database.sqlite /tmp/database.sqlite');
        // 检查下载是否成功
        if (!file_exists('/tmp/panel.zip') || filesize('/tmp/panel.zip') < 4096) {
            $this->error('检测到面板新版本下载失败，已终止更新，请加QQ群：12370907 反馈处理');
            return;
        }
        $this->info('正在备份插件...');
        shell_exec('rm -rf /tmp/plugins');
        shell_exec('mkdir /tmp/plugins');
        shell_exec('\cp -r /www/panel/plugins/* /tmp/plugins');
        // 检查备份是否成功
        if (!file_exists('/tmp/database.sqlite') || !is_dir('/tmp/plugins/Openresty')) {
            $this->error('检测到面板旧配置备份失败，已终止更新，请加QQ群：12370907 反馈处理');
            return;
        }
        $this->info('正在删除旧版本...');
        shell_exec('rm -rf /www/panel');
        shell_exec('mkdir /www/panel');
        $this->info('正在解压新版本...');
        shell_exec('unzip -o /tmp/panel.zip -d /www/panel');
        // 检查解压是否成功
        if (!file_exists('/www/panel/artisan')) {
            $this->error('检测到面板新版本解压失败，请加QQ群：12370907 反馈处理');
            return;
        }
        $this->info('正在恢复数据库...');
        shell_exec('\cp /tmp/database.sqlite /www/panel/database/database.sqlite');
        // 检查恢复是否成功
        if (!file_exists('/www/panel/database/database.sqlite')) {
            $this->error('检测到面板数据库恢复失败，请加QQ群：12370907 反馈处理');
            return;
        }
        $this->info('正在恢复插件...');
        shell_exec('\cp -r /tmp/plugins/* /www/panel/plugins');
        $this->info('正在更新面板数据库...');
        shell_exec('cd /www/panel && php-panel artisan migrate');
        $this->info('正在设置面板权限...');
        shell_exec('chown -R root:root /www/panel');
        shell_exec('chmod -R 600 /www/panel');
        shell_exec('chmod 755 /www/panel');
        shell_exec('chmod -R 755 /www/panel/public');
        shell_exec('chown -R root:root /www/server/cron');
        shell_exec('chmod -R 700 /www/server/cron');
        shell_exec('chmod -R 600 /www/server/cron/logs');
        shell_exec('chown -R root:root /www/server/vhost');
        shell_exec('chmod -R 644 /www/server/vhost');
        $this->info('正在重载面板服务...');
        $reloadCheck = shell_exec('systemctl reload panel.service 2>&1');
        // 检查重启是否成功
        if (!empty($reloadCheck)) {
            $this->error('检测到面板服务重载失败，请加QQ群：12370907 反馈处理');
        } else {
            $this->info('正在清理临时文件...');
            shell_exec('rm -rf /tmp/panel.zip');
            shell_exec('rm -rf /tmp/database.sqlite');
            shell_exec('rm -rf /tmp/plugins');
            $this->info('面板更新成功');
        }
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

        if (isset($matches[1])) {
            $this->info('面板访问端口为：'.$matches[1]);
        } else {
            $this->error('获取面板端口失败，请检查nginx主配置文件');
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
     * @return void
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
        Setting::query()->where('name', 'mysql_root_password')->updateOrCreate([
            'name' => 'mysql_root_password',
        ], [
            'value' => $password,
        ]);
        $this->info('成功');
    }

    /**
     * 清理所有运行中和等待中的任务
     * @return void
     */
    private function cleanRunningTask(): void
    {
        // 更新任务状态
        Task::query()->update(['status' => 'finished']);
        // 将所有队列任务清空
        shell_exec('php-panel /www/panel/artisan queue:clear');
        $this->info('成功');
    }

    /**
     * 备份网站/MySQL数据库/PostgreSQL数据库到指定目录
     * @return void
     */
    private function backup(): void
    {
        $type = $this->argument('a1');
        $name = $this->argument('a2');
        $path = $this->argument('a3');

        // 判空
        if (empty($type) || empty($name) || empty($path)) {
            $this->error('参数错误');
            return;
        }

        // 判断目录是否存在
        if (!is_dir($path)) {
            $this->error('目录不存在');
            return;
        }

        // 判断目录是否可写
        if (!is_writable($path)) {
            $this->error('目录不可写');
            return;
        }

        // 判断备份目录是否以/结尾，有则去掉
        if (str_ends_with($path, '/')) {
            $path = substr($path, 0, -1);
        }

        // 判断类型
        if ($type == 'website') {
            // 备份网站
            // 从数据库中获取网站目录
            $sitePath = Website::query()->where('name', $name)->value('path');
            if (empty($sitePath)) {
                $this->error('网站不存在');
                return;
            }
            $backupFile = $path.'/'.$name.'_'.date('YmdHis').'.zip';
            shell_exec('zip -r '.$backupFile.' '.escapeshellarg($sitePath).' 2>&1');
            $this->info('成功');
        } elseif ($type == 'mysql') {
            // 备份MySQL数据库
            $password = Setting::query()->where('name', 'mysql_root_password')->value('value');
            $backupFile = $path.'/'.$name.'_'.date('YmdHis').'.sql';
            // 判断数据库是否存在
            $name = escapeshellarg($name);
            $check = shell_exec("mysql -u root -p".$password." -e 'use ".$name."' 2>&1");
            if (str_contains($check, 'ERROR')) {
                $this->error('数据库不存在');
                return;
            }
            shell_exec("mysqldump -u root -p".$password." ".$name." > ".$backupFile." 2>&1");
            // zip压缩
            shell_exec('zip -r '.$backupFile.'.zip '.escapeshellarg($backupFile).' 2>&1');
            // 删除sql文件
            unlink($backupFile);
            $this->info('成功');
        } elseif ($type == 'postgresql') {
            // 备份PostgreSQL数据库
            $backupFile = $path.'/'.$name.'_'.date('YmdHis').'.sql';
            // 判断数据库是否存在
            $check = shell_exec('su - postgres -c "psql -l" 2>&1');
            if (!str_contains($check, $name)) {
                $this->error('数据库不存在');
                return;
            }
            $name = escapeshellarg($name);
            shell_exec('su - postgres -c "pg_dump '.$name.'" > '.$backupFile.' 2>&1');
            // zip压缩
            shell_exec('zip -r '.$backupFile.'.zip '.escapeshellarg($backupFile).' 2>&1');
            // 删除sql文件
            unlink($backupFile);
            $this->info('成功');
        } else {
            $this->error('参数错误');
        }
    }

    /**
     * 写入网站数据到面板
     * @return void
     */
    private function writeSite(): void
    {
        $name = $this->argument('a1');
        $status = $this->argument('a2');
        $path = $this->argument('a3');
        $php = $this->argument('a4');
        $ssl = $this->argument('a5');

        // 判空
        if (empty($name) || empty($status) || empty($path) || empty($php) || empty($ssl)) {
            $this->error('参数错误');
            return;
        }

        // 判断网站是否存在
        if (Website::query()->where('name', $name)->exists()) {
            $this->error('网站已存在');
            return;
        }

        // 判断目录是否存在
        if (!is_dir($path)) {
            $this->error('目录不存在');
            return;
        }

        // 写入网站
        Website::query()->create([
            'name' => $name,
            'status' => $status,
            'path' => $path,
            'php' => $php,
            'ssl' => $ssl,
        ]);
    }

    /**
     * 删除面板网站数据
     * @return void
     */
    private function deleteSite(): void
    {
        $name = $this->argument('a1');

        // 判空
        if (empty($name)) {
            $this->error('参数错误');
            return;
        }

        // 判断网站是否存在
        if (!Website::query()->where('name', $name)->exists()) {
            $this->error('网站不存在');
            return;
        }

        // 删除网站
        Website::query()->where('name', $name)->delete();
    }

    /**
     * 写入/更新面板设置数据
     * @return void
     */
    private function writeSetting(): void
    {
        $name = $this->argument('a1');
        $value = $this->argument('a2');

        // 判空
        if (empty($name) || empty($value)) {
            $this->error('参数错误');
            return;
        }
        Setting::query()->updateOrCreate(['name' => $name], ['value' => $value]);
    }

    /**
     * 删除面板设置数据
     * @return void
     */
    private function deleteSetting(): void
    {
        $name = $this->argument('a1');

        // 判空
        if (empty($name)) {
            $this->error('参数错误');
            return;
        }

        // 判断设置是否存在
        if (!Setting::query()->where('name', $name)->exists()) {
            $this->error('设置不存在');
            return;
        }

        // 删除设置
        Setting::query()->where('name', $name)->delete();
    }
}
