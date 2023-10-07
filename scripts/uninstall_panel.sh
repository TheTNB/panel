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

LOGO="+----------------------------------------------------\n| 耗子Linux面板卸载脚本\n+----------------------------------------------------\n| Copyright © 2022-"$(date +%Y)" 耗子科技 All rights reserved.\n+----------------------------------------------------"
HR="+----------------------------------------------------"
download_Url=""
setup_Path="/www"

Prepare_System() {
    if [ $(whoami) != "root" ]; then
        echo -e $HR
        echo "错误：请使用root用户运行卸载命令。"
        exit 1
    fi

    isInstalled=$(systemctl status panel 2>&1 | grep "Active")
    if [ "${isInstalled}" == "" ]; then
        echo -e $HR
        echo "错误：耗子Linux面板未安装，无需卸载。"
        exit 1
    fi

    if ! id -u "www" > /dev/null 2>&1; then
        groupadd www
        useradd -s /sbin/nologin -g www www
    fi
}

Remove_Swap() {
    swapFile="${setup_Path}/swap"
    if [ -f "${swapFile}" ]; then
        swapoff ${swapFile}
        rm -f ${swapFile}
        sed -i '/swap/d' /etc/fstab
    fi

    mount -a
    if [ "$?" != "0" ]; then
        echo -e $HR
        echo "错误：检测到系统的 /etc/fstab 文件配置有误，请检查排除后重试，问题解决前勿重启系统。"
        exit 1
    fi
}

Remove_Panel() {
    systemctl stop panel
    systemctl disable panel
    rm -f /etc/systemd/system/panel.service
    rm -rf ${setup_Path}
}

clear
echo -e "${LOGO}"

# 卸载确认
echo -e "高危操作，卸载面板前请务必备份好所有数据，提前卸载面板所有插件！"
echo -e "卸载面板后，所有数据将被清空，无法恢复！"
read -r -p "输入 y 并回车以确认卸载面板：" uninstall
if [ "${uninstall}" != 'y' ]; then
    echo "输入不正确，已退出卸载。"
    exit
fi

echo -e "${LOGO}"
echo '正在卸载耗子Linux面板...'
echo -e $HR

Prepare_System
Remove_Swap
Remove_Panel

clear

echo -e "${LOGO}"
echo '耗子Linux面板卸载完成。'
echo '感谢您的使用，欢迎您再次使用耗子Linux面板。'
echo -e $HR
