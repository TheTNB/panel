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
OS=$(source /etc/os-release && { [[ "$ID" == "debian" ]] && echo "debian"; } || { [[ "$ID" == "centos" ]] || [[ "$ID" == "rhel" ]] || [[ "$ID" == "rocky" ]] || [[ "$ID" == "almalinux" ]] && echo "centos"; } || echo "unknown")
downloadUrl="https://dl.cdn.haozi.net/panel/redis"
setupPath="/www"
redisPath="${setupPath}/server/redis"
redisVersion="7.0.12"
cpuCore=$(cat /proc/cpuinfo | grep "processor" | wc -l)

# 准备目录
cd ${redisPath}

# 下载源码
wget -T 120 -t 3 -O ${redisPath}/redis-${redisVersion}.tar.gz ${downloadUrl}/redis-${redisVersion}.tar.gz
tar -zxvf redis-${redisVersion}.tar.gz
rm -f redis-${redisVersion}.tar.gz
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
systemctl restart redis

panel writePlugin redis ${redisVersion}

echo -e "${HR}\nRedis 升级完成\n${HR}"
