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
OS=$(source /etc/os-release && { [[ "$ID" == "debian" ]] && echo "debian"; } || { [[ "$ID" == "centos" ]] || [[ "$ID" == "rhel" ]] || [[ "$ID" == "rocky" ]] || [[ "$ID" == "almalinux" ]] && echo "centos"; } || echo "unknown")

downloadUrl="https://jihulab.com/haozi-team/download/-/raw/main/panel/php_extensions"
action="$1"
phpVersion="$2"
imagickVersion="3.7.0"

Install() {
    # 检查是否已经安装
    isInstall=$(cat /www/server/php/${phpVersion}/etc/php.ini | grep '^extension=imagick$')
    if [ "${isInstall}" != "" ]; then
        echo -e $HR
        echo "PHP-${phpVersion} 已安装 imagick"
        exit 1
    fi

    # 安装依赖
    if [ "${OS}" == "centos" ]; then
        dnf install ImageMagick ImageMagick-devel -y
    elif [ "${OS}" == "debian" ]; then
        apt-get install imagemagick libmagickwand-dev -y
    else
        echo -e $HR
        echo "错误：耗子Linux面板不支持该系统"
        exit 1
    fi

    cd /www/server/php/${phpVersion}/src/ext
    rm -rf imagick
    rm -rf imagick.tar.gz
    wget -T 60 -t 3 -O imagick.tar.gz ${downloadUrl}/imagick-${imagickVersion}.tar.gz
    tar -zxvf imagick.tar.gz
    mv imagick-${imagickVersion} imagick
    cd imagick
    /www/server/php/${phpVersion}/bin/phpize
    ./configure --with-php-config=/www/server/php/${phpVersion}/bin/php-config
    make
    if [ "$?" != "0" ]; then
        echo -e $HR
        echo "PHP-${phpVersion} imagick 编译失败"
        exit 1
    fi
    make install
    if [ "$?" != "0" ]; then
        echo -e $HR
        echo "PHP-${phpVersion} imagick 安装失败"
        exit 1
    fi

    sed -i '/;haozi/a\extension=imagick' /www/server/php/${phpVersion}/etc/php.ini

    # 重载PHP
    systemctl reload php-fpm-${phpVersion}.service
    echo -e $HR
    echo "PHP-${phpVersion} imagick 安装成功"
}

Uninstall() {
    # 检查是否已经安装
    isInstall=$(cat /www/server/php/${phpVersion}/etc/php.ini | grep '^extension=imagick$')
    if [ "${isInstall}" == "" ]; then
        echo -e $HR
        echo "PHP-${phpVersion} 未安装 imagick"
        exit 1
    fi

    sed -i '/extension=imagick/d' /www/server/php/${phpVersion}/etc/php.ini

    # 重载PHP
    systemctl reload php-fpm-${phpVersion}.service
    echo -e $HR
    echo "PHP-${phpVersion} imagick 卸载成功"
}

if [ "$action" == 'install' ]; then
    Install
fi
if [ "$action" == 'uninstall' ]; then
    Uninstall
fi
