package asynqque

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

func HandleWelcomeEmailTask(c context.Context, t *asynq.Task) error {
	body := WelcomeEmail{}
	err := json.Unmarshal(t.Payload(), &body)
	_ = err

	if err != nil {
		return err
	}
	fmt.Printf("Send Welcome Email to User ID %d\n", body.UserId)
	return nil
}

func HandleReminderEmailTask(c context.Context, t *asynq.Task) error {
	body := ReminderEmail{}
	err := json.Unmarshal(t.Payload(), &body)
	_ = err

	fmt.Printf("Send Reminder Email to User ID %d\n", body.UserId)
	fmt.Printf("Reason: time is up (%v)\n", body.SentAt)
	return nil
}
