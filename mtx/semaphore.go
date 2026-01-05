package mtx

import (
	"context"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

type TimeoutSemaphore struct {
	sem    *semaphore.Weighted
	mu     sync.Mutex
	timers map[uint64]*time.Timer
}

func NewTimeoutSemaphore(maxCount int64) *TimeoutSemaphore {
	return &TimeoutSemaphore{
		sem:    semaphore.NewWeighted(maxCount),
		timers: make(map[uint64]*time.Timer),
	}
}

func (ts *TimeoutSemaphore) AcquireWithAutoRelease(ctx context.Context, weight int64, timeout time.Duration, id uint64,
	fnTimeoutCallback func(id uint64)) error {
	if err := ts.sem.Acquire(ctx, weight); err != nil {
		return err
	}

	ts.mu.Lock()
	timer := time.AfterFunc(timeout, func() {
		ts.mu.Lock()
		_, ok := ts.timers[id]
		if ok {
			delete(ts.timers, id)
		}
		ts.mu.Unlock()

		if !ok {
			return
		}

		if fnTimeoutCallback != nil {
			fnTimeoutCallback(id)
		}

		ts.sem.Release(weight)
	})
	ts.timers[id] = timer
	ts.mu.Unlock()

	return nil
}

func (ts *TimeoutSemaphore) Release(id uint64, weight int64) {
	ts.mu.Lock()
	timer, ok := ts.timers[id]
	if ok {
		delete(ts.timers, id)
	}
	ts.mu.Unlock()

	if !ok {
		return
	}

	timer.Stop()
	ts.sem.Release(weight)
}
