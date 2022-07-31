package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gitlab.com/g6834/team41/analytics/internal/app"
)

type businessTestSuite struct {
	suite.Suite

	app *app.App
}

func TestBusinessSuite(t *testing.T) {
	suite.Run(t, &businessTestSuite{})
}

func (s *businessTestSuite) SetupSuite() {
	s.app = app.NewApp()
}

func (s *businessTestSuite) SetupTest() {
	_, err := s.app.Events.CreateTask(objectId)
	s.Require().NoError(err, "CreateTask fail")

	_, err = s.app.Events.CreateLetter(objectId, emailA)
	s.Assert().NoError(err, "fail CreateLetter emailA")

	_, err = s.app.Events.CreateLetter(objectId, emailB)
	s.Assert().NoError(err, "fail CreateLetter emailB")
}

func (s *businessTestSuite) TearDownTest() {
	err := s.app.Events.DeleteTask(objectId)
	s.Assert().NoError(err, "DeleteTask fail")
}

func (s *businessTestSuite) getStatistics() (completed int, uncompleted int, sumReaction int) {
	var err error

	completed, err = s.app.Statistics.GetCountAcceptedTask()
	s.Assert().NoError(err, "GetCountAcceptedTask fail")

	uncompleted, err = s.app.Statistics.GetCountDeclinedTask()
	s.Assert().NoError(err, "GetCountDeclinedTask fail")

	sumReaction, err = s.app.Statistics.GetSumReaction(objectId)
	s.Assert().NoError(err, "GetSumReaction fail")

	return
}

func (s *businessTestSuite) TestCounters() {

	initialCompleted, initialUncompleted, lastSumReaction := s.getStatistics()

	var err error

	time.Sleep(time.Millisecond)

	err = s.app.Events.AcceptedLetter(objectId, emailA)
	s.Assert().NoError(err, "fail AcceptedLetter emailA")

	s.Run("after emailA clicks accepted", func() {
		var completed, uncompleted, sumReaction int
		expectedCompleted := initialCompleted
		expectedUncompleted := initialUncompleted

		completed, uncompleted, sumReaction = s.getStatistics()

		s.Assert().Equal(expectedCompleted, completed, "wrong 'completed'")
		s.Assert().Equal(expectedUncompleted, uncompleted, "wrong 'uncompleted'")
		s.Assert().Less(lastSumReaction, sumReaction, "wrong 'sumReaction'")

		lastSumReaction = sumReaction
	})

	time.Sleep(time.Millisecond)

	err = s.app.Events.DeclinedLetter(objectId, emailB)
	s.Assert().NoError(err, "fail DeclinedLetter emailB")

	s.Run("after emailB clicks declined", func() {
		var completed, uncompleted, sumReaction int
		expectedCompleted := initialCompleted
		expectedUncompleted := initialUncompleted

		completed, uncompleted, sumReaction = s.getStatistics()

		s.Assert().Equal(expectedCompleted, completed, "wrong 'completed'")
		s.Assert().Equal(expectedUncompleted, uncompleted, "wrong 'uncompleted'")
		s.Assert().Equal(lastSumReaction, sumReaction, "wrong 'sumReaction'")

		lastSumReaction = sumReaction
	})

	time.Sleep(time.Millisecond)

	err = s.app.Events.FinishTask(objectId)
	s.Require().NoError(err, "FinishTask fail")

	s.Run("after task finished", func() {
		var completed, uncompleted, sumReaction int
		expectedCompleted := initialCompleted + 1
		expectedUncompleted := initialUncompleted - 1

		completed, uncompleted, sumReaction = s.getStatistics()

		s.Assert().Equal(expectedCompleted, completed, "wrong 'completed'")
		s.Assert().Equal(expectedUncompleted, uncompleted, "wrong 'uncompleted'")
		s.Assert().Equal(lastSumReaction, sumReaction, "wrong 'sumReaction'")
	})
}
