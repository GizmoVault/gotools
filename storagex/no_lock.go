package storagex

type NoLock struct {
}

func (NoLock) RLock() {

}

func (NoLock) RUnlock() {

}

func (NoLock) Lock() {

}

func (NoLock) Unlock() {

}
