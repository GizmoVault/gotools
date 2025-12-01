package aes

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test1(t *testing.T) {
	d, e := hex.DecodeString("fd0627f6b2f9499e74a425ffdb67c774")
	assert.Nil(t, e)
	d, e = CBCDecrypt(d, []byte("8973aa0f19b98a2f"))
	assert.Nil(t, e)
	assert.Equal(t, string(d), "800|1467792422")

	d, e = CBCEncrypt([]byte("800|1467792422"), []byte("8973aa0f19b98a2f"))
	assert.Nil(t, e)
	assert.Equal(t, hex.EncodeToString(d), "fd0627f6b2f9499e74a425ffdb67c774")

	d, e = CBCEncrypt([]byte("800|1467792422"), []byte("8973aa0f19b98a2f28f70a37f80229d1"))
	assert.Nil(t, e)
	assert.Equal(t, hex.EncodeToString(d), "39653a811d14a79c5e80e3536179c459")

	//
	d, e = hex.DecodeString("fd0627f6b2f9499e74a425ffdb67c774")
	assert.Nil(t, e)
	_, e = CBCDecrypt(d, []byte("8973aa0f19b98a2f28f70a37f80229d1"))
	assert.NotNil(t, e)
}
