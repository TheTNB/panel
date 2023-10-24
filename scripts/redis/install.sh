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
OS=$(source /etc/os-release && { [[ "$ID" == "debian" ]] && echo "debian"; } || { [[ "$ID" == "centos" ]] || [[ "$ID" == "rhel" ]] || [[ "$ID" == "rocky" ]] || [[ "$ID" == "almalinux" ]] && echo "centos"; } || echo "unknown")
downloadUrl="https://jihulab.com/haozi-team/download/-/raw/main/panel/redis"
setupPath="/www"
redisPath="${setupPath}/server/redis"
redisVersion="7.0.12"
cpuCore=$(cat /proc/cpuinfo | grep "processor" | wc -l)

if ! id -u "redis" > /dev/null 2>&1; then
    groupadd redis
    useradd -s /sbin/nologin -g redis redis
fi

# 安装依赖
if [ "${OS}" == "centos" ]; then
    dnf makecache -y
    dnf groupinstall "Development Tools" -y
    dnf install systemd-devel openssl-devel -y
elif [ "${OS}" == "debian" ]; then
    apt-get update
    apt-get install build-essential libsystemd-dev libssl-dev -y
else
    echo -e $HR
    echo "错误：耗子Linux面板不支持该系统"
    exit 1
fi

# 准备目录
rm -rf ${redisPath}
mkdir -p ${redisPath}
cd ${redisPath}

# 下载源码
wget -T 120 -t 3 -O ${redisPath}/redis-${redisVersion}.tar.gz ${downloadUrl}/redis-${redisVersion}.tar.gz
wget -T 20 -t 3 -O ${redisPath}/redis-${redisVersion}.tar.gz.checksum.txt ${downloadUrl}/redis-${redisVersion}.tar.gz.checksum.txt

if ! sha256sum --status -c redis-${redisVersion}.tar.gz.checksum.txt; then
    echo -e $HR
    echo "错误：Redis源码 checksum 校验失败，文件可能被篡改或不完整，已终止操作"
    rm -rf ${redisPath}
    exit 1
fi

tar -zxvf redis-${redisVersion}.tar.gz
rm -f redis-${redisVersion}.tar.gz
rm -f redis-${redisVersion}.tar.gz.checksum.txt
mv redis-${redisVersion}/* ./ && rm -rf redis-${redisVersion}
mkdir -p ${redisPath}/bin

make BUILD_TLS=yes USE_SYSTEMD=yes -j${cpuCore}
if [ "$?" != "0" ]; then
    echo -e $HR
    echo "错误：Redis编译失败，请截图错误信息寻求帮助。"
    rm -rf ${redisPath}
    exit 1
fi
make PREFIX=${redisPath} install
if [ ! -f "${redisPath}/bin/redis-server" ]; then
    echo -e $HR
    echo "错误：Redis安装失败，请截图错误信息寻求帮助。"
    rm -rf ${redisPath}
    exit 1
fi

# 设置软链接
ln -sf ${redisPath}/bin/redis-cli /usr/bin/redis-cli
ln -sf ${redisPath}/bin/redis-server /usr/bin/redis-server
ln -sf ${redisPath}/bin/redis-sentinel /usr/bin/redis-sentinel
ln -sf ${redisPath}/bin/redis-benchmark /usr/bin/redis-benchmark
ln -sf ${redisPath}/bin/redis-check-aof /usr/bin/redis-check-aof
ln -sf ${redisPath}/bin/redis-check-rdb /usr/bin/redis-check-rdb

# 设置配置文件
VM_OVERCOMMIT_MEMORY=$(cat /etc/sysctl.conf|grep vm.overcommit_memory)
NET_CORE_SOMAXCONN=$(cat /etc/sysctl.conf|grep net.core.somaxconn)
if [ -z "${VM_OVERCOMMIT_MEMORY}" ] && [ -z "${NET_CORE_SOMAXCONN}" ];then
		echo "vm.overcommit_memory = 1" >> /etc/sysctl.conf
		echo "net.core.somaxconn = 1024" >> /etc/sysctl.conf
		sysctl -p
fi

sed -i 's/dir .\//dir \/www\/server\/redis\//g' ${redisPath}/redis.conf
sed -i 's/# supervised.*/supervised systemd/g' ${redisPath}/redis.conf
sed -i 's/daemonize.*/daemonize no/g' ${redisPath}/redis.conf

if [ "${ARCH}" == "aarch64" ]; then
    echo "ignore-warnings ARM64-COW-BUG" >> ${redisPath}/redis.conf
fi

chown -R redis:redis ${redisPath}
chmod -R 755 ${redisPath}

# 设置服务
cp -r utils/systemd-redis_server.service /etc/systemd/system/redis.service
sed -i "s!ExecStart=.*!ExecStart=${redisPath}/bin/redis-server ${redisPath}/redis.conf!g" /etc/systemd/system/redis.service
sed -i "s!#User=.*!User=redis!g" /etc/systemd/system/redis.service
sed -i "s!#Group=.*!Group=redis!g" /etc/systemd/system/redis.service
sed -i "s!#WorkingDirectory=.*!WorkingDirectory=${redisPath}!g" /etc/systemd/system/redis.service

systemctl daemon-reload
systemctl enable redis
systemctl start redis

panel writePlugin redis ${redisVersion}

echo -e "${HR}\nRedis 安装完成\n${HR}"
