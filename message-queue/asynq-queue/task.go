package asynqque

import (
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
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

func NewWelcomeEmailTask(id int) *asynq.Task {
	payload := WelcomeEmail{UserId: id}
	bPayload, err := json.Marshal(payload)
	_ = err

	return asynq.NewTask(TypeWelcomeEmail, bPayload)
}

func NewReminderEmailTask(id int, ts time.Time) *asynq.Task {
	payload := ReminderEmail{
		UserId: id,
		SentAt: ts,
	}

	bPayload, err := json.Marshal(payload)
	_ = err

	return asynq.NewTask(TypeReminderEmail, bPayload)
}
