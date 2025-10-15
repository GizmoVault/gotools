package storagex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFsLoadSave(t *testing.T) {
	t.SkipNow()

	stg := NewRawFSStorage("tmp")

	d, err := stg.ReadFile("test.txt")
	assert.NoError(t, err)
	t.Log(string(d))

	assert.NoError(t, stg.WriteFile("test.txt", []byte("hello world")))
	assert.NoError(t, stg.WriteFile("test.txt", []byte("hello world2")))
}
