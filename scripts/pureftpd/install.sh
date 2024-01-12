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
downloadUrl="https://git.haozi.net/opensource/download/-/raw/main/panel/pure-ftpd"
setupPath="/www"
pureftpdPath="${setupPath}/server/pure-ftpd"
pureftpdVersion="1.0.50"

source ${setupPath}/panel/scripts/calculate_j.sh
j=$(calculate_j)

# 准备安装目录
rm -rf ${pureftpdPath}
mkdir -p ${pureftpdPath}
cd ${pureftpdPath}

wget -T 120 -t 3 -O ${pureftpdPath}/pure-ftpd-${pureftpdVersion}.tar.gz ${downloadUrl}/pure-ftpd-${pureftpdVersion}.tar.gz
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

make "-j${j}"
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
sed -i "s!# PIDFile\s*/var/run/pure-ftpd.pid!PIDFile ${pureftpdPath}/etc/pure-ftpd.pid!" ${pureftpdPath}/etc/pure-ftpd.conf
touch ${pureftpdPath}/etc/pureftpd.passwd
touch ${pureftpdPath}/etc/pureftpd.pdb

openssl dhparam -out ${pureftpdPath}/etc/pure-ftpd-dhparams.pem 2048
openssl req -x509 -nodes -days 36500 -newkey rsa:2048 -sha256 -keyout ${pureftpdPath}/etc/pure-ftpd.pem -out ${pureftpdPath}/etc/pure-ftpd.pem -subj "/C=CN/ST=Tianjin/L=Tianjin/O=HaoZi Technology Co., Ltd./OU=HaoZi Panel/CN=Panel"
chmod 600 ${pureftpdPath}/etc/*.pem

# 添加系统服务
ln -sf ${pureftpdPath}/bin/pure-pw /usr/bin/pure-pw

cat > /etc/systemd/system/pure-ftpd.service << EOF
[Unit]
Description=Pure-FTPd FTP server
After=syslog.target network.target

[Service]
Type=forking
PIDFile=${pureftpdPath}/etc/pure-ftpd.pid
ExecStart=${pureftpdPath}/sbin/pure-ftpd ${pureftpdPath}/etc/pure-ftpd.conf
ExecStartPost=/bin/sleep 2
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
