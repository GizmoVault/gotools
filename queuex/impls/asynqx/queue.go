package asynqx

import (
	"context"
	"errors"
	"time"

	"github.com/GizmoVault/gotools/queuex"
	"github.com/hibiken/asynq"
)

type RedisClientOpt struct {
	// Network type to use, either tcp or unix.
	// Default is tcp.
	Network string

	// Redis server address in "host:port" format.
	Addr string

	// Username to authenticate the current connection when Redis ACLs are used.
	// See: https://redis.io/commands/auth.
	Username string

	// Password to authenticate the current connection.
	// See: https://redis.io/commands/auth.
	Password string

	// Redis DB to select after connecting to a server.
	// See: https://redis.io/commands/select.
	DB int

	// Dial timeout for establishing new connections.
	// Default is 5 seconds.
	DialTimeout time.Duration

	// Timeout for socket reads.
	// If timeout is reached, read commands will fail with a timeout error
	// instead of blocking.
	//
	// Use value -1 for no timeout and 0 for default.
	// Default is 3 seconds.
	ReadTimeout time.Duration

	// Timeout for socket writes.
	// If timeout is reached, write commands will fail with a timeout error
	// instead of blocking.
	//
	// Use value -1 for no timeout and 0 for default.
	// Default is ReadTimout.
	WriteTimeout time.Duration

	// Maximum number of socket connections.
	// Default is 10 connections per every CPU as reported by runtime.NumCPU.
	PoolSize int
}

//nolint:gocritic // follow asynq
func (opt RedisClientOpt) ToAsyncQRedisClientOpt() asynq.RedisClientOpt {
	return asynq.RedisClientOpt{
		Network:      opt.Network,
		Addr:         opt.Addr,
		Username:     opt.Username,
		Password:     opt.Password,
		DB:           opt.DB,
		DialTimeout:  opt.DialTimeout,
		ReadTimeout:  opt.ReadTimeout,
		WriteTimeout: opt.WriteTimeout,
		PoolSize:     opt.PoolSize,
	}
}

//
//
//

//nolint:gocritic // follow asynq
func NewConsumerQueue(redisClientOpt RedisClientOpt) (queuex.ConsumerQueue, error) {
	return NewConsumerQueueWithServer(asynq.NewServer(
		redisClientOpt.ToAsyncQRedisClientOpt(),
		asynq.Config{
			Concurrency: 10,
		},
	))
}

func NewConsumerQueueWithServer(server *asynq.Server) (q queuex.ConsumerQueue, err error) {
	err = server.Ping()
	if err != nil {
		return
	}

	return &serverQueueImpl{
		server: server,
		mux:    asynq.NewServeMux(),
	}, nil
}

type serverQueueImpl struct {
	server *asynq.Server
	mux    *asynq.ServeMux
}

func (impl *serverQueueImpl) Run() error {
	return impl.server.Run(impl.mux)
}

func (impl *serverQueueImpl) HandleFunc(key string, h queuex.Handler) {
	impl.mux.HandleFunc(key, func(ctx context.Context, task *asynq.Task) error {
		taskID, _ := asynq.GetTaskID(ctx)

		err := h(ctx, taskID, &queuex.Task{
			Key:     task.Type(),
			Payload: task.Payload(),
		})

		if err != nil && errors.Is(err, queuex.ErrorSkipRetry) {
			err = asynq.SkipRetry
		}

		return err
	})
}

//
//
//

//nolint:gocritic // follow asynq
func NewProducerQueue(redisClientOpt RedisClientOpt) (queuex.ProducerQueue, error) {
	return NewProducerQueueWithClient(asynq.NewClient(redisClientOpt.ToAsyncQRedisClientOpt()))
}

func NewProducerQueueWithClient(client *asynq.Client) (q queuex.ProducerQueue, err error) {
	err = client.Ping()
	if err != nil {
		return
	}

	return &clientQueueImpl{
		client: client,
	}, nil
}

type clientQueueImpl struct {
	client *asynq.Client
}

func (impl *clientQueueImpl) Enqueue(task *queuex.Task, delay time.Duration) (id string, err error) {
	var options []asynq.Option

	options = append(options, asynq.Retention(time.Hour))

	if delay > 0 {
		options = append(options, asynq.ProcessIn(delay))
	}

	taskInfo, err := impl.client.Enqueue(asynq.NewTask(task.Key, task.Payload), options...)
	if err != nil {
		return
	}

	id = taskInfo.ID

	return
}
