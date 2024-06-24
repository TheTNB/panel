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

panel deletePlugin mysql${1}

echo -e "${HR}\nMySQL-${1} 卸载完成\n${HR}"
