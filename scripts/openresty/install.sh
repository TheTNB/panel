#!/bin/bash
export PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:$PATH

: '
Copyright 2022 HaoZi Technology Co., Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
'

HR="+----------------------------------------------------"
ARCH=$(uname -m)
OS=$(source /etc/os-release && { [[ "$ID" == "debian" ]] && echo "debian"; } || { [[ "$ID" == "centos" ]] || [[ "$ID" == "rhel" ]] || [[ "$ID" == "rocky" ]] || [[ "$ID" == "almalinux" ]] && echo "centos"; } || echo "unknown")
downloadUrl="https://dl.cdn.haozi.net/panel/openresty"
setupPath="/www"
openrestyPath="${setupPath}/server/openresty"
openrestyVersion="1.21.4.1"
cpuCore=$(cat /proc/cpuinfo | grep "processor" | wc -l)

# 安装依赖
if [ "${OS}" == "centos" ]; then
    dnf install dnf-plugins-core -y
    dnf install epel-release -y
    dnf config-manager --set-enabled epel
    dnf config-manager --set-enabled PowerTools
    dnf config-manager --set-enabled powertools
    dnf config-manager --set-enabled CRB
    dnf config-manager --set-enabled Crb
    dnf config-manager --set-enabled crb
    /usr/bin/crb enable
    dnf makecache
    dnf groupinstall "Development Tools" -y
    dnf install tar unzip gd gd-devel git-core flex perl oniguruma oniguruma-devel libsodium-devel libxml2-devel libxslt-devel GeoIP-devel bison yajl yajl-devel curl curl-devel libtermcap-devel ncurses-devel libevent-devel readline-devel libuuid-devel brotli-devel icu libicu libicu-devel openssl openssl-devel -y
elif [ "${OS}" == "debian" ]; then
    apt update
    apt install build-essential tar unzip libgd3 libgd-dev git flex perl libonig-dev libsodium-dev libxml2-dev libxslt1-dev libgeoip-dev bison libyajl-dev curl libcurl4-openssl-dev libncurses5-dev libevent-dev libreadline-dev uuid-dev libbrotli-dev icu-devtools libicu-dev openssl libssl-dev -y
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
wget -T 120 -O ${openrestyPath}/openresty-${openrestyVersion}.tar.gz ${downloadUrl}/openresty-${openrestyVersion}.tar.gz
tar -zxvf openresty-${openrestyVersion}.tar.gz
rm -f openresty-${openrestyVersion}.tar.gz
mv openresty-${openrestyVersion} src
cd src

# openssl
wget -T 120 -O openssl.tar.gz ${downloadUrl}/openssl/openssl-1.1.1u.tar.gz
tar -zxvf openssl.tar.gz
rm -f openssl.tar.gz
mv openssl-1.1.1u openssl
rm -f openssl.tar.gz

# pcre
wget -T 60 -O pcre-8.45.tar.gz ${downloadUrl}/pcre/pcre-8.45.tar.gz
tar -zxvf pcre-8.45.tar.gz
rm -f pcre-8.45.tar.gz
mv pcre-8.45 pcre
rm -f pcre-8.45.tar.gz

# ngx_cache_purge
wget -T 20 -O ngx_cache_purge.tar.gz ${downloadUrl}/modules/ngx_cache_purge-2.3.tar.gz
tar -zxvf ngx_cache_purge.tar.gz
rm -f ngx_cache_purge.tar.gz
mv ngx_cache_purge-2.3 ngx_cache_purge
rm -f ngx_cache_purge.tar.gz

# nginx-sticky-module
wget -T 20 -O nginx-sticky-module.zip ${downloadUrl}/modules/nginx-sticky-module.zip
unzip -o nginx-sticky-module.zip
rm -f nginx-sticky-module.zip

# nginx-dav-ext-module
wget -T 20 -O nginx-dav-ext-module-3.0.0.tar.gz ${downloadUrl}/modules/nginx-dav-ext-module-3.0.0.tar.gz
tar -xvf nginx-dav-ext-module-3.0.0.tar.gz
rm -f nginx-dav-ext-module-3.0.0.tar.gz
mv nginx-dav-ext-module-3.0.0 nginx-dav-ext-module

# waf
wget -T 60 -O uthash.zip ${downloadUrl}/modules/uthash-2.3.0.zip
unzip -o uthash.zip
mv uthash-2.3.0 uthash
rm -f uthash.zip
cd ../
wget -T 20 -O ngx_waf.zip ${downloadUrl}/modules/ngx_waf-6.1.9.zip
unzip -o ngx_waf.zip
mv ngx_waf-6.1.9 ngx_waf
rm -f ngx_waf.zip
cd ngx_waf/inc
wget -T 60 -O libinjection.zip ${downloadUrl}/modules/libinjection-3.10.0.zip
unzip -o libinjection.zip
mv libinjection-3.10.0 libinjection
rm -f libinjection.zip
cd ../
make -j$(nproc)
if [ "$?" != "0" ]; then
    echo -e $HR
    echo "错误：OpenResty waf拓展初始化失败，请截图错误信息寻求帮助。"
    rm -rf ${openrestyPath}
    exit 1
fi
cd ${openrestyPath}/src

# brotli
wget -T 20 -O ngx_brotli.zip ${downloadUrl}/modules/ngx_brotli-1.0.0rc.zip
unzip -o ngx_brotli.zip
mv ngx_brotli-1.0.0rc ngx_brotli
rm -f ngx_brotli.zip
cd ngx_brotli/deps
rm -rf brotli
wget -T 20 -O brotli.zip ${downloadUrl}/modules/brotli-1.0.9.zip
unzip -o brotli.zip
mv brotli-1.0.9 brotli
rm -f brotli.zip
cd ${openrestyPath}/src

cd ${openrestyPath}/src
export LD_LIBRARY_PATH=/usr/local/lib/:$LD_LIBRARY_PATH
export LIB_UTHASH=${openrestyPath}/src/uthash

./configure --user=www --group=www --prefix=${openrestyPath} --with-luajit --add-module=${openrestyPath}/src/ngx_cache_purge --add-module=${openrestyPath}/src/nginx-sticky-module --with-openssl=${openrestyPath}/src/openssl --with-pcre=${openrestyPath}/src/pcre --with-http_v2_module --with-http_slice_module --with-threads --with-stream --with-stream_ssl_module --with-stream_realip_module --with-stream_geoip_module --with-stream_ssl_preread_module --with-http_stub_status_module --with-http_ssl_module --with-http_image_filter_module --with-http_gzip_static_module --with-http_gunzip_module --with-ipv6 --with-http_geoip_module --with-http_sub_module --with-http_flv_module --with-http_addition_module --with-http_realip_module --with-http_mp4_module --with-ld-opt="-Wl,-E" --with-cc-opt="-O2 -std=gnu99" --with-cpu-opt="amd64" --with-http_dav_module --add-module=${openrestyPath}/src/nginx-dav-ext-module --add-module=${openrestyPath}/src/ngx_brotli --add-module=${openrestyPath}/ngx_waf
if [[ "${cpuCore}" -gt "1" ]]; then
    make -j2
else
    make
fi
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
cat >${openrestyPath}/conf/nginx.conf <<EOF
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
cat >${openrestyPath}/conf/pathinfo.conf <<EOF
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
cat >${openrestyPath}/html/index.html <<EOF
<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="utf-8">
<title>耗子Linux面板</title>
</head>
<body>
<h1>耗子Linux面板</h1>
<p>这是耗子Linux面板的OpenResty默认页面！</p>
<p>当您看到此页面，说明该域名尚未与站点绑定。</p>
</body>
</html>
EOF

# 写入站点停止页
cat >${openrestyPath}/html/stop.html <<EOF
<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="utf-8">
<title>网站已停止 - 耗子Linux面板</title>
</head>
<body>
<h1>耗子Linux面板</h1>
<p>该网站已被管理员停止访问！</p>
<p>当您看到此页面，说明该网站已被管理员停止对外访问，请联系管理员了解详情。</p>
</body>
</html>
EOF

# 处理文件权限
chmod 755 ${openrestyPath}
chmod 644 ${openrestyPath}/html
chmod -R 755 /www/wwwroot
chown -R www:www /www/wwwroot
chmod -R 644 /www/server/vhost

# 写入无php配置文件
echo "" >${openrestyPath}/conf/enable-php-0.conf
# 写入代理默认配置文件
cat >${openrestyPath}/conf/proxy.conf <<EOF
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

# 建立日志目录
mkdir -p /www/wwwlogs/waf
chown www:www /www/wwwlogs/waf
chmod 755 /www/wwwlogs/waf

# 写入服务文件
cat >/etc/systemd/system/openresty.service <<EOF
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
systemctl start openresty.service

panel writePlugin openresty ${openrestyVersion}

echo -e "${HR}\nOpenResty 安装完成\n${HR}"
