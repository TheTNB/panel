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
downloadUrl="https://jihulab.com/haozi-team/download/-/raw/main/panel/pure-ftpd"
setupPath="/www"
pureftpdPath="${setupPath}/server/pure-ftpd"
pureftpdVersion="1.0.50"

# 准备安装目录
cp ${pureftpdPath}/etc/pureftpd.passwd /tmp/pureftpd.passwd
cp ${pureftpdPath}/etc/pureftpd.pdb /tmp/pureftpd.pdb
cp ${pureftpdPath}/etc/pureftpd.conf /tmp/pureftpd.conf
systemctl stop pure-ftpd.service
rm -rf ${pureftpdPath}
mkdir -p ${pureftpdPath}
cd ${pureftpdPath}

wget -T 60 -t 3 -O ${pureftpdPath}/pure-ftpd-${pureftpdVersion}.tar.gz ${downloadUrl}/pure-ftpd-${pureftpdVersion}.tar.gz
wget -T 20 -t 3 -O ${pureftpdPath}/pure-ftpd-${pureftpdVersion}.tar.gz.checksum.txt ${downloadUrl}/pure-ftpd-${pureftpdVersion}.tar.gz.checksum.txt

if ! sha256sum --status -c pure-ftpd-${pureftpdVersion}.tar.gz.checksum.txt; then
    echo -e $HR
    echo "错误：Pure-Ftpd-${pureftpdVersion}源码 checksum 校验失败，文件可能被篡改或不完整，已终止操作"
    rm -rf ${pureftpdPath}
    exit 1
fi

tar -xvf pure-ftpd-${pureftpdVersion}.tar.gz
rm -f pure-ftpd-${pureftpdVersion}.tar.gz
rm -f pure-ftpd-${pureftpdVersion}.tar.gz.checksum.txt
mv pure-ftpd-${pureftpdVersion} src
cd src

./configure --prefix=${pureftpdPath} CFLAGS=-O2 --with-puredb --with-quotas --with-cookie --with-virtualhosts --with-diraliases --with-sysquotas --with-ratios --with-altlog --with-paranoidmsg --with-shadow --with-welcomemsg --with-throttling --with-uploadscript --with-language=simplified-chinese --with-rfc2640 --with-ftpwho --with-tls
if [ "$?" != "0" ]; then
    echo -e $HR
    echo "错误：Pure-Ftpd-${pureftpdVersion}编译配置失败，请截图错误信息寻求帮助。"
    exit 1
fi

make
if [ "$?" != "0" ]; then
    echo -e $HR
    echo "错误：Pure-Ftpd-${pureftpdVersion}编译失败，请截图错误信息寻求帮助。"
    exit 1
fi

make install
if [ ! -f "${pureftpdPath}/bin/pure-pw" ]; then
    echo -e $HR
    echo "错误：Pure-Ftpd-${pureftpdVersion}安装失败，请截图错误信息寻求帮助。"
    exit 1
fi

# 还原配置
cp /tmp/pureftpd.passwd ${pureftpdPath}/etc/pureftpd.passwd
cp /tmp/pureftpd.pdb ${pureftpdPath}/etc/pureftpd.pdb
cp /tmp/pureftpd.conf ${pureftpdPath}/etc/pureftpd.conf

systemctl start pure-ftpd.service

panel writePlugin pureftpd 1.0.50

echo -e "${HR}\nPure-Ftpd-${pureftpdVersion} 升级完成\n${HR}"
