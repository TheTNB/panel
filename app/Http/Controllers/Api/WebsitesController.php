<?php
/**
 * 耗子Linux面板 - 网站控制器
 * @author 耗子
 */

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Models\Website;
use App\Models\Setting;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Validation\ValidationException;
use skoerfgen\ACMECert\ACME_Exception;
use skoerfgen\ACMECert\ACMECert;

class WebsitesController extends Controller
{
    /**
     * 获取面板网站
     * @return JsonResponse
     */
    public function getList(): JsonResponse
    {
        $websiteList = Website::query()->get();
        // 判空
        if ($websiteList->isEmpty()) {
            return response()->json([
                'code' => 0,
                'msg' => '无数据',
                'data' => []
            ]);
        }
        foreach ($websiteList as $website) {
            // 如果PHP是0，将其设置为字符串的00
            if ($website->php == '0') {
                $website->php = '00';
            }
        }
        return response()->json([
            'code' => 0,
            'msg' => '获取成功',
            'data' => $websiteList
        ]);
    }

    /**
     * 添加面板网站
     * @param  Request  $request
     * @return JsonResponse
     */
    public function add(Request $request): JsonResponse
    {
        // 消毒数据
        try {
            $credentials = $this->validate($request, [
                'name' => 'required|max:255',
                'domain' => 'required',
                'path' => 'string|nullable|max:255',
                'php' => 'required',
                'note' => 'string|nullable|max:255',
                'db' => 'required|boolean',
                'db_type' => 'required_if:db,true|max:10',
                'db_name' => 'required_if:db,true|max:255',
                'db_username' => 'required_if:db,true|max:255',
                'db_password' => ['required_if:db,true', 'max:255'],
            ]);
        } catch (ValidationException $e) {
            return response()->json([
                'code' => 1,
                'msg' => '参数错误：'.$e->getMessage(),
                'errors' => $e->errors()
            ], 200);
        }

        // 对db_password单独验证
        if ($credentials['db']) {
            if (!preg_match('/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*(_|[^\w])).+$/',
                $credentials['db_password'])) {
                return response()->json([
                    'code' => 1,
                    'msg' => '数据库密码必须包含大小写字母、数字、特殊字符'
                ], 200);
            } elseif (strlen($credentials['db_password']) < 8) {
                return response()->json([
                    'code' => 1,
                    'msg' => '数据库密码长度不能小于8位'
                ], 200);
            }
        }

        // 禁止添加重复网站
        $website = Website::query()->where('name', $credentials['name'])->first();
        if ($website) {
            return response()->json([
                'code' => 1,
                'msg' => '网站已存在'
            ]);
        }
        // 禁止phpmyadmin作为名称
        if ($credentials['name'] == 'phpmyadmin') {
            return response()->json([
                'code' => 1,
                'msg' => '该名称为保留名称，请更换'
            ]);
        }
        // path为空时，设置默认值
        if (empty($credentials['path'])) {
            $credentials['path'] = '/www/wwwroot/'.$credentials['name'];
        }
        // 如果path不以/开头，则返回错误
        if (!str_starts_with($credentials['path'], '/')) {
            return response()->json([
                'code' => 1,
                'msg' => '网站路径必须以/开头'
            ]);
        }
        // ssl默认设置为0
        $credentials['ssl'] = 0;
        // 运行状态默认设置为1
        $credentials['status'] = 1;
        $domain = trim($credentials['domain']);
        // 入库
        Website::query()->create($credentials);
        // 创建网站目录
        shell_exec("mkdir -p ".$credentials['path']);
        // 创建index.html
        shell_exec("touch ".$credentials['path']."/index.html");
        // 写入到index.html
        $index_html = <<<EOF
<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="utf-8">
<title>耗子Linux面板</title>
</head>
<body>
<h1>耗子Linux面板</h1>
<p>这是耗子Linux面板的网站默认页面！</p>
<p>当您看到此页面，说明您的网站已创建成功。</p>
</body>
</html>

EOF;
        file_put_contents($credentials['path']."/index.html", $index_html);

        // 创建nginx配置
        $port_list = "";
        $domain_list = "";
        $domain_arr = explode(PHP_EOL, $domain);
        foreach ($domain_arr as $key => $value) {
            $temp = explode(":", $value);
            $domain_list .= " ".$temp[0];
            if (!isset($temp[1])) {
                if ($key == count($domain_arr) - 1) {
                    $port_list .= "    listen 80;";
                } else {
                    $port_list .= "    listen 80;".PHP_EOL;
                }
            } else {
                if ($key == count($domain_arr) - 1) {
                    $port_list .= "    listen ".$temp[1].";";
                } else {
                    $port_list .= "    listen ".$temp[1].";".PHP_EOL;
                }
            }

        }
        $nginx_config = <<<EOF
# 配置文件中的标记位请勿随意修改，改错将导致面板无法识别！
# 有自定义配置需求的，请将自定义的配置写在各标记位下方。
server
{
    # port标记位开始
$port_list
    # port标记位结束
    # server_name标记位开始
    server_name$domain_list;
    # server_name标记位结束
    # index标记位开始
    index index.php index.html;
    # index标记位结束
    # root标记位开始
    root $credentials[path];
    # root标记位结束

    # ssl标记位开始
    # ssl标记位结束

    # php标记位开始
    include enable-php-$credentials[php].conf;
    # php标记位结束

    # waf标记位开始
    waf on;
    waf_rule_path /www/server/nginx/ngx_waf/assets/rules/;
    waf_mode DYNAMIC;
    waf_cc_deny rate=1000r/m duration=60m;
    waf_cache capacity=50;
    # waf标记位结束

    # 错误页配置，可自行设置
    #error_page 404 /404.html;
    #error_page 502 /502.html;

    # 伪静态规则引入，修改后将导致面板设置的伪静态规则失效
    include /www/server/vhost/rewrite/$credentials[name].conf;

    # 面板默认禁止访问部分敏感目录，可自行修改
    location ~ ^/(\.user.ini|\.htaccess|\.git|\.svn)
    {
        return 404;
    }
    # 面板默认不记录静态资源的访问日志并开启1小时浏览器缓存，可自行修改
    location ~ .*\.(js|css)$
    {
        expires 1h;
        error_log /dev/null;
        access_log /dev/null;
    }
    access_log /www/wwwlogs/$credentials[name].log;
    error_log /www/wwwlogs/$credentials[name].log;
}
EOF;
        // 写入nginx配置
        file_put_contents('/www/server/vhost/'.$credentials['name'].'.conf', $nginx_config);
        shell_exec('echo "" > /www/server/vhost/rewrite/'.$credentials['name'].'.conf');
        shell_exec('echo "" > /www/server/vhost/ssl/'.$credentials['name'].'.pem');
        shell_exec('echo "" > /www/server/vhost/ssl/'.$credentials['name'].'.key');
        shell_exec("systemctl reload nginx");

        // 创建数据库
        if ($credentials['db']) {
            if ($credentials['db_type'] == 'mysql') {
                $password = Setting::query()->where('name', 'mysql_root_password')->value('value');
                shell_exec("mysql -u root -p".$password." -e \"CREATE DATABASE IF NOT EXISTS ".$credentials['db_name']." DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;\" 2>&1");
                shell_exec("mysql -u root -p".$password." -e \"CREATE USER '".$credentials['db_username']."'@'localhost' IDENTIFIED BY '".$credentials['db_password']."';\"");
                shell_exec("mysql -u root -p".$password." -e \"GRANT ALL PRIVILEGES ON ".$credentials['db_name'].".* TO '".$credentials['db_username']."'@'localhost';\"");
                shell_exec("mysql -u root -p".$password." -e \"flush privileges;\"");
            } elseif ($credentials['db_type'] == 'postgresql') {
                shell_exec('echo "CREATE DATABASE '.$credentials['db_name'].';"|su - postgres -c "psql"');
                shell_exec('echo "CREATE USER '.$credentials['db_username'].' WITH PASSWORD \''.$credentials['db_password'].'\';"|su - postgres -c "psql"');
                shell_exec('echo "GRANT ALL PRIVILEGES ON DATABASE '.$credentials['db_name'].' TO '.$credentials['db_username'].';"|su - postgres -c "psql"');
                // 写入用户配置
                shell_exec('echo "host    '.$credentials['db_name'].'    '.$credentials['db_username'].'    127.0.0.1/32    scram-sha-256" >> /www/server/postgresql/15/pg_hba.conf');
                // 重载
                shell_exec('systemctl reload postgresql-15');
            }
        }
        $res['code'] = 0;
        $res['msg'] = 'success';
        return response()->json($res);
    }


    /**
     * 删除面板网站
     * @param  Request  $request
     * @return JsonResponse
     */
    public function delete(Request $request): JsonResponse
    {
        $name = $request->input('name');
        // 从数据库删除
        Website::query()->where('name', $name)->delete();
        // 删除站点目录
        shell_exec("rm -rf /www/wwwroot/$name");
        // 删除nginx配置
        shell_exec("rm -rf /www/server/vhost/$name.conf");
        // 删除rewrite配置
        shell_exec("rm -rf /www/server/vhost/rewrite/$name.conf");
        // 删除ssl配置
        shell_exec("rm -rf /www/server/vhost/ssl/$name.pem");
        shell_exec("rm -rf /www/server/vhost/ssl/$name.key");

        $res['code'] = 0;
        $res['msg'] = 'success';
        return response()->json($res);
    }

    /**
     * 获取网站全局设置
     * @return JsonResponse
     */
    public function getDefaultSettings(): JsonResponse
    {
        $index = @file_get_contents('/www/server/nginx/html/index.html') ? file_get_contents('/www/server/nginx/html/index.html') : '';
        $stop = @file_get_contents('/www/server/nginx/html/stop.html') ? file_get_contents('/www/server/nginx/html/stop.html') : '';
        $res['code'] = 0;
        $res['msg'] = 'success';
        $res['data'] = [
            'index' => $index,
            'stop' => $stop,
        ];
        return response()->json($res);
    }

    /**
     * 保存网站全局设置
     * @return JsonResponse
     */
    public function saveDefaultSettings(): JsonResponse
    {
        $index = request()->input('index');
        $stop = request()->input('stop');
        file_put_contents('/www/server/nginx/html/index.html', $index);
        file_put_contents('/www/server/nginx/html/stop.html', $stop);
        $res['code'] = 0;
        $res['msg'] = 'success';
        return response()->json($res);
    }

    /**
     * 获取面板网站设置
     * @param  Request  $request
     * @return JsonResponse
     */
    public function getSiteSettings(Request $request): JsonResponse
    {
        $name = $request->input('name');
        $website = Website::query()->where('name', $name)->first();
        // 通过name读取相应的nginx配置
        $nginx_config = file_get_contents('/www/server/vhost/'.$name.'.conf');
        // 从nginx配置中port标记位提取全部端口
        $port_raw = cut('# port标记位开始', '# port标记位结束', $nginx_config);
        preg_match_all('/listen\s+(.*);/', $port_raw, $matches);
        foreach ($matches[1] as $k => $v) {
            if ($k == 0) {
                $website['port'] = $v;
            } else {
                $website['port'] .= PHP_EOL.$v;
            }
        }
        // 从nginx配置中server_name标记位提取全部域名
        $server_name_raw = cut('# server_name标记位开始', '# server_name标记位结束', $nginx_config);
        preg_match_all('/server_name\s+(.+);/', $server_name_raw, $matches1);
        $domain_arr = explode(" ", $matches1[1][0]);
        foreach ($domain_arr as $k => $v) {
            if ($k == 0) {
                $website['domain'] = $v;
            } else {
                $website['domain'] .= PHP_EOL.$v;
            }
        }
        // 从nginx配置中root标记位提取运行目录
        $root_raw = cut('# root标记位开始', '# root标记位结束', $nginx_config);
        preg_match_all('/root\s+(.+);/', $root_raw, $matches2);
        $website['root'] = $matches2[1][0];
        // 从nginx配置中index标记位提取全部默认文件
        $index_raw = cut('# index标记位开始', '# index标记位结束', $nginx_config);
        preg_match_all('/index\s+(.+);/', $index_raw, $matches3);
        $website['index'] = $matches3[1][0];

        // 检查网站目录下是否存在.user.ini文件且设置了open_basedir
        if (file_exists($website['path'].'/.user.ini')) {
            $user_ini = file_get_contents($website['path'].'/.user.ini');
            if (str_contains($user_ini, 'open_basedir')) {
                $website['open_basedir'] = 1;
            } else {
                $website['open_basedir'] = 0;
            }
        } else {
            $website['open_basedir'] = 0;
        }

        if ($website['ssl'] == '1') {
            $ssl_certificate_raw = cut('# ssl标记位开始', '# ssl标记位结束', $nginx_config);
            // 从nginx配置中ssl_certificate标记位提取全部证书路径
            preg_match_all('/ssl_certificate\s+(.+);/', $ssl_certificate_raw, $matches4);
            $website['ssl_certificate'] = file_get_contents($matches4[1][0]);
            // 从nginx配置中ssl_certificate_key标记位提取全部证书密钥路径
            preg_match_all('/ssl_certificate_key\s+(.+);/', $ssl_certificate_raw, $matches5);
            $website['ssl_certificate_key'] = file_get_contents($matches5[1][0]);
            $website['http_redirect'] = str_contains($nginx_config, '# http重定向标记位');
            $website['hsts'] = str_contains($nginx_config, '# hsts标记位');
            try {
                $sslDate = (new ACMECert())->getRemainingDays($website['ssl_certificate']);
                $sslDate = round($sslDate, 2);
            } catch (Exception $e) {
                $sslDate = '未知';
            }
            $website['ssl_date'] = $sslDate;
        } else {
            $website['ssl_certificate'] = @file_get_contents('/www/server/vhost/ssl/'.$name.'.pem');
            $website['ssl_certificate_key'] = @file_get_contents('/www/server/vhost/ssl/'.$name.'.key');
            $website['http_redirect'] = 0;
            $website['hsts'] = 0;
        }

        // 从nginx配置中ssl标记位提取waf配置
        $waf_raw = cut('# waf标记位开始', '# waf标记位结束', $nginx_config);
        if (str_contains($waf_raw, 'waf on;')) {
            $website['waf'] = 1;
        } else {
            $website['waf'] = 0;
        }
        preg_match_all('/waf_mode\s+(.+);/', $waf_raw, $matches6);
        $website['waf_mode'] = $matches6[1][0];
        preg_match_all('/waf_cc_deny\s+(.+);/', $waf_raw, $matches7);
        $website['waf_cc_deny'] = $matches7[1][0];
        preg_match_all('/waf_cache\s+(.+);/', $waf_raw, $matches8);
        $website['waf_cache'] = $matches8[1][0];

        // 读取伪静态文件的内容
        $website['rewrite'] = file_get_contents('/www/server/vhost/rewrite/'.$name.'.conf');

        // 读取配置原文
        $website['config_raw'] = file_get_contents('/www/server/vhost/'.$name.'.conf');

        // 读取访问日志
        $website['log'] = shell_exec('tail -n 100 /www/wwwlogs/'.$name.'.log');
        // log需要转义实体
        $website['log'] = htmlspecialchars($website['log']);

        // 如果PHP是0，将其设置为字符串的00
        if ($website['php'] == '0') {
            $website['php'] = '00';
        }

        $res['code'] = 0;
        $res['msg'] = 'success';
        $res['data'] = $website;
        return response()->json($res);
    }

    /**
     * 保存网站设置
     * @param  Request  $request
     * @return JsonResponse
     */
    public function saveSiteSettings(Request $request): JsonResponse
    {
        // 获取前端传递过来的数据
        $name = $request->input('name');
        $config = $request->input('config');

        $website = Website::query()->where('name', $name)->first();
        if (!$website) {
            return response()->json([
                'code' => 1,
                'msg' => '网站不存在',
            ], 200);
        }
        if ($website->status != 1) {
            return response()->json([
                'code' => 1,
                'msg' => '网站已停用，请先启用',
            ], 200);
        }

        $res['code'] = 0;
        $res['msg'] = 'success';

        // 如果config_raw与本地配置文件不一致，则更新配置文件，然后直接返回
        $configRaw = shell_exec('cat /www/server/vhost/'.$name.'.conf');
        if (trim($configRaw) != trim($config['config_raw'])) {
            file_put_contents('/www/server/vhost/'.$name.'.conf', $config['config_raw']);
            return response()->json($res);
        }

        // 检查网站目录是否存在
        if (!is_dir($config['path'])) {
            $res['code'] = 1;
            $res['msg'] = '网站目录不存在';
            return response()->json($res);
        }

        // 域名
        $domain = "server_name";
        $domain_arr = explode(PHP_EOL, $config['domain']);
        foreach ($domain_arr as $v) {
            $domain .= " ".$v;
        }
        $domain .= ';';
        $domain_config_old = cut('# server_name标记位开始', '# server_name标记位结束', $configRaw);
        if (!empty(trim($domain_config_old)) && $domain_config_old != PHP_EOL) {
            $configRaw = str_replace($domain_config_old, PHP_EOL."    ".$domain.PHP_EOL.'    ', $configRaw);
        }

        // 端口
        $port = "";
        $portArr = explode(PHP_EOL, $config['port']);
        foreach ($portArr as $k => $v) {
            // 检查端口是否均为数字
            if (!is_numeric($v) && $v != '443 ssl http2') {
                $res['code'] = 1;
                $res['msg'] = '端口必须为数字';
                return response()->json($res);
            }
            // 检查是否443端口
            if ($v == '443' && $config['ssl'] == '1') {
                $v = '443 ssl http2';
            }
            if ($k != count($portArr) - 1) {
                $port .= "    listen ".$v.';'.PHP_EOL;
            } else {
                $port .= "    listen ".$v.';';
            }
        }
        $port_config_old = cut('# port标记位开始', '# port标记位结束', $configRaw);
        if (!empty(trim($port_config_old)) && $port_config_old != PHP_EOL) {
            $configRaw = str_replace($port_config_old, PHP_EOL.$port.PHP_EOL.'    ', $configRaw);
        }

        // 运行目录
        $pathConfig = cut('# root标记位开始', '# root标记位结束', $configRaw);
        preg_match_all('/root\s+(.+);/', $pathConfig, $matches1);
        $pathConfigOld = $matches1[1][0];
        if (!empty(trim($pathConfigOld)) && $pathConfigOld != PHP_EOL) {
            $pathConfigNew = str_replace($pathConfigOld, $config['root'], $pathConfig);
            $configRaw = str_replace($pathConfig, $pathConfigNew, $configRaw);
        }

        // 默认文件
        $indexConfig = cut('# index标记位开始', '# index标记位结束', $configRaw);
        preg_match_all('/index\s+(.+);/', $indexConfig, $matches2);
        $indexConfigOld = $matches2[1][0];
        if (!empty(trim($indexConfigOld)) && $indexConfigOld != PHP_EOL) {
            $indexConfigNew = str_replace($indexConfigOld, $config['index'], $indexConfig);
            $configRaw = str_replace($indexConfig, $indexConfigNew, $configRaw);
        }

        // open_basedir
        if ($config['open_basedir'] == 1) {
            // 判断$config['path']是否为'/'结尾
            if (str_ends_with($config['path'], '/')) {
                $open_basedir = "open_basedir=".$config['path'].":/tmp/";
            } else {
                $open_basedir = "open_basedir=".$config['path']."/:/tmp/";
            }
            // 写入open_basedir配置到.user.ini文件
            if (is_dir($config['path'])) {
                file_put_contents($config['path'].'/.user.ini', $open_basedir);
                // 为.user.ini文件添加i权限
                // shell_exec('chattr +i '.$config['path'].'/.user.ini');
            }
        } else {
            // 移除.user.ini文件的i权限
            shell_exec('chattr -i '.$config['path'].'/.user.ini');
            // 删除.user.ini文件
            if (file_exists($config['path'].'/.user.ini')) {
                unlink($config['path'].'/.user.ini');
            }
        }

        // waf
        $waf = $config['waf'] == 1 ? 'on' : 'off';
        $wafMode = empty($config['waf_mode']) ? 'DYNAMIC' : $config['waf_mode'];
        $wafCcDeny = empty($config['waf_cc_deny']) ? 'rate=1000r/m duration=60m' : $config['waf_cc_deny'];
        $wafCache = empty($config['waf_cache']) ? 'capacity=50' : $config['waf_cache'];

        $wafConfig = <<<EOF
# waf标记位开始
    waf $waf;
    waf_rule_path /www/server/nginx/ngx_waf/assets/rules/;
    waf_mode $wafMode;
    waf_cc_deny $wafCcDeny;
    waf_cache $wafCache;
EOF;
        $wafConfig .= PHP_EOL.'    ';
        $wafConfigOld = cut('# waf标记位开始', '# waf标记位结束', $configRaw);
        if (!empty(trim($wafConfigOld)) && $wafConfigOld != PHP_EOL) {
            $configRawClean = str_replace($wafConfigOld, "", $configRaw);
        } else {
            $configRawClean = $configRaw;
        }
        $configRaw = str_replace('# waf标记位开始', $wafConfig, $configRawClean);

        // ssl
        if ($config['ssl'] == '1') {
            // 写入证书
            file_put_contents("/www/server/vhost/ssl/".$name.'.pem', $config['ssl_certificate']);
            file_put_contents("/www/server/vhost/ssl/".$name.'.key', $config['ssl_certificate_key']);
            $ssl_config = <<<EOF
# ssl标记位开始
    ssl_certificate /www/server/vhost/ssl/$name.pem;
    ssl_certificate_key /www/server/vhost/ssl/$name.key;
    ssl_session_timeout 1d;
    ssl_session_cache shared:SSL:10m;
    ssl_session_tickets off;
    ssl_dhparam /etc/ssl/certs/dhparam.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
EOF;
            if ($config['http_redirect'] == '1') {
                $ssl_config .= PHP_EOL;
                $ssl_config .= <<<EOF
    # http重定向标记位开始
    if (\$server_port !~ 443){
        return 301 https://\$host\$request_uri;
    }
    error_page 497  https://\$host\$request_uri;
    # http重定向标记位结束
EOF;
            }
            if ($config['hsts'] == '1') {
                $ssl_config .= PHP_EOL;
                $ssl_config .= <<<EOF
    # hsts标记位开始
    add_header Strict-Transport-Security "max-age=63072000" always;
    # hsts标记位结束
EOF;
            }
            $ssl_config .= PHP_EOL.'    ';
            $ssl_config_old = cut('# ssl标记位开始', '# ssl标记位结束', $configRaw);
            if (!empty(trim($ssl_config_old)) && $ssl_config_old != PHP_EOL) {
                $configRaw_clean = str_replace($ssl_config_old, "", $configRaw);
            } else {
                $configRaw_clean = $configRaw;
            }
            $configRaw = str_replace('# ssl标记位开始', $ssl_config, $configRaw_clean);

        } else {
            // 更新nginx配置文件
            $ssl_config_old = cut('# ssl标记位开始', '# ssl标记位结束', $configRaw);
            if (!empty(trim($ssl_config_old)) && $ssl_config_old != PHP_EOL) {
                $configRaw = str_replace($ssl_config_old, PHP_EOL.'    ', $configRaw);
            }
        }

        // 如果PHP版本不一致，则更新PHP版本
        $php_old = Website::query()->where('name', $name)->value('php');
        if ($config['php'] != $php_old) {
            $php_config_old = cut('# php标记位开始', '# php标记位结束', $configRaw);
            $php_config_new = PHP_EOL;
            $php_config_new .= <<<EOL
    include enable-php-$config[php].conf;
EOL;
            $php_config_new .= PHP_EOL.'    ';

            if (!empty(trim($php_config_old)) && $php_config_old != PHP_EOL) {
                $configRaw = str_replace($php_config_old, $php_config_new, $configRaw);
            }
        }

        // 将数据入库
        $website->php = $config['php'];
        $website->ssl = $config['ssl'];
        $website->path = $config['path'];
        $website->save();
        file_put_contents('/www/server/vhost/'.$name.'.conf', $configRaw);
        file_put_contents('/www/server/vhost/rewrite/'.$name.'.conf', $config['rewrite']);
        shell_exec('systemctl reload nginx');
        return response()->json($res);
    }

    /**
     * 清理网站日志
     * @param  Request  $request
     * @return JsonResponse
     */
    public function clearSiteLog(Request $request): JsonResponse
    {
        $name = $request->input('name');
        shell_exec('echo "" > /www/wwwlogs/'.$name.'.log');
        $res['code'] = 0;
        $res['msg'] = 'success';
        return response()->json($res);
    }

    /**
     * 修改网站备注
     * @param  Request  $request
     * @return JsonResponse
     */
    public function updateSiteNote(Request $request): JsonResponse
    {
        $name = $request->input('name');
        $note = $request->input('note');
        Website::query()->where('name', $name)->update(['note' => $note]);
        $res['code'] = 0;
        $res['msg'] = 'success';
        return response()->json($res);
    }

    /**
     * 获取备份列表
     */
    public function getBackupList(): JsonResponse
    {
        $backupPath = '/www/backup/website';
        // 判断备份目录是否存在
        if (!is_dir($backupPath)) {
            mkdir($backupPath, 0644, true);
        }
        $backupFiles = scandir($backupPath);
        $backupFiles = array_diff($backupFiles, ['.', '..']);
        $backupFiles = array_values($backupFiles);
        $backupFiles = array_map(function ($backupFile) {
            return [
                'backup' => $backupFile,
                'size' => formatBytes(filesize('/www/backup/website/'.$backupFile)),
            ];
        }, $backupFiles);
        $res['code'] = 0;
        $res['msg'] = 'success';
        $res['data'] = $backupFiles;
        return response()->json($res);
    }

    /**
     * 创建备份
     */
    public function createBackup(Request $request): JsonResponse
    {
        // 消毒数据
        try {
            $credentials = $this->validate($request, [
                'name' => 'required|max:255',
            ]);
        } catch (ValidationException $e) {
            return response()->json([
                'code' => 1,
                'msg' => '参数错误：'.$e->getMessage(),
                'errors' => $e->errors()
            ], 200);
        }

        $backupPath = '/www/backup/website';
        // 判断备份目录是否存在
        if (!is_dir($backupPath)) {
            mkdir($backupPath, 0644, true);
        }

        // 从数据库中获取网站目录
        $sitePath = Website::query()->where('name', $credentials['name'])->value('path');
        $backupFile = $backupPath.'/'.$credentials['name'].'_'.date('YmdHis').'.zip';
        shell_exec('zip -r '.$backupFile.' '.$sitePath.' 2>&1');

        $res['code'] = 0;
        $res['msg'] = 'success';
        return response()->json($res);
    }

    /**
     * 上传备份
     */
    public function uploadBackup(Request $request): JsonResponse
    {
        // 消毒数据
        try {
            $credentials = $this->validate($request, [
                'file' => 'required|file',
            ]);
        } catch (ValidationException $e) {
            return response()->json([
                'code' => 1,
                'msg' => '参数错误：'.$e->getMessage(),
                'errors' => $e->errors()
            ], 200);
        }

        $file = $request->file('file');
        $backupPath = '/www/backup/website';

        // 判断备份目录是否存在
        if (!is_dir($backupPath)) {
            mkdir($backupPath, 0644, true);
        }
        $backupFile = $backupPath.'/'.$file->getClientOriginalName();
        $file->move($backupPath, $file->getClientOriginalName());

        // 返回文件名
        $res['code'] = 0;
        $res['msg'] = 'success';
        $res['data'] = $file->getClientOriginalName();
        return response()->json($res);
    }

    /**
     * 恢复备份
     */
    public function restoreBackup(Request $request): JsonResponse
    {
        // 消毒数据
        try {
            $credentials = $this->validate($request, [
                'name' => 'required|max:255',
                'backup' => 'required|max:255',
            ]);
        } catch (ValidationException $e) {
            return response()->json([
                'code' => 1,
                'msg' => '参数错误：'.$e->getMessage(),
                'errors' => $e->errors()
            ], 200);
        }

        $backupPath = '/www/backup/website';
        // 判断备份目录是否存在
        if (!is_dir($backupPath)) {
            mkdir($backupPath, 0644, true);
        }
        $backupFile = $backupPath.'/'.$credentials['backup'];
        // 判断备份文件是否存在
        if (!is_file($backupFile)) {
            return response()->json([
                'code' => 1,
                'msg' => '备份文件不存在',
            ], 200);
        }

        shell_exec('rm -rf /www/wwwroot/'.$credentials['name'].'/*');
        shell_exec('unzip -o '.$backupFile.' -d /www/wwwroot/'.$credentials['name'].' 2>&1');
        // 设置权限
        shell_exec('chown -R www:www /www/wwwroot/'.$credentials['name']);
        shell_exec('chmod -R 755 /www/wwwroot/'.$credentials['name']);

        $res['code'] = 0;
        $res['msg'] = 'success';
        return response()->json($res);
    }

    /**
     * 删除备份
     */
    public function deleteBackup(Request $request): JsonResponse
    {
        // 消毒数据
        try {
            $credentials = $this->validate($request, [
                'backup' => 'required|max:255',
            ]);
        } catch (ValidationException $e) {
            return response()->json([
                'code' => 1,
                'msg' => '参数错误：'.$e->getMessage(),
                'errors' => $e->errors()
            ], 200);
        }

        $backupPath = '/www/backup/website';
        // 判断备份目录是否存在
        if (!is_dir($backupPath)) {
            mkdir($backupPath, 0644, true);
        }
        $backupFile = $backupPath.'/'.$credentials['backup'];
        // 判断备份文件是否存在
        if (!is_file($backupFile)) {
            return response()->json([
                'code' => 1,
                'msg' => '备份文件不存在',
            ], 200);
        }

        unlink($backupFile);
        $res['code'] = 0;
        $res['msg'] = 'success';
        return response()->json($res);
    }

    /**
     * 重置网站配置文件
     */
    public function resetSiteConfig(Request $request): JsonResponse
    {
        try {
            $credentials = $this->validate($request, [
                'name' => 'required|max:255',
            ]);
        } catch (ValidationException $e) {
            return response()->json([
                'code' => 1,
                'msg' => '参数错误：'.$e->getMessage(),
                'errors' => $e->errors()
            ], 200);
        }

        $website = Website::query()->where('name', $credentials['name'])->first();
        if (!$website) {
            return response()->json([
                'code' => 1,
                'msg' => '网站不存在',
            ], 200);
        }

        // 如果PHP是0，将其设置为字符串的00
        if ($website['php'] == '0') {
            $website['php'] = '00';
        }

        // 更新网站状态为运行
        $website->status = 1;
        // 更新网站ssl状态为关闭
        $website->ssl = 0;
        $website->save();

        $nginxConfig = <<<EOF
# 配置文件中的标记位请勿随意修改，改错将导致面板无法识别！
# 有自定义配置需求的，请将自定义的配置写在各标记位下方。
server
{
    # port标记位开始
    listen 80;
    # port标记位结束
    # server_name标记位开始
    server_name localhost;
    # server_name标记位结束
    # index标记位开始
    index index.php index.html;
    # index标记位结束
    # root标记位开始
    root $website[path];
    # root标记位结束

    # ssl标记位开始
    # ssl标记位结束

    # php标记位开始
    include enable-php-$website[php].conf;
    # php标记位结束

    # waf标记位开始
    waf on;
    waf_rule_path /www/server/nginx/ngx_waf/assets/rules/;
    waf_mode DYNAMIC;
    waf_cc_deny rate=1000r/m duration=60m;
    waf_cache capacity=50;
    # waf标记位结束

    # 错误页配置，可自行设置
    #error_page 404 /404.html;
    #error_page 502 /502.html;

    # 伪静态规则引入，修改后将导致面板设置的伪静态规则失效
    include /www/server/vhost/rewrite/$website[name].conf;

    # 面板默认禁止访问部分敏感目录，可自行修改
    location ~ ^/(\.user.ini|\.htaccess|\.git|\.svn)
    {
        return 404;
    }
    # 面板默认不记录静态资源的访问日志并开启1小时浏览器缓存，可自行修改
    location ~ .*\.(js|css)$
    {
        expires 1h;
        error_log /dev/null;
        access_log /dev/null;
    }
    access_log /www/wwwlogs/$website[name].log;
    error_log /www/wwwlogs/$website[name].log;
}
EOF;
        file_put_contents('/www/server/vhost/'.$website['name'].'.conf', $nginxConfig);
        // 重置伪静态规则
        shell_exec('echo "" > /www/server/vhost/rewrite/'.$website['name'].'.conf');
        // 重载nginx
        shell_exec('systemctl reload nginx');

        // 返回
        $res['code'] = 0;
        $res['msg'] = 'success';
        return response()->json($res);
    }

    /**
     * 设置网站运行状态
     */
    public function setSiteStatus(Request $request): JsonResponse
    {
        try {
            $credentials = $this->validate($request, [
                'name' => 'required|max:255',
                'status' => 'required|in:0,1',
            ]);
        } catch (ValidationException $e) {
            return response()->json([
                'code' => 1,
                'msg' => '参数错误：'.$e->getMessage(),
                'errors' => $e->errors()
            ], 200);
        }

        $website = Website::query()->where('name', $credentials['name'])->first();
        if (!$website) {
            return response()->json([
                'code' => 1,
                'msg' => '网站不存在',
            ], 200);
        }

        $nginxConfig = file_get_contents('/www/server/vhost/'.$website['name'].'.conf');

        // 运行目录
        $pathConfig = cut('# root标记位开始', '# root标记位结束', $nginxConfig);
        preg_match_all('/root\s+(.+);/', $pathConfig, $matches1);
        $pathConfigOld = $matches1[1][0];
        if (!empty(trim($pathConfigOld)) && $pathConfigOld != PHP_EOL) {
            if ($credentials['status'] == 0) {
                $pathConfigNew = str_replace($pathConfigOld, '/www/server/nginx/html', $pathConfig);
                // 将旧配置追加到新配置中
                $pathConfigNew .= '# '.$pathConfigOld.PHP_EOL.'    ';
            } else {
                // 匹配旧配置
                preg_match_all('/# (.+)/', $pathConfig, $matches2);
                // 还原旧配置
                $pathConfigNew = str_replace($pathConfigOld, $matches2[1][0], $pathConfig);
                // 删除旧配置
                $pathConfigNew = str_replace(PHP_EOL.'    # '.$matches2[1][0], '', $pathConfigNew);
            }
            $nginxConfig = str_replace($pathConfig, $pathConfigNew, $nginxConfig);
        }

        // 默认文件
        $indexConfig = cut('# index标记位开始', '# index标记位结束', $nginxConfig);
        preg_match_all('/index\s+(.+);/', $indexConfig, $matches2);
        $indexConfigOld = $matches2[1][0];
        if (!empty(trim($indexConfigOld)) && $indexConfigOld != PHP_EOL) {
            if ($credentials['status'] == 0) {
                $indexConfigNew = str_replace($indexConfigOld, 'stop.html', $indexConfig);
                // 将旧配置追加到新配置中
                $indexConfigNew .= '# '.$indexConfigOld.PHP_EOL.'    ';
            } else {
                // 匹配旧配置
                preg_match_all('/# (.+)/', $indexConfig, $matches2);
                // 还原旧配置
                $indexConfigNew = str_replace($indexConfigOld, $matches2[1][0], $indexConfig);
                // 删除旧配置
                $indexConfigNew = str_replace(PHP_EOL.'    # '.$matches2[1][0], '', $indexConfigNew);
            }
            $nginxConfig = str_replace($indexConfig, $indexConfigNew, $nginxConfig);
        }

        // 写入配置文件
        file_put_contents('/www/server/vhost/'.$website['name'].'.conf', $nginxConfig);

        $website->status = $credentials['status'];
        $website->save();

        // 重载nginx
        shell_exec('systemctl reload nginx');

        // 返回
        $res['code'] = 0;
        $res['msg'] = 'success';
        return response()->json($res);
    }

    /**
     * 签发SSL证书
     * @param  Request  $request
     * @return JsonResponse
     * @throws ACME_Exception|Exception
     */
    public function issueSsl(Request $request): JsonResponse
    {
        try {
            $input = $this->validate($request, [
                'type' => 'required|in:lets,buypass,google,sslcom,zerossl',
                'name' => 'required',
            ]);
        } catch (ValidationException $e) {
            return response()->json([
                'code' => 1,
                'msg' => '参数错误：'.$e->getMessage(),
                'errors' => $e->errors()
            ], 200);
        }

        $user = $request->user();

        // 检查网站是否存在
        $website = Website::query()->where('name', $input['name'])->first();
        if (!$website) {
            return response()->json([
                'code' => 1,
                'msg' => '网站不存在',
            ], 200);
        }
        // 从配置文件中获取网站域名
        $nginxConfig = file_get_contents('/www/server/vhost/'.$website['name'].'.conf');
        $domainConfig = cut('# server_name标记位开始', '# server_name标记位结束', $nginxConfig);
        preg_match_all('/server_name\s+(.+);/', $domainConfig, $matches1);
        $domains = explode(" ", $matches1[1][0]);
        // 从配置文件中获取网站目录
        $pathConfig = cut('# root标记位开始', '# root标记位结束', $nginxConfig);
        preg_match_all('/root\s+(.+);/', $pathConfig, $matches2);
        $path = $matches2[1][0];

        /**
         * 对域名需要进行一下处理，如果域名是泛域名，返回暂不支持泛域名
         */
        foreach ($domains as $domain) {
            if (str_contains($domain, '*')) {
                return response()->json([
                    'code' => 1,
                    'msg' => '暂不支持泛域名',
                ], 200);
            }
        }

        switch ($input['type']) {
            case 'lets':
                $ac = new ACMECert('https://acme-v02.api.letsencrypt.org/directory');
                break;
            case 'buypass':
                $ac = new ACMECert('https://api.buypass.com/acme/directory');
                break;
            case 'google':
                $ac = new ACMECert('https://dv.acme-v02.api.pki.goog/directory');
                break;
            case 'sslcom':
                $ac = new ACMECert('https://acme.ssl.com/sslcom-dv-rsa');
                break;
            case 'zerossl':
                $ac = new ACMECert('https://acme.zerossl.com/v2/DV90');
                break;
            default:
                $res = [
                    'code' => 1,
                    'msg' => '参数错误：type',
                ];
                return response()->json($res);
                break;
        }

        try {
            $accountKey = $ac->generateECKey('P-384');
            $certKey = $ac->generateECKey('P-384');
        } catch (Exception $e) {
            return response()->json([
                'code' => 1,
                'msg' => '生成密钥失败：'.$e->getMessage(),
            ], 200);
        }
        try {
            $ac->loadAccountKey($accountKey);
        } catch (Exception $e) {
            return response()->json([
                'code' => 1,
                'msg' => '加载密钥失败：'.$e->getMessage(),
            ], 200);
        }
        try {
            $ac->register(true, $user->email);
        } catch (Exception $e) {
            return response()->json([
                'code' => 1,
                'msg' => '注册CA账户失败：'.$e->getMessage(),
            ], 200);
        }

        // 初始化域名数组
        $domainConfig = [];
        foreach ($domains as $domain) {
            $domainConfig[$domain] = [
                'challenge' => 'http-01',
                'docroot' => $path
            ];
        }

        $handler = function ($opts) {
            $fn = $opts['config']['docroot'].$opts['key'];
            @mkdir(dirname($fn), 0777, true);
            file_put_contents($fn, $opts['value']);
            return function ($opts) {
                unlink($opts['config']['docroot'].$opts['key']);
            };
        };

        // 申请证书
        try {
            $fullchain = $ac->getCertificateChain($certKey, $domainConfig, $handler);
        } catch (ACME_Exception $e) {
            return response()->json([
                'code' => 1,
                'msg' => '申请证书失败：'.$e->getMessage(),
            ], 200);
        }

        // 写入证书
        $sslDir = '/www/server/vhost/ssl/';
        file_put_contents($sslDir.$website['name'].'.key', $certKey);
        file_put_contents($sslDir.$website['name'].'.pem', $fullchain);

        // 返回
        $res = [
            'code' => 0,
            'msg' => 'success',
        ];
        return response()->json($res);
    }

}
