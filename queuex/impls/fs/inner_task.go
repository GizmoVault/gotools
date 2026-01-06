package fs

import (
	"time"

	"github.com/GizmoVault/gotools/queuex"
)

type innerTask struct {
	ID      string
	Key     string
	Payload []byte
	At      int64
}

func (it *innerTask) GetTask() *queuex.Task {
	if it == nil {
		return nil
	}

	task := &queuex.Task{
		Key: it.Key,
	}

	if it.Payload != nil {
		task.Payload = make([]byte, len(it.Payload))
		copy(task.Payload, it.Payload)
	}

	return task
}

func fromTask(id string, task *queuex.Task, delay time.Duration) *innerTask {
	at := time.Now()
	if delay > 0 {
		at = at.Add(delay)
	}

	newTask := &innerTask{
		ID:  id,
		Key: task.Key,
		At:  at.Unix(),
	}

	if task.Payload != nil {
		newTask.Payload = make([]byte, len(task.Payload))
		copy(newTask.Payload, task.Payload)
	}

	return newTask
}
