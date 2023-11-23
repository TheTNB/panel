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
    dnf install -y rsync
elif [ "${OS}" == "debian" ]; then
    apt-get install -y rsync
else
    echo -e $HR
    echo "错误：不支持的操作系统"
    exit 1
fi

# 写入配置
cat > /etc/rsyncd.conf << EOF
uid = root
gid = root
port = 873
use chroot = no
read only = no
dont compress = *.jpg *.jpeg *.png *.gif *.webp *.avif *.mp4 *.avi *.mov *.mkv *.mp3 *.wav *.aac *.flac *.zip *.rar *.7z *.gz *.tgz *.tar *.pdf *.epub *.iso *.exe *.apk *.dmg *.rpm *.deb *.msi
hosts allow = 127.0.0.1/32 ::1/128
# hosts deny =
max connections = 100
timeout = 1800
lock file = /var/run/rsync.lock
pid file = /var/run/rsyncd.pid
log file = /var/log/rsyncd.log

EOF

touch /etc/rsyncd.secrets
chmod 644 /etc/rsyncd.conf
chmod 600 /etc/rsyncd.secrets

# 写入服务文件
cat > /etc/systemd/system/rsyncd.service << EOF
[Unit]
Description=fast remote file copy program daemon
After=network-online.target remote-fs.target nss-lookup.target
Wants=network-online.target
ConditionPathExists=/etc/rsyncd.conf

[Service]
ExecStart=/usr/bin/rsync --daemon --no-detach "\$OPTIONS"
ExecReload=/bin/kill -HUP \$MAINPID
KillMode=process
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable rsyncd.service
systemctl restart rsyncd.service

panel writePlugin rsync 3.2.7
