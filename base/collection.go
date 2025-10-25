package base

import (
	"reflect"

	"github.com/GizmoVault/gotools/base/commerrx"
	"github.com/GizmoVault/gotools/base/syncx"
)

type Entity interface {
	GetID() uint64
}

type Collection[T Entity, L syncx.RWLocker] struct {
	lock  L
	items map[uint64]T
}

func NewCollection[T Entity, L syncx.RWLocker](l L) *Collection[T, L] {
	return &Collection[T, L]{
		lock:  l,
		items: make(map[uint64]T),
	}
}

func (c *Collection[T, L]) Add(item T) error {
	if reflect.ValueOf(item).IsNil() { // 第25行
		return commerrx.ErrInvalidArgument
	}

	id := item.GetID()
	if id == 0 {
		return commerrx.ErrInvalidArgument
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	if _, exists := c.items[id]; exists {
		return commerrx.ErrAlreadyExists
	}

	c.items[id] = item

	return nil
}

func (c *Collection[T, L]) Remove(item T) error {
	if reflect.ValueOf(item).IsNil() { // 第25行
		return commerrx.ErrInvalidArgument
	}

	return c.RemoveByID(item.GetID())
}

func (c *Collection[T, L]) RemoveByID(id uint64) error {
	if id == 0 {
		return commerrx.ErrInvalidArgument
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	if _, exists := c.items[id]; !exists {
		return commerrx.ErrNotFound
	}

	delete(c.items, id)

	return nil
}

func (c *Collection[T, L]) Get(id uint64) (T, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	var zero T

	if item, exists := c.items[id]; exists {
		return item, true
	}

	return zero, false
}

func (c *Collection[T, L]) Len() int {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return len(c.items)
}

func (c *Collection[T, L]) GetItems() []T {
	c.lock.RLock()
	defer c.lock.RUnlock()

	items := make([]T, 0, len(c.items))
	for _, item := range c.items {
		items = append(items, item)
	}

	return items
}

func (c *Collection[T, L]) Walk(callback func(T) error) (err error) {
	if callback == nil {
		return commerrx.ErrInvalidArgument
	}

	for _, item := range c.GetItems() {
		err = callback(item)
		if err != nil {
			break
		}
	}

	return
}

func (c *Collection[T, L]) WalkASync(callback func(T)) {
	if callback == nil {
		return
	}

	for _, item := range c.GetItems() {
		go func(item T) {
			callback(item)
		}(item)
	}
}
