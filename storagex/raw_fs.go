package storagex

import (
	"os"
	"path"
	"path/filepath"

	"github.com/GizmoVault/gotools/base/constx"
	"github.com/GizmoVault/gotools/pathx"
)

func NewRawFSStorage(rootPath string) FileStorage {
	if rootPath == "" {
		rootPath, _ = os.Getwd()
	}

	return &fsStorageImpl{
		rootPath: rootPath,
	}
}

type fsStorageImpl struct {
	rootPath string
}

func (impl *fsStorageImpl) WriteFile(name string, data []byte) error {
	if !path.IsAbs(name) {
		name = filepath.Join(impl.rootPath, name)
	}

	_ = pathx.MustDirOfFileExists(name)

	if impl.trySafeWriteFile(name, data) {
		return nil
	}

	return os.WriteFile(name, data, constx.PermReadWrite)
}

func (impl *fsStorageImpl) ReadFile(name string) ([]byte, error) {
	if !path.IsAbs(name) {
		name = filepath.Join(impl.rootPath, name)
	}

	return os.ReadFile(name)
}

func (*fsStorageImpl) trySafeWriteFile(name string, data []byte) (ok bool) {
	exists, err := pathx.IsFileExists(name)
	if err != nil {
		return
	}

	if !exists {
		return
	}

	nameBak := name + ".bak"

	if o, e := pathx.IsFileExists(nameBak); e == nil && o {
		_ = os.Remove(nameBak)
	}

	err = os.Rename(name, nameBak)
	if err != nil {
		return
	}

	err = os.WriteFile(name, data, constx.PermReadWrite)
	if err != nil {
		return
	}

	_ = os.Remove(nameBak)

	ok = true

	return
}
