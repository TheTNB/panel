#!/bin/bash
export PATH=/bin:/sbin:/usr/bin:/usr/sbin:/www/server/bin:/www/server/sbin:$PATH

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
downloadUrl="https://dl.cdn.haozi.net/panel/gitea"
giteaPath="/www/server/gitea"
giteaVersion="1.22.0"

if [ ! -d "${giteaPath}" ]; then
    mkdir -p ${giteaPath}
fi

# 架构判断
if [ "${ARCH}" == "x86_64" ]; then
    giteaFile="gitea-${giteaVersion}-linux-amd64.7z"
elif [ "${ARCH}" == "aarch64" ]; then
    giteaFile="gitea-${giteaVersion}-linux-arm64.7z"
else
    echo -e $HR
    echo "错误：不支持的架构"
    exit 1
fi

# 安装依赖
if [ "${OS}" == "centos" ]; then
    dnf makecache -y
    dnf install git git-lfs -y
elif [ "${OS}" == "debian" ]; then
    apt-get update
    apt-get install git git-lfs -y
else
    echo -e $HR
    echo "错误：耗子 Linux 面板不支持该系统"
    exit 1
fi

git lfs install
git lfs version

# 下载
cd ${giteaPath}
wget -T 120 -t 3 -O ${giteaPath}/${giteaFile} ${downloadUrl}/${giteaFile}
wget -T 20 -t 3 -O ${giteaPath}/${giteaFile}.checksum.txt ${downloadUrl}/${giteaFile}.checksum.txt
if ! sha256sum --status -c ${giteaPath}/${giteaFile}.checksum.txt; then
    echo -e $HR
    echo "错误：gitea checksum 校验失败，文件可能被篡改或不完整，已终止操作"
    rm -rf ${giteaPath}
    exit 1
fi

# 解压
cd ${giteaPath}
7z x ${giteaFile}
rm -f ${giteaFile} ${giteaFile}.checksum.txt
mv gitea-${giteaVersion}-linux-* gitea
if [ ! -f "${giteaPath}/gitea" ]; then
    echo -e $HR
    echo "错误：gitea 解压失败"
    rm -rf ${giteaPath}
    exit 1
fi

# 初始化目录
mkdir -p ${giteaPath}/{custom,data,log}
chown -R www:www ${giteaPath}
chmod -R 750 ${giteaPath}
ln -sf ${giteaPath}/gitea /usr/local/bin/gitea

# 配置systemd
cat >/etc/systemd/system/gitea.service <<EOF
[Unit]
Description=Gitea (Git with a cup of tea)
After=network.target
###
# 可以自行添加数据库服务依赖
# Can add database service dependencies yourself
###
#
#Wants=mysqld.service
#After=mysqld.service
#
#Wants=postgresql.service
#After=postgresql.service
#
#Wants=redis.service
#After=redis.service
#

[Service]
LimitNOFILE=524288:524288
RestartSec=2s
Type=simple
User=www
Group=www
WorkingDirectory=/www/server/gitea/
ExecStart=/usr/local/bin/gitea web --config /www/server/gitea/app.ini
Restart=always
Environment=USER=www HOME=/home/www GITEA_WORK_DIR=/www/server/gitea
CapabilityBoundingSet=CAP_NET_BIND_SERVICE
AmbientCapabilities=CAP_NET_BIND_SERVICE
PrivateUsers=false

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable gitea
systemctl start gitea

# 防火墙
if [ "${OS}" == "centos" ]; then
    firewall-cmd --zone=public --add-port=3000/tcp --permanent
    firewall-cmd --reload
elif [ "${OS}" == "debian" ]; then
    ufw allow 3000/tcp
    ufw reload
fi

panel writePlugin gitea ${giteaVersion}
echo -e $HR
echo "gitea 安装完成，请访问 IP:3000 完成初始化向导"
echo "安装后建议修改 systemd 配置 /etc/systemd/system/gitea.service 中的数据库依赖"
echo -e $HR
