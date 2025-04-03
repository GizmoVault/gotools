package hashx

import (
	"hash"

	"crypto/hmac"
	"crypto/md5"  //nolint:gosec // .
	"crypto/sha1" //nolint:gosec // .
	"encoding/hex"
)

func MD5(s string) string {
	sum := md5.Sum([]byte(s)) //nolint:gosec // .

	return hex.EncodeToString(sum[:])
}

func SHA1(s string) string {
	sum := sha1.Sum([]byte(s)) //nolint:gosec // .

	return hex.EncodeToString(sum[:])
}

func HMac(key []byte, s string) string {
	return HMacEx(sha1.New, key, s)
}

func HMacEx(h func() hash.Hash, key []byte, s string) string {
	mac := hmac.New(h, key)
	mac.Write([]byte(s))
	sum := mac.Sum(nil)

	return hex.EncodeToString(sum)
}
