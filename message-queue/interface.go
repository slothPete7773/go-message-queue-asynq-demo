package messagequeue

type MessageQueue interface {
	Close()
}

type Publisher interface {
	Publish(Payload) error
	MessageQueue
}

type Consumer interface {
	Run() error
	MessageQueue
}

type Payload struct {
	Body          []byte
	QueueName     string
	TaskTypeLabel string
}
