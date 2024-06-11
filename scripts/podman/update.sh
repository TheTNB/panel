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
OS=$(source /etc/os-release && { [[ "$ID" == "debian" ]] && echo "debian"; } || { [[ "$ID" == "centos" ]] || [[ "$ID" == "rhel" ]] || [[ "$ID" == "rocky" ]] || [[ "$ID" == "almalinux" ]] && echo "centos"; } || echo "unknown")
podmanVersion="4.0.0"

if [ "${OS}" == "centos" ]; then
    dnf makecache -y
    dnf update podman -y
elif [ "${OS}" == "debian" ]; then
    apt-get update
    apt-get upgrade podman -y
else
    echo -e $HR
    echo "错误：耗子面板不支持该系统"
    exit 1
fi

systemctl restart podman

panel writePlugin podman ${podmanVersion}
echo -e ${HR}
echo "podman 安装完成"
echo -e ${HR}
