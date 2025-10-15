package logx

import "testing"

func TestConsoleLoggerWrapper(_ *testing.T) {
	logger := NewConsoleLoggerWrapper()
	logger.GetLogger().SetLevel(LevelWarn)

	logger.Debug("a")
	logger.Info("b")
	logger.Warn("c")
	logger.Error("d")
}

func TestFileLoggerWrapper(_ *testing.T) {
	logger := NewFileLoggerWrapper("tmp/file.log")
	logger.GetLogger().SetLevel(LevelWarn)

	logger.Debug("a")
	logger.Info("b")
	logger.Warn("c")
	logger.Error("d")
}

func TestFileLoggerWrapper2(_ *testing.T) {
	logger := NewCommLoggerEx(true, NewFileRecorder("tmp/file.log"))
	logger.SetLevel(LevelInfo)

	lw := NewWrapper(logger)
	lw.Debug("a")
	lw.Info("b")
	lw.Warn("c")
	lw.Error("d")
	lw.Errorf("e:%d", 43)
}
