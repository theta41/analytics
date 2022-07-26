package tests

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"gitlab.com/g6834/team41/analytics/internal/app"
)

type statisticssTestSuite struct {
	suite.Suite

	app *app.App
}

func TestStatisticsSuite(t *testing.T) {
	suite.Run(t, &statisticssTestSuite{})
}

func (s *statisticssTestSuite) SetupSuite() {
	s.app = app.NewApp()
}

func (s *statisticssTestSuite) Test_StatisticsError() {
	var err error

	_, err = s.app.Statistics.GetCountAcceptedTask()
	s.Assert().NoError(err, "GetCountAcceptedTask fail")

	_, err = s.app.Statistics.GetCountDeclinedTask()
	s.Assert().NoError(err, "GetCountDeclinedTask fail")

	_, err = s.app.Statistics.GetSumReaction(objectId)
	s.Assert().NoError(err, "GetSumReaction fail")
}
