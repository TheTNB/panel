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
rm -rf ${pureftpdPath}
mkdir -p ${pureftpdPath}
cd ${pureftpdPath}

wget -T 120 -t 3 -O ${pureftpdPath}/pure-ftpd-${pureftpdVersion}.tar.gz ${downloadUrl}/pure-ftpd-${pureftpdVersion}.tar.gz
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

# 修改 pure-ftpd 配置文件
sed -i "s!# PureDB\s*/etc/pureftpd.pdb!PureDB ${pureftpdPath}/etc/pureftpd.pdb!" ${pureftpdPath}/etc/pure-ftpd.conf
sed -i 's!# ChrootEveryone\s*yes!ChrootEveryone yes!' ${pureftpdPath}/etc/pure-ftpd.conf
sed -i 's!NoAnonymous\s*no!NoAnonymous yes!' ${pureftpdPath}/etc/pure-ftpd.conf
sed -i 's!AnonymousCanCreateDirs\s*yes!AnonymousCanCreateDirs no!' ${pureftpdPath}/etc/pure-ftpd.conf
sed -i 's!AnonymousCantUpload\s*yes!AnonymousCantUpload no!' ${pureftpdPath}/etc/pure-ftpd.conf
sed -i 's!PAMAuthentication\s*yes!PAMAuthentication no!' ${pureftpdPath}/etc/pure-ftpd.conf
sed -i 's!UnixAuthentication\s*yes!UnixAuthentication no!' ${pureftpdPath}/etc/pure-ftpd.conf
sed -i 's!# PassivePortRange\s*30000 50000!PassivePortRange 39000 40000!' ${pureftpdPath}/etc/pure-ftpd.conf
sed -i 's!PassivePortRange\s*30000 50000!PassivePortRange 39000 40000!' ${pureftpdPath}/etc/pure-ftpd.conf
sed -i 's!LimitRecursion\s*10000 8!LimitRecursion 20000 8!' ${pureftpdPath}/etc/pure-ftpd.conf
sed -i 's!# TLS!TLS!' ${pureftpdPath}/etc/pure-ftpd.conf
sed -i "s!# CertFile\s*/etc/ssl/private/pure-ftpd.pem!CertFile ${pureftpdPath}/etc/pure-ftpd.pem!" ${pureftpdPath}/etc/pure-ftpd.conf
sed -i 's!# Bind\s*127.0.0.1,21!Bind 0.0.0.0,21!' ${pureftpdPath}/etc/pure-ftpd.conf
touch ${pureftpdPath}/etc/pureftpd.passwd
touch ${pureftpdPath}/etc/pureftpd.pdb

openssl dhparam -out ${pureftpdPath}/etc/pure-ftpd-dhparams.pem 2048
openssl req -x509 -nodes -days 3560 -newkey rsa:2048 -sha256 -keyout ${pureftpdPath}/etc/pure-ftpd.pem -out ${pureftpdPath}/etc/pure-ftpd.pem << EOF
CN
Beijing
Beijing
HaoZi Technology Co., Ltd
HaoZi Panel
github.com/haozi-team/panel
panel@haozi.net
EOF
chmod 600 ${pureftpdPath}/etc/*.pem

# 添加系统服务
ln -sf ${pureftpdPath}/bin/pure-pw /usr/bin/pure-pw

cat > /etc/systemd/system/pure-ftpd.service << EOF
[Unit]
Description=Pure-FTPd FTP server
After=network.target

[Service]
Type=forking
PIDFile=/var/run/pure-ftpd.pid
ExecStart=${pureftpdPath}/sbin/pure-ftpd ${pureftpdPath}/etc/pure-ftpd.conf
ExecReload=/bin/kill -HUP \$MAINPID
ExecStop=/bin/kill -TERM \$MAINPID

[Install]
WantedBy=multi-user.target
EOF

# 添加防火墙规则
if [ "${OS}" == "centos" ]; then
    firewall-cmd --zone=public --add-port=21/tcp --permanent
    firewall-cmd --zone=public --add-port=39000-40000/tcp --permanent
    firewall-cmd --reload
elif [ "${OS}" == "debian" ]; then
    ufw allow 21/tcp
    ufw allow 39000:40000/tcp
    ufw reload
fi

systemctl daemon-reload
systemctl enable pure-ftpd.service
systemctl start pure-ftpd.service

panel writePlugin pureftpd 1.0.50

echo -e "${HR}\nPure-Ftpd-${pureftpdVersion} 安装完成\n${HR}"
