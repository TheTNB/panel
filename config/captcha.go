package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("captcha", map[string]any{
		// 验证码图片高度
		"height": 60,

		// 验证码图片宽度
		"width": 120,

		// 验证码的长度
		"length": 6,

		// 数字的最大倾斜角度
		"maxskew": 0.6,

		// 图片背景里的混淆点数量
		"dotcount": 40,

		// 过期时间，单位是分钟
		"expire_time": 5,

		// debug 模式下的过期时间，方便本地开发调试
		"debug_expire_time": 3600,

		// 非 production 环境，使用此 key 可跳过验证，方便测试
		"testing_key": "captcha_skip_test",
	})
}
