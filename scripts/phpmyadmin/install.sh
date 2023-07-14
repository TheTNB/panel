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
downloadUrl="https://dl.cdn.haozi.net/panel/phpmyadmin"
setupPath="/www"
phpmyadminPath="${setupPath}/wwwroot/phpmyadmin"
phpmyadminVersion="5.2.1"
randomDir="$(cat /dev/urandom | head -n 16 | md5sum | head -c 10)"

# 准备安装目录
rm -rf ${phpmyadminPath}
mkdir -p ${phpmyadminPath}
cd ${phpmyadminPath}

wget -O phpmyadmin.zip ${downloadUrl}/phpMyAdmin-${phpmyadminVersion}-all-languages.zip
if [ "$?" != "0" ]; then
    echo -e $HR
    echo "错误：phpMyAdmin 下载失败"
    rm -rf ${phpmyadminPath}
    exit 1
fi
unzip -o phpmyadmin.zip
mv phpMyAdmin-${phpmyadminVersion}-all-languages phpmyadmin_${randomDir}
chown -R www:www ${phpmyadminPath}
chmod -R 755 ${phpmyadminPath}
rm -rf phpmyadmin.zip

# 判断PHP版本
phpVersion="74"
if [ -d "/www/server/php/80" ]; then
    phpVersion="80"
fi
if [ -d "/www/server/php/81" ]; then
    phpVersion="81"
fi
if [ -d "/www/server/php/82" ]; then
    phpVersion="82"
fi

# 写入 phpMyAdmin 配置文件
cat >/www/server/vhost/openresty/phpmyadmin.conf <<EOF
# 配置文件中的标记位请勿随意修改，改错将导致面板无法识别！
# 有自定义配置需求的，请将自定义的配置写在各标记位下方。
server
{
    # port标记位开始
    listen 888;
    # port标记位结束
    # server_name标记位开始
    server_name phpmyadmin;
    # server_name标记位结束
    # index标记位开始
    index index.php;
    # index标记位结束
    # root标记位开始
    root /www/wwwroot/phpmyadmin;
    # root标记位结束

    # php标记位开始
    include enable-php-${phpVersion}.conf;
    # php标记位结束

    # 面板默认禁止访问部分敏感目录，可自行修改
    location ~ ^/(\.user.ini|\.htaccess|\.git|\.svn)
    {
        return 404;
    }
    location ~ /tmp/ {
        return 403;
    }
    # 面板默认不记录静态资源的访问日志并开启1小时浏览器缓存，可自行修改
    location ~ .*\.(js|css)$
    {
        expires 1h;
        error_log /dev/null;
        access_log /dev/null;
    }

    access_log /www/wwwlogs/phpmyadmin.log;
    error_log /www/wwwlogs/phpmyadmin.log;
}
EOF
# 设置文件权限
chown -R root:root /www/server/vhost/openresty/phpmyadmin.conf
chmod -R 644 /www/server/vhost/openresty/phpmyadmin.conf

# 放行端口
firewall-cmd --permanent --zone=public --add-port=888/tcp >/dev/null 2>&1
firewall-cmd --reload

panel writePlugin phpmyadmin
systemctl reload openresty
