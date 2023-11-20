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
redisVersion="7.2.3"
cpuCore=$(cat /proc/cpuinfo | grep "processor" | wc -l)

# 读取信息
redisConfig="${redisPath}/redis.conf"
redisPort=$(cat ${redisConfig} | grep 'port ' | grep -v '#' | awk '{print $2}')
redisPass=$(cat ${redisConfig} | grep 'requirepass ' | grep -v '#' | awk '{print $2}')
redisHost=$(cat ${redisConfig} | grep 'bind ' | grep -v '#' | awk '{print $2}')
redisDir=$(cat ${redisConfig} | grep 'dir ' | grep -v '#' | awk '{print $2}')

# 备份
if [ -z "${redisPass}" ]; then
    redis-cli -p ${redisPort} << EOF
SAVE
EOF
else
    redis-cli -p ${redisPort} -a ${redisPass} << EOF
SAVE
EOF
fi
mv ${redisDir}/dump.rdb /tmp/dump.rdb.bak

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
    exit 1
fi
make PREFIX=${redisPath} install
if [ ! -f "${redisPath}/bin/redis-server" ]; then
    echo -e $HR
    echo "错误：Redis升级失败，请截图错误信息寻求帮助。"
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

# 恢复配置
if [ -n "${redisPass}" ]; then
    sed -i "s!# requirepass .*!requirepass ${redisPass}!g" ${redisPath}/redis.conf
fi
if [ -n "${redisHost}" ]; then
    sed -i "s!bind .*!bind ${redisHost}!g" ${redisPath}/redis.conf
fi
if [ -n "${redisPort}" ]; then
    sed -i "s!port .*!port ${redisPort}!g" ${redisPath}/redis.conf
fi
if [ -n "${redisDir}" ]; then
    sed -i "s!dir .*!dir ${redisDir}!g" ${redisPath}/redis.conf
fi

# 恢复数据
if [ -f "/tmp/dump.rdb.bak" ]; then
    mv /tmp/dump.rdb.bak ${redisDir}/dump.rdb
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
