package asynqque

import (
	"log"
	"sync"

	"github.com/hibiken/asynq"
)

var (
	once     = sync.Once{}
	mqClient = &asynq.Client{}
)

func SetupMessageQueue() {
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

	log.Println("listen to message queue")

	mux.HandleFunc(TypeWelcomeEmail, HandleWelcomeEmailTask)
	mux.HandleFunc(TypeReminderEmail, HandleReminderEmailTask)

	if err := worker.Run(mux); err != nil {
		log.Fatal(err)
	}

}

func NewMessageQueueClient() *asynq.Client {
	once.Do(func() {
		redisConn := asynq.RedisClientOpt{Addr: "localhost:6379"}
		mqClient = asynq.NewClient(redisConn)

	})

	return mqClient
}
