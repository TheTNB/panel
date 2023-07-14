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

action="$1"      # 操作
phpVersion="$2" # PHP版本

Install() {
    # 检查是否已经安装
    isInstall=$(cat /www/server/php/${phpVersion}/etc/php.ini | grep '^zend_extension=opcache$')
    if [ "${isInstall}" != "" ]; then
        echo -e $HR
        echo "PHP-${phpVersion} 已安装 Zend OPcache"
        exit 1
    fi

    if [ "${phpVersion}" -ge "80" ]; then
        sed -i '/;haozi/a\zend_extension=opcache\nopcache.enable = 1\nopcache.enable_cli=1\nopcache.memory_consumption=128\nopcache.interned_strings_buffer=32\nopcache.max_accelerated_files=100000\nopcache.revalidate_freq=3\nopcache.save_comments=0\nopcache.jit_buffer_size=128m\nopcache.jit=1205' /www/server/php/${phpVersion}/etc/php.ini
    else
        sed -i '/;haozi/a\zend_extension=opcache\nopcache.enable = 1\nopcache.enable_cli=1\nopcache.memory_consumption=128\nopcache.interned_strings_buffer=32\nopcache.max_accelerated_files=100000\nopcache.revalidate_freq=3\nopcache.save_comments=0' /www/server/php/${phpVersion}/etc/php.ini
    fi
    # 重载PHP
    systemctl reload php-fpm-${phpVersion}.service
    echo -e $HR
    echo "PHP-${phpVersion} Zend OPcache 安装成功"
}

Uninstall() {
    # 检查是否已经安装
    isInstall=$(cat /www/server/php/${phpVersion}/etc/php.ini | grep '^zend_extension=opcache$')
    if [ "${isInstall}" == "" ]; then
        echo -e $HR
        echo "PHP-${phpVersion} 未安装 Zend OPcache"
        exit 1
    fi

    sed -i '/^opcache.*$/d' /www/server/php/${phpVersion}/etc/php.ini
    sed -i '/zend_extension=opcache/d' /www/server/php/${phpVersion}/etc/php.ini

    # 重载PHP
    systemctl reload php-fpm-${phpVersion}.service
    echo -e $HR
    echo "PHP-${phpVersion} Zend OPcache 卸载成功"
}

if [ "$action" == 'install' ]; then
    Install
fi
if [ "$action" == 'uninstall' ]; then
    Uninstall
fi
