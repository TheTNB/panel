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

LOGO="+----------------------------------------------------\n| 耗子面板安装脚本\n+----------------------------------------------------\n| Copyright © 2022-"$(date +%Y)" 耗子科技 All rights reserved.\n+----------------------------------------------------"
HR="+----------------------------------------------------"
download_Url=""
setup_Path="/www"
sshPort=$(cat /etc/ssh/sshd_config | grep 'Port ' | awk '{print $2}')
ipLocation=$(curl -s https://ip.ping0.cc/geo)

Prepare_system() {
    if [ $(whoami) != "root" ]; then
        echo -e $HR
        echo "错误：请使用root用户运行安装命令。"
        exit 1
    fi

    ARCH=$(uname -m)
    OS=$(source /etc/os-release && { [[ "$ID" == "debian" ]] && echo "debian"; } || { [[ "$ID" == "centos" ]] || [[ "$ID" == "rhel" ]] || [[ "$ID" == "rocky" ]] || [[ "$ID" == "almalinux" ]] && echo "centos"; } || echo "unknown")
    if [ "${OS}" == "unknown" ]; then
        echo -e $HR
        echo "错误：该系统不支持安装耗子面板，请更换Debian12/RHEL9安装。"
        exit 1
    fi
    if [ "${ARCH}" != "x86_64" ] && [ "${ARCH}" != "aarch64" ]; then
        echo -e $HR
        echo "错误：该系统架构不支持安装耗子面板，请更换x86_64/aarch64架构安装。"
        exit 1
    fi

    is64bit=$(getconf LONG_BIT)
    if [ "${is64bit}" != '64' ]; then
        echo -e $HR
        echo "错误：32位系统不支持安装耗子面板，请更换64位系统安装。"
        exit 1
    fi

    isInstalled=$(systemctl status panel 2>&1 | grep "Active")
    if [ "${isInstalled}" != "" ]; then
        echo -e $HR
        echo "错误：耗子面板已安装，请勿重复安装。"
        exit 1
    fi

    wwwUserCheck=$(cat /etc/passwd | grep www)
    if [ "${wwwUserCheck}" == "" ]; then
        groupadd www
        useradd -s /sbin/nologin -g www www
    fi

    rm -rf /etc/localtime
    ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

    [ -s /etc/selinux/config ] && sed -i 's/SELINUX=enforcing/SELINUX=disabled/g' /etc/selinux/config
    setenforce 0 >/dev/null 2>&1

    ulimit -n 204800
    echo 6553560 >/proc/sys/fs/file-max
    checkSoftNofile=$(cat /etc/security/limits.conf | grep '^* soft nofile .*$')
    checkHardNofile=$(cat /etc/security/limits.conf | grep '^* hard nofile .*$')
    checkSoftNproc=$(cat /etc/security/limits.conf | grep '^* soft nproc .*$')
    checkHardNproc=$(cat /etc/security/limits.conf | grep '^* hard nproc .*$')
    checkFsFileMax=$(cat /etc/sysctl.conf | grep '^fs.file-max.*$')
    if [ "${checkSoftNofile}" == "" ]; then
        echo "* soft nofile 204800" >>/etc/security/limits.conf
    fi
    if [ "${checkHardNofile}" == "" ]; then
        echo "* hard nofile 204800" >>/etc/security/limits.conf
    fi
    if [ "${checkSoftNproc}" == "" ]; then
        echo "* soft nproc 204800" >>/etc/security/limits.conf
    fi
    if [ "${checkHardNproc}" == "" ]; then
        echo "* hard nproc 204800 " >>/etc/security/limits.conf
    fi
    if [ "${checkFsFileMax}" == "" ]; then
        echo fs.file-max = 6553560 >>/etc/sysctl.conf
    fi

    if [ "${OS}" == "centos" ]; then
        if [[ ${ipLocation} =~ "中国" ]]; then
            sed -e 's|^mirrorlist=|#mirrorlist=|g' \
                -e 's|^#baseurl=http://dl.rockylinux.org/$contentdir|baseurl=https://mirrors.aliyun.com/rockylinux|g' \
                -i.bak \
                /etc/yum.repos.d/[Rr]ocky*.repo
            sed -e 's|^mirrorlist=|#mirrorlist=|g' \
                -e 's|^# baseurl=https://repo.almalinux.org|baseurl=https://mirrors.aliyun.com|g' \
                -i.bak \
                /etc/yum.repos.d/[Aa]lmalinux*.repo

            dnf makecache
        fi
        dnf install dnf-plugins-core -y
        dnf install epel-release -y
        dnf config-manager --set-enabled epel
        if [[ ${ipLocation} =~ "中国" ]]; then
            sed -i 's|^#baseurl=https://download.example/pub|baseurl=https://mirrors.aliyun.com|' /etc/yum.repos.d/epel*
            sed -i 's|^metalink|#metalink|' /etc/yum.repos.d/epel*
            dnf makecache
        fi
        dnf config-manager --set-enabled PowerTools
        dnf config-manager --set-enabled powertools
        dnf config-manager --set-enabled CRB
        dnf config-manager --set-enabled Crb
        dnf config-manager --set-enabled crb
        /usr/bin/crb enable
        dnf makecache
        dnf install -y curl wget zip unzip tar git jq git-core
    elif [ "${OS}" == "debian" ]; then
        if [[ ${ipLocation} =~ "中国" ]]; then
            sed -i 's/deb.debian.org/mirrors.aliyun.com/g' /etc/apt/sources.list
            sed -i 's/security.debian.org/mirrors.aliyun.com/g' /etc/apt/sources.list
        fi
        apt update
        apt install -y curl wget zip unzip tar git jq git
    else
        echo -e $HR
        echo "错误：该系统不支持安装耗子面板，请更换Debian12/RHEL9安装。"
        exit 1
    fi
}

Auto_Swap() {
    # 判断是否有swap
    swap=$(LC_ALL=C free | grep Swap | awk '{print $2}')
    if [ "${swap}" -gt 1 ]; then
        return
    fi

    if [ ! -d ${setup_Path} ]; then
        mkdir ${setup_Path}
    fi

    # 设置swap
    swapFile="${setup_Path}/swap"
    dd if=/dev/zero of=$swapFile bs=1M count=2048
    chmod 600 $swapFile
    mkswap -f $swapFile
    swapon $swapFile
    echo "$swapFile    swap    swap    defaults    0 0" >>/etc/fstab
}

Init_Panel() {
    mkdir ${setup_Path}/server/cron
    mkdir ${setup_Path}/server/cron/logs
    chmod -R 644 ${setup_Path}/server/cron
    chmod -R 644 ${setup_Path}/server/cron/logs
    mkdir ${setup_Path}/panel
    rm -rf ${setup_Path}/panel/*
    # 下载面板zip包并解压
    if [ "${ARCH}" == "x86_64" ]; then
        panelZip=$(curl "https://api.github.com/repos/HaoZi-Team/Panel/releases/latest" | jq -r '.assets[] | select(.name | contains("amd64v3")) | .browser_download_url')
    elif [ "${ARCH}" == "aarch64" ]; then
        panelZip=$(curl "https://api.github.com/repos/HaoZi-Team/Panel/releases/latest" | jq -r '.assets[] | select(.name | contains("arm64")) | .browser_download_url')
    else
        echo -e $HR
        echo "错误：该系统架构不支持安装耗子面板，请更换x86_64/aarch64架构安装。"
        exit 1
    fi
    wget -O ${setup_Path}/panel/panel.zip "${download_Url}${panelZip}"
    cd ${setup_Path}/panel
    unzip -o panel.zip
    rm -rf panel.zip
    cp panel-example.conf panel.conf
    ${setup_Path}/panel/panel --env="panel.conf" artisan key:generate
    ${setup_Path}/panel/panel --env="panel.conf" artisan jwt:secret
    chmod -R 700 ${setup_Path}/panel
    cp scripts/panel.sh /usr/bin/panel
    chmod -R 700 /usr/bin/panel
    # 防火墙放行
    if [ "${OS}" == "centos" ]; then
        yum install firewalld -y
        systemctl enable firewalld
        systemctl start firewalld
        firewall-cmd --set-default-zone=public >/dev/null 2>&1
        firewall-cmd --permanent --zone=public --add-port=22/tcp >/dev/null 2>&1
        firewall-cmd --permanent --zone=public --add-port=80/tcp >/dev/null 2>&1
        firewall-cmd --permanent --zone=public --add-port=443/tcp >/dev/null 2>&1
        firewall-cmd --permanent --zone=public --add-port=8888/tcp >/dev/null 2>&1
        firewall-cmd --permanent --zone=public --add-port=${sshPort}/tcp >/dev/null 2>&1
        firewall-cmd --reload
    elif [ "${OS}" == "debian" ]; then
        apt install ufw -y
        sudo ufw enable
        sudo ufw allow 22/tcp
        sudo ufw allow 80/tcp
        sudo ufw allow 443/tcp
        sudo ufw allow 8888/tcp
        sudo ufw allow ${sshPort}/tcp
    fi
    # 写入服务文件
    cat >/etc/systemd/system/panel.service <<EOF
[Unit]
Description=HaoZi Panel
After=syslog.target network.target
Wants=network.target

[Service]
Type=simple
WorkingDirectory=${setup_Path}/panel/
ExecStart=./panel --env="panel.conf"
ExecReload=kill -s HUP \$MAINPID
ExecStop=kill -s QUIT \$MAINPID
User=root
Restart=always

[Install]
WantedBy=multi-user.target
EOF
    systemctl daemon-reload
    systemctl enable panel.service
    systemctl start panel.service

    clear
    echo -e $LOGO
    echo '面板安装成功！'
    echo -e $HR
    panel init
    panel getInfo
}

clear
echo -e $LOGO

# 安装确认
read -p "面板将安装至${setup_Path}目录，请输入 y 并回车以开始安装：" install
if [ "$install" != 'y' ]; then
    echo "输入不正确，已退出安装。"
    exit
fi

#代理设置
read -p "是否使用GitHub代理安装？(y/n)" proxy
if [ "$proxy" == 'y' ]; then
    download_Url="https://ghproxy.com/"
fi

clear
echo -e $LOGO
echo '安装面板依赖软件（如报错请检查 Dnf/Yum 源是否正常）'
echo -e $HR
sleep 3s
Prepare_system
Auto_Swap

clear
echo -e $LOGO
echo '安装面板运行环境（视网络情况可能需要较长时间）'
echo -e $HR
sleep 3s
Init_Panel
