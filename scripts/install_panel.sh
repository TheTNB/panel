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

LOGO="+----------------------------------------------------\n| 耗子Linux面板安装脚本\n+----------------------------------------------------\n| Copyright © 2022-"$(date +%Y)" 耗子科技 All rights reserved.\n+----------------------------------------------------"
HR="+----------------------------------------------------"
setup_Path="/www"
sshPort=$(cat /etc/ssh/sshd_config | grep 'Port ' | awk '{print $2}')
inChina=$(curl --retry 2 -m 10 -L https://www.cloudflare-cn.com/cdn-cgi/trace 2> /dev/null | grep -qx 'loc=CN' && echo "true" || echo "false")

Prepare_System() {
    if [ $(whoami) != "root" ]; then
        echo -e $HR
        echo "错误：请使用root用户运行安装命令。"
        exit 1
    fi

    ARCH=$(uname -m)
    OS=$(source /etc/os-release && { [[ "$ID" == "debian" ]] && echo "debian"; } || { [[ "$ID" == "centos" ]] || [[ "$ID" == "rhel" ]] || [[ "$ID" == "rocky" ]] || [[ "$ID" == "almalinux" ]] && echo "centos"; } || echo "unknown")
    if [ "${OS}" == "unknown" ]; then
        echo -e $HR
        echo "错误：该系统不支持安装耗子Linux面板，请更换Debian12/RHEL9安装。"
        exit 1
    fi
    if [ "${ARCH}" != "x86_64" ] && [ "${ARCH}" != "aarch64" ]; then
        echo -e $HR
        echo "错误：该系统架构不支持安装耗子Linux面板，请更换x86_64/aarch64架构安装。"
        exit 1
    fi

    is64bit=$(getconf LONG_BIT)
    if [ "${is64bit}" != '64' ]; then
        echo -e $HR
        echo "错误：32位系统不支持安装耗子Linux面板，请更换64位系统安装。"
        exit 1
    fi

    isInstalled=$(systemctl status panel 2>&1 | grep "Active")
    if [ "${isInstalled}" != "" ]; then
        echo -e $HR
        echo "错误：耗子Linux面板已安装，请勿重复安装。"
        exit 1
    fi

    if ! id -u "www" > /dev/null 2>&1; then
        groupadd www
        useradd -s /sbin/nologin -g www www
    fi

    if [ ! -d ${setup_Path} ]; then
        mkdir ${setup_Path}
    fi

    timedatectl set-timezone Asia/Shanghai

    [ -s /etc/selinux/config ] && sed -i 's/SELINUX=enforcing/SELINUX=disabled/g' /etc/selinux/config
    setenforce 0 > /dev/null 2>&1

    ulimit -n 204800
    echo 6553560 > /proc/sys/fs/file-max
    checkSoftNofile=$(cat /etc/security/limits.conf | grep '^* soft nofile .*$')
    checkHardNofile=$(cat /etc/security/limits.conf | grep '^* hard nofile .*$')
    checkSoftNproc=$(cat /etc/security/limits.conf | grep '^* soft nproc .*$')
    checkHardNproc=$(cat /etc/security/limits.conf | grep '^* hard nproc .*$')
    checkFsFileMax=$(cat /etc/sysctl.conf | grep '^fs.file-max.*$')
    if [ "${checkSoftNofile}" == "" ]; then
        echo "* soft nofile 204800" >> /etc/security/limits.conf
    fi
    if [ "${checkHardNofile}" == "" ]; then
        echo "* hard nofile 204800" >> /etc/security/limits.conf
    fi
    if [ "${checkSoftNproc}" == "" ]; then
        echo "* soft nproc 204800" >> /etc/security/limits.conf
    fi
    if [ "${checkHardNproc}" == "" ]; then
        echo "* hard nproc 204800 " >> /etc/security/limits.conf
    fi
    if [ "${checkFsFileMax}" == "" ]; then
        echo fs.file-max = 6553560 >> /etc/sysctl.conf
    fi

    if [ "${OS}" == "centos" ]; then
        if ${inChina}; then
            sed -e 's|^mirrorlist=|#mirrorlist=|g' \
                -e 's|^#baseurl=http://dl.rockylinux.org/$contentdir|baseurl=https://mirrors.cloud.tencent.com/rocky|g' \
                -i.bak \
                /etc/yum.repos.d/[Rr]ocky*.repo
            sed -e 's|^mirrorlist=|#mirrorlist=|g' \
                -e 's|^# baseurl=https://repo.almalinux.org|baseurl=https://mirrors.cloud.tencent.com|g' \
                -i.bak \
                /etc/yum.repos.d/[Aa]lmalinux*.repo

            dnf makecache -y
        fi
        dnf install dnf-plugins-core -y
        dnf install epel-release -y
        dnf config-manager --set-enabled epel
        if ${inChina}; then
            sed -i 's|^#baseurl=https://download.example/pub|baseurl=https://mirrors.cloud.tencent.com|' /etc/yum.repos.d/epel*
            sed -i 's|^metalink|#metalink|' /etc/yum.repos.d/epel*
            dnf makecache -y
        fi
        dnf config-manager --set-enabled PowerTools
        dnf config-manager --set-enabled powertools
        dnf config-manager --set-enabled CRB
        dnf config-manager --set-enabled Crb
        dnf config-manager --set-enabled crb
        /usr/bin/crb enable
        dnf makecache -y
        dnf install -y curl wget zip unzip tar git jq git-core dos2unix
    elif [ "${OS}" == "debian" ]; then
        if ${inChina}; then
            sed -i 's/deb.debian.org/mirrors.aliyun.com/g' /etc/apt/sources.list
            sed -i 's/security.debian.org/mirrors.aliyun.com/g' /etc/apt/sources.list
        fi
        apt-get update -y
        apt-get install -y curl wget zip unzip tar git jq git dos2unix rsyslog
    else
        echo -e $HR
        echo "错误：该系统不支持安装耗子Linux面板，请更换Debian12/RHEL9安装。"
        exit 1
    fi

    if [ "$?" != "0" ]; then
        echo -e $HR
        echo "错误：安装面板依赖软件失败，请截图错误信息寻求帮助。"
        exit 1
    fi
}

Auto_Swap() {
    # 判断是否有swap
    swap=$(LC_ALL=C free | grep Swap | awk '{print $2}')
    if [ "${swap}" -gt 1 ]; then
        return
    fi

    # 设置swap
    swapFile="${setup_Path}/swap"
    btrfsCheck=$(df -T /www | awk '{print $2}' | tail -n 1)
    if [ "${btrfsCheck}" == "btrfs" ]; then
        btrfs filesystem mkswapfile --size 4G --uuid clear ${swapFile}
    else
        dd if=/dev/zero of=$swapFile bs=1M count=4096
    fi
    chmod 600 $swapFile
    mkswap -f $swapFile
    swapon $swapFile
    echo "$swapFile    swap    swap    defaults    0 0" >> /etc/fstab

    mount -a
    if [ "$?" != "0" ]; then
        echo -e $HR
        echo "错误：检测到系统的 /etc/fstab 文件配置有误，请检查排除后重试，问题解决前勿重启系统。"
        exit 1
    fi
}

Init_Panel() {
    mkdir ${setup_Path}/server
    mkdir ${setup_Path}/server/cron
    mkdir ${setup_Path}/server/cron/logs
    chmod -R 755 ${setup_Path}/server
    mkdir ${setup_Path}/panel
    rm -rf ${setup_Path}/panel/*
    # 下载面板zip包并解压
    if [ "${ARCH}" == "x86_64" ]; then
        if ${inChina}; then
            panelZip=$(curl -sSL "https://jihulab.com/api/v4/projects/haozi-team%2Fpanel/releases/permalink/latest" | jq -r '.assets.links[] | select(.name | contains("amd64v2")) | .direct_asset_url')
            panelZipName=$(curl -sSL "https://jihulab.com/api/v4/projects/haozi-team%2Fpanel/releases/permalink/latest" | jq -r '.assets.links[] | select(.name | contains("amd64v2")) | .name')
        else
            panelZip=$(curl -sSL "https://api.github.com/repos/haozi-team/panel/releases/latest" | jq -r '.assets[] | select(.name | contains("amd64v2")) | .browser_download_url')
            panelZipName=$(curl -sSL "https://api.github.com/repos/haozi-team/panel/releases/latest" | jq -r '.assets[] | select(.name | contains("amd64v2")) | .name')
        fi
    elif [ "${ARCH}" == "aarch64" ]; then
        if ${inChina}; then
            panelZip=$(curl -sSL "https://jihulab.com/api/v4/projects/haozi-team%2Fpanel/releases/permalink/latest" | jq -r '.assets.links[] | select(.name | contains("arm64")) | .direct_asset_url')
            panelZipName=$(curl -sSL "https://jihulab.com/api/v4/projects/haozi-team%2Fpanel/releases/permalink/latest" | jq -r '.assets.links[] | select(.name | contains("arm64")) | .name')
        else
            panelZip=$(curl -sSL "https://api.github.com/repos/haozi-team/panel/releases/latest" | jq -r '.assets[] | select(.name | contains("arm64")) | .browser_download_url')
            panelZipName=$(curl -sSL "https://api.github.com/repos/haozi-team/panel/releases/latest" | jq -r '.assets[] | select(.name | contains("arm64")) | .name')
        fi
    else
        echo -e $HR
        echo "错误：该系统架构不支持安装耗子Linux面板，请更换x86_64/aarch64架构安装。"
        exit 1
    fi
    if [ "$?" != "0" ] || [ "${panelZip}" == "" ] || [ "${panelZipName}" == "" ]; then
        echo -e $HR
        echo "错误：获取面板下载链接失败，请截图错误信息寻求帮助。"
        exit 1
    fi
    wget -T 120 -t 3 -O ${setup_Path}/panel/${panelZipName} "${panelZip}"

    # 下载 checksums 文件
    if ${inChina}; then
        checksumsFile=$(curl -sSL "https://jihulab.com/api/v4/projects/haozi-team%2Fpanel/releases/permalink/latest" | jq -r '.assets.links[] | select(.name | contains("checksums")) | .direct_asset_url')
        checksumsFileName=$(curl -sSL "https://jihulab.com/api/v4/projects/haozi-team%2Fpanel/releases/permalink/latest" | jq -r '.assets.links[] | select(.name | contains("checksums")) | .name')
    else
        checksumsFile=$(curl -sSL "https://api.github.com/repos/haozi-team/panel/releases/latest" | jq -r '.assets[] | select(.name | contains("checksums")) | .browser_download_url')
        checksumsFileName=$(curl -sSL "https://api.github.com/repos/haozi-team/panel/releases/latest" | jq -r '.assets[] | select(.name | contains("checksums")) | .name')
    fi
    wget -T 20 -t 3 -O ${setup_Path}/panel/${checksumsFileName} "${checksumsFile}"
    cd ${setup_Path}/panel
    if ! sha256sum --status -c ${checksumsFileName}; then
        echo -e $HR
        echo "错误：面板压缩包 checksum 校验失败，文件可能被篡改或不完整，已终止操作"
        exit 1
    fi
    unzip -o ${panelZipName}
    if [ "$?" != "0" ]; then
        echo -e $HR
        echo "错误：解压面板失败，请截图错误信息寻求帮助。"
        exit 1
    fi
    rm -rf ${panelZipName}
    rm -rf ${checksumsFileName}
    cp panel-example.conf panel.conf
    ${setup_Path}/panel/panel --env="panel.conf" artisan key:generate
    ${setup_Path}/panel/panel --env="panel.conf" artisan jwt:secret
    ${setup_Path}/panel/panel --env="panel.conf" artisan migrate
    chmod -R 700 ${setup_Path}/panel
    cp scripts/panel.sh /usr/bin/panel
    chmod -R 700 /usr/bin/panel
    # 防火墙放行
    if [ "${OS}" == "centos" ]; then
        yum install firewalld -y
        systemctl enable firewalld
        systemctl start firewalld
        firewall-cmd --set-default-zone=public > /dev/null 2>&1
        firewall-cmd --permanent --zone=public --add-port=22/tcp > /dev/null 2>&1
        firewall-cmd --permanent --zone=public --add-port=80/tcp > /dev/null 2>&1
        firewall-cmd --permanent --zone=public --add-port=443/tcp > /dev/null 2>&1
        firewall-cmd --permanent --zone=public --add-port=8888/tcp > /dev/null 2>&1
        firewall-cmd --permanent --zone=public --add-port=${sshPort}/tcp > /dev/null 2>&1
        firewall-cmd --reload
    elif [ "${OS}" == "debian" ]; then
        apt-get install ufw -y
        echo y | ufw enable
        ufw allow 22/tcp
        ufw allow 80/tcp
        ufw allow 443/tcp
        ufw allow 8888/tcp
        ufw allow ${sshPort}/tcp
        ufw reload
    fi
    if [ "$?" != "0" ]; then
        echo -e $HR
        echo "错误：防火墙放行失败，请截图错误信息寻求帮助。"
        exit 1
    fi
    # 写入服务文件
    cat > /etc/systemd/system/panel.service << EOF
[Unit]
Description=HaoZi Panel
After=syslog.target network.target
Wants=network.target

[Service]
Type=simple
WorkingDirectory=${setup_Path}/panel/
ExecStart=/www/panel/panel --env="/www/panel/panel.conf"
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
    if [ "$?" != "0" ]; then
        echo -e $HR
        echo "错误：面板启动失败，请截图错误信息寻求帮助。"
        exit 1
    fi

    clear
    echo -e $LOGO
    echo '面板安装成功！'
    echo -e $HR
    panel init
    panel getInfo
    rm -f install_panel.sh
    rm -f install_panel.sh.checksum.txt
}

clear
echo -e $LOGO

# 安装确认
read -p "面板将安装至 ${setup_Path} 目录，请输入 y 并回车以开始安装：" install
if [ "$install" != 'y' ]; then
    echo "输入不正确，已退出安装。"
    exit
fi

clear
echo -e $LOGO
echo '安装面板依赖软件（如报错请检查 APT/Yum 源是否正常）'
echo -e $HR
sleep 2s
Prepare_System
Auto_Swap

echo -e $LOGO
echo '安装面板运行环境（视网络情况可能需要较长时间）'
echo -e $HR
sleep 2s
Init_Panel
