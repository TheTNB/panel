package captcha

import (
	"time"

	"github.com/goravel/framework/facades"
)

// CacheStore 实现 base64Captcha.Store interface
type CacheStore struct {
	KeyPrefix string
}

// Set 实现 base64Captcha.Store interface 的 Set 方法
func (s *CacheStore) Set(key string, value string) error {

	ExpireTime := time.Minute * time.Duration(facades.Config().GetInt("captcha.expire_time"))
	// 方便本地开发调试
	if facades.Config().GetBool("app.debug") {
		ExpireTime = time.Minute * time.Duration(facades.Config().GetInt("captcha.debug_expire_time"))
	}

	err := facades.Cache().Put(s.KeyPrefix+key, value, ExpireTime)
	if err != nil {
		return err
	}

	return nil
}

// Get 实现 base64Captcha.Store interface 的 Get 方法
func (s *CacheStore) Get(key string, clear bool) string {
	key = s.KeyPrefix + key
	val := facades.Cache().Get(key, "").(string)
	if clear {
		facades.Cache().Forget(key)
	}
	return val
}

// Verify 实现 base64Captcha.Store interface 的 Verify 方法
func (s *CacheStore) Verify(key, answer string, clear bool) bool {
	v := s.Get(key, clear)
	return v == answer
}
