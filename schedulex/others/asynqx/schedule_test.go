package asynqx

import (
	"testing"
	"time"

	"github.com/GizmoVault/gotools/base/logx"
	"github.com/GizmoVault/gotools/configx/ut"
	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/require"
)

func TestSchedule(t *testing.T) {
	cfg := ut.SetupAndCheckFlags(t, ut.FlagRedis)

	s, err := NewScheduleTaskPool(asynq.RedisClientOpt{
		Addr:     cfg.RedisOpt.Addr,
		Password: cfg.RedisOpt.Password,
		DB:       1,
	}, "testQ", func(id, key string, payload []byte) string {
		t.Log("callback", "id", id, "key", key, "payload", string(payload))

		return "success"
	}, true, logx.NewConsoleLoggerWrapper())
	require.NoError(t, err)

	err = s.ClearQueue()
	require.NoError(t, err)

	id1, err := s.AddTask("task1", time.Now().Add(time.Second), []byte("task 1"))
	require.NoError(t, err)
	t.Log("id1", id1)

	id21, err := s.AddTask("task2", time.Now().Add(time.Second*3), []byte("task 2 -1"))
	require.NoError(t, err)
	t.Log("id21", id21)

	id22, err := s.AddTask("task2", time.Now().Add(time.Second*3), []byte("task 2 -2"))
	require.NoError(t, err)
	t.Log("id22", id22)

	id31, err := s.AddTask("task3", time.Now().Add(time.Second*3), []byte("task 3 -1"))
	require.NoError(t, err)
	t.Log("id31", id31)

	id32, err := s.AddTask("task3", time.Now().Add(time.Second*3), []byte("task 3 -2"))
	require.NoError(t, err)
	t.Log("id32", id32)

	id41, err := s.AddTask("task4", time.Now().Add(time.Second*10), []byte("task 4 -1"))
	require.NoError(t, err)
	t.Log("id41", id41)

	err = s.RemoveTask(id21)
	require.NoError(t, err)

	s.RemoveTaskByKey("task3")

	time.Sleep(time.Second * 20)
}
