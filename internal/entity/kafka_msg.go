package entity

import (
	"fmt"
	"strconv"
	"time"
)

type KafkaMsgAboutTaskEvent struct {
	TaskID          int32
	ApproversNumber int32
	User            string     // требования бизнес-логики сервиса аналитики - пользователь (если создание, то email автора, если согласование, то email согласующего)
	Type            KafkaTypes // требования бизнес-логики сервиса аналитики - тип события (создание, согласование и т.д.)
	Time            time.Time  // требования бизнес-логики сервиса аналитики - т.к. время события (транзакции) и стандартное время msg в кафке может различаться
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
	ApproveLink string     `json:"approve_link"`
	RejectLink  string     `json:"reject_link"`
	MailType    KafkaTypes `json:"mail_type"`
}

func NewKafkaMsgToMailService(task *Task, taskType KafkaTypes, userEmail, link string) *KafkaMsgToMailService {
	kafkaMsg := KafkaMsgToMailService{
		TaskID:      int32(task.ID),
		Description: task.Descr,
		Body:        task.Body,
		Addressee:   userEmail,
		MailType:    taskType, // либо Approved отсылки очередному или в финале, либо Rejected для оповещения всех
		ApproveLink: fmt.Sprintf("%s/task/approve/%s/%s", link, userEmail, strconv.Itoa(task.ID)),
		RejectLink:  fmt.Sprintf("%s/task/reject/%s/%s", link, userEmail, strconv.Itoa(task.ID)),
	}

	return &kafkaMsg
}
