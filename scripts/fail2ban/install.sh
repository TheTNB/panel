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

if [ "${OS}" == "centos" ]; then
    dnf install -y fail2ban
elif [ "${OS}" == "debian" ]; then
    apt-get install -y fail2ban
else
    echo -e $HR
    echo "错误：不支持的操作系统"
    exit 1
fi

if [ "$?" != "0" ]; then
    echo -e $HR
    echo "错误：fail2ban安装失败，请截图错误信息寻求帮助。"
    exit 1
fi

# 修改 fail2ban 配置文件
sed -i 's!# logtarget.*!logtarget = /var/log/fail2ban.log!' /etc/fail2ban/fail2ban.conf
sed -i 's!logtarget\s*=.*!logtarget = /var/log/fail2ban.log!' /etc/fail2ban/jail.conf
cat > /etc/fail2ban/jail.local << EOF
[DEFAULT]
ignoreip = 127.0.0.1/8
bantime = 600
findtime = 300
maxretry = 5
banaction = firewallcmd-ipset
action = %(action_mwl)s

# ssh-START
[ssh]
enabled = true
filter = sshd
port = 22
maxretry = 5
findtime = 300
bantime = 86400
action = %(action_mwl)s
logpath = /var/log/secure
# ssh-END
EOF
# 替换端口
sshPort=$(cat /etc/ssh/sshd_config | grep 'Port ' | awk '{print $2}')
if [ "${sshPort}" == "" ]; then
    sshPort="22"
fi
sed -i "s/port = 22/port = ${sshPort}/g" /etc/fail2ban/jail.local

# Debian 的特殊处理
if [ "${OS}" == "debian" ]; then
    sed -i "s/\/var\/log\/secure/\/var\/log\/auth.log/g" /etc/fail2ban/jail.local
    sed -i "s/banaction = firewallcmd-ipset/banaction = ufw/g" /etc/fail2ban/jail.local
fi

# 启动 fail2ban
systemctl unmask fail2ban
systemctl daemon-reload
systemctl enable fail2ban
systemctl restart fail2ban

panel writePlugin fail2ban 1.0.0
