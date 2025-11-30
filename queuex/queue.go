package queuex

import (
	"context"
	"time"

	"github.com/GizmoVault/gotools/storagex"
	"github.com/google/uuid"
)

type Task struct {
	ID      string
	Key     string
	Payload []byte
}

func MarshalTask[S storagex.Serial](key string, s storagex.Serial, d any) *Task {
	task := &Task{
		ID:      uuid.NewString(),
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
	Enqueue(task *Task, delay time.Duration) error
}

type Handler func(ctx context.Context, task *Task)

type ConsumerQueue interface {
	HandleFunc(key string, h Handler)
}

type Queue interface {
	ProducerQueue
	ConsumerQueue

	Stop()
}
