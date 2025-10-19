package main

import (
	"asynq-demo/task"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
)

func HandleWelcomeEmailTask(c context.Context, t *asynq.Task) error {
	body := task.WelcomeEmail{}
	err := json.Unmarshal(t.Payload(), &body)
	_ = err

	if err != nil {
		return err
	}
	fmt.Printf("Send Welcome Email to User ID %d\n", body.UserId)
	return nil
}

func HandleReminderEmailTask(c context.Context, t *asynq.Task) error {
	body := task.ReminderEmail{}
	err := json.Unmarshal(t.Payload(), &body)
	_ = err

	fmt.Printf("Send Reminder Email to User ID %d\n", body.UserId)
	fmt.Printf("Reason: time is up (%v)\n", body.SentAt)
	return nil
}

func main() {
	redisConn := asynq.RedisClientOpt{Addr: "localhost:6379"}
	worker := asynq.NewServer(redisConn, asynq.Config{
		Concurrency: 10,
		Queues: map[string]int{
			"critical": 6,
			"default":  3,
			"low":      1,
		},
	})
	mux := asynq.NewServeMux()
	mux.HandleFunc(task.TypeWelcomeEmail, HandleWelcomeEmailTask)
	mux.HandleFunc(task.TypeReminderEmail, HandleReminderEmailTask)
	if err := worker.Run(mux); err != nil {
		log.Fatal(err)
	}

}
