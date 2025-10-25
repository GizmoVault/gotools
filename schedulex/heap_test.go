package schedulex

import (
	"testing"
	"time"

	"container/heap"
	"github.com/stretchr/testify/assert"
)

func Test_taskHeap(t *testing.T) {
	h := &taskHeap{}
	heap.Init(h)

	heap.Push(h, &taskItem{
		at: time.Now().Add(10 * time.Second),
	})
	heap.Push(h, &taskItem{
		at: time.Now().Add(8 * time.Second),
	})
	heap.Push(h, &taskItem{
		at: time.Now().Add(12 * time.Second),
	})

	v := (*h)[0]
	t.Log(v.at)

	vi := heap.Pop(h)
	t.Log(vi.(*taskItem).at)

	v = (*h)[0]
	t.Log(v.at)

	vi = heap.Pop(h)
	t.Log(vi.(*taskItem).at)

	v = (*h)[0]
	t.Log(v.at)

	vi = heap.Pop(h)
	t.Log(vi.(*taskItem).at)
}

func Test_HeapTaskPool(t *testing.T) {
	hp := CreateHeapTaskPool()

	err := hp.AddTask("1", time.Now().Add(time.Second*3), TaskFunc(func(key string, _ ...any) {
		t.Log("cb ", key)
	}))
	assert.NoError(t, err)

	err = hp.AddTask("2", time.Now().Add(time.Second), TaskFunc(func(key string, _ ...any) {
		t.Log("cb ", key)
	}))
	assert.NoError(t, err)

	err = hp.AddTask("3", time.Now().Add(time.Second*2), TaskFunc(func(key string, _ ...any) {
		t.Log("cb ", key)
	}))
	assert.NoError(t, err)

	time.Sleep(time.Second * 4)

	hp.Stop()
}
