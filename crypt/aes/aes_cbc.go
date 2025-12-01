package aes

import (
	"crypto/aes"
	"crypto/cipher"

	"github.com/GizmoVault/gotools/base/errorx"
)

func CBCEncrypt(origData, key []byte) (crypted []byte, err error) {
	defer func() {
		if errR := recover(); errR != nil {
			err = errorx.ErrCrashed
		}
	}()

	block, err := aes.NewCipher(key)
	if err != nil {
		err = errorx.FromErrorAndMessage(err, "aes:NewCipher")

		return
	}

	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted = make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)

	return
}

func CBCDecrypt(encryptedData, key []byte) (decryptedData []byte, err error) {
	defer func() {
		if errR := recover(); errR != nil {
			err = errorx.ErrCrashed
		}
	}()

	block, err := aes.NewCipher(key)
	if err != nil {
		err = errorx.FromErrorAndMessage(err, "aes:NewCipher")

		return
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(encryptedData))
	blockMode.CryptBlocks(origData, encryptedData)

	decryptedData, err = PKCS5UnPadding(origData)

	return
}
