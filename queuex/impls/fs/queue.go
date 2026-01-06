package fs

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/GizmoVault/gotools/base"
	"github.com/GizmoVault/gotools/base/errorx"
	"github.com/GizmoVault/gotools/base/logx"
	"github.com/GizmoVault/gotools/base/syncx"
	"github.com/GizmoVault/gotools/queuex"
	"github.com/GizmoVault/gotools/schedulex"
	"github.com/GizmoVault/gotools/storagex"
	"github.com/google/uuid"
)

func NewFsQueue(ctx context.Context, fileName string, logger logx.Wrapper) (queuex.Queue, error) {
	return NewFsQueueWithFNNow(ctx, fileName, nil, logger)
}

func NewFsQueueWithFNNow(ctx context.Context, fileName string, now base.FNNow, logger logx.Wrapper) (queuex.Queue, error) {
	if logger == nil {
		logger = logx.NewConsoleLoggerWrapper()
	}

	logger = logger.WithFields(logx.StringField(logx.ClsKey, "queueImpl"))

	expiredStg, err := storagex.NewMemWithFile[map[string]*innerTask, storagex.Serial, syncx.RWLocker](
		make(map[string]*innerTask), &storagex.JSONSerial{}, &sync.RWMutex{}, fileName+".expired", nil)
	if err != nil {
		return nil, err
	}

	impl := &queueImpl{
		logger:     logger,
		ctx:        ctx,
		expiredStg: expiredStg,
		taskPool:   schedulex.NewHeapTaskPool(now),
		m:          make(map[string]queuex.Handler),
	}

	impl.stg, err = storagex.NewMemWithFileEx[map[string]*innerTask, storagex.Serial, syncx.RWLocker](
		make(map[string]*innerTask), &storagex.JSONSerial{}, &sync.RWMutex{}, fileName, nil, impl)
	if err != nil {
		return nil, err
	}

	return impl, nil
}

type queueImpl struct {
	logger     logx.Wrapper
	ctx        context.Context
	stg        *storagex.MemWithFile[map[string]*innerTask, storagex.Serial, syncx.RWLocker]
	expiredStg *storagex.MemWithFile[map[string]*innerTask, storagex.Serial, syncx.RWLocker]
	taskPool   schedulex.ScheduleTaskPool

	mLock sync.Mutex
	m     map[string]queuex.Handler
}

func (*queueImpl) Run() error {
	<-make(chan any)

	return errorx.ErrLogic
}

func (impl *queueImpl) Stop() {
	impl.taskPool.Stop()
}

func (*queueImpl) BeforeLoad() {

}

func (impl *queueImpl) AfterLoad(m map[string]*innerTask, err error) {
	if err != nil {
		impl.logger.WithFields(logx.ErrorField(err)).Error("AfterLoad failed")

		return
	}

	for _, task := range m {
		if e := impl.taskPool.AddTask(task.Key, time.Unix(task.At, 0), impl.taskCallback); e != nil {
			impl.logger.WithFields(logx.ErrorField(err)).Errorf("taskPool AddTask failed")
		}
	}
}

func (*queueImpl) BeforeSave() {

}

func (*queueImpl) AfterSave(_ map[string]*innerTask, _ error) {

}

//
//
//

func (impl *queueImpl) Enqueue(task *queuex.Task, delay time.Duration) (id string, err error) {
	if task == nil || task.Key == "" {
		err = errorx.ErrInvalidArgs

		return
	}

	id = uuid.NewString()

	var at int64

	err = impl.stg.Change(func(oldM map[string]*innerTask) (newM map[string]*innerTask, err error) {
		newM = oldM

		if len(newM) == 0 {
			newM = make(map[string]*innerTask)
		}

		if _, ok := newM[id]; ok {
			err = errorx.ErrExists

			return
		}

		newM[id] = fromTask(id, task, delay)

		at = newM[id].At

		return
	})

	if err != nil {
		return
	}

	err = impl.taskPool.AddTask(id, time.Unix(at, 0), impl.taskCallback)

	return
}

func (impl *queueImpl) HandleFunc(key string, h queuex.Handler) {
	impl.mLock.Lock()
	defer impl.mLock.Unlock()

	if h == nil {
		delete(impl.m, key)

		return
	}

	impl.m[key] = h

	go func() {
		time.Sleep(time.Second * 5)

		for {
			var task *innerTask
			var handler queuex.Handler

			impl.expiredStg.Read(func(m map[string]*innerTask) {
				for _, task = range m {
					handler = impl.getHandler(task.Key)
					if handler == nil {
						continue
					}

					break
				}
			})

			if handler == nil {
				break
			}

			_ = h(impl.ctx, task.ID, task.GetTask())

			_ = impl.expiredStg.Change(func(oldM map[string]*innerTask) (newM map[string]*innerTask, err error) {
				newM = oldM
				if len(newM) == 0 {
					err = errorx.ErrSkip

					return
				}

				delete(newM, task.ID)

				return
			})
		}
	}()
}

func (impl *queueImpl) getHandler(key string) queuex.Handler {
	impl.mLock.Lock()
	defer impl.mLock.Unlock()

	h, ok := impl.m[key]
	if ok {
		return h
	}

	for s, handler := range impl.m {
		if strings.HasPrefix(key, s) {
			return handler
		}
	}

	return nil
}

//
//
//

func (impl *queueImpl) taskCallback(key string, _ ...any) {
	var task *innerTask
	var ok bool

	impl.stg.Read(func(m map[string]*innerTask) {
		task, ok = m[key]
	})

	if !ok {
		return
	}

	h := impl.getHandler(task.Key)
	if h == nil {
		_ = impl.expiredStg.Change(func(oldM map[string]*innerTask) (newM map[string]*innerTask, err error) {
			newM = oldM
			if len(newM) == 0 {
				newM = make(map[string]*innerTask)
			}

			newM[task.ID] = task

			return
		})
	} else {
		_ = h(impl.ctx, task.ID, task.GetTask())
	}

	_ = impl.stg.Change(func(oldM map[string]*innerTask) (newM map[string]*innerTask, err error) {
		newM = oldM
		if len(newM) == 0 {
			newM = make(map[string]*innerTask)
		}

		delete(newM, task.ID)

		return
	})
}
