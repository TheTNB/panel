// Package captcha 处理图片验证码逻辑
package captcha

import (
	"sync"

	"github.com/goravel/framework/facades"
	"github.com/mojocn/base64Captcha"
)

type Captcha struct {
	Base64Captcha *base64Captcha.Captcha
}

// once 确保 internalCaptcha 对象只初始化一次
var once sync.Once

// internalCaptcha 内部使用的 Captcha 对象
var internalCaptcha *Captcha

// NewCaptcha 单例模式获取
func NewCaptcha() *Captcha {
	once.Do(func() {
		// 初始化 Captcha 对象
		internalCaptcha = &Captcha{}

		// 使用 Cache facade 进行存储，并配置存储 Key 的前缀
		store := CacheStore{
			KeyPrefix: facades.Config().GetString("app.name") + ":captcha:",
		}

		// 配置 base64Captcha 驱动信息
		driver := base64Captcha.NewDriverDigit(
			facades.Config().GetInt("captcha.height"),         // 宽
			facades.Config().GetInt("captcha.width"),          // 高
			facades.Config().GetInt("captcha.length"),         // 长度
			facades.Config().Get("captcha.maxskew").(float64), // 数字的最大倾斜角度
			facades.Config().GetInt("captcha.dotcount"),       // 图片背景里的混淆点数量
		)

		// 实例化 base64Captcha 并赋值给内部使用的 internalCaptcha 对象
		internalCaptcha.Base64Captcha = base64Captcha.NewCaptcha(driver, &store)
	})

	return internalCaptcha
}

// GenerateCaptcha 生成图片验证码
func (c *Captcha) GenerateCaptcha() (id string, b64s string, err error) {
	return c.Base64Captcha.Generate()
}

// VerifyCaptcha 验证验证码是否正确
func (c *Captcha) VerifyCaptcha(id string, answer string, clear bool) (match bool) {

	// 方便本地和 API 自动测试
	if facades.Config().GetBool("app.debug") && id == facades.Config().GetString("captcha.testing_key") {
		return true
	}
	// 这样方便用户多次提交，防止表单提交错误需要多次输入图片验证码
	return c.Base64Captcha.Verify(id, answer, clear)
}
