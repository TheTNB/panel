package captcha

import (
	"testing"
	"time"

	testifymock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/goravel/framework/testing/mock"
)

type CaptchaTestSuite struct {
	suite.Suite
	captcha *Captcha
}

func TestCaptchaTestSuite(t *testing.T) {
	mockConfig := mock.Config()
	mockConfig.On("GetString", "app.name").Return("HaoZiPanel").Once()
	mockConfig.On("GetInt", "captcha.height").Return(80).Once()
	mockConfig.On("GetInt", "captcha.width").Return(240).Once()
	mockConfig.On("GetInt", "captcha.length").Return(4).Once()
	mockConfig.On("Get", "captcha.maxskew").Return(0.7).Once()
	mockConfig.On("GetInt", "captcha.dotcount").Return(80).Once()
	mockConfig.On("GetInt", "captcha.expire_time").Return(5).Once()
	mockConfig.On("GetInt", "captcha.debug_expire_time").Return(10).Once()
	mockConfig.On("GetBool", "app.debug").Return(true).Once()
	mockCache, _, _ := mock.Cache()
	mockCache.On("Put", testifymock.Anything, testifymock.Anything, time.Minute*time.Duration(10)).Return(nil).Once()
	suite.Run(t, &CaptchaTestSuite{
		captcha: NewCaptcha(),
	})
	mockConfig.AssertExpectations(t)
}

func (s *CaptchaTestSuite) TestGenerateCaptcha() {
	id, base64, err := s.captcha.GenerateCaptcha()
	s.NotEmpty(id)
	s.NotEmpty(base64)
	s.Nil(err)
}
