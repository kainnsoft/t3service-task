package entity

type KafkaTypes int

// event types
const (
	Created KafkaTypes = iota
	Approved
	Rejected
	SendMail
	ToApprove
	ToReject
	AboutTaskEvent
	ToMailService
)

func (kt KafkaTypes) String() string {
	return [...]string{"created", "approved", "rejected", "send_mail", "to_approve", "to_reject", "about_task_event", "to_mail_service"}[kt]
}
