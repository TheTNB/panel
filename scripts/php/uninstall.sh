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
setupPath="/www"
phpVersion="${1}"
phpPath="${setupPath}/server/php/${phpVersion}"

systemctl stop php-fpm-${phpVersion}
systemctl disable php-fpm-${phpVersion}
rm -rf /lib/systemd/system/php-fpm-${phpVersion}.service
systemctl daemon-reload

rm -rf ${phpPath}
rm -f /usr/bin/php-${phpVersion}

panel deletePlugin openresty

echo -e "${HR}\nPHP uninstall completed.\n${HR}"
