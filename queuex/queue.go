package queuex

import (
	"context"
	"errors"
	"time"

	"github.com/GizmoVault/gotools/storagex"
)

type Task struct {
	Key     string
	Payload []byte
}

func MarshalTask[S storagex.Serial](key string, s storagex.Serial, d any) *Task {
	task := &Task{
		Key:     key,
		Payload: nil,
	}

	if d != nil {
		data, _ := s.Marshal(d)
		task.Payload = data
	}

	return task
}

func UnMarshalTaskPayload[S storagex.Serial](s storagex.Serial, payload []byte, v any) error {
	return s.Unmarshal(payload, v)
}

type ProducerQueue interface {
	Enqueue(task *Task, delay time.Duration) (id string, err error)
}

var ErrorSkipRetry error = errors.New("skip retry")

type Handler func(ctx context.Context, id string, task *Task) error

type ConsumerQueue interface {
	HandleFunc(key string, h Handler)

	Run() error
}

type Queue interface {
	ProducerQueue
	ConsumerQueue

	Stop()
}
