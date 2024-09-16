package ntp

import (
	"testing"
	"time"

	"github.com/go-rat/utils/env"
	"github.com/stretchr/testify/suite"
)

type NTPTestSuite struct {
	suite.Suite
}

func TestNTPTestSuite(t *testing.T) {
	suite.Run(t, &NTPTestSuite{})
}

func (suite *NTPTestSuite) TestNowWithDefaultAddresses() {
	now, _ := Now()
	suite.WithinDuration(time.Now(), now, time.Minute)
}

func (suite *NTPTestSuite) TestNowWithCustomAddress() {
	now, err := Now("time.windows.com")
	suite.NoError(err)
	suite.WithinDuration(time.Now(), now, time.Minute)
}

func (suite *NTPTestSuite) TestNowWithInvalidAddress() {
	_, err := Now("invalid.address")
	suite.Error(err)
}

func (suite *NTPTestSuite) TestUpdateSystemTime() {
	if env.IsWindows() {
		suite.T().Skip("Skipping on Windows")
	}
	err := UpdateSystemTime(time.Now())
	suite.NoError(err)
}

func (suite *NTPTestSuite) TestUpdateSystemTimeZone() {
	if env.IsWindows() {
		suite.T().Skip("Skipping on Windows")
	}
	err := UpdateSystemTimeZone("UTC")
	suite.NoError(err)
}
