package asynqx

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/GizmoVault/gotools/base/errorx"
	"github.com/GizmoVault/gotools/base/logx"
	"github.com/hibiken/asynq"
)

type TaskFunc func(id, key string, payload []byte) string

//nolint:gocritic // follow asynq
func NewScheduleTaskPool(opt asynq.RedisClientOpt, queueName string, callback TaskFunc, clearAllTasks bool,
	logger logx.Wrapper) (st *ScheduleTask, err error) {
	if logger == nil {
		logger = logx.NewNopLoggerWrapper()
	}

	logger = logger.WithFields(logx.StringField(logx.ClsKey, "ScheduleTask"))

	if queueName == "" {
		err = errorx.ErrInvalidArgs

		return
	}

	st = &ScheduleTask{
		logger:    logger,
		queueName: queueName,
		server: asynq.NewServer(opt, asynq.Config{
			Queues: map[string]int{
				queueName: 1,
			},
		}),
		client:    asynq.NewClient(opt),
		inspector: asynq.NewInspector(opt),
		callback:  callback,
	}

	err = st.init(clearAllTasks)
	if err != nil {
		return
	}

	return
}

type ScheduleTask struct {
	wg        sync.WaitGroup
	logger    logx.Wrapper
	queueName string
	server    *asynq.Server
	client    *asynq.Client
	inspector *asynq.Inspector

	callback TaskFunc
}

func (impl *ScheduleTask) ClearQueue() (err error) {
	_, err = impl.inspector.GetQueueInfo(impl.queueName)
	if err != nil {
		return
	}

	methods := []func(string) (int, error){
		impl.inspector.DeleteAllPendingTasks,
		impl.inspector.DeleteAllScheduledTasks,
		impl.inspector.DeleteAllRetryTasks,
		impl.inspector.DeleteAllCompletedTasks,
		impl.inspector.DeleteAllArchivedTasks,
	}

	for _, m := range methods {
		_, err = m(impl.queueName)
		if err != nil {
			if !errors.Is(err, asynq.ErrQueueNotFound) {
				return
			}

			err = nil
		}
	}

	return
}

func (impl *ScheduleTask) AddTask(key string, t time.Time, payload []byte) (id string, err error) {
	taskInfo, err := impl.client.Enqueue(asynq.NewTask(key, payload), asynq.Queue(impl.queueName), asynq.ProcessAt(t),
		asynq.Retention(time.Hour*24))
	if err != nil {
		return
	}

	id = taskInfo.ID

	return
}

func (impl *ScheduleTask) RemoveTaskByKey(key string) {
	methods := []func(queue string, opts ...asynq.ListOption) ([]*asynq.TaskInfo, error){
		impl.inspector.ListPendingTasks,
		impl.inspector.ListScheduledTasks,
		impl.inspector.ListRetryTasks,
	}

	for _, method := range methods {
		page := 1
		for {
			tasks, err := method(impl.queueName, asynq.Page(page), asynq.PageSize(100))
			if err != nil {
				impl.logger.WithFields(logx.ErrorField(err)).Error("fetch task list failed")

				continue
			}

			if len(tasks) == 0 {
				break
			}

			for _, task := range tasks {
				if task.Type == key {
					err = impl.inspector.DeleteTask(impl.queueName, task.ID)
					if err != nil {
						impl.logger.WithFields(logx.ErrorField(err)).Error("delete task failed")
					}
				}
			}

			page++
		}
	}
}

func (impl *ScheduleTask) RemoveTask(id string) error {
	return impl.inspector.DeleteTask(impl.queueName, id)
}

func (impl *ScheduleTask) Stop() {
	impl.server.Stop()
	impl.server.Shutdown()
}

func (impl *ScheduleTask) Start(callback TaskFunc) error {
	if impl.callback != nil {
		return errorx.ErrLogic
	}

	impl.callback = callback
	impl.start()

	return nil
}

func (impl *ScheduleTask) init(clearAllTasks bool) (err error) {
	err = impl.server.Ping()
	if err != nil {
		return
	}

	if clearAllTasks {
		_ = impl.ClearQueue()
	}

	if impl.callback != nil {
		impl.start()
	}

	return
}

func (impl *ScheduleTask) start() {
	impl.wg.Add(1)

	go func() {
		defer impl.wg.Done()

		_ = impl.server.Run(asynq.HandlerFunc(func(ctx context.Context, task *asynq.Task) error {
			id, _ := asynq.GetTaskID(ctx)

			result := impl.callback(id, task.Type(), task.Payload())

			_, _ = task.ResultWriter().Write([]byte(result))

			return nil
		}))
	}()
}
