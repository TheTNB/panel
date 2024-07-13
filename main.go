/*
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
*/
package main

import (
	"github.com/goravel/framework/facades"

	"github.com/TheTNB/panel/v2/bootstrap"
)

//	@title			耗子面板 API
//	@version		2
//	@description	耗子面板的 API 信息

//	@contact.name	耗子科技
//	@contact.email	admin@haozi.net

//	@license.name	GNU Affero General Public License v3
//	@license		url   https://www.gnu.org/licenses/agpl-3.0.html

//	@securityDefinitions.apikey	BearerToken
//	@in							header
//	@name						Authorization

// @BasePath	/api
func main() {
	// 启动框架
	bootstrap.Boot()

	// 启动 HTTP 服务
	if facades.Config().GetBool("panel.ssl") {
		go func() {
			if err := facades.Route().RunTLS(); err != nil {
				facades.Log().Infof("Route run error: %v", err)
			}
		}()
	} else {
		go func() {
			if err := facades.Route().Run(); err != nil {
				facades.Log().Infof("Route run error: %v", err)
			}
		}()
	}

	// 启动计划任务
	go facades.Schedule().Run()

	// 启动队列
	go func() {
		if err := facades.Queue().Worker(nil).Run(); err != nil {
			facades.Log().Errorf("Queue run error: %v", err)
		}
	}()

	select {}
}
