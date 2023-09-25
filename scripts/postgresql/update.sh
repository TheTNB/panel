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
memTotal=$(LC_ALL=C free -m | grep Mem | awk '{print  $2}')
OS=$(source /etc/os-release && { [[ "$ID" == "debian" ]] && echo "debian"; } || { [[ "$ID" == "centos" ]] || [[ "$ID" == "rhel" ]] || [[ "$ID" == "rocky" ]] || [[ "$ID" == "almalinux" ]] && echo "centos"; } || echo "unknown")
downloadUrl="https://dl.cdn.haozi.net/panel/postgresql"
setupPath="/www"
postgresqlPath="${setupPath}/server/postgresql"
postgresqlVersion=""
cpuCore=$(cat /proc/cpuinfo | grep "processor" | wc -l)

if [[ "${1}" == "15" ]]; then
    postgresqlVersion="15.4"
elif [[ "${1}" == "16" ]]; then
    postgresqlVersion="16.0"
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
    echo "错误：耗子Linux面板不支持该系统"
    exit 1
fi

# 停止已有服务
systemctl stop postgresql

# 准备目录
rm -rf ${postgresqlPath}/src
cd ${postgresqlPath}

# 下载源码
wget -T 120 -O ${postgresqlPath}/postgresql-${postgresqlVersion}.tar.gz ${downloadUrl}/postgresql-${postgresqlVersion}.tar.gz
tar -zxvf postgresql-${postgresqlVersion}.tar.gz
rm -f postgresql-${postgresqlVersion}.tar.gz
mv postgresql-${postgresqlVersion} src

# 编译
cd src
./configure --prefix=${postgresqlPath} --enable-nls='zh_CN en' --with-icu --with-ssl=openssl --with-systemd --with-libxml --with-libxslt
if [ "$?" != "0" ]; then
    echo -e $HR
    echo "错误：PostgreSQL 编译初始化失败，请截图错误信息寻求帮助。"
    exit 1
fi
make -j${cpuCore}
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
