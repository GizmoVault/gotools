package storagex

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

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

	return impl.trySafeWriteFile(name, data)
}

func (impl *fsStorageImpl) ReadFile(name string) ([]byte, error) {
	if !path.IsAbs(name) {
		name = filepath.Join(impl.rootPath, name)
	}

	return impl.trySafeReadFile(name)
}

func (*fsStorageImpl) fileNameBackup(name string) string {
	return name + ".bak"
}

func (*fsStorageImpl) fileNameBackupDone(name string) string {
	return name + ".bak.done"
}

func (impl *fsStorageImpl) trySafeReadFile(name string) ([]byte, error) {
	fileBackup := impl.fileNameBackup(name)
	fileBackupDone := impl.fileNameBackupDone(name)

	ok, err := pathx.IsFileExists(fileBackup)
	if err == nil && ok {
		if ok2, err2 := pathx.IsFileExists(fileBackupDone); err2 == nil && ok2 {
			_ = os.Rename(name, fmt.Sprintf("%s.r.%d", name, time.Now().UnixMilli()))

			_ = os.Rename(fileBackup, name)
		}
	}

	_ = os.Remove(fileBackup)
	_ = os.Remove(fileBackupDone)

	return os.ReadFile(name)
}

func (impl *fsStorageImpl) trySafeWriteFile(name string, data []byte) (err error) {
	fileBackup := impl.fileNameBackup(name)
	fileBackupDone := impl.fileNameBackupDone(name)

	ok, err := pathx.IsFileExists(fileBackup)
	if err == nil && ok {
		if ok2, err2 := pathx.IsFileExists(fileBackupDone); err2 == nil && ok2 {
			_ = os.Rename(name, fmt.Sprintf("%s.r.w.%d", name, time.Now().UnixMilli()))

			_ = os.Rename(fileBackup, name)
		}
	}

	_ = os.Remove(fileBackup)
	_ = os.Remove(fileBackupDone)

	exists, err := pathx.IsFileExists(name)
	if err == nil && exists {
		err = os.Rename(name, fileBackup)
		if err != nil {
			return
		}

		_ = os.WriteFile(fileBackupDone, []byte(time.Now().String()), constx.PermReadWrite)
	}

	err = os.WriteFile(name, data, constx.PermReadWrite)
	if err != nil {
		return
	}

	_ = os.Remove(fileBackupDone)
	_ = os.Remove(fileBackup)

	return
}
