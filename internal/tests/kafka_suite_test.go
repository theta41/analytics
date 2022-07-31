package tests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gitlab.com/g6834/team41/analytics/internal/app"
	"gitlab.com/g6834/team41/analytics/internal/env"
	mq "gitlab.com/g6834/team41/analytics/internal/kafka"
)

type kafkaTestSuite struct {
	suite.Suite

	app *app.App
}

func TestKafkaSuite(t *testing.T) {
	suite.Run(t, &kafkaTestSuite{})
}

func (s *kafkaTestSuite) SetupSuite() {
	s.app = app.NewApp()
}

func (s *kafkaTestSuite) TestConsumerConnection() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	mq.StartConsumer(
		ctx,
		env.E.C.Kafka.Brokers,
		env.E.C.Kafka.Topic,
		env.E.C.Kafka.GroupId,
		s.app.Events,
	)
}
