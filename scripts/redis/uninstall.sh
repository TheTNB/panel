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

systemctl stop redis
systemctl disable redis
rm -rf /etc/systemd/system/redis.service
systemctl daemon-reload
pkill -9 redis
rm -rf /www/server/redis

rm -rf /usr/bin/redis-cli
rm -rf /usr/bin/redis-server
rm -rf /usr/bin/redis-benchmark
rm -rf /usr/bin/redis-check-aof
rm -rf /usr/bin/redis-check-rdb
rm -rf /usr/bin/redis-sentinel

panel deletePlugin redis

echo -e "${HR}\nRedis 卸载完成\n${HR}"
