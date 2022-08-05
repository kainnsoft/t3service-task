package entity

import (
	"fmt"
	"time"
)

type KafkaMsgAboutTaskEvent struct {
	TaskID          int32
	ApproversNumber int32
	Type            KafkaTypes
	User            string
	Time            time.Time
}

func NewKafkaMsgAboutTaskEvent(task *Task, taskType KafkaTypes, taskUserEmail string) *KafkaMsgAboutTaskEvent {
	kafkaMsg := KafkaMsgAboutTaskEvent{
		TaskID:          int32(task.ID),
		Time:            time.Now(),
		Type:            taskType,
		User:            taskUserEmail,
		ApproversNumber: int32(len(task.Approvers)),
	}

	return &kafkaMsg
}

type KafkaMsgToMailService struct {
	TaskID      int32      `json:"task_id"`
	Description string     `json:"description"`
	Body        string     `json:"body"`
	Addressee   string     `json:"addressee"`
	MailType    KafkaTypes `json:"mail_type"`
	ApproveLink string     `json:"approve_link"`
	RejectLink  string     `json:"reject_link"`
}

func NewKafkaMsgToMailService(task *Task, taskType KafkaTypes, userEmail string) *KafkaMsgToMailService {
	kafkaMsg := KafkaMsgToMailService{
		TaskID:      int32(task.ID),
		Description: task.Descr,
		Body:        task.Body,
		Addressee:   userEmail,
		MailType:    taskType,
		ApproveLink: fmt.Sprintf("/task/approve/%s", userEmail),
		RejectLink:  fmt.Sprintf("/task/reject/%s", userEmail),
	}

	return &kafkaMsg
}
