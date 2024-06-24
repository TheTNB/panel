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
    dnf install cmake bison ncurses-devel libtirpc-devel openssl-devel pkg-config openldap-devel libudev-devel cyrus-sasl-devel patchelf rpcgen rpcsvc-proto-devel -y
    dnf install gcc-toolset-12-gcc gcc-toolset-12-gcc-c++ gcc-toolset-12-binutils gcc-toolset-12-annobin-annocheck gcc-toolset-12-annobin-plugin-gcc -y
elif [ "${OS}" == "debian" ]; then
    apt-get update
    apt-get install build-essential cmake bison libncurses5-dev libtirpc-dev libssl-dev pkg-config libldap2-dev libudev-dev libsasl2-dev patchelf -y
else
    echo -e $HR
    echo "错误：耗子面板不支持该系统"
    exit 1
fi

mysqlUserCheck=$(cat /etc/passwd | grep mysql)
if [ "${mysqlUserCheck}" == "" ]; then
    groupadd mysql
    useradd -s /sbin/nologin -g mysql mysql
fi

# 准备目录
rm -rf ${mysqlPath}
mkdir -p ${mysqlPath}
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
rm mysql-test/CMakeLists.txt
sed -i 's/ADD_SUBDIRECTORY(mysql-test)//g' CMakeLists.txt
mkdir build
cd build

# 5.7 和 8.0 需要 boost
if [[ "${1}" == "57" ]] || [[ "${1}" == "80" ]]; then
    MAYBE_WITH_BOOST="-DWITH_BOOST=../boost"
fi

cmake .. -DCMAKE_INSTALL_PREFIX=${mysqlPath} -DMYSQL_DATADIR=${mysqlPath}/data -DSYSCONFDIR=${mysqlPath}/conf -DWITH_MYISAM_STORAGE_ENGINE=1 -DWITH_INNOBASE_STORAGE_ENGINE=1 -DWITH_ARCHIVE_STORAGE_ENGINE=1 -DWITH_FEDERATED_STORAGE_ENGINE=1 -DWITH_BLACKHOLE_STORAGE_ENGINE=1 -DDEFAULT_CHARSET=utf8mb4 -DDEFAULT_COLLATION=utf8mb4_general_ci -DENABLED_LOCAL_INFILE=1 -DWITH_DEBUG=0 -DWITH_UNIT_TESTS=OFF -DINSTALL_MYSQLTESTDIR= -DCMAKE_BUILD_TYPE=Release -DWITH_SYSTEMD=1 -DSYSTEMD_PID_DIR=${mysqlPath} ${MAYBE_WITH_BOOST}
if [ "$?" != "0" ]; then
    echo -e $HR
    echo "错误：MySQL 编译初始化失败，请截图错误信息寻求帮助。"
    rm -rf ${mysqlPath}
    exit 1
fi

make "-j${j}"
if [ "$?" != "0" ]; then
    echo -e $HR
    echo "错误：MySQL 编译失败，请截图错误信息寻求帮助。"
    rm -rf ${mysqlPath}
    exit 1
fi

# 安装
make install
if [ "$?" != "0" ]; then
    echo -e $HR
    echo "错误：MySQL 安装失败，请截图错误信息寻求帮助。"
    rm -rf ${mysqlPath}
    exit 1
fi

# 配置
mkdir ${mysqlPath}/conf
cat > ${mysqlPath}/conf/my.cnf << EOF
[client]
port = 3306
socket = /tmp/mysql.sock

[mysqld]
port = 3306
socket = /tmp/mysql.sock
datadir = ${mysqlPath}/data
default_storage_engine = InnoDB
skip-external-locking
table_definition_cache = 400
performance_schema_max_table_instances = 400
key_buffer_size = 8M
max_allowed_packet = 1G
table_open_cache = 32
sort_buffer_size = 256K
net_buffer_length = 4K
read_buffer_size = 128K
read_rnd_buffer_size = 256K
myisam_sort_buffer_size = 4M
thread_cache_size = 4
query_cache_size = 4M
tmp_table_size = 8M
explicit_defaults_for_timestamp = 1
#skip-name-resolve
max_connections = 500
max_connect_errors = 100
open_files_limit = 65535
early-plugin-load = ""

log-bin = mysql-bin
server-id = 1
slow_query_log = 1
slow-query-log-file = ${mysqlPath}/mysql-slow.log
long_query_time = 3
log-error = ${mysqlPath}/mysql-error.log

innodb_data_home_dir = ${mysqlPath}/data
innodb_data_file_path = ibdata1:10M:autoextend
innodb_log_group_home_dir = ${mysqlPath}/data
innodb_buffer_pool_size = 16M
innodb_redo_log_capacity = 5M
innodb_log_buffer_size = 8M
innodb_flush_log_at_trx_commit = 1
innodb_lock_wait_timeout = 50
innodb_max_dirty_pages_pct = 90
innodb_read_io_threads = 4
innodb_write_io_threads = 4

[mysqldump]
quick
max_allowed_packet = 500M

[myisamchk]
key_buffer_size = 20M
sort_buffer_size = 20M
read_buffer = 2M
write_buffer = 2M

[mysqlhotcopy]
interactive-timeout
EOF

# 根据CPU核心数确定写入线程数
sed -i 's/innodb_write_io_threads = 4/innodb_write_io_threads = '${cpuCore}'/g' ${mysqlPath}/conf/my.cnf
sed -i 's/innodb_read_io_threads = 4/innodb_read_io_threads = '${cpuCore}'/g' ${mysqlPath}/conf/my.cnf

if [[ "${1}" == "84" ]] || [[ "${1}" == "80" ]]; then
    sed -i '/query_cache_size/d' ${mysqlPath}/conf/my.cnf
fi
if [[ "${1}" == "57" ]]; then
    sed -i '/innodb_redo_log_capacity/d' ${mysqlPath}/conf/my.cnf
fi

# 根据内存大小调参
if [[ ${memTotal} -gt 1024 && ${memTotal} -lt 2048 ]]; then
    sed -i "s#^key_buffer_size.*#key_buffer_size = 32M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^table_open_cache.*#table_open_cache = 128#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^sort_buffer_size.*#sort_buffer_size = 768K#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^read_buffer_size.*#read_buffer_size = 768K#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^myisam_sort_buffer_size.*#myisam_sort_buffer_size = 8M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^thread_cache_size.*#thread_cache_size = 16#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^query_cache_size.*#query_cache_size = 16M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^tmp_table_size.*#tmp_table_size = 32M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^innodb_buffer_pool_size.*#innodb_buffer_pool_size = 128M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^innodb_redo_log_capacity.*#innodb_redo_log_capacity = 64M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^innodb_log_buffer_size.*#innodb_log_buffer_size = 16M#" ${mysqlPath}/conf/my.cnf
elif [[ ${memTotal} -ge 2048 && ${memTotal} -lt 4096 ]]; then
    sed -i "s#^key_buffer_size.*#key_buffer_size = 64M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^table_open_cache.*#table_open_cache = 256#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^sort_buffer_size.*#sort_buffer_size = 1M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^read_buffer_size.*#read_buffer_size = 1M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^myisam_sort_buffer_size.*#myisam_sort_buffer_size = 16M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^thread_cache_size.*#thread_cache_size = 32#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^query_cache_size.*#query_cache_size = 32M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^tmp_table_size.*#tmp_table_size = 64M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^innodb_buffer_pool_size.*#innodb_buffer_pool_size = 256M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^innodb_redo_log_capacity.*#innodb_redo_log_capacity = 128M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^innodb_log_buffer_size.*#innodb_log_buffer_size = 32M#" ${mysqlPath}/conf/my.cnf
elif [[ ${memTotal} -ge 4096 && ${memTotal} -lt 8192 ]]; then
    sed -i "s#^key_buffer_size.*#key_buffer_size = 128M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^table_open_cache.*#table_open_cache = 512#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^sort_buffer_size.*#sort_buffer_size = 2M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^read_buffer_size.*#read_buffer_size = 2M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^myisam_sort_buffer_size.*#myisam_sort_buffer_size = 32M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^thread_cache_size.*#thread_cache_size = 64#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^query_cache_size.*#query_cache_size = 64M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^tmp_table_size.*#tmp_table_size = 64M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^innodb_buffer_pool_size.*#innodb_buffer_pool_size = 512M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^innodb_redo_log_capacity.*#innodb_redo_log_capacity = 256M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^innodb_log_buffer_size.*#innodb_log_buffer_size = 64M#" ${mysqlPath}/conf/my.cnf
elif [[ ${memTotal} -ge 8192 && ${memTotal} -lt 16384 ]]; then
    sed -i "s#^key_buffer_size.*#key_buffer_size = 256M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^table_open_cache.*#table_open_cache = 1024#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^sort_buffer_size.*#sort_buffer_size = 4M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^read_buffer_size.*#read_buffer_size = 4M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^myisam_sort_buffer_size.*#myisam_sort_buffer_size = 64M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^thread_cache_size.*#thread_cache_size = 128#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^query_cache_size.*#query_cache_size = 128M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^tmp_table_size.*#tmp_table_size = 128M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^innodb_buffer_pool_size.*#innodb_buffer_pool_size = 1024M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^innodb_redo_log_capacity.*#innodb_redo_log_capacity = 512M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^innodb_log_buffer_size.*#innodb_log_buffer_size = 128M#" ${mysqlPath}/conf/my.cnf
elif [[ ${memTotal} -ge 16384 && ${memTotal} -lt 32768 ]]; then
    sed -i "s#^key_buffer_size.*#key_buffer_size = 512M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^table_open_cache.*#table_open_cache = 2048#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^sort_buffer_size.*#sort_buffer_size = 8M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^read_buffer_size.*#read_buffer_size = 8M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^myisam_sort_buffer_size.*#myisam_sort_buffer_size = 128M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^thread_cache_size.*#thread_cache_size = 256#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^query_cache_size.*#query_cache_size = 256M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^tmp_table_size.*#tmp_table_size = 256M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^innodb_buffer_pool_size.*#innodb_buffer_pool_size = 2048M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^innodb_redo_log_capacity.*#innodb_redo_log_capacity = 1G#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^innodb_log_buffer_size.*#innodb_log_buffer_size = 256M#" ${mysqlPath}/conf/my.cnf
elif [[ ${memTotal} -ge 32768 ]]; then
    sed -i "s#^key_buffer_size.*#key_buffer_size = 1024M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^table_open_cache.*#table_open_cache = 4096#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^sort_buffer_size.*#sort_buffer_size = 16M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^read_buffer_size.*#read_buffer_size = 16M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^myisam_sort_buffer_size.*#myisam_sort_buffer_size = 256M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^thread_cache_size.*#thread_cache_size = 512#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^query_cache_size.*#query_cache_size = 512M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^tmp_table_size.*#tmp_table_size = 512M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^innodb_buffer_pool_size.*#innodb_buffer_pool_size = 4096M#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^innodb_redo_log_capacity.*#innodb_redo_log_capacity = 2G#" ${mysqlPath}/conf/my.cnf
    sed -i "s#^innodb_log_buffer_size.*#innodb_log_buffer_size = 512M#" ${mysqlPath}/conf/my.cnf
fi

# 初始化
cd ${mysqlPath}
rm -rf ${mysqlPath}/src
rm -rf ${mysqlPath}/data
mkdir -p ${mysqlPath}/data
chown -R mysql:mysql ${mysqlPath}
chmod -R 755 ${mysqlPath}
chmod 644 ${mysqlPath}/conf/my.cnf

${mysqlPath}/bin/mysqld --initialize-insecure --user=mysql --basedir=${mysqlPath} --datadir=${mysqlPath}/data

# 软链接
ln -sf ${mysqlPath}/bin/* /usr/bin/

# 写入 systemd 配置
cat > /etc/systemd/system/mysqld.service << EOF
[Unit]
Description=MySQL Server
Documentation=http://dev.mysql.com/doc/refman/en/using-systemd.html
After=network.target
After=syslog.target

[Install]
WantedBy=multi-user.target

[Service]
User=mysql
Group=mysql

Type=forking

PIDFile=/www/server/mysql/mysqld.pid

# Disable service start and stop timeout logic of systemd for mysqld service.
TimeoutSec=0

# Execute pre and post scripts as root
PermissionsStartOnly=true

# Start main service
ExecStart=/www/server/mysql/bin/mysqld --daemonize --pid-file=/www/server/mysql/mysqld.pid \$MYSQLD_OPTS

# Use this to switch malloc implementation
EnvironmentFile=-/etc/sysconfig/mysql

# Sets open_files_limit
LimitNOFILE = 500000

Restart=on-failure

RestartPreventExitStatus=1

PrivateTmp=false
EOF

systemctl daemon-reload
systemctl enable mysqld
systemctl start mysqld

mysqlPassword=$(cat /dev/urandom | head -n 16 | md5sum | head -c 16)
${mysqlPath}/bin/mysqladmin -u root password ${mysqlPassword}
${mysqlPath}/bin/mysql -uroot -p${mysqlPassword} -e "DROP DATABASE test;"
${mysqlPath}/bin/mysql -uroot -p${mysqlPassword} -e "DELETE FROM mysql.user WHERE user='';"
${mysqlPath}/bin/mysql -uroot -p${mysqlPassword} -e "FLUSH PRIVILEGES;"

panel writePlugin mysql${1} ${mysqlVersion}
panel writeMysqlPassword ${mysqlPassword}

echo -e "${HR}\nMySQL-${1} 安装完成\n${HR}"
