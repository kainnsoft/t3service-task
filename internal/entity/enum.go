package entity

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
