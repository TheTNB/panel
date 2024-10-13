package nginx

var order = []string{"listen", "server_name", "index", "root",
	"ssl_certificate", "ssl_certificate_key", "ssl_session_timeout", "ssl_session_cache", "ssl_protocols", "ssl_ciphers", "ssl_prefer_server_ciphers", "ssl_early_data", "ssl_stapling", "ssl_stapling_verify", "ssl_trusted_certificate",
	"resolver", "error_page", "include", "if", "location", "access_log", "error_log"}

const defaultConf = `server
{
    listen 80;
    server_name localhost;
    index index.php index.html;
    root /www/wwwroot/default;
    # 错误页配置，可自行修改
    #error_page 502 /502.html;
    error_page 404 /404.html;
    include enable-php-0.conf;
    # acme证书签发配置，不可修改
    include /www/server/vhost/acme/test.conf;
    # 伪静态规则引入，修改后将导致面板设置的伪静态规则失效
    include /www/server/vhost/rewrite/test.conf;
    # 禁止访问部分敏感目录，可自行修改
    location ~ ^/(\.user.ini|\.htaccess|\.git|\.svn)
    {
        return 404;
    }
    # 不记录静态资源的访问日志，可自行修改
    location ~ .*\.(js|css|ttf|otf|woff|woff2|eot)$
    {
        expires 1h;
        error_log /dev/null;
        access_log /dev/null;
    }
    access_log /www/wwwlogs/default.log;
    error_log /www/wwwlogs/default.log;
}
`
