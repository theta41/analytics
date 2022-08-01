package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"gitlab.com/g6834/team41/analytics/internal/ports"

	"github.com/sirupsen/logrus"

	"github.com/segmentio/kafka-go"
)

const (
	keyCreateTask    = "create-task"
	keyFinishTask    = "finish"
	keyDeleteTask    = "delete"
	keyCreateLetter  = "create-letter"
	keyAcceptLetter  = "accept"
	keyDeclineLetter = "decline"
)

type mqCreateTask struct {
	ObjectId uint32 `json:"task_id"`
}

type mqFinishTask mqCreateTask

type mqCreateLetter struct {
	ObjectId uint32 `json:"task_id"`
	Email    string `json:"email"`
}

type mqAccecptLetter mqCreateLetter
type mqDeclinedLetter mqCreateLetter

func StartConsumer(ctx context.Context, brokers []string, topic, groupId string, events ports.Events) {
	c, err := newClient(brokers, topic, groupId)
	if err != nil {
		logrus.Error("failed to connect to kafka: ", err)
		return
	}

	log := logrus.WithField("mq", "kafka-consumer")

	err = c.consume(ctx, func(binKey, binValue []byte) {
		key := string(binKey)
		log := log.WithField("key", key)

		switch key {

		default:
			log.Error("unknown message key")

		case keyCreateTask:
			var msg mqCreateTask
			err := json.Unmarshal(binValue, &msg)
			if err != nil {
				log.Error("unmarshal value error: ", err)
			} else {
				log.Info("got message: ", msg)
				events.CreateTask(msg.ObjectId)
			}

		case keyFinishTask:
			var msg mqFinishTask
			err := json.Unmarshal(binValue, &msg)
			if err != nil {
				log.Error("unmarshal value error: ", err)
			} else {
				log.Info("got message: ", msg)
				events.FinishTask(msg.ObjectId)
			}

		case keyDeleteTask:
			objectId, err := strconv.Atoi(string(binValue))
			if err != nil {
				log.Error("strconv.Atoi error: ", err)
			} else {
				log.Info("got objectId: ", objectId)
				events.DeleteTask(uint32(objectId))
			}

		case keyCreateLetter:
			var msg mqCreateLetter
			err := json.Unmarshal(binValue, &msg)
			if err != nil {
				log.Error("unmarshal value error: ", err)
			} else {
				log.Info("got message: ", msg)
				events.CreateLetter(msg.ObjectId, msg.Email)
			}

		case keyAcceptLetter:
			var msg mqAccecptLetter
			err := json.Unmarshal(binValue, &msg)
			if err != nil {
				log.Error("unmarshal value error: ", err)
			} else {
				log.Info("got message: ", msg)
				events.AcceptedLetter(msg.ObjectId, msg.Email)
			}

		case keyDeclineLetter:
			var msg mqDeclinedLetter
			err := json.Unmarshal(binValue, &msg)
			if err != nil {
				log.Error("unmarshal value error: ", err)
			} else {
				log.Info("got message: ", msg)
				events.DeclinedLetter(msg.ObjectId, msg.Email)
			}
		}
	})

	if err != nil {
		log.Error("consume error: ", err)
	}
}

type client struct {
	reader *kafka.Reader
}

func newClient(brokers []string, topic, groupId string) (*client, error) {
	if len(brokers) == 0 || brokers[0] == "" {
		return nil, errors.New("invalid kafka configuration")
	}

	c := client{}

	log := logrus.WithField("mq", "kafka-reader")

	c.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupId,
		Topic:    topic,
		MinBytes: 1,
		MaxBytes: 10e6,
		//Logger:      kafka.LoggerFunc(log.Infof),
		ErrorLogger: kafka.LoggerFunc(log.Errorf),
	})

	return &c, nil
}

func (c *client) consume(ctx context.Context, onMsg func(key, value []byte)) error {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			return err
		}

		onMsg(m.Key, m.Value)
	}
}
