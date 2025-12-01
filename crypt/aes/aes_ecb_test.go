package aes

import (
	"encoding/base64"
	"encoding/hex"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestECB1(t *testing.T) {
	rawData := "0.56"
	key := "0123456789abcdef"

	encryptedData, err := ECBEncrypt([]byte(rawData), []byte(key))
	assert.Nil(t, err)

	encryptedDataBase64 := base64.RawURLEncoding.EncodeToString(encryptedData)
	t.Log(encryptedDataBase64)

	dd, err := base64.RawURLEncoding.DecodeString(encryptedDataBase64)
	assert.Nil(t, err)
	rawDataB, err := ECBDecrypt(dd, []byte(key))
	assert.Nil(t, err)
	t.Log(string(rawDataB))
}

func TestECB2(t *testing.T) {
	_, err := ECBDecrypt([]byte("11"), []byte("0123456789abcdef"))
	assert.NotNil(t, err)
}

func TestECB3(t *testing.T) {
	encryptedData, err := ECBEncrypt([]byte("300"), []byte("c4c7cb4150306b77aa6580f54a077e1c"))
	assert.Nil(t, err)

	encryptedDataBase64 := base64.RawURLEncoding.EncodeToString(encryptedData)
	t.Log(encryptedDataBase64)
}

func TestECB4(t *testing.T) {
	dd, err := base64.RawURLEncoding.DecodeString("I_4tmKEIbSBsfONPi6dfGQ")
	assert.Nil(t, err)
	ori, err := ECBDecrypt(dd, []byte("1234123412341234123412341234abcd"))
	assert.Nil(t, err)
	v, err := strconv.Atoi(string(ori))
	assert.Nil(t, err)
	assert.Equal(t, v, 1000)
}

func TestECB5(t *testing.T) {
	hexS := "dfd5bb7564ddfd204e08cf146e71baa8"
	decodeData := make([]byte, hex.DecodedLen(len(hexS)))
	_, err := hex.Decode(decodeData, []byte(hexS))
	assert.Nil(t, err)
	d, err := ECBDecrypt(decodeData, []byte("0123456789abcdef"))
	assert.Nil(t, err)
	t.Log(string(d))
}
