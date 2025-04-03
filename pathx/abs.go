package pathx

import (
	"path"
	"path/filepath"
)

func Abs(filePath string) string {
	r, _ := filepath.Abs(filePath)

	return r
}

func AbsEx(baseRoot, filePath string) string {
	if path.IsAbs(filePath) {
		return filePath
	}

	return filepath.Join(baseRoot, filePath)
}
