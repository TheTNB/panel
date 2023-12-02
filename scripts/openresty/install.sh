#!/bin/bash
export PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:$PATH

: '
Copyright (C) 2022 - now  HaoZi Technology Co., Ltd.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
'

HR="+----------------------------------------------------"
ARCH=$(uname -m)
OS=$(source /etc/os-release && { [[ "$ID" == "debian" ]] && echo "debian"; } || { [[ "$ID" == "centos" ]] || [[ "$ID" == "rhel" ]] || [[ "$ID" == "rocky" ]] || [[ "$ID" == "almalinux" ]] && echo "centos"; } || echo "unknown")
downloadUrl="https://git.haozi.net/opensource/download/-/raw/main/panel/openresty"
setupPath="/www"
openrestyPath="${setupPath}/server/openresty"
openrestyVersion="1.21.4.3"
cpuCore=$(cat /proc/cpuinfo | grep "processor" | wc -l)

source ${setupPath}/panel/scripts/calculate_j.sh
j=$(calculate_j)

# 安装依赖
if [ "${OS}" == "centos" ]; then
    dnf makecache -y
    dnf groupinstall "Development Tools" -y
    dnf install tar unzip gd gd-devel git-core flex perl oniguruma oniguruma-devel libsodium-devel libxml2-devel libxslt-devel bison yajl yajl-devel curl curl-devel ncurses-devel libevent-devel readline-devel libuuid-devel brotli-devel icu libicu libicu-devel openssl openssl-devel -y
elif [ "${OS}" == "debian" ]; then
    apt-get update
    apt-get install build-essential tar unzip libgd3 libgd-dev git flex perl libonig-dev libsodium-dev libxml2-dev libxslt1-dev bison libyajl-dev curl libcurl4-openssl-dev libncurses5-dev libevent-dev libreadline-dev uuid-dev libbrotli-dev icu-devtools libicu-dev openssl libssl-dev -y
else
    echo -e $HR
    echo "错误：耗子Linux面板不支持该系统"
    exit 1
fi

# 准备目录
rm -rf ${openrestyPath}
mkdir -p ${openrestyPath}
cd ${openrestyPath}

# 下载源码
wget -T 120 -t 3 -O ${openrestyPath}/openresty-${openrestyVersion}.tar.gz ${downloadUrl}/openresty-${openrestyVersion}.tar.gz
wget -T 20 -t 3 -O ${openrestyPath}/openresty-${openrestyVersion}.tar.gz.checksum.txt ${downloadUrl}/openresty-${openrestyVersion}.tar.gz.checksum.txt

if ! sha256sum --status -c openresty-${openrestyVersion}.tar.gz.checksum.txt; then
    echo -e $HR
    echo "错误：OpenResty 源码 checksum 校验失败，文件可能被篡改或不完整，已终止操作"
    rm -rf ${openrestyPath}
    exit 1
fi

tar -zxvf openresty-${openrestyVersion}.tar.gz
rm -f openresty-${openrestyVersion}.tar.gz
rm -f openresty-${openrestyVersion}.tar.gz.checksum.txt
mv openresty-${openrestyVersion} src
cd src

# openssl
wget -T 120 -t 3 -O openssl-1.1.1w.tar.gz ${downloadUrl}/openssl/openssl-1.1.1w.tar.gz
wget -T 20 -t 3 -O openssl-1.1.1w.tar.gz.checksum.txt ${downloadUrl}/openssl/openssl-1.1.1w.tar.gz.checksum.txt

if ! sha256sum --status -c openssl-1.1.1w.tar.gz.checksum.txt; then
    echo -e $HR
    echo "错误：OpenSSL 源码 checksum 校验失败，文件可能被篡改或不完整，已终止操作"
    rm -rf ${openrestyPath}
    exit 1
fi

tar -zxvf openssl-1.1.1w.tar.gz
rm -f openssl-1.1.1w.tar.gz
rm -f openssl-1.1.1w.tar.gz.checksum.txt
mv openssl-1.1.1w openssl

# patch openssl
cd openssl
wget -T 20 -t 3 -O openssl-1.1.1f-sess_set_get_cb_yield.patch ${downloadUrl}/openssl/openssl-1.1.1f-sess_set_get_cb_yield.patch
wget -T 20 -t 3 -O openssl-1.1.1f-sess_set_get_cb_yield.patch.checksum.txt ${downloadUrl}/openssl/openssl-1.1.1f-sess_set_get_cb_yield.patch.checksum.txt

if ! sha256sum --status -c openssl-1.1.1f-sess_set_get_cb_yield.patch.checksum.txt; then
    echo -e $HR
    echo "错误：OpenSSL 补丁文件 checksum 校验失败，文件可能被篡改或不完整，已终止操作"
    rm -rf ${openrestyPath}
    exit 1
fi

patch -p1 < openssl-1.1.1f-sess_set_get_cb_yield.patch
rm -f openssl-1.1.1f-sess_set_get_cb_yield.patch
rm -f openssl-1.1.1f-sess_set_get_cb_yield.patch.checksum.txt
cd ../

# pcre
wget -T 60 -t 3 -O pcre-8.45.tar.gz ${downloadUrl}/pcre/pcre-8.45.tar.gz
wget -T 20 -t 3 -O pcre-8.45.tar.gz.checksum.txt ${downloadUrl}/pcre/pcre-8.45.tar.gz.checksum.txt

if ! sha256sum --status -c pcre-8.45.tar.gz.checksum.txt; then
    echo -e $HR
    echo "错误：pcre 源码 checksum 校验失败，文件可能被篡改或不完整，已终止操作"
    rm -rf ${openrestyPath}
    exit 1
fi

tar -zxvf pcre-8.45.tar.gz
rm -f pcre-8.45.tar.gz
rm -f pcre-8.45.tar.gz.checksum.txt
mv pcre-8.45 pcre

# ngx_cache_purge
wget -T 20 -t 3 -O ngx_cache_purge-2.3.tar.gz ${downloadUrl}/modules/ngx_cache_purge-2.3.tar.gz
wget -T 20 -t 3 -O ngx_cache_purge-2.3.tar.gz.checksum.txt ${downloadUrl}/modules/ngx_cache_purge-2.3.tar.gz.checksum.txt

if ! sha256sum --status -c ngx_cache_purge-2.3.tar.gz.checksum.txt; then
    echo -e $HR
    echo "错误：ngx_cache_purge 源码 checksum 校验失败，文件可能被篡改或不完整，已终止操作"
    rm -rf ${openrestyPath}
    exit 1
fi

tar -zxvf ngx_cache_purge-2.3.tar.gz
rm -f ngx_cache_purge-2.3.tar.gz
rm -f ngx_cache_purge-2.3.tar.gz.checksum.txt
mv ngx_cache_purge-2.3 ngx_cache_purge

# nginx-sticky-module
wget -T 20 -t 3 -O nginx-sticky-module.zip ${downloadUrl}/modules/nginx-sticky-module.zip
wget -T 20 -t 3 -O nginx-sticky-module.zip.checksum.txt ${downloadUrl}/modules/nginx-sticky-module.zip.checksum.txt

if ! sha256sum --status -c nginx-sticky-module.zip.checksum.txt; then
    echo -e $HR
    echo "错误：nginx-sticky-module 源码 checksum 校验失败，文件可能被篡改或不完整，已终止操作"
    rm -rf ${openrestyPath}
    exit 1
fi

unzip -o nginx-sticky-module.zip
rm -f nginx-sticky-module.zip
rm -f nginx-sticky-module.zip.checksum.txt

# nginx-dav-ext-module
wget -T 20 -t 3 -O nginx-dav-ext-module-3.0.0.tar.gz ${downloadUrl}/modules/nginx-dav-ext-module-3.0.0.tar.gz
wget -T 20 -t 3 -O nginx-dav-ext-module-3.0.0.tar.gz.checksum.txt ${downloadUrl}/modules/nginx-dav-ext-module-3.0.0.tar.gz.checksum.txt

if ! sha256sum --status -c nginx-dav-ext-module-3.0.0.tar.gz.checksum.txt; then
    echo -e $HR
    echo "错误：nginx-dav-ext-module 源码 checksum 校验失败，文件可能被篡改或不完整，已终止操作"
    rm -rf ${openrestyPath}
    exit 1
fi

tar -xvf nginx-dav-ext-module-3.0.0.tar.gz
rm -f nginx-dav-ext-module-3.0.0.tar.gz
rm -f nginx-dav-ext-module-3.0.0.tar.gz.checksum.txt
mv nginx-dav-ext-module-3.0.0 nginx-dav-ext-module

# waf
wget -T 60 -t 3 -O uthash-2.3.0.zip ${downloadUrl}/modules/uthash-2.3.0.zip
wget -T 20 -t 3 -O uthash-2.3.0.zip.checksum.txt ${downloadUrl}/modules/uthash-2.3.0.zip.checksum.txt

if ! sha256sum --status -c uthash-2.3.0.zip.checksum.txt; then
    echo -e $HR
    echo "错误：uthash 源码 checksum 校验失败，文件可能被篡改或不完整，已终止操作"
    rm -rf ${openrestyPath}
    exit 1
fi

unzip -o uthash-2.3.0.zip
mv uthash-2.3.0 uthash
rm -f uthash-2.3.0.zip
rm -f uthash-2.3.0.zip.checksum.txt
cd ../

wget -T 20 -t 3 -O ngx_waf-6.1.9.zip ${downloadUrl}/modules/ngx_waf-6.1.9.zip
wget -T 20 -t 3 -O ngx_waf-6.1.9.zip.checksum.txt ${downloadUrl}/modules/ngx_waf-6.1.9.zip.checksum.txt

if ! sha256sum --status -c ngx_waf-6.1.9.zip.checksum.txt; then
    echo -e $HR
    echo "错误：ngx_waf 源码 checksum 校验失败，文件可能被篡改或不完整，已终止操作"
    rm -rf ${openrestyPath}
    exit 1
fi

unzip -o ngx_waf-6.1.9.zip
mv ngx_waf-6.1.9 ngx_waf
rm -f ngx_waf-6.1.9.zip
rm -f ngx_waf-6.1.9.zip.checksum.txt

cd ngx_waf/inc
wget -T 60 -t 3 -O libinjection-3.10.0.zip ${downloadUrl}/modules/libinjection-3.10.0.zip
wget -T 20 -t 3 -O libinjection-3.10.0.zip.checksum.txt ${downloadUrl}/modules/libinjection-3.10.0.zip.checksum.txt

if ! sha256sum --status -c libinjection-3.10.0.zip.checksum.txt; then
    echo -e $HR
    echo "错误：libinjection 源码 checksum 校验失败，文件可能被篡改或不完整，已终止操作"
    rm -rf ${openrestyPath}
    exit 1
fi

unzip -o libinjection-3.10.0.zip
mv libinjection-3.10.0 libinjection
rm -f libinjection-3.10.0.zip
rm -f libinjection-3.10.0.zip.checksum.txt

cd ../
make "-j${j}"
if [ "$?" != "0" ]; then
    echo -e $HR
    echo "错误：OpenResty waf拓展初始化失败，请截图错误信息寻求帮助。"
    rm -rf ${openrestyPath}
    exit 1
fi
cd ${openrestyPath}/src

# brotli
wget -T 20 -t 3 -O ngx_brotli-1.0.0rc.zip ${downloadUrl}/modules/ngx_brotli-1.0.0rc.zip
wget -T 20 -t 3 -O ngx_brotli-1.0.0rc.zip.checksum.txt ${downloadUrl}/modules/ngx_brotli-1.0.0rc.zip.checksum.txt

if ! sha256sum --status -c ngx_brotli-1.0.0rc.zip.checksum.txt; then
    echo -e $HR
    echo "错误：ngx_brotli 源码 checksum 校验失败，文件可能被篡改或不完整，已终止操作"
    rm -rf ${openrestyPath}
    exit 1
fi

unzip -o ngx_brotli-1.0.0rc.zip
mv ngx_brotli-1.0.0rc ngx_brotli
rm -f ngx_brotli-1.0.0rc.zip
rm -f ngx_brotli-1.0.0rc.zip.checksum.txt
cd ngx_brotli/deps
rm -rf brotli

wget -T 20 -t 3 -O brotli-1.0.9.zip ${downloadUrl}/modules/brotli-1.0.9.zip
wget -T 20 -t 3 -O brotli-1.0.9.zip.checksum.txt ${downloadUrl}/modules/brotli-1.0.9.zip.checksum.txt

if ! sha256sum --status -c brotli-1.0.9.zip.checksum.txt; then
    echo -e $HR
    echo "错误：brotli 源码 checksum 校验失败，文件可能被篡改或不完整，已终止操作"
    rm -rf ${openrestyPath}
    exit 1
fi

unzip -o brotli-1.0.9.zip
mv brotli-1.0.9 brotli
rm -f brotli-1.0.9.zip
rm -f brotli-1.0.9.zip.checksum.txt

cd ${openrestyPath}/src
export LD_LIBRARY_PATH=/usr/local/lib/:$LD_LIBRARY_PATH
export LIB_UTHASH=${openrestyPath}/src/uthash

./configure --user=www --group=www --prefix=${openrestyPath} --with-luajit --add-module=${openrestyPath}/src/ngx_cache_purge --add-module=${openrestyPath}/src/nginx-sticky-module --with-openssl=${openrestyPath}/src/openssl --with-pcre=${openrestyPath}/src/pcre --with-pcre-jit --with-http_v2_module --with-http_slice_module --with-threads --with-stream --with-stream_ssl_module --with-stream_realip_module --with-stream_ssl_preread_module --with-http_stub_status_module --with-http_ssl_module --with-http_image_filter_module --with-http_gzip_static_module --with-http_gunzip_module --with-ipv6 --with-http_sub_module --with-http_flv_module --with-http_addition_module --with-http_realip_module --with-http_mp4_module --with-http_auth_request_module --with-http_secure_link_module --with-http_random_index_module --with-ld-opt="-Wl,-E" --with-cc-opt="-DNGX_LUA_ABORT_AT_PANIC" --with-luajit-xcflags="-DLUAJIT_NUMMODE=2 -DLUAJIT_ENABLE_LUA52COMPAT" --with-compat --with-http_dav_module --add-module=${openrestyPath}/src/nginx-dav-ext-module --add-module=${openrestyPath}/src/ngx_brotli --add-module=${openrestyPath}/ngx_waf
make "-j${j}"
if [ "$?" != "0" ]; then
    echo -e $HR
    echo "错误：OpenResty编译失败，请截图错误信息寻求帮助。"
    rm -rf ${openrestyPath}
    exit 1
fi
make install
if [ ! -f "${openrestyPath}/nginx/sbin/nginx" ]; then
    echo -e $HR
    echo "错误：OpenResty安装失败，请截图错误信息寻求帮助。"
    rm -rf ${openrestyPath}
    exit 1
fi

# 设置软链接
ln -sf ${openrestyPath}/nginx/html ${openrestyPath}/html
ln -sf ${openrestyPath}/nginx/conf ${openrestyPath}/conf
ln -sf ${openrestyPath}/nginx/logs ${openrestyPath}/logs
ln -sf ${openrestyPath}/nginx/sbin ${openrestyPath}/sbin
ln -sf ${openrestyPath}/nginx/sbin/nginx /usr/bin/openresty
rm -f ${openrestyPath}/conf/nginx.conf

# 创建配置目录
cd ${openrestyPath}
rm -f openresty-${openrestyVersion}.tar.gz
rm -rf src
mkdir -p /www/wwwroot/default
mkdir -p /www/wwwlogs
mkdir -p /www/server/vhost
mkdir -p /www/server/vhost
mkdir -p /www/server/vhost/rewrite
mkdir -p /www/server/vhost/ssl

# 写入主配置文件
cat > ${openrestyPath}/conf/nginx.conf << EOF
# 该文件为OpenResty主配置文件，不建议随意修改！
user www www;
worker_processes auto;
error_log /www/wwwlogs/openresty_error.log crit;
pid /www/server/openresty/nginx.pid;
worker_rlimit_nofile 51200;

stream {
    log_format tcp_format '\$time_local|\$remote_addr|\$protocol|\$status|\$bytes_sent|\$bytes_received|\$session_time|\$upstream_addr|\$upstream_bytes_sent|\$upstream_bytes_received|\$upstream_connect_time';

    access_log /www/wwwlogs/tcp-access.log tcp_format;
    error_log /www/wwwlogs/tcp-error.log;
}

events {
    use epoll;
    worker_connections 51200;
    multi_accept on;
}

http {
    include mime.types;
    include proxy.conf;
    include default.conf;
    default_type application/octet-stream;

    server_names_hash_bucket_size 512;
    client_header_buffer_size 32k;
    large_client_header_buffers 4 32k;
    client_max_body_size 200m;
    client_body_buffer_size 10M;
    client_body_in_file_only off;

    variables_hash_max_size 2048;
    variables_hash_bucket_size 128;

    sendfile on;
    tcp_nopush on;

    keepalive_timeout 60;

    tcp_nodelay on;

    fastcgi_connect_timeout 300;
    fastcgi_send_timeout 300;
    fastcgi_read_timeout 300;
    fastcgi_buffer_size 64k;
    fastcgi_buffers 8 64k;
    fastcgi_busy_buffers_size 256k;
    fastcgi_temp_file_write_size 256k;
    fastcgi_intercept_errors on;

    gzip on;
    gzip_min_length 1k;
    gzip_buffers 32 4k;
    gzip_http_version 1.1;
    gzip_comp_level 6;
    gzip_types *;
    gzip_vary on;
    gzip_proxied any;
    gzip_disable "MSIE [1-6]\.";
    brotli on;
    brotli_comp_level 6;
    brotli_min_length 10;
    brotli_window 1m;
    brotli_types *;
    brotli_static on;

    limit_conn_zone \$binary_remote_addr zone=perip:10m;
    limit_conn_zone \$server_name zone=perserver:10m;

    server_tokens off;
    access_log off;

    # 服务状态页
    server {
        listen 80;
        server_name 127.0.0.1;
        allow 127.0.0.1;

        location /nginx_status {
            stub_status on;
            access_log off;
        }
        location ~ ^/phpfpm_status/(?<version>\d+)$ {
            fastcgi_pass unix:/tmp/php-cgi-\$version.sock;
            include fastcgi_params;
            fastcgi_param SCRIPT_FILENAME \$fastcgi_script_name;
        }
    }
    include /www/server/vhost/*.conf;
}
EOF
# 写入pathinfo配置文件
cat > ${openrestyPath}/conf/pathinfo.conf << EOF
set \$real_script_name \$fastcgi_script_name;
if (\$fastcgi_script_name ~ "^(.+?\.php)(/.+)$") {
    set \$real_script_name \$1;
    set \$path_info \$2;
 }
fastcgi_param SCRIPT_FILENAME \$document_root\$real_script_name;
fastcgi_param SCRIPT_NAME \$real_script_name;
fastcgi_param PATH_INFO \$path_info;
EOF
# 写入默认站点页
cat > ${openrestyPath}/html/index.html << EOF
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>未找到网站 - 耗子Linux面板</title>
    <style>
        body {
            background-color: #f9f9f9;
            margin: 0;
            padding: 0;
        }
        .container {
            max-width: 800px;
            margin: 2em auto;
            background-color: #ffffff;
            padding: 20px;
            border-radius: 12px;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
        }
        h1 {
            font-size: 2.5em;
            margin-top: 0;
            margin-bottom: 20px;
            text-align: center;
            color: #333;
            border-bottom: 2px solid #ddd;
            padding-bottom: 0.5em;
        }
        p {
            color: #555;
            line-height: 1.8;
        }
        @media screen and (max-width: 768px) {
            .container {
                padding: 15px;
                margin: 2em 15px;
            }
            h1 {
                font-size: 1.8em;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>耗子Linux面板</h1>
        <p>这是耗子Linux面板的OpenResty默认页面！</p>
        <p>当您看到此页面，说明无法在服务器上找到该域名对应的站点。</p>
    </div>
</body>
</html>
EOF

# 写入站点停止页
cat > ${openrestyPath}/html/stop.html << EOF
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>网站已停止 - 耗子Linux面板</title>
    <style>
        body {
            background-color: #f9f9f9;
            margin: 0;
            padding: 0;
        }
        .container {
            max-width: 800px;
            margin: 2em auto;
            background-color: #ffffff;
            padding: 20px;
            border-radius: 12px;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
        }
        h1 {
            font-size: 2.5em;
            margin-top: 0;
            margin-bottom: 20px;
            text-align: center;
            color: #333;
            border-bottom: 2px solid #ddd;
            padding-bottom: 0.5em;
        }
        p {
            color: #555;
            line-height: 1.8;
        }
        @media screen and (max-width: 768px) {
            .container {
                padding: 15px;
                margin: 2em 15px;
            }
            h1 {
                font-size: 1.8em;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>耗子Linux面板</h1>
        <p>该网站已被管理员停止访问！</p>
        <p>当您看到此页面，说明该网站已被服务器管理员停止对外访问。</p>
    </div>
</body>
</html>
EOF

# 处理文件权限
chmod -R 755 ${openrestyPath}
chmod -R 755 /www/wwwroot
chown -R www:www /www/wwwroot
chmod -R 644 /www/server/vhost

# 写入无php配置文件
echo "" > ${openrestyPath}/conf/enable-php-0.conf

# 自动为所有PHP版本创建配置文件
if [ -d "${setupPath}/server/php" ]; then
    cd ${setupPath}/server/php
    phpList=$(ls -l | grep ^d | awk '{print $NF}')
    for phpVersion in ${phpList}; do
        if [ -d "${setupPath}/server/php/${phpVersion}" ]; then
            # 写入PHP配置文件
            cat > ${openrestyPath}/conf/enable-php-${phpVersion}.conf << EOF
location ~ \.php$ {
    try_files \$uri =404;
    fastcgi_pass unix:/tmp/php-cgi-${phpVersion}.sock;
    fastcgi_index index.php;
    include fastcgi.conf;
    include pathinfo.conf;
}
EOF
        fi
    done
fi

# 写入代理默认配置文件
cat > ${openrestyPath}/conf/proxy.conf << EOF
proxy_temp_path ${openrestyPath}/proxy_temp_dir;
proxy_cache_path ${openrestyPath}/proxy_cache_dir levels=1:2 keys_zone=cache_one:20m inactive=1d max_size=5g;
proxy_connect_timeout 60;
proxy_read_timeout 60;
proxy_send_timeout 60;
proxy_buffer_size 32k;
proxy_buffers 4 64k;
proxy_busy_buffers_size 128k;
proxy_temp_file_write_size 128k;
proxy_next_upstream error timeout invalid_header http_500 http_503 http_404;
proxy_cache cache_one;
EOF

# 写入默认站点配置文件
cat > ${openrestyPath}/conf/default.conf << EOF
server
{
    listen 80 default_server;
    listen 443 ssl default_server;
    server_name _;
    index index.html;
    root /www/server/openresty/html;
    ssl_reject_handshake on;
}
EOF

# 建立日志目录
mkdir -p /www/wwwlogs/waf
chown www:www /www/wwwlogs/waf
chmod 755 /www/wwwlogs/waf

# 写入服务文件
cat > /etc/systemd/system/openresty.service << EOF
[Unit]
Description=The OpenResty Application Platform
After=syslog.target network-online.target remote-fs.target nss-lookup.target
Wants=network-online.target

[Service]
Type=forking
PIDFile=/www/server/openresty/nginx.pid
ExecStartPre=/www/server/openresty/sbin/nginx -t -c /www/server/openresty/conf/nginx.conf
ExecStart=/www/server/openresty/sbin/nginx -c /www/server/openresty/conf/nginx.conf
ExecReload=/www/server/openresty/sbin/nginx -s reload
ExecStop=/www/server/openresty/sbin/nginx -s quit

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable openresty.service
systemctl restart openresty.service

panel writePlugin openresty ${openrestyVersion}

echo -e "${HR}\nOpenResty 安装完成\n${HR}"
