package entity

import (
	"fmt"
	"time"
)

type KafkaTypes int

// event types
const (
	Created KafkaTypes = iota
	Approved
	Rejected
	SendMail
	AboutTaskEvent
	ToMailService
)

func (kt KafkaTypes) String() string {
	return [...]string{"created", "approved", "rejected", "send_mail", "about_task_event", "to_mail_service"}[kt]
}

type KafkaMsgAboutTaskEvent struct {
	TaskId          int32
	Time            time.Time
	Type            string
	User            string
	ApproversNumber int32
}

func NewKafkaMsgAboutTaskEvent(task *Task, taskType, taskUserEmail string) *KafkaMsgAboutTaskEvent {
	kafkaMsg := KafkaMsgAboutTaskEvent{
		TaskId:          int32(task.Id),
		Time:            time.Now(),
		Type:            taskType,
		User:            taskUserEmail,
		ApproversNumber: int32(len(task.Approvers)),
	}
	return &kafkaMsg
}

type KafkaMsgToMailService struct {
	TaskId      int32  `json:"task_id"`
	Description string `json:"description"`
	Body        string `json:"body"`
	Addressee   string `json:"addressee"`
	MailType    string `json:"mail_type"`
	ApproveLink string `json:"approve_link"`
	RejectLink  string `json:"reject_link"`
}

func NewKafkaMsgToMailService(task *Task, taskType, userEmail string) *KafkaMsgToMailService {
	kafkaMsg := KafkaMsgToMailService{
		TaskId:      int32(task.Id),
		Description: task.Descr,
		Body:        task.Body,
		Addressee:   userEmail,
		MailType:    taskType,
		ApproveLink: fmt.Sprintf("/task/approve/%s", userEmail),
		RejectLink:  fmt.Sprintf("/task/reject/%s", userEmail),
	}

	return &kafkaMsg
}
