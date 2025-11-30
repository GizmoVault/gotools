package storagex_test

import (
	"os"
	"testing"

	"github.com/GizmoVault/gotools/storagex"
	"github.com/stretchr/testify/assert"
)

type utKVItem struct {
	N int
	S string
}

func TestSimpleKV(t *testing.T) {
	_ = os.Remove("tmp/ut.txt")

	kv, err := storagex.NewKV("tmp/ut.txt")
	assert.NoError(t, err)

	var item utKVItem

	ok, err := kv.Get("key", &item)
	assert.Nil(t, err)
	assert.False(t, ok)

	err = kv.Set("key", &utKVItem{
		N: 10,
		S: "S20S",
	})
	assert.Nil(t, err)

	ok, err = kv.Get("key", &item)
	assert.Nil(t, err)
	assert.True(t, ok)
	assert.EqualValues(t, 10, item.N)
	assert.EqualValues(t, "S20S", item.S)

	kv2, err := storagex.NewKV("tmp/ut.txt")
	assert.NoError(t, err)

	ok, err = kv2.Get("key", &item)
	assert.Nil(t, err)
	assert.True(t, ok)
	assert.EqualValues(t, 10, item.N)
	assert.EqualValues(t, "S20S", item.S)
}

func TestSimpleKV2(t *testing.T) {
	_ = os.Remove("tmp/ut.txt")

	kv, err := storagex.NewKV("tmp/ut.txt")
	assert.NoError(t, err)

	items, err := kv.GetAll([]string{"key", "key1"}, &utKVItem{})
	assert.Nil(t, err)
	assert.True(t, items[0] == nil)

	err = kv.SetAll([]string{"key", "key1"}, &utKVItem{
		N: 10,
		S: "S10S",
	}, &utKVItem{
		N: 20,
		S: "S20S",
	})
	assert.Nil(t, err)

	items, err = kv.GetAll([]string{"key", "key1", "key2"}, &utKVItem{}, &utKVItem{}, &utKVItem{})
	assert.Nil(t, err)
	assert.True(t, len(items) == 3)
	assert.NotNil(t, items[0])
	assert.NotNil(t, items[1])
	assert.Nil(t, items[2])
	assert.EqualValues(t, 10, items[0].(*utKVItem).N)
	assert.EqualValues(t, "S10S", items[0].(*utKVItem).S)
	assert.EqualValues(t, 20, items[1].(*utKVItem).N)
	assert.EqualValues(t, "S20S", items[1].(*utKVItem).S)

	kv2, err := storagex.NewKV("tmp/ut.txt")
	assert.NoError(t, err)

	items, err = kv2.GetAll([]string{"key", "key1", "key2"}, &utKVItem{}, &utKVItem{}, &utKVItem{})
	assert.Nil(t, err)
	assert.True(t, len(items) == 3)
	assert.NotNil(t, items[0])
	assert.NotNil(t, items[1])
	assert.Nil(t, items[2])
	assert.EqualValues(t, 10, items[0].(*utKVItem).N)
	assert.EqualValues(t, "S10S", items[0].(*utKVItem).S)
	assert.EqualValues(t, 20, items[1].(*utKVItem).N)
	assert.EqualValues(t, "S20S", items[1].(*utKVItem).S)
}
