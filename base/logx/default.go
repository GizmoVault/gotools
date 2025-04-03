package logx

func NewNopLoggerWrapper() Wrapper {
	return NewWrapper(&NopLogger{})
}

func NewConsoleLoggerWrapper() Wrapper {
	logger := NewCommLogger(&ConsoleRecorder{})
	logger.SetLevel(LevelInfo)

	return NewWrapper(logger)
}

func NewFileLoggerWrapper(filePath string) Wrapper {
	logger := NewCommLogger(NewFileRecorder(filePath))
	logger.SetLevel(LevelInfo)

	return NewWrapper(logger)
}
