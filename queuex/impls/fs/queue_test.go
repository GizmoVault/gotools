package fs

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/GizmoVault/gotools/base/logx"
	"github.com/GizmoVault/gotools/queuex"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestQueue(t *testing.T) {
	_ = os.RemoveAll("tmp")

	queue, err := NewFsQueue(t.Context(), "tmp/ut_queue.dat", logx.NewConsoleLoggerWrapper())
	assert.NoError(t, err)

	queue.HandleFunc("email:user", func(_ context.Context, task *queuex.Task) {
		t.Log(time.Now(), "email:user callback", "task.id", task.ID, "task.key", task.Key, "task.payload", task.Payload)
	})

	t.Log(time.Now(), "start enqueue")
	err = queue.Enqueue(&queuex.Task{
		ID:      uuid.NewString(),
		Key:     "email:user",
		Payload: []byte{1, 2, 3},
	}, time.Second)
	assert.NoError(t, err)

	err = queue.Enqueue(&queuex.Task{
		ID:      uuid.NewString(),
		Key:     "email:user",
		Payload: []byte{4, 5, 6},
	}, time.Second)
	assert.NoError(t, err)

	err = queue.Enqueue(&queuex.Task{
		ID:      uuid.NewString(),
		Key:     "email:vip",
		Payload: []byte{7, 8, 9},
	}, time.Second)
	assert.NoError(t, err)

	err = queue.Enqueue(&queuex.Task{
		ID:      uuid.NewString(),
		Key:     "email:vip",
		Payload: []byte{10, 11, 12},
	}, time.Second)
	assert.NoError(t, err)
	t.Log(time.Now(), "end enqueue")

	time.Sleep(time.Second * 3)
}

func TestQueue2(t *testing.T) {
	queue, err := NewFsQueue(t.Context(), "tmp/ut_queue.dat", logx.NewConsoleLoggerWrapper())
	assert.NoError(t, err)

	queue.HandleFunc("email:vip", func(_ context.Context, task *queuex.Task) {
		t.Log(time.Now(), "email:vip callback", "task.id", task.ID, "task.key", task.Key, "task.payload", task.Payload)
	})

	time.Sleep(time.Second * 10)
}
