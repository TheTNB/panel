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
if [ "${OS}" == "unknown" ]; then
    echo -e $HR
    echo "错误：该系统不支持安装耗子面板，请更换 Debian 12.x / RHEL 9.x 安装。"
    exit 1
fi

oldVersion=$(panel getSetting version)
oldVersion=${oldVersion#v}
panelPath="/www/panel"

# 大于
function version_gt() { test "$(echo -e "$1\n$2" | tr " " "\n" | sort -V | head -n 1)" != "$1"; }
# 小于
function version_lt() { test "$(echo -e "$1\n$2" | tr " " "\n" | sort -rV | head -n 1)" != "$1"; }
# 大于等于
function version_ge() { test "$(echo -e "$1\n$2" | tr " " "\n" | sort -rV | head -n 1)" == "$1"; }
# 小于等于
function version_le() { test "$(echo -e "$1\n$2" | tr " " "\n" | sort -V | head -n 1)" == "$1"; }

if [ -z "$oldVersion" ]; then
    if [ -f "$panelPath/database/panel.db" ]; then
        echo "DB_FILE=$panelPath/database/panel.db" >> $panelPath/panel.conf
        oldVersion=$(panel getSetting version)
        oldVersion=${oldVersion#v}
        sed -i '/DB_FILE/d' $panelPath/panel.conf
    else
        echo "错误：无法获取面板版本"
        echo "Error: can't get panel version"
        exit 1
    fi
fi

# 判断版本号是否合法
versionPattern="^[0-9]+\.[0-9]+\.[0-9]+$"
if [[ ! $oldVersion =~ $versionPattern ]]; then
    if [ -f "$panelPath/database/panel.db" ]; then
        echo "DB_FILE=$panelPath/database/panel.db" >> $panelPath/panel.conf
        oldVersion=$(panel getSetting version)
        oldVersion=${oldVersion#v}
        sed -i '/DB_FILE/d' $panelPath/panel.conf
    else
        echo "错误：面板版本号不合法"
        echo "Error: panel version is illegal"
        exit 1
    fi
fi

echo $HR

if version_lt "$oldVersion" "2.1.8"; then
    echo "更新面板到 v2.1.8 ..."
    echo "Update panel to v2.1.8 ..."
    oldEntrance=$(panel getSetting entrance)
    echo "APP_ENTRANCE=$oldEntrance" >> $panelPath/panel.conf
    panel deleteSetting entrance
fi

if version_lt "$oldVersion" "2.1.30"; then
    echo "更新面板到 v2.1.30 ..."
    echo "Update panel to v2.1.30 ..."
    sed -i '/APP_HOST/d' $panelPath/panel.conf
    echo "APP_SSL=false" >> $panelPath/panel.conf
    mv $panelPath/database/panel.db $panelPath/storage/panel.db
    openssl req -x509 -nodes -days 36500 -newkey ec:<(openssl ecparam -name secp384r1) -keyout $panelPath/storage/ssl.key -out $panelPath/storage/ssl.crt -subj "/C=CN/ST=Tianjin/L=Tianjin/O=HaoZi Technology Co., Ltd./OU=HaoZi Panel/CN=Panel"
fi

if version_lt "$oldVersion" "2.2.0"; then
    echo "更新面板到 v2.2.0 ..."
    echo "Update panel to v2.2.0 ..."
    echo "APP_LOCALE=zh_CN" >> $panelPath/panel.conf
fi

if version_lt "$oldVersion" "2.2.4"; then
    echo "更新面板到 v2.2.4 ..."
    echo "Update panel to v2.2.4 ..."
    if [ "${OS}" == "centos" ]; then
        dnf makecache
        dnf install -y p7zip p7zip-plugins rsyslog
        systemctl enable rsyslog
        systemctl start rsyslog
    else
        apt-get update -y
        apt-get install -y p7zip p7zip-full
    fi
fi

if version_lt "$oldVersion" "2.2.10"; then
    echo "更新面板到 v2.2.10 ..."
    echo "Update panel to v2.2.10 ..."
    if [ -f "/usr/bin/podman" ]; then
        panel writePlugin podman 4.0.0
        if [ "${OS}" == "debian" ]; then
            apt-get install containers-storage -y
            cp /usr/share/containers/storage.conf /etc/containers/storage.conf
        fi
        systemctl enable podman
        systemctl enable podman.socket
        systemctl enable podman-restart
        systemctl start podman
        systemctl start podman.socket
        systemctl start podman-restart
    fi
fi

if version_lt "$oldVersion" "2.2.14"; then
    echo "更新面板到 v2.2.14 ..."
    echo "Update panel to v2.2.14 ..."
    if [ -f "/www/server/openresty/bin/openresty" ]; then
        mkdir -p /www/server/vhost/acme
        chmod -R 644 /www/server/vhost/acme
    fi
fi

if version_lt "$oldVersion" "2.2.16"; then
    echo "更新面板到 v2.2.16 ..."
    echo "Update panel to v2.2.16 ..."
    if [ -f "/www/server/mysql/bin/mysql" ]; then
        ln -sf /www/server/mysql/bin/* /usr/bin/
        rm -f /etc/profile.d/mysql.sh
        source /etc/profile
    fi
fi

if version_lt "$oldVersion" "2.2.20"; then
    echo "更新面板到 v2.2.20 ..."
    echo "Update panel to v2.2.20 ..."
    echo "SESSION_LIFETIME=120" >> $panelPath/panel.conf
fi

echo $HR
echo "更新结束"
echo "Update finished"
