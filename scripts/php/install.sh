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
downloadUrl="https://jihulab.com/haozi-team/download/-/raw/main/panel/php"
setupPath="/www"
phpVersion="${1}"
phpVersionCode=""
phpPath="${setupPath}/server/php/${phpVersion}"
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
    dnf install autoconf glibc-headers gdbm-devel gd gd-devel perl oniguruma-devel libsodium-devel libxml2-devel sqlite-devel libzip-devel bzip2-devel xz-devel libpng-devel libjpeg-devel libwebp-devel libavif-devel freetype-devel gmp-devel openssl-devel readline-devel libxslt-devel libcurl-devel pkgconfig libedit-devel zlib-devel pcre-devel crontabs libicu libicu-devel c-ares -y
elif [ "${OS}" == "debian" ]; then
    apt-get update
    apt-get install build-essential autoconf libc6-dev libgdbm-dev libgd-tools libgd-dev perl libonig-dev libsodium-dev libxml2-dev libsqlite3-dev libzip-dev libbz2-dev liblzma-dev libpng-dev libjpeg-dev libwebp-dev libavif-dev libfreetype6-dev libgmp-dev libssl-dev libreadline-dev libxslt1-dev libcurl4-openssl-dev pkg-config libedit-dev zlib1g-dev libpcre3-dev cron libicu-dev libc-ares2 libc-ares-dev -y
else
    echo -e $HR
    echo "错误：耗子Linux面板不支持该系统"
    exit 1
fi

# 准备安装目录
rm -rf ${phpPath}
mkdir -p ${phpPath}
cd ${phpPath}

# 下载源码
if [ "${phpVersion}" == "74" ]; then
    phpVersionCode="7.4.33"
elif [ "${phpVersion}" == "80" ]; then
    phpVersionCode="8.0.30"
elif [ "${phpVersion}" == "81" ]; then
    phpVersionCode="8.1.23"
elif [ "${phpVersion}" == "82" ]; then
    phpVersionCode="8.2.10"
else
    echo -e $HR
    echo "错误：PHP-${phpVersion}不支持，请检查版本号是否正确。"
    exit 1
fi

wget -T 120 -t 3 -O ${phpPath}/php-${phpVersionCode}.tar.gz ${downloadUrl}/php-${phpVersionCode}.tar.gz
wget -T 20 -t 3 -O ${phpPath}/php-${phpVersionCode}.tar.gz.checksum.txt ${downloadUrl}/php-${phpVersionCode}.tar.gz.checksum.txt

if ! sha256sum --status -c php-${phpVersionCode}.tar.gz.checksum.txt; then
    echo -e $HR
    echo "错误：PHP-${phpVersion}源码 checksum 校验失败，文件可能被篡改或不完整，已终止操作"
    rm -rf ${phpPath}
    exit 1
fi

tar -xvf php-${phpVersionCode}.tar.gz
rm -f php-${phpVersionCode}.tar.gz
rm -f php-${phpVersionCode}.tar.gz.checksum.txt
mv php-* src

if [ "${phpVersion}" -le "80" ]; then
    wget -T 120 -t 3 -O ${phpPath}/openssl-1.1.1u.tar.gz ${downloadUrl}/openssl/openssl-1.1.1u.tar.gz
    wget -T 20 -t 3 -O ${phpPath}/openssl-1.1.1u.tar.gz.checksum.txt ${downloadUrl}/openssl/openssl-1.1.1u.tar.gz.checksum.txt

    if ! sha256sum --status -c openssl-1.1.1u.tar.gz.checksum.txt; then
        echo -e $HR
        echo "错误：PHP-${phpVersion} OpenSSL 源码 checksum 校验失败，文件可能被篡改或不完整，已终止操作"
        rm -rf ${phpPath}
        exit 1
    fi

    tar -zxvf openssl-1.1.1u.tar.gz
    rm -f openssl-1.1.1u.tar.gz
    rm -f openssl-1.1.1u.tar.gz.checksum.txt
    mv openssl-1.1.1u openssl
    cd openssl
    ./config --prefix=/usr/local/openssl-1.1 --openssldir=/usr/local/openssl-1.1
    make -j$(nproc)
    make install
    echo "/usr/local/openssl-1.1/lib" > /etc/ld.so.conf.d/openssl-1.1.conf
    ldconfig
    cd ..
    rm -rf openssl

    export CFLAGS="-I/usr/local/openssl-1.1/include -I/usr/local/curl/include"
    export LIBS="-L/usr/local/openssl-1.1/lib -L/usr/local/curl/lib"
fi

# 配置
cd src
if [ "${phpVersion}" == "81" ] || [ "${phpVersion}" == "82" ]; then
    ./configure --prefix=${phpPath} --with-config-file-path=${phpPath}/etc --enable-fpm --with-fpm-user=www --with-fpm-group=www --enable-mysqlnd --with-mysqli=mysqlnd --with-pdo-mysql=mysqlnd --with-freetype --with-jpeg --with-zlib --with-libxml-dir=/usr --enable-xml --disable-rpath --enable-bcmath --enable-shmop --enable-sysvsem --enable-inline-optimization --with-curl --enable-mbregex --enable-mbstring --enable-intl --enable-pcntl --enable-ftp --enable-gd --with-openssl --with-mhash --enable-pcntl --enable-sockets --with-xmlrpc --enable-soap --with-gettext --enable-fileinfo --enable-opcache --with-sodium --with-webp --with-avif
else
    ./configure --prefix=${phpPath} --with-config-file-path=${phpPath}/etc --enable-fpm --with-fpm-user=www --with-fpm-group=www --enable-mysqlnd --with-mysqli=mysqlnd --with-pdo-mysql=mysqlnd --with-freetype --with-jpeg --with-zlib --with-libxml-dir=/usr --enable-xml --disable-rpath --enable-bcmath --enable-shmop --enable-sysvsem --enable-inline-optimization --with-curl --enable-mbregex --enable-mbstring --enable-intl --enable-pcntl --enable-ftp --enable-gd --with-openssl --with-mhash --enable-pcntl --enable-sockets --with-xmlrpc --enable-soap --with-gettext --enable-fileinfo --enable-opcache --with-sodium --with-webp
fi

# 编译安装
if [[ "${cpuCore}" -gt "1" ]]; then
    make -j2
else
    make
fi
make install
if [ ! -f "${phpPath}/bin/php" ]; then
    echo -e $HR
    echo "错误：PHP-${phpVersion}安装失败，请截图错误信息寻求帮助！"
    rm -rf ${phpPath}
    exit 1
fi

# 创建php配置
mkdir -p ${phpPath}/etc
\cp php.ini-production ${phpPath}/etc/php.ini

# 安装zip拓展
cd ${phpPath}/src/ext/zip
${phpPath}/bin/phpize
./configure --with-php-config=${phpPath}/bin/php-config
make -j${cpuCore}
make install
if [ "$?" != "0" ]; then
    echo -e $HR
    echo "错误：PHP-${phpVersion} zip拓展安装失败，请截图错误信息寻求帮助。"
    exit 1
fi
cd ../../

# 写入拓展标记位
echo ";下方标记位禁止删除，否则将导致PHP拓展无法正常安装！" >> ${phpPath}/etc/php.ini
echo ";haozi" >> ${phpPath}/etc/php.ini
# 写入zip拓展到php配置
echo "extension=zip" >> ${phpPath}/etc/php.ini

# 设置软链接
rm -f /usr/bin/php-${phpVersion}
rm -f /usr/bin/pear
rm -f /usr/bin/pecl
ln -sf ${phpPath}/bin/php /usr/bin/php
ln -sf ${phpPath}/bin/php /usr/bin/php-${phpVersion}
ln -sf ${phpPath}/bin/phpize /usr/bin/phpize
ln -sf ${phpPath}/bin/pear /usr/bin/pear
ln -sf ${phpPath}/bin/pecl /usr/bin/pecl
ln -sf ${phpPath}/sbin/php-fpm /usr/bin/php-fpm-${phpVersion}

# 设置fpm
cat > ${phpPath}/etc/php-fpm.conf << EOF
[global]
pid = ${phpPath}/var/run/php-fpm.pid
error_log = ${phpPath}/var/log/php-fpm.log
log_level = notice

[www]
listen = /tmp/php-cgi-${phpVersion}.sock
listen.backlog = -1
listen.allowed_clients = 127.0.0.1
listen.owner = www
listen.group = www
listen.mode = 0666
user = www
group = www
pm = dynamic
pm.max_children = 30
pm.start_servers = 5
pm.min_spare_servers = 5
pm.max_spare_servers = 10
request_terminate_timeout = 100
request_slowlog_timeout = 30
pm.status_path = /phpfpm_status/${phpVersion}
slowlog = var/log/slow.log
EOF

# 设置PHP进程数
memTotal=$(free -m | grep Mem | awk '{print  $2}')
if [[ ${memTotal} -gt 1024 && ${memTotal} -le 2048 ]]; then
    sed -i "s#pm.max_children.*#pm.max_children = 50#" ${phpPath}/etc/php-fpm.conf
    sed -i "s#pm.start_servers.*#pm.start_servers = 5#" ${phpPath}/etc/php-fpm.conf
    sed -i "s#pm.min_spare_servers.*#pm.min_spare_servers = 5#" ${phpPath}/etc/php-fpm.conf
    sed -i "s#pm.max_spare_servers.*#pm.max_spare_servers = 10#" ${phpPath}/etc/php-fpm.conf
elif [[ ${memTotal} -gt 2048 && ${memTotal} -le 4096 ]]; then
    sed -i "s#pm.max_children.*#pm.max_children = 80#" ${phpPath}/etc/php-fpm.conf
    sed -i "s#pm.start_servers.*#pm.start_servers = 5#" ${phpPath}/etc/php-fpm.conf
    sed -i "s#pm.min_spare_servers.*#pm.min_spare_servers = 5#" ${phpPath}/etc/php-fpm.conf
    sed -i "s#pm.max_spare_servers.*#pm.max_spare_servers = 20#" ${phpPath}/etc/php-fpm.conf
elif [[ ${memTotal} -gt 4096 && ${memTotal} -le 8192 ]]; then
    sed -i "s#pm.max_children.*#pm.max_children = 150#" ${phpPath}/etc/php-fpm.conf
    sed -i "s#pm.start_servers.*#pm.start_servers = 10#" ${phpPath}/etc/php-fpm.conf
    sed -i "s#pm.min_spare_servers.*#pm.min_spare_servers = 10#" ${phpPath}/etc/php-fpm.conf
    sed -i "s#pm.max_spare_servers.*#pm.max_spare_servers = 30#" ${phpPath}/etc/php-fpm.conf
elif [[ ${memTotal} -gt 8192 && ${memTotal} -le 16384 ]]; then
    sed -i "s#pm.max_children.*#pm.max_children = 200#" ${phpPath}/etc/php-fpm.conf
    sed -i "s#pm.start_servers.*#pm.start_servers = 15#" ${phpPath}/etc/php-fpm.conf
    sed -i "s#pm.min_spare_servers.*#pm.min_spare_servers = 15#" ${phpPath}/etc/php-fpm.conf
    sed -i "s#pm.max_spare_servers.*#pm.max_spare_servers = 30#" ${phpPath}/etc/php-fpm.conf
elif [[ ${memTotal} -gt 16384 ]]; then
    sed -i "s#pm.max_children.*#pm.max_children = 300#" ${phpPath}/etc/php-fpm.conf
    sed -i "s#pm.start_servers.*#pm.start_servers = 20#" ${phpPath}/etc/php-fpm.conf
    sed -i "s#pm.min_spare_servers.*#pm.min_spare_servers = 20#" ${phpPath}/etc/php-fpm.conf
    sed -i "s#pm.max_spare_servers.*#pm.max_spare_servers = 50#" ${phpPath}/etc/php-fpm.conf
fi
sed -i "s#listen.backlog.*#listen.backlog = 8192#" ${phpPath}/etc/php-fpm.conf
# 最大上传限制100M
sed -i 's/post_max_size =.*/post_max_size = 100M/g' ${phpPath}/etc/php.ini
sed -i 's/upload_max_filesize =.*/upload_max_filesize = 100M/g' ${phpPath}/etc/php.ini
# 时区PRC
sed -i 's/;date.timezone =.*/date.timezone = PRC/g' ${phpPath}/etc/php.ini
sed -i 's/short_open_tag =.*/short_open_tag = On/g' ${phpPath}/etc/php.ini
sed -i 's/;cgi.fix_pathinfo=.*/cgi.fix_pathinfo=1/g' ${phpPath}/etc/php.ini
# 最大运行时间
sed -i 's/max_execution_time =.*/max_execution_time = 86400/g' ${phpPath}/etc/php.ini
sed -i 's/;sendmail_path =.*/sendmail_path = \/usr\/sbin\/sendmail -t -i/g' ${phpPath}/etc/php.ini
# 禁用函数
sed -i 's/disable_functions =.*/disable_functions = passthru,exec,system,putenv,chroot,chgrp,chown,shell_exec,popen,proc_open,pcntl_exec,ini_alter,ini_restore,dl,openlog,syslog,readlink,symlink,popepassthru,pcntl_alarm,pcntl_fork,pcntl_waitpid,pcntl_wait,pcntl_wifexited,pcntl_wifstopped,pcntl_wifsignaled,pcntl_wifcontinued,pcntl_wexitstatus,pcntl_wtermsig,pcntl_wstopsig,pcntl_signal,pcntl_signal_dispatch,pcntl_get_last_error,pcntl_strerror,pcntl_sigprocmask,pcntl_sigwaitinfo,pcntl_sigtimedwait,pcntl_exec,pcntl_getpriority,pcntl_setpriority,imap_open,apache_setenv/g' ${phpPath}/etc/php.ini
sed -i 's/display_errors = Off/display_errors = On/g' ${phpPath}/etc/php.ini
sed -i 's/error_reporting =.*/error_reporting = E_ALL \& \~E_NOTICE/g' ${phpPath}/etc/php.ini

# 设置SSL根证书
#sed -i "s#;openssl.cafile=#openssl.cafile=/etc/pki/tls/certs/ca-bundle.crt#" ${phpPath}/etc/php.ini
#sed -i "s#;curl.cainfo =#curl.cainfo = /etc/pki/tls/certs/ca-bundle.crt#" ${phpPath}/etc/php.ini

# 关闭php外显
sed -i 's/expose_php = On/expose_php = Off/g' ${phpPath}/etc/php.ini

# 写入openresty 调用php配置文件
cat > /www/server/openresty/conf/enable-php-${phpVersion}.conf << EOF
location ~ \.php$ {
    try_files \$uri =404;
    fastcgi_pass unix:/tmp/php-cgi-${phpVersion}.sock;
    fastcgi_index index.php;
    include fastcgi.conf;
    include pathinfo.conf;
}
EOF

# 添加php-fpm到服务
\cp ${phpPath}/src/sapi/fpm/php-fpm.service /lib/systemd/system/php-fpm-${phpVersion}.service
sed -i "/PrivateTmp/d" /lib/systemd/system/php-fpm-${phpVersion}.service
systemctl daemon-reload

# 启动php
systemctl enable php-fpm-${phpVersion}.service
systemctl start php-fpm-${phpVersion}.service

panel writePlugin php${phpVersion} ${phpVersionCode}

echo -e "${HR}\nPHP-${phpVersion} 安装完成\n${HR}"
