package tests

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"gitlab.com/g6834/team41/analytics/internal/app"
)

type eventsTestSuite struct {
	suite.Suite

	app *app.App
}

func TestEventsSuite(t *testing.T) {
	suite.Run(t, &eventsTestSuite{})
}

func (s *eventsTestSuite) SetupSuite() {
	s.app = app.NewApp()
}

func (s *eventsTestSuite) Test_TaskErrors() {
	defer func() {
		err := s.app.Events.DeleteTask(objectId)
		s.Assert().NoError(err, "fail when defer DeleteTask")
	}()

	var err error

	err = s.app.Events.CreateTask(objectId)
	s.Require().NoError(err, "CreateTask fail")

	err = s.app.Events.FinishTask(objectId)
	s.Require().NoError(err, "FinishTask fail")
}

func (s *eventsTestSuite) Test_LetterErrors() {
	defer func() {
		err := s.app.Events.DeleteTask(objectId)
		s.Assert().NoError(err, "fail when defer DeleteTask")
	}()

	var err error

	err = s.app.Events.CreateLetter(objectId, emailA)
	s.Require().Error(err, "CreateLetter without task should fail")

	err = s.app.Events.CreateTask(objectId)
	s.Require().NoError(err, "CreateTask fail")

	err = s.app.Events.CreateLetter(objectId, emailA)
	s.Assert().NoError(err, "fail CreateLetter emailA")

	if err == nil {
		err = s.app.Events.AcceptedLetter(objectId, emailA)
		s.Assert().NoError(err, "fail AcceptedLetter emailA")
	}

	err = s.app.Events.CreateLetter(objectId, emailB)
	s.Assert().NoError(err, "fail CreateLetter emailB")

	if err == nil {
		err = s.app.Events.DeclinedLetter(objectId, emailB)
		s.Assert().NoError(err, "fail DeclinedLetter emailB")
	}
}
