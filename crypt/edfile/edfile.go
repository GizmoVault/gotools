package edfile

import (
	erand "crypto/rand"
	"encoding/binary"
	"github.com/GizmoVault/gotools/crypt/hash"
	"math/rand"
	"os"

	"github.com/GizmoVault/gotools/base/errorx"
	"github.com/GizmoVault/gotools/crypt/aes"
	"github.com/GizmoVault/gotools/pathx"
)

func EncodePlainFile(d []byte) []byte {
	r := make([]byte, len(d)+100)

	_, _ = erand.Read(r)

	startPos := rand.Intn(60) //nolint:gosec // fixme

	binary.LittleEndian.PutUint32(r[10:], uint32(len(d)))   //nolint:gosec // fixme
	binary.LittleEndian.PutUint32(r[14:], uint32(startPos)) //nolint:gosec // fixme

	copy(r[startPos+20:], d)

	return r
}

func DecodePlainFile(d []byte) (dd []byte, ok bool) {
	if len(d) < 100 {
		return
	}

	dLen := binary.LittleEndian.Uint32(d[10:])
	dPos := binary.LittleEndian.Uint32(d[14:]) + 20

	//nolint:gosec // todo
	if dPos+dLen >= uint32(len(d)) {
		return
	}

	dd = make([]byte, dLen)
	copy(dd, d[dPos:])

	ok = true

	return
}

func deriveSecKeyFromKeyS(key string) []byte {
	sum := hash.MD5Sum([]byte(key))

	return sum[:]
}

func WriteSecFile(name, key string, data []byte) (err error) {
	ed, err := aes.CBCEncrypt(EncodePlainFile(data), deriveSecKeyFromKeyS(key))
	if err != nil {
		return
	}

	_ = pathx.MustDirOfFileExists(name)

	err = os.WriteFile(name, ed, 0600)

	return
}

func ReadSecFile(name, key string) (data []byte, err error) {
	d, err := os.ReadFile(name)
	if err != nil {
		return
	}

	dd, err := aes.CBCDecrypt(d, deriveSecKeyFromKeyS(key))
	if err != nil {
		return
	}

	data, ok := DecodePlainFile(dd)
	if !ok {
		err = errorx.ErrFail

		return
	}

	return
}
