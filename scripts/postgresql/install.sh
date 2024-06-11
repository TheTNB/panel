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
    echo "错误：耗子面板不支持该系统"
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
wget -T 120 -t 3 -O ${postgresqlPath}/postgresql-${postgresqlVersion}.7z ${downloadUrl}/postgresql-${postgresqlVersion}.7z
wget -T 20 -t 3 -O ${postgresqlPath}/postgresql-${postgresqlVersion}.7z.checksum.txt ${downloadUrl}/postgresql-${postgresqlVersion}.7z.checksum.txt

if ! sha256sum --status -c postgresql-${postgresqlVersion}.7z.checksum.txt; then
    echo -e $HR
    echo "错误：PostgreSQL 源码 checksum 校验失败，文件可能被篡改或不完整，已终止操作"
    rm -rf ${postgresqlPath}
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
    rm -rf ${postgresqlPath}
    exit 1
fi
make "-j${j}"
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

cd ${postgresqlPath}
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
LimitNOFILE=500000

[Install]
WantedBy=multi-user.target
EOF

# 在 /etc/systemd/logind.conf 设置 RemoveIPC=no，不然会删除 /dev/shm 下的共享内存文件
checkRemoveIPC=$(cat /etc/systemd/logind.conf | grep '^RemoveIPC=no.*$')
if [ "${checkRemoveIPC}" == "" ]; then
    echo "RemoveIPC=no" >> /etc/systemd/logind.conf
    systemctl restart systemd-logind
fi

# 启动服务
systemctl daemon-reload
systemctl enable postgresql
systemctl start postgresql

panel writePlugin postgresql${1} ${postgresqlVersion}

echo -e "${HR}\nPostgreSQL-${1} 安装完成\n${HR}"
