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
OS=$(source /etc/os-release && { [[ "$ID" == "debian" ]] && echo "debian"; } || { [[ "$ID" == "centos" ]] || [[ "$ID" == "rhel" ]] || [[ "$ID" == "rocky" ]] || [[ "$ID" == "almalinux" ]] && echo "centos"; } || echo "unknown")

if [ "${OS}" == "centos" ]; then
    dnf install -y s3fs-fuse
elif [ "${OS}" == "debian" ]; then
    apt-get install -y s3fs
else
    echo -e $HR
    echo "错误：不支持的操作系统"
    exit 1
fi
if [ "$?" != "0" ]; then
    echo -e $HR
    echo "错误：安装软件失败，请截图错误信息寻求帮助。"
    exit 1
fi

panel writePlugin s3fs 1.9
