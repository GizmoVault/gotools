package wrapper

import "sync"

func DoInLock(locker sync.Locker, fn func(tag any), fnTag any) {
	locker.Lock()
	defer locker.Unlock()

	fn(fnTag)
}
