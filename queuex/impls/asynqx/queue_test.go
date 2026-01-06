package asynqx

import (
	"context"
	"testing"
	"time"

	"github.com/GizmoVault/gotools/queuex"
	"github.com/stretchr/testify/assert"
)

func TestQueue(t *testing.T) {
	t.SkipNow()

	queue, err := NewProducerQueue(RedisClientOpt{
		Addr:     "192.168.31.11:6379",
		Password: "repass",
		DB:       1,
	})
	assert.NoError(t, err)

	time.Sleep(1 * time.Second)

	go func() {
		consumeQueue, e := NewConsumerQueue(RedisClientOpt{
			Addr:     "192.168.31.11:6379",
			Password: "repass",
			DB:       1,
		})
		assert.NoError(t, e)

		consumeQueue.HandleFunc("email:user", func(_ context.Context, id string, task *queuex.Task) error {
			t.Log(time.Now(), "email:user callback", "task.id", id, "task.key", task.Key, "task.payload", task.Payload)

			return nil
		})

		_ = consumeQueue.Run()
	}()

	t.Log(time.Now(), "start enqueue")
	id1, err := queue.Enqueue(&queuex.Task{
		Key:     "email:user",
		Payload: []byte{1, 2, 3},
	}, time.Second)
	assert.NoError(t, err)
	t.Log("start enqueue", id1)

	id2, err := queue.Enqueue(&queuex.Task{
		Key:     "email:user",
		Payload: []byte{4, 5, 6},
	}, time.Second)
	assert.NoError(t, err)
	t.Log("start enqueue", id2)

	id3, err := queue.Enqueue(&queuex.Task{
		Key:     "email:vip",
		Payload: []byte{7, 8, 9},
	}, time.Second)
	assert.NoError(t, err)
	t.Log("start enqueue", id3)

	id4, err := queue.Enqueue(&queuex.Task{
		Key:     "email:vip",
		Payload: []byte{10, 11, 12},
	}, time.Second)
	assert.NoError(t, err)
	t.Log(time.Now(), "end enqueue")
	t.Log("start enqueue", id4)

	time.Sleep(time.Second * 1000)
}

func TestQueue2(t *testing.T) {
	t.SkipNow()

	consumeQueue, err := NewConsumerQueue(RedisClientOpt{
		Addr:     "192.168.31.11:6379",
		Password: "repass",
		DB:       1,
	})
	assert.NoError(t, err)

	consumeQueue.HandleFunc("email:vip", func(_ context.Context, id string, task *queuex.Task) error {
		t.Log(time.Now(), "email:vip callback", "task.id", id, "task.key", task.Key, "task.payload", task.Payload)

		return nil
	})

	_ = consumeQueue.Run()

	time.Sleep(time.Second * 10)
}
