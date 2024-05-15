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
downloadUrl="https://dl.cdn.haozi.net/panel/phpmyadmin"
setupPath="/www"
phpmyadminPath="${setupPath}/server/phpmyadmin"
phpmyadminVersion="5.2.1"
randomDir="$(cat /dev/urandom | head -n 16 | md5sum | head -c 10)"

# 准备安装目录
rm -rf ${phpmyadminPath}
mkdir -p ${phpmyadminPath}
cd ${phpmyadminPath}

wget -T 60 -t 3 -O phpMyAdmin-${phpmyadminVersion}-all-languages.zip ${downloadUrl}/phpMyAdmin-${phpmyadminVersion}-all-languages.zip
wget -T 20 -t 3 -O phpMyAdmin-${phpmyadminVersion}-all-languages.zip.checksum.txt ${downloadUrl}/phpMyAdmin-${phpmyadminVersion}-all-languages.zip.checksum.txt

if ! sha256sum --status -c phpMyAdmin-${phpmyadminVersion}-all-languages.zip.checksum.txt; then
    echo -e $HR
    echo "错误：phpMyAdmin 源码 checksum 校验失败，文件可能被篡改或不完整，已终止操作"
    rm -rf ${phpmyadminPath}
    exit 1
fi

unzip -o phpMyAdmin-${phpmyadminVersion}-all-languages.zip
mv phpMyAdmin-${phpmyadminVersion}-all-languages phpmyadmin_${randomDir}
chown -R www:www ${phpmyadminPath}
chmod -R 755 ${phpmyadminPath}
rm -rf phpMyAdmin-${phpmyadminVersion}-all-languages.zip
rm -rf phpMyAdmin-${phpmyadminVersion}-all-languages.zip.checksum.txt

# 判断PHP版本
phpVersion=""
if [ -d "/www/server/php/74" ]; then
    phpVersion="74"
fi
if [ -d "/www/server/php/80" ]; then
    phpVersion="80"
fi
if [ -d "/www/server/php/81" ]; then
    phpVersion="81"
fi
if [ -d "/www/server/php/82" ]; then
    phpVersion="82"
fi

if [ "${phpVersion}" == "" ]; then
    echo -e $HR
    echo "错误：未安装 PHP"
    rm -rf ${phpmyadminPath}
    exit 1
fi

# 写入 phpMyAdmin 配置文件
cat > /www/server/vhost/phpmyadmin.conf << EOF
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
    root /www/server/phpmyadmin;
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
chown -R root:root /www/server/vhost/phpmyadmin.conf
chmod -R 644 /www/server/vhost/phpmyadmin.conf
chmod -R 755 ${phpmyadminPath}
chown -R www:www ${phpmyadminPath}

# 放行端口
if [ "${OS}" == "centos" ]; then
    firewall-cmd --permanent --zone=public --add-port=888/tcp > /dev/null 2>&1
    firewall-cmd --reload
elif [ "${OS}" == "debian" ]; then
    ufw allow 888/tcp > /dev/null 2>&1
    ufw reload
fi

panel writePlugin phpmyadmin 5.2.1
systemctl reload openresty

echo -e "${HR}\phpMyAdmin 安装完成\n${HR}"
