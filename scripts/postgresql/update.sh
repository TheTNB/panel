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
memTotal=$(LC_ALL=C free -m | grep Mem | awk '{print  $2}')
OS=$(source /etc/os-release && { [[ "$ID" == "debian" ]] && echo "debian"; } || { [[ "$ID" == "centos" ]] || [[ "$ID" == "rhel" ]] || [[ "$ID" == "rocky" ]] || [[ "$ID" == "almalinux" ]] && echo "centos"; } || echo "unknown")
downloadUrl="https://dl.cdn.haozi.net/panel/postgresql"
setupPath="/www"
postgresqlPath="${setupPath}/server/postgresql"
postgresqlVersion=""

source ${setupPath}/panel/scripts/calculate_j.sh
j=$(calculate_j)

if [[ "${1}" == "15" ]]; then
    postgresqlVersion="15.7"
elif [[ "${1}" == "16" ]]; then
    postgresqlVersion="16.3"
else
    echo -e $HR
    echo "错误：不支持的 PostgreSQL 版本！"
    exit 1
fi

# 安装依赖
if [ "${OS}" == "centos" ]; then
    dnf makecache -y
    dnf groupinstall "Development Tools" -y
    dnf install make gettext zlib-devel readline-devel libicu-devel libxml2-devel libxslt-devel openssl-devel systemd-devel -y
elif [ "${OS}" == "debian" ]; then
    apt-get update
    apt-get install build-essential make gettext zlib1g-dev libreadline-dev libicu-dev libxml2-dev libxslt-dev libssl-dev libsystemd-dev -y
else
    echo -e $HR
    echo "错误：耗子 Linux 面板不支持该系统"
    exit 1
fi

# 停止已有服务
systemctl stop postgresql

# 准备目录
rm -rf ${postgresqlPath}/src
cd ${postgresqlPath}

# 下载源码
wget -T 120 -t 3 -O ${postgresqlPath}/postgresql-${postgresqlVersion}.7z ${downloadUrl}/postgresql-${postgresqlVersion}.7z
wget -T 20 -t 3 -O ${postgresqlPath}/postgresql-${postgresqlVersion}.7z.checksum.txt ${downloadUrl}/postgresql-${postgresqlVersion}.7z.checksum.txt

if ! sha256sum --status -c postgresql-${postgresqlVersion}.7z.checksum.txt; then
    echo -e $HR
    echo "错误：PostgreSQL 源码 checksum 校验失败，文件可能被篡改或不完整，已终止操作"
    exit 1
fi

7z x postgresql-${postgresqlVersion}.7z
rm -f postgresql-${postgresqlVersion}.7z
rm -f postgresql-${postgresqlVersion}.7z.checksum.txt
mv postgresql-${postgresqlVersion} src
chmod -R 755 src

# 编译
cd src
./configure --prefix=${postgresqlPath} --enable-nls='zh_CN en' --with-icu --with-ssl=openssl --with-systemd --with-libxml --with-libxslt
if [ "$?" != "0" ]; then
    echo -e $HR
    echo "错误：PostgreSQL 编译初始化失败，请截图错误信息寻求帮助。"
    exit 1
fi
make "-j${j}"
if [ "$?" != "0" ]; then
    echo -e $HR
    echo "错误：PostgreSQL 编译失败，请截图错误信息寻求帮助。"
    exit 1
fi
make install
if [ "$?" != "0" ]; then
    echo -e $HR
    echo "错误：PostgreSQL 安装失败，请截图错误信息寻求帮助。"
    exit 1
fi

cd ${postgresqlPath}
rm -rf ${postgresqlPath}/src

# 配置
chown -R postgres:postgres ${postgresqlPath}
chmod -R 700 ${postgresqlPath}

panel writePlugin postgresql${1} ${postgresqlVersion}

systemctl daemon-reload
systemctl restart postgresql

echo -e "${HR}\nPostgreSQL-${1} 升级完成\n${HR}"
