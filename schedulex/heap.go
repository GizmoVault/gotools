package schedulex

import (
	"container/heap"
	"sync"
	"time"

	"github.com/GizmoVault/gotools/base"
	"github.com/GizmoVault/gotools/base/commerrx"
	"github.com/google/uuid"
)

const (
	taskOpAdd = iota
	taskOpDel
	taskOpUpdate
)

type taskOpInfo struct {
	opType int
	key    string
	at     time.Time
	exec   TaskFunc
	params []interface{}
}

type taskItem struct {
	at       time.Time
	exec     TaskFunc
	params   []any
	canceled bool
	key      string
}

type taskHeap []*taskItem

func (th *taskHeap) Len() int {
	return len(*th)
}

func (th *taskHeap) Less(i, j int) bool {
	return (*th)[i].at.Unix() < (*th)[j].at.Unix()
}

func (th *taskHeap) Swap(i, j int) {
	(*th)[i], (*th)[j] = (*th)[j], (*th)[i]
}

func (th *taskHeap) Push(x any) {
	*th = append(*th, x.(*taskItem))
}

func (th *taskHeap) Pop() any {
	old := *th
	n := len(old)
	x := old[n-1]
	*th = old[0 : n-1]
	return x
}

type HeapTaskPool struct {
	sync.WaitGroup
	fnNow  base.FNNow
	closed chan bool
	taskOp chan *taskOpInfo
	keys   map[string]*taskItem

	tasks *taskHeap
}

func NewHeapTaskPool(now base.FNNow) *HeapTaskPool {
	pool := &HeapTaskPool{
		fnNow: now,
	}
	pool.Start()
	return pool
}

func (tp *HeapTaskPool) Start() {
	tp.Wait()

	tp.closed = make(chan bool)
	tp.taskOp = make(chan *taskOpInfo, 100)
	tp.keys = make(map[string]*taskItem)

	tp.tasks = &taskHeap{}
	heap.Init(tp.tasks)

	tp.Add(1)
	go tp.loop()
}

func (tp *HeapTaskPool) Stop() {
	close(tp.closed)

	tp.Wait()
}

func execTask(key string, exec TaskFunc, params []any) {
	go exec(key, params...)
}

func (tp *HeapTaskPool) now() time.Time {
	if tp.fnNow != nil {
		return tp.fnNow()
	}

	return time.Now()
}

func (tp *HeapTaskPool) process() time.Duration {
	for {
		if tp.tasks.Len() <= 0 {
			return 24 * time.Hour
		}

		timeNow := tp.now()

		t := (*tp.tasks)[0]
		if t.canceled {
			heap.Pop(tp.tasks)

			continue
		}

		if t.at.After(timeNow) {
			return t.at.Sub(timeNow)
		}

		heap.Pop(tp.tasks)

		delete(tp.keys, t.key)

		execTask(t.key, t.exec, t.params)
	}
}

func (tp *HeapTaskPool) addOrUpdateTask(t time.Time, key string, exec TaskFunc, params []interface{}) {
	if key == "" {
		key = uuid.NewString()
	} else {
		tp.cancelTask(key)
	}

	ti := &taskItem{
		at:     t,
		exec:   exec,
		params: params,
		key:    key,
	}

	tp.keys[key] = ti

	heap.Push(tp.tasks, ti)
}

func (tp *HeapTaskPool) cancelTask(key string) {
	if t, ok := tp.keys[key]; ok {
		t.canceled = true
	}
}

func (tp *HeapTaskPool) loop() {
	defer tp.Done()

	nextTaskInterval := time.Second

	for {
		select {
		case <-tp.closed:
			return
		case taskI := <-tp.taskOp:
			switch taskI.opType {
			case taskOpAdd, taskOpUpdate:
				tp.addOrUpdateTask(taskI.at, taskI.key, taskI.exec, taskI.params)
			case taskOpDel:
				tp.cancelTask(taskI.key)
			}

			nextTaskInterval = tp.process()
		case <-time.After(nextTaskInterval):
			nextTaskInterval = tp.process()
		}
	}
}

func (tp *HeapTaskPool) AddTask(key string, t time.Time, exec TaskFunc, params ...interface{}) error {
	if exec == nil {
		return commerrx.ErrInvalidArgument
	}

	tp.taskOp <- &taskOpInfo{
		opType: taskOpAdd,
		key:    key,
		at:     t,
		exec:   exec,
		params: params,
	}

	return nil
}

func (tp *HeapTaskPool) RemoveTask(key string) error {
	if key == "" {
		return commerrx.ErrInvalidArgument
	}

	tp.taskOp <- &taskOpInfo{
		opType: taskOpDel,
		key:    key,
	}

	return nil
}
