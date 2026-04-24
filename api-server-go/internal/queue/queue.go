package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"mochat-api-server/internal/pkg/logger"
)

const (
	QueueDefault      = "default"
	QueueEmployee     = "employee"
	QueueContact      = "contact"
	QueueWelcome      = "welcome"
	QueueRoom         = "room"
	QueueChat         = "chat"
	QueueFile         = "file"
	QueueCallback     = "callback"
	QueueRemind       = "remind"
	QueueMessageMedia = "message_media"
)

var Client *asynq.Client

func InitClient(redisAddr string, redisPassword string, redisDB int) error {
	Client = asynq.NewClient(asynq.RedisClientOpt{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})
	return nil
}

func CloseClient() error {
	if Client != nil {
		return Client.Close()
	}
	return nil
}

func Enqueue(taskType string, payload interface{}, queue string, opts ...asynq.Option) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload failed: %w", err)
	}

	defaultOpts := []asynq.Option{
		asynq.Queue(queue),
		asynq.MaxRetry(3),
		asynq.Timeout(60 * time.Second),
		asynq.Retention(24 * time.Hour),
	}

	defaultOpts = append(defaultOpts, opts...)

	task := asynq.NewTask(taskType, data)
	info, err := Client.Enqueue(task, defaultOpts...)
	if err != nil {
		return fmt.Errorf("enqueue task failed: %w", err)
	}

	logger.Sugar.Debugf("enqueued task: id=%s type=%s queue=%s", info.ID, taskType, queue)
	return nil
}

type TaskHandler func(ctx context.Context, payload []byte) error

func NewServer(redisAddr string, redisPassword string, redisDB int) *asynq.Server {
	return asynq.NewServer(asynq.RedisClientOpt{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	}, asynq.Config{
		Concurrency: 20,
		Queues: map[string]int{
			QueueDefault:      6,
			QueueEmployee:     8,
			QueueContact:      4,
			QueueWelcome:      6,
			QueueRoom:         4,
			QueueChat:         6,
			QueueFile:         6,
			QueueCallback:     6,
			QueueRemind:       6,
			QueueMessageMedia: 10,
		},
		RetryDelayFunc: func(n int, err error, task *asynq.Task) time.Duration {
			return time.Duration(n) * 10 * time.Second
		},
	})
}
