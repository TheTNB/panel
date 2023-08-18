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
postgresqlPassword=$(cat /dev/urandom | head -n 16 | md5sum | head -c 16)
cpuCore=$(cat /proc/cpuinfo | grep "processor" | wc -l)

if [[ "${1}" == "15" ]]; then
    postgresqlVersion="15.3"
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

postgresqlUserCheck=$(cat /etc/passwd | grep postgres)
if [ "${postgresqlUserCheck}" == "" ]; then
    groupadd postgres
    useradd -g postgres postgres
fi

# 准备目录
rm -rf ${postgresqlPath}
mkdir -p ${postgresqlPath}
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
    rm -rf ${postgresqlPath}
    exit 1
fi
make -j${cpuCore}
if [ "$?" != "0" ]; then
    echo -e $HR
    echo "错误：PostgreSQL 编译失败，请截图错误信息寻求帮助。"
    rm -rf ${postgresqlPath}
    exit 1
fi
make install
if [ "$?" != "0" ]; then
    echo -e $HR
    echo "错误：PostgreSQL 安装失败，请截图错误信息寻求帮助。"
    rm -rf ${postgresqlPath}
    exit 1
fi

rm -rf ${postgresqlPath}/src

# 配置
mkdir -p ${postgresqlPath}/data
mkdir -p ${postgresqlPath}/logs
chown -R postgres:postgres ${postgresqlPath}
chmod -R 700 ${postgresqlPath}

echo "export PATH=${postgresqlPath}/bin:\$PATH" >> /etc/profile
source /etc/profile

mkdir -p /home/postgres
cd /home/postgres
if [ -f /home/postgres/.bash_profile ]; then
        echo "export PGHOME=${postgresqlPath}" >> /home/postgres/.bash_profile
        echo "export PGDATA=${postgresqlPath}/data" >> /home/postgres/.bash_profile
        echo "export PATH=${postgresqlPath}/bin:\$PATH " >> /home/postgres/.bash_profile
        echo "MANPATH=$PGHOME/share/man:$MANPATH" >> /home/postgres/.bash_profile
        echo "LD_LIBRARY_PATH=$PGHOME/lib:$LD_LIBRARY_PATH" >> /home/postgres/.bash_profile
        source /home/postgres/.bash_profile
fi
if [ -f /home/postgres/.profile ]; then
        echo "export PGHOME=${postgresqlPath}" >> /home/postgres/.profile
        echo "export PGDATA=${postgresqlPath}/data" >> /home/postgres/.profile
        echo "export PATH=${postgresqlPath}/bin:\$PATH " >> /home/postgres/.profile
        echo "MANPATH=$PGHOME/share/man:$MANPATH" >> /home/postgres/.profile
        echo "LD_LIBRARY_PATH=$PGHOME/lib:$LD_LIBRARY_PATH" >> /home/postgres/.profile
        source /home/postgres/.profile
fi

# 初始化
su - postgres -c "${postgresqlPath}/bin/initdb -D ${postgresqlPath}/data"
if [ "$?" != "0" ]; then
    echo -e $HR
    echo "错误：PostgreSQL 初始化失败，请截图错误信息寻求帮助。"
    rm -rf ${postgresqlPath}
    exit 1
fi

# 配置慢查询日志
cat >> ${postgresqlPath}/data/postgresql.conf << EOF
logging_collector = on
log_destination = 'stderr'
log_directory = '${postgresqlPath}/logs'
log_filename = 'postgresql-%Y-%m-%d.log'
log_statement = all
log_min_duration_statement = 5000
EOF

# 写入服务
cat > /etc/systemd/system/postgresql.service << EOF
[Unit]
Description=PostgreSQL database server
Documentation=man:postgres(1)
After=network-online.target
Wants=network-online.target

[Service]
Type=notify
User=postgres
ExecStart=${postgresqlPath}/bin/postgres -D ${postgresqlPath}/data
ExecReload=/bin/kill -HUP \$MAINPID
KillMode=mixed
KillSignal=SIGINT
TimeoutSec=infinity

[Install]
WantedBy=multi-user.target
EOF

# 启动服务
systemctl daemon-reload
systemctl enable postgresql
systemctl start postgresql

panel writePlugin postgresql${1} ${postgresqlVersion}

echo -e "${HR}\nPostgreSQL-${1} 安装完成\n${HR}"
