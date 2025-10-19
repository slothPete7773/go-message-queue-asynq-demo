package main

import (
	"asynq-demo/task"
	"log"
	"math/rand"
	"time"

	"github.com/hibiken/asynq"
)

func main() {
	redisConn := asynq.RedisClientOpt{Addr: "localhost:6379"}
	client := asynq.NewClient(redisConn)
	defer client.Close()
	// for {

	userID := rand.Intn(1000) + 10
	delay := 2 * time.Minute
	task1 := task.NewWelcomeEmailTask(userID)
	task2 := task.NewReminderEmailTask(userID, time.Now().Add(delay))
	_, err := client.Enqueue(task1, asynq.Queue("critical"))
	if err != nil {
		log.Fatal(err)
	}
	_, err = client.Enqueue(task2, asynq.Queue("low"), asynq.ProcessIn(delay))
	if err != nil {
		log.Fatal(err)
	}
	// }
}
