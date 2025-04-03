package logx

import (
	"fmt"
	"os"
	"path"
)

type FileRecorder interface {
}

func NewFileRecorder(filePath string) Recorder {
	dir, _ := path.Split(filePath)

	_ = os.MkdirAll(dir, os.ModePerm)

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666) //nolint: mnd // .
	if err != nil {
		return nil
	}

	return &fileRecorderImpl{
		f: f,
	}
}

type fileRecorderImpl struct {
	f *os.File
}

func (impl *fileRecorderImpl) Log(_ Level, a ...interface{}) {
	_, _ = fmt.Fprintln(impl.f, a...)
}

func (impl *fileRecorderImpl) Logf(_ Level, format string, a ...interface{}) {
	_, _ = impl.f.WriteString(fmt.Sprintf(format, a...) + "\n")
}
