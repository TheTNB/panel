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

downloadUrl="https://jihulab.com/haozi-team/download/-/raw/main/panel/php_extensions"
action="$1"
phpVersion="$2"

Install() {
    # 检查是否已经安装
    isInstall=$(cat /www/server/php/${phpVersion}/etc/php.ini | grep 'ioncube_loader_lin')
    if [ "${isInstall}" != "" ]; then
        echo -e $HR
        echo "PHP-${phpVersion} 已安装 ionCube"
        exit 1
    fi

    mkdir /usr/local/ioncube
    wget -T 60 -t 3 -O /usr/local/ioncube/ioncube_loader_lin_${phpVersion}.so ${downloadUrl}/ioncube_loader_lin_${phpVersion}.so
    wget -T 20 -t 3 -O /usr/local/ioncube/ioncube_loader_lin_${phpVersion}.so.checksum.txt ${downloadUrl}/ioncube_loader_lin_${phpVersion}.so.checksum.txt

    if ! sha256sum --status -c /usr/local/ioncube/ioncube_loader_lin_${phpVersion}.so.checksum.txt; then
        echo -e $HR
        echo "错误：PHP-${phpVersion} ionCube 源码 checksum 校验失败，文件可能被篡改或不完整，已终止操作"
        exit 1
    fi

    rm -f /usr/local/ioncube/ioncube_loader_lin_${phpVersion}.so.checksum.txt

    sed -i -e "/;haozi/a\zend_extension=/usr/local/ioncube/ioncube_loader_lin_${phpVersion}.so" /www/server/php/${phpVersion}/etc/php.ini

    # 重载PHP
    systemctl reload php-fpm-${phpVersion}.service
    echo -e $HR
    echo "PHP-${phpVersion} ionCube 安装成功"
}

Uninstall() {
    # 检查是否已经安装
    isInstall=$(cat /www/server/php/${phpVersion}/etc/php.ini | grep 'ioncube_loader_lin')
    if [ "${isInstall}" == "" ]; then
        echo -e $HR
        echo "PHP-${phpVersion} 未安装 ionCube"
        exit 1
    fi

    rm -f /usr/local/ioncube/ioncube_loader_lin_${phpVersion}.so
    sed -i '/ioncube_loader_lin/d' /www/server/php/${phpVersion}/etc/php.ini

    # 重载PHP
    systemctl reload php-fpm-${phpVersion}.service
    echo -e $HR
    echo "PHP-${phpVersion} ionCube 卸载成功"
}

if [ "$action" == 'install' ]; then
    Install
fi
if [ "$action" == 'uninstall' ]; then
    Uninstall
fi
