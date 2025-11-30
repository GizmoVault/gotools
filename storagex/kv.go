package storagex

import (
	"encoding/json"
	"sync"

	"github.com/GizmoVault/gotools/base/errorx"
	"github.com/GizmoVault/gotools/base/syncx"
)

func NewKV(file string) (StorageTiny2, error) {
	return NewKVEx(file, nil)
}

func NewKVEx(file string, storage FileStorage) (StorageTiny2, error) {
	stg, err := NewMemWithFile[map[string]string, Serial, syncx.RWLocker](make(map[string]string), &JSONSerial{},
		&sync.RWMutex{}, file, storage)
	if err != nil {
		return nil, err
	}

	return &kvImpl{
		d: stg,
	}, nil
}

type kvImpl struct {
	d *MemWithFile[map[string]string, Serial, syncx.RWLocker]
}

func (impl *kvImpl) GetList(itemGen func(key string) interface{}) (items []interface{}, err error) {
	if itemGen == nil {
		err = errorx.ErrInvalidArgs

		return
	}

	impl.d.Read(func(values map[string]string) {
		for key, value := range values {
			item := itemGen(key)
			if item == nil {
				continue
			}

			err = json.Unmarshal([]byte(value), &item)
			if err != nil {
				continue
			}

			items = append(items, item)
		}
	})

	return
}

func (impl *kvImpl) GetMap(itemGen func(key string) interface{}) (items map[string]interface{}, err error) {
	if itemGen == nil {
		err = errorx.ErrInvalidArgs

		return
	}

	items = make(map[string]interface{})

	impl.d.Read(func(values map[string]string) {
		for key, value := range values {
			item := itemGen(key)
			if item == nil {
				continue
			}

			err = json.Unmarshal([]byte(value), &item)
			if err != nil {
				continue
			}

			items[key] = item
		}
	})

	return
}

func (impl *kvImpl) Set(key string, v interface{}) error {
	return impl.SetAll([]string{key}, v)
}

func (impl *kvImpl) Get(key string, v interface{}) (ok bool, err error) {
	vs, err := impl.GetAll([]string{key}, v)
	if err != nil {
		return
	}

	ok = vs[0] != nil

	return
}

func (impl *kvImpl) Del(key string) error {
	return impl.DelAll([]string{key})
}

func (impl *kvImpl) SetAll(keys []string, vs ...interface{}) error {
	if len(keys) != len(vs) {
		return errorx.ErrInvalidArgs
	}

	ds := make([][]byte, 0, len(keys))

	for _, v := range vs {
		d, err := json.Marshal(v)
		if err != nil {
			return err
		}

		ds = append(ds, d)
	}

	return impl.d.Change(func(v map[string]string) (newV map[string]string, err error) {
		newV = v

		if newV == nil {
			newV = make(map[string]string)
		}

		for idx := range keys {
			newV[keys[idx]] = string(ds[idx])
		}

		return
	})
}

func (impl *kvImpl) GetAll(keys []string, vsi ...interface{}) (vs []interface{}, err error) {
	var ds []string

	impl.d.Read(func(v map[string]string) {
		for _, key := range keys {
			ds = append(ds, v[key])
		}
	})

	vs = make([]interface{}, len(keys))

	for idx := range ds {
		if ds[idx] == "" {
			vs[idx] = nil

			continue
		}

		if idx >= len(vsi) || vsi[idx] == nil {
			vs[idx] = ds[idx]

			continue
		}

		err = json.Unmarshal([]byte(ds[idx]), vsi[idx])
		if err != nil {
			return
		}

		vs[idx] = vsi[idx]
	}

	return
}

func (impl *kvImpl) DelAll(keys []string) error {
	return impl.d.Change(func(v map[string]string) (newV map[string]string, err error) {
		newV = v

		if newV == nil {
			newV = make(map[string]string)
		}

		for _, key := range keys {
			delete(newV, key)
		}

		return
	})
}
