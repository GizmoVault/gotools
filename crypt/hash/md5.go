package hash

import "crypto/md5" //nolint:gosec // no problem

func MD5Sum(data []byte) [16]byte {
	return md5.Sum(data) //nolint:gosec // no problem
}
