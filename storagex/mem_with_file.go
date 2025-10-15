package storagex

import (
	"errors"
	"os"
	"time"
)

type Serial interface {
	Marshal(t any) ([]byte, error)
	Unmarshal(d []byte, t any) error
}

type Lock interface {
	RLock()
	RUnlock()

	Lock()
	Unlock()
}

type EventObserver[T any] interface {
	BeforeLoad()
	AfterLoad(memD T, err error)
	BeforeSave()
	AfterSave(memD T, err error)
}

type MemWithFile[T any, S Serial, L Lock] struct {
	memD   T
	serial S
	lock   L

	fileName string
	storage  FileStorage
	ob       EventObserver[T]

	changedFlag      bool
	autoSaveInterval time.Duration
}

func NewMemWithFile[T any, S Serial, L Lock](d T, serial S, lock L, fileName string, storage FileStorage) *MemWithFile[T, S, L] {
	return NewMemWithFileEx(d, serial, lock, fileName, storage, nil)
}

func NewMemWithFileEx[T any, S Serial, L Lock](d T, serial S, lock L, fileName string, storage FileStorage,
	ob EventObserver[T]) *MemWithFile[T, S, L] {
	return NewMemWithFileEx1(d, serial, lock, fileName, storage, ob, 0)
}

func NewMemWithFileEx1[T any, S Serial, L Lock](d T, serial S, lock L, fileName string, storage FileStorage,
	ob EventObserver[T], autoSaveInterval time.Duration) *MemWithFile[T, S, L] {
	if storage == nil && fileName != "" {
		storage = NewRawFSStorage("")
	}

	mwf := &MemWithFile[T, S, L]{
		memD:             d,
		serial:           serial,
		lock:             lock,
		fileName:         fileName,
		storage:          storage,
		ob:               ob,
		autoSaveInterval: autoSaveInterval,
	}

	_ = mwf.load()

	if autoSaveInterval > 0 {
		go mwf.autoSaveRoutine()
	}

	return mwf
}

func (mwf *MemWithFile[T, S, L]) autoSaveRoutine() {
	for {
		time.Sleep(mwf.autoSaveInterval)

		mwf.lock.Lock()

		if mwf.changedFlag {
			mwf.changedFlag = false

			_ = mwf.save()
		}

		mwf.lock.Unlock()
	}
}

func (mwf *MemWithFile[T, S, L]) Read(proc func(memD T)) {
	mwf.lock.RLock()
	defer mwf.lock.RUnlock()

	proc(mwf.memD)
}

func (mwf *MemWithFile[T, S, L]) Change(proc func(memD T) (newMemD T, err error)) error {
	mwf.lock.Lock()
	defer mwf.lock.Unlock()

	newMemD, err := proc(mwf.memD)
	if err != nil {
		return err
	}

	mwf.memD = newMemD

	if mwf.autoSaveInterval <= 0 {
		return mwf.save()
	}

	mwf.changedFlag = true

	return nil
}

func (mwf *MemWithFile[T, S, L]) load() error {
	if mwf.fileName == "" {
		return nil
	}

	if mwf.ob != nil {
		mwf.ob.BeforeLoad()
	}

	d, err := mwf.storage.ReadFile(mwf.fileName)
	if err != nil {
		var pathError *os.PathError
		if errors.As(err, &pathError) {
			err = nil
		}

		if mwf.ob != nil {
			mwf.ob.AfterLoad(mwf.memD, err)
		}

		return err
	}

	var m T

	err = mwf.serial.Unmarshal(d, &m)
	if err != nil {
		if mwf.ob != nil {
			mwf.ob.AfterLoad(mwf.memD, err)
		}

		return err
	}

	mwf.memD = m

	if mwf.ob != nil {
		mwf.ob.AfterLoad(mwf.memD, nil)
	}

	return nil
}

func (mwf *MemWithFile[T, S, L]) save() error {
	if mwf.fileName == "" {
		return nil
	}

	if mwf.ob != nil {
		mwf.ob.BeforeSave()
	}

	d, err := mwf.serial.Marshal(mwf.memD)
	if err != nil {
		if mwf.ob != nil {
			mwf.ob.AfterSave(mwf.memD, err)
		}

		return err
	}

	err = mwf.storage.WriteFile(mwf.fileName, d)
	if err != nil {
		if mwf.ob != nil {
			mwf.ob.AfterSave(mwf.memD, err)
		}

		return err
	}

	if mwf.ob != nil {
		mwf.ob.AfterSave(mwf.memD, nil)
	}

	return nil
}
