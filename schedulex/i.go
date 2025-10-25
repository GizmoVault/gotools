package schedulex

import (
	"time"

	"github.com/GizmoVault/gotools/base"
)

type TaskFunc func(key string, args ...any)

type ScheduleTaskPool interface {
	AddTask(key string, t time.Time, exec TaskFunc, params ...any) error
	RemoveTask(key string) error

	Stop()
}

func CreateHeapTaskPool() ScheduleTaskPool {
	return NewHeapTaskPool(nil)
}

func CreateHeapTaskPoolWithFNNow(now base.FNNow) ScheduleTaskPool {
	return NewHeapTaskPool(now)
}
