package asynqque

import (
	messagequeue "asynq-demo/message-queue"
	"fmt"
	"log"
	"sync"

	"github.com/hibiken/asynq"
)

type AsynqMq struct {
	asynqClient *asynq.Client
}

var (
	client = AsynqMq{}
	once   = sync.Once{}
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

func NewMessageQueueClient() messagequeue.Publisher {
	once.Do(func() {
		redisConn := asynq.RedisClientOpt{Addr: "localhost:6379"}
		client.asynqClient = asynq.NewClient(redisConn)

	})

	return &client
}

func (mq *AsynqMq) Close() {
	err := mq.asynqClient.Close()
	if err != nil {
		log.Fatal("failed to close asynq client: %w", err)
	}
}

func (mq *AsynqMq) Publish(payload messagequeue.Payload) error {
	task := asynq.NewTask(payload.TaskTypeLabel, payload.Body)

	taskInfo, err := mq.asynqClient.Enqueue(
		task,
		asynq.Queue(payload.QueueName),
	)
	if err != nil {
		return fmt.Errorf("error-enqueue-job: %w", err)
	}

	log.Println("task info: ", taskInfo)

	return nil
}
