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

# 计算 j 值（通用）
calculate_j() {
    export LC_ALL=C
    total_mem=$(free -m | awk '/^Mem:/{print $2}')
    total_swap=$(free -m | awk '/^Swap:/{print $2}')
    total=$((total_mem + total_swap))
    j_value=$((total / 1024))
    cpu_cores=$(nproc)

    if [ $j_value -eq 0 ]; then
        j_value=1
    fi

    if [ $j_value -gt "$cpu_cores" ]; then
        j_value=$cpu_cores
    fi

    echo "$j_value"
}

# 计算 j 值（2倍内存）
calculate_j2() {
    export LC_ALL=C
    total_mem=$(free -m | awk '/^Mem:/{print $2}')
    total_swap=$(free -m | awk '/^Swap:/{print $2}')
    total=$((total_mem + total_swap))
    j_value=$((total / 2024))
    cpu_cores=$(nproc)

    if [ $j_value -eq 0 ]; then
        j_value=1
    fi

    if [ $j_value -gt "$cpu_cores" ]; then
        j_value=$cpu_cores
    fi

    echo "$j_value"
}
