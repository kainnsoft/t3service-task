package repository

import (
	"encoding/json"
	"fmt"
	"strconv"
	"team3-task/internal/entity"
	kafkaPkg "team3-task/pkg/kafka"

	"github.com/segmentio/kafka-go"
)

type KafkaProducers struct {
	KafProducerAboutTaskEvent *kafkaPkg.Client // for analityc topic
	KafProducerToMailService  *kafkaPkg.Client // for mail topic
}

func SendMessagesToKafka(c *kafkaPkg.Client,
	task *entity.Task,
	taskType entity.KafkaTypes,
	taskUser string,
	msgType entity.KafkaTypes) error {

	var (
		msgValue []byte
		err      error
	)

	switch msgType {
	case entity.AboutTaskEvent:
		msgNew := entity.NewKafkaMsgAboutTaskEvent(task, taskType, taskUser)
		msgValue, err = json.MarshalIndent(msgNew, "", " ")

		if err != nil {
			return fmt.Errorf("repository.SendMessagesToKafka json.Marshal error: %v", err)
		}
	case entity.ToMailService:
		msgNew := entity.NewKafkaMsgToMailService(task, taskType, taskUser)
		msgValue, err = json.MarshalIndent(msgNew, "", " ")

		if err != nil {
			return fmt.Errorf("repository.SendMessagesToKafka json.Marshal error: %v", err)
		}
	}

	keyTaskID := strconv.Itoa(task.ID)
	msg := []kafka.Message{
		{
			Key:   []byte(keyTaskID),
			Value: msgValue,
		},
	}
	err = c.SendMessages(msg)

	if err != nil {
		return fmt.Errorf("repository.SendMessagesToKafka c.SendMessages(msg) error: %v", err)
	}

	return nil
}
