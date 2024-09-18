package api

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/TheTNB/panel/internal/app"
)

type APITestSuite struct {
	suite.Suite
	api *API
}

func TestAPITestSuite(t *testing.T) {
	app.Version = "2.3.0"
	suite.Run(t, &APITestSuite{
		api: NewAPI(),
	})
}

func (s *APITestSuite) TestGetLatestVersion() {
	_, err := s.api.GetLatestVersion()
	s.NoError(err)
}

func (s *APITestSuite) TestGetVersionsLog() {
	_, err := s.api.GetIntermediateVersions()
	s.NoError(err)
}
