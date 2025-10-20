# go-message-queue-asynq-demo

Run main app

```sh
go run main.go
```

Run message queue

```sh
go run cmd/message-queue/queue.go
```

Request that trigger message queue job

```sh
curl localhost:8088/users/1
```