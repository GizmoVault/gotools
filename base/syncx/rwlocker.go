package syncx

import "sync"

type RLocker interface {
	RLock()
	RUnlock()
}

type RWLocker interface {
	sync.Locker
	RLocker
}
