package asynqque

import (
	messagequeue "asynq-demo/message-queue"
	"encoding/json"
	"time"
)

const (
	TypeWelcomeEmail  = "email:welcome"
	TypeReminderEmail = "email:reminder"
)

type WelcomeEmail struct {
	UserId int `json:"userId"`
}

type ReminderEmail struct {
	UserId int       `json:"userId"`
	SentAt time.Time `json:"sentAt"`
}

func NewWelcomeEmailTask(id int) *messagequeue.Payload {
	payload := WelcomeEmail{UserId: id}
	bPayload, err := json.Marshal(payload)
	_ = err

	return &messagequeue.Payload{
		Body:          bPayload,
		TaskTypeLabel: TypeWelcomeEmail,
		QueueName:     "low",
	}
}

func NewReminderEmailTask(id int, ts time.Time) *messagequeue.Payload {
	payload := ReminderEmail{
		UserId: id,
		SentAt: ts,
	}

	bPayload, err := json.Marshal(payload)
	_ = err

	return &messagequeue.Payload{
		Body:          bPayload,
		TaskTypeLabel: TypeReminderEmail,
		QueueName:     "critical",
	}
}
