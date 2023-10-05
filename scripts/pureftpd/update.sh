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
OS=$(source /etc/os-release && { [[ "$ID" == "debian" ]] && echo "debian"; } || { [[ "$ID" == "centos" ]] || [[ "$ID" == "rhel" ]] || [[ "$ID" == "rocky" ]] || [[ "$ID" == "almalinux" ]] && echo "centos"; } || echo "unknown")
downloadUrl="https://dl.cdn.haozi.net/panel/pure-ftpd"
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
if [ "$?" != "0" ]; then
    echo -e $HR
    echo "错误：Pure-Ftpd-${pureftpdVersion}下载失败，请检查网络是否正常。"
    exit 1
fi

tar -xvf pure-ftpd-${pureftpdVersion}.tar.gz
rm -f pure-ftpd-${pureftpdVersion}.tar.gz
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
