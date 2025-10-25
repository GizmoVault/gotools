package logx

import "github.com/GizmoVault/gotools/base"

func NewNopLoggerWrapper() Wrapper {
	return NewWrapper(&NopLogger{})
}

func NewConsoleLoggerWrapper() Wrapper {
	logger := NewCommLogger(nil, &ConsoleRecorder{})
	logger.SetLevel(LevelInfo)

	return NewWrapper(logger)
}

func NewConsoleLoggerWrapperWithFNNow(now base.FNNow) Wrapper {
	logger := NewCommLogger(now, &ConsoleRecorder{})
	logger.SetLevel(LevelInfo)

	return NewWrapper(logger)
}

func NewFileLoggerWrapper(filePath string) Wrapper {
	logger := NewCommLogger(nil, NewFileRecorder(filePath))
	logger.SetLevel(LevelInfo)

	return NewWrapper(logger)
}
