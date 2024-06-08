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
downloadUrl="https://dl.cdn.haozi.net/panel/gitea"
giteaPath="/www/server/gitea"
giteaVersion="1.22.0"

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

# 替换文件
systemctl stop gitea
rm -f gitea
mv gitea-${giteaVersion}-linux-* gitea
if [ ! -f "${giteaPath}/gitea" ]; then
    echo -e $HR
    echo "错误：gitea 解压失败"
    rm -rf ${giteaPath}
    exit 1
fi

chown -R www:www ${giteaPath}
chmod -R 750 ${giteaPath}
systemctl start gitea

panel writePlugin gitea ${giteaVersion}
echo -e $HR
echo "gitea 升级完成"
echo -e $HR
