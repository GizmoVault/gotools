package storagex_test

import (
	"testing"

	"github.com/GizmoVault/gotools/base/syncx"
	"github.com/GizmoVault/gotools/storagex"
	"github.com/stretchr/testify/assert"
)

func TestMemAndFile1(t *testing.T) {
	mm, err := storagex.NewMemWithFile[map[int]string, storagex.Serial, syncx.RWLocker](make(map[int]string),
		&storagex.JSONSerial{}, &storagex.NoLock{}, "tmp/utStorage.txt", nil)
	assert.NoError(t, err)
	t.Log(mm)

	memWithFile, err := storagex.NewMemWithFile(make(map[int]string), &storagex.JSONSerial{}, &storagex.NoLock{},
		"tmp/utStorage.txt", nil)
	assert.NoError(t, err)
	assert.NotNil(t, memWithFile)

	_ = memWithFile.Change(func(m map[int]string) (map[int]string, error) {
		if m == nil {
			m = make(map[int]string)
		}

		m[1] = "1xx"

		return m, nil
	})

	memWithFile.Read(func(m map[int]string) {
		t.Log(m[1])
	})
}
