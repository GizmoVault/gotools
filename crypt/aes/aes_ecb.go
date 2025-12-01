package aes

import (
	"crypto/aes"
	"github.com/GizmoVault/gotools/base/errorx"
)

type PaddingType int

const (
	PaddingTypeNone PaddingType = iota
	PaddingTypePKCS5
	PaddingTypePKCS7
)

func ECBEncrypt(origData, key []byte) (encryptedData []byte, err error) {
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

	en := NewECBEncryptor(block)

	origData = PKCS7Padding(origData, block.BlockSize())
	encryptedData = make([]byte, len(origData))

	en.CryptBlocks(encryptedData, origData)

	return
}

func ECBEncryptEx(origData, key []byte, paddingType PaddingType) (encryptedData []byte, err error) {
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

	en := NewECBEncryptor(block)

	switch paddingType {
	case PaddingTypePKCS5:
		origData = PKCS5Padding(origData)
	case PaddingTypePKCS7:
		origData = PKCS7Padding(origData, block.BlockSize())
	default:
	}

	encryptedData = make([]byte, len(origData))

	en.CryptBlocks(encryptedData, origData)

	return
}

func ECBDecrypt(encryptedData, key []byte) (origData []byte, err error) {
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

	en := NewECBDecrypter(block)
	origData = make([]byte, len(encryptedData))

	en.CryptBlocks(origData, encryptedData)

	origData, err = PKCS5UnPadding(origData)

	return
}

func ECBDecryptEx(encryptedData, key []byte, paddingType PaddingType) (origData []byte, err error) {
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

	en := NewECBDecrypter(block)
	origData = make([]byte, len(encryptedData))
	en.CryptBlocks(origData, encryptedData)

	switch paddingType {
	case PaddingTypePKCS5:
		origData, err = PKCS5UnPadding(origData)
	case PaddingTypePKCS7:
		origData, err = PKCS7UnPadding(origData)
	default:
	}

	return
}
