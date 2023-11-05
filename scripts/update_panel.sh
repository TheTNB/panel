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
oldVersion=$(panel getSetting version)
oldVersion=${oldVersion#v}
panelPath="/www/panel"

# 大于
function version_gt() { test "$(echo -e "$1\n$2" | tr " " "\n" | sort -V | head -n 1)" != "$1"; }
# 小于
function version_lt() { test "$(echo -e "$1\n$2" | tr " " "\n" | sort -rV | head -n 1)" != "$1"; }
# 大于等于
function version_ge() { test "$(echo -e "$1\n$2" | tr " " "\n" | sort -rV | head -n 1)" == "$1"; }
# 小于等于
function version_le() { test "$(echo -e "$1\n$2" | tr " " "\n" | sort -V | head -n 1)" == "$1"; }

if [ -z "$oldVersion" ]; then
    echo "错误：无法获取面板版本"
    echo "Error: can't get panel version"
    exit 1
fi

echo $HR

if version_lt "$oldVersion" "2.1.8"; then
    echo "更新面板到 v2.1.8 ..."
    echo "Update panel to v2.1.8 ..."
    oldEntrance=$(panel getSetting entrance)
    echo "APP_ENTRANCE=$oldEntrance" >> $panelPath/panel.conf
    panel deleteSetting entrance
fi

echo $HR
echo "更新结束"
echo "Update finished"
