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
downloadUrl="https://dl.cdn.haozi.net/panel/frp"
frpPath="/www/server/frp"
frpVersion="0.58.0"

if [ ! -d "${frpPath}" ]; then
    mkdir -p ${frpPath}
fi

# 架构判断
if [ "${ARCH}" == "x86_64" ]; then
    frpFile="frp_${frpVersion}_linux_amd64.7z"
elif [ "${ARCH}" == "aarch64" ]; then
    frpFile="frp_${frpVersion}_linux_arm64.7z"
else
    echo -e $HR
    echo "错误：不支持的架构"
    exit 1
fi

# 备份配置
if [ -f "${frpPath}/frps.toml" ]; then
    cp -f ${frpPath}/frps.toml ${frpPath}/frps.toml.bak
fi
if [ -f "${frpPath}/frpc.toml" ]; then
    cp -f ${frpPath}/frpc.toml ${frpPath}/frpc.toml.bak
fi

# 下载frp
cd ${frpPath}
wget -T 120 -t 3 -O ${frpPath}/${frpFile} ${downloadUrl}/${frpFile}
wget -T 20 -t 3 -O ${frpPath}/${frpFile}.checksum.txt ${downloadUrl}/${frpFile}.checksum.txt
if ! sha256sum --status -c ${frpPath}/${frpFile}.checksum.txt; then
    echo -e $HR
    echo "错误：frp checksum 校验失败，文件可能被篡改或不完整，已终止操作"
    rm -rf ${frpPath}
    exit 1
fi

# 解压frp
cd ${frpPath}
7z x ${frpFile}
chmod -R 700 ${frpPath}
rm -f ${frpFile} ${frpFile}.checksum.txt

# 还原配置
if [ -f "${frpPath}/frps.toml.bak" ]; then
    cp -f ${frpPath}/frps.toml.bak ${frpPath}/frps.toml
fi
if [ -f "${frpPath}/frpc.toml.bak" ]; then
    cp -f ${frpPath}/frpc.toml.bak ${frpPath}/frpc.toml
fi

systemctl restart frps
systemctl restart frpc

panel writePlugin frp ${frpVersion}
echo -e ${HR}
echo "frp 安装完成"
echo -e ${HR}

