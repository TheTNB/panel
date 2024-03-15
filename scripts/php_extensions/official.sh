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
OS=$(source /etc/os-release && { [[ "$ID" == "debian" ]] && echo "debian"; } || { [[ "$ID" == "centos" ]] || [[ "$ID" == "rhel" ]] || [[ "$ID" == "rocky" ]] || [[ "$ID" == "almalinux" ]] && echo "centos"; } || echo "unknown")

action="$1"      # 操作
phpVersion="$2" # PHP版本
extensionName="$3" # 扩展名称
addArgs="" # 附加参数

Install() {
    # 检查是否已经安装
    isInstall=$(cat /www/server/php/${phpVersion}/etc/php.ini | grep "^extension=${extensionName}$")
    if [ "${isInstall}" != "" ]; then
        echo -e $HR
        echo "PHP-${phpVersion} 已安装 ${extensionName}"
        exit 1
    fi

    # 安装依赖
    if [ "${extensionName}" == "snmp" ]; then
        if [ "${OS}" == "centos" ]; then
            dnf install -y net-snmp-devel
        elif [ "${OS}" == "debian" ]; then
            apt-get install -y libsnmp-dev
        fi
    fi
    if [ "${extensionName}" == "ldap" ]; then
        if [ "${OS}" == "centos" ]; then
            dnf install -y openldap-devel
            ln -sf /usr/lib64/libldap* /usr/lib
        elif [ "${OS}" == "debian" ]; then
            apt-get install -y libldap2-dev
            ln -sf /usr/lib/x86_64-linux-gnu/libldap* /usr/lib
        fi
    fi
    if [ "${extensionName}" == "imap" ]; then
        if [ "${OS}" == "centos" ]; then
            # RHEL 9 的仓库中没有 libc-client-devel，待考虑
            dnf install -y libc-client-devel
        elif [ "${OS}" == "debian" ]; then
            apt-get install -y libc-client-dev
        fi
        addArgs="--with-imap --with-imap-ssl --with-kerberos"
    fi
    if [ "${extensionName}" == "enchant" ]; then
        if [ "${OS}" == "centos" ]; then
            dnf install -y enchant-devel
        elif [ "${OS}" == "debian" ]; then
            apt-get install -y libenchant-2-dev
        fi
    fi
    if [ "${extensionName}" == "pspell" ]; then
        if [ "${OS}" == "centos" ]; then
            dnf install -y aspell-devel
        elif [ "${OS}" == "debian" ]; then
            apt-get install -y libpspell-dev
        fi
    fi
    if [ "${extensionName}" == "gmp" ]; then
        if [ "${OS}" == "centos" ]; then
            dnf install -y gmp-devel
        elif [ "${OS}" == "debian" ]; then
            apt-get install -y libgmp-dev
        fi
    fi
    if [ "${extensionName}" == "gettext" ]; then
        if [ "${OS}" == "centos" ]; then
            dnf install -y gettext-devel
        elif [ "${OS}" == "debian" ]; then
            apt-get install -y libgettextpo-dev
        fi
    fi
    if [ "${extensionName}" == "bz2" ]; then
        if [ "${OS}" == "centos" ]; then
            dnf install -y bzip2-devel
        elif [ "${OS}" == "debian" ]; then
            apt-get install -y libbz2-dev
        fi
    fi
    if [ "${extensionName}" == "zip" ]; then
        if [ "${OS}" == "centos" ]; then
            dnf install -y libzip-devel
        elif [ "${OS}" == "debian" ]; then
            apt-get install -y libzip-dev
        fi
    fi
    if [ "${extensionName}" == "pdo_pgsql" ]; then
        addArgs="--with-pdo-pgsql=/www/server/postgresql"
    fi

    # 安装扩展
    if [ ! -d /www/server/php/${phpVersion}/src/ext/${extensionName} ]; then
        echo -e $HR
        echo "PHP-${phpVersion} ${extensionName} 源码不存在"
        exit 1
    fi
    cd /www/server/php/${phpVersion}/src/ext/${extensionName}
    /www/server/php/${phpVersion}/bin/phpize
    ./configure --with-php-config=/www/server/php/${phpVersion}/bin/php-config ${addArgs}
    make
    if [ "$?" != "0" ]; then
        echo -e $HR
        echo "PHP-${phpVersion} ${extensionName} 编译失败"
        exit 1
    fi
    make install
    if [ "$?" != "0" ]; then
        echo -e $HR
        echo "PHP-${phpVersion} ${extensionName} 安装失败"
        exit 1
    fi

    sed -i "/;haozi/a\extension=${extensionName}" /www/server/php/${phpVersion}/etc/php.ini

    # 重载PHP
    systemctl reload php-fpm-${phpVersion}.service
    echo -e $HR
    echo "PHP-${phpVersion} ${extensionName} 安装成功"
}

Uninstall() {
    # 检查是否已经安装
    isInstall=$(cat /www/server/php/${phpVersion}/etc/php.ini | grep "^extension=${extensionName}$")
    if [ "${isInstall}" == "" ]; then
        echo -e $HR
        echo "PHP-${phpVersion} 未安装 ${extensionName}"
        exit 1
    fi

    sed -i "/extension=${extensionName}/d" /www/server/php/${phpVersion}/etc/php.ini

    # 重载PHP
    systemctl reload php-fpm-${phpVersion}.service
    echo -e $HR
    echo "PHP-${phpVersion} ${extensionName} 卸载成功"
}

if [ "$action" == 'install' ]; then
    Install
fi
if [ "$action" == 'uninstall' ]; then
    Uninstall
fi
