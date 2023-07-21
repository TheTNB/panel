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

systemctl stop mysqld
systemctl disable mysqld
rm -rf /etc/systemd/system/mysqld.service
systemctl daemon-reload
pkill -9 mysqld
rm -rf /www/server/mysql

rm -f /usr/bin/mysql*
rm -f /usr/lib/libmysql*
rm -f /usr/lib64/libmysql*

userdel -r mysql
groupdel mysql

sed -i '#export PATH=/www/server/mysql#d' /etc/profile
source /etc/profile

panel deletePlugin mysql${1}

echo -e "${HR}\nMySQL-${1} 卸载完成\n${HR}"
