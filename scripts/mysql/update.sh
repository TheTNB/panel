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
memTotal=$(LC_ALL=C free -m | grep Mem | awk '{print  $2}')
OS=$(source /etc/os-release && { [[ "$ID" == "debian" ]] && echo "debian"; } || { [[ "$ID" == "centos" ]] || [[ "$ID" == "rhel" ]] || [[ "$ID" == "rocky" ]] || [[ "$ID" == "almalinux" ]] && echo "centos"; } || echo "unknown")
downloadUrl="https://dl.cdn.haozi.net/panel/mysql"
setupPath="/www"
mysqlPath="${setupPath}/server/mysql"
mysqlVersion=""
mysqlPassword=$(panel getSetting mysql_root_password)
cpuCore=$(cat /proc/cpuinfo | grep "processor" | wc -l)

source ${setupPath}/panel/scripts/calculate_j.sh
j=$(calculate_j)

if [[ "${1}" == "84" ]]; then
    mysqlVersion="8.4.0"
    j=$(calculate_j2)
elif [[ "${1}" == "80" ]]; then
    mysqlVersion="8.0.37"
    j=$(calculate_j2)
elif [[ "${1}" == "57" ]]; then
    mysqlVersion="5.7.44"
else
    echo -e $HR
    echo "错误：不支持的 MySQL 版本！"
    exit 1
fi

# 安装依赖
if [ "${OS}" == "centos" ]; then
    dnf makecache -y
    dnf groupinstall "Development Tools" -y
    dnf install cmake bison ncurses-devel libtirpc-devel openssl-devel pkg-config openldap-devel libudev-devel cyrus-sasl-devel patchelf rpcgen rpcsvc-proto-devel p7zip p7zip-plugins -y
    dnf install gcc-toolset-12-gcc gcc-toolset-12-gcc-c++ gcc-toolset-12-binutils gcc-toolset-12-annobin-annocheck gcc-toolset-12-annobin-plugin-gcc -y
elif [ "${OS}" == "debian" ]; then
    apt-get update
    apt-get install build-essential cmake bison libncurses5-dev libtirpc-dev libssl-dev pkg-config libldap2-dev libudev-dev libsasl2-dev patchelf p7zip p7zip-full -y
else
    echo -e $HR
    echo "错误：耗子 Linux 面板不支持该系统"
    exit 1
fi

# 停止已有服务
systemctl stop mysqld

# 准备目录
cd ${mysqlPath}

# 下载源码
wget -T 120 -t 3 -O ${mysqlPath}/mysql-${mysqlVersion}.7z ${downloadUrl}/mysql-${mysqlVersion}.7z
wget -T 20 -t 3 -O ${mysqlPath}/mysql-${mysqlVersion}.7z.checksum.txt ${downloadUrl}/mysql-${mysqlVersion}.7z.checksum.txt
if ! sha256sum --status -c mysql-${mysqlVersion}.7z.checksum.txt; then
    echo -e $HR
    echo "错误：MySQL 源码 checksum 校验失败，文件可能被篡改或不完整，已终止操作"
    rm -rf ${mysqlPath}
    exit 1
fi

7z x mysql-${mysqlVersion}.7z
rm -f mysql-${mysqlVersion}.7z
rm -f mysql-${mysqlVersion}.7z.checksum.txt

# 编译
mv mysql-${mysqlVersion} src
chmod -R 755 src
cd src
mkdir build
cd build

# 5.7 需要 boost
if [[ "${1}" == "57" ]]; then
    MAYBE_WITH_BOOST="-DWITH_BOOST=../boost"
fi

cmake .. -DCMAKE_INSTALL_PREFIX=${mysqlPath} -DMYSQL_DATADIR=${mysqlPath}/data -DSYSCONFDIR=${mysqlPath}/conf -DWITH_MYISAM_STORAGE_ENGINE=1 -DWITH_INNOBASE_STORAGE_ENGINE=1 -DWITH_ARCHIVE_STORAGE_ENGINE=1 -DWITH_FEDERATED_STORAGE_ENGINE=1 -DWITH_BLACKHOLE_STORAGE_ENGINE=1 -DDEFAULT_CHARSET=utf8mb4 -DDEFAULT_COLLATION=utf8mb4_general_ci -DENABLED_LOCAL_INFILE=1 -DWITH_DEBUG=0 -DWITH_UNIT_TESTS=OFF -DINSTALL_MYSQLTESTDIR= -DCMAKE_BUILD_TYPE=Release -DWITH_SYSTEMD=1 -DSYSTEMD_PID_DIR=${mysqlPath} ${MAYBE_WITH_BOOST}
if [ "$?" != "0" ]; then
    echo -e $HR
    echo "错误：MySQL 编译初始化失败，请截图错误信息寻求帮助。"
    exit 1
fi

make "-j${j}"
if [ "$?" != "0" ]; then
    echo -e $HR
    echo "错误：MySQL 编译失败，请截图错误信息寻求帮助。"
    exit 1
fi

# 安装
make install
if [ "$?" != "0" ]; then
    echo -e $HR
    echo "错误：MySQL 安装失败，请截图错误信息寻求帮助。"
    exit 1
fi

# 设置权限
chown -R mysql:mysql ${mysqlPath}
chmod -R 755 ${mysqlPath}
chmod 644 ${mysqlPath}/conf/my.cnf

# 启动服务
systemctl start mysqld

# 执行更新后的初始化
if [[ "${1}" == "57" ]]; then
    ${mysqlPath}/bin/mysql_upgrade -uroot -p${mysqlPassword}
fi
${mysqlPath}/bin/mysql -uroot -p${mysqlPassword} -e "DROP DATABASE test;"
${mysqlPath}/bin/mysql -uroot -p${mysqlPassword} -e "DELETE FROM mysql.user WHERE user='';"
${mysqlPath}/bin/mysql -uroot -p${mysqlPassword} -e "FLUSH PRIVILEGES;"

panel writePlugin mysql${1} ${mysqlVersion}

echo -e "${HR}\nMySQL-${1} 升级完成\n${HR}"
