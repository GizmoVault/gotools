package impls

import (
	"github.com/GizmoVault/gotools/base/logx"
	"github.com/sirupsen/logrus"
)

func NewLogrus() logx.Logger {
	return NewLogrusEx(nil)
}

func NewLogrusEx(logger *logrus.Logger) logx.Logger {
	if logger == nil {
		logger = logrus.New()
		logger.SetFormatter(new(logrus.JSONFormatter))
	}

	logger.SetLevel(logrus.TraceLevel)

	return &logrusImpl{
		rl:    logrus.NewEntry(logger),
		level: logx.LevelInfo,
	}
}

type logrusImpl struct {
	rl    *logrus.Entry
	level logx.Level
}

func (*logrusImpl) mapLevel(level logx.Level) logrus.Level {
	switch level {
	case logx.LevelFatal:
		return logrus.FatalLevel
	case logx.LevelError:
		return logrus.ErrorLevel
	case logx.LevelWarn:
		return logrus.WarnLevel
	case logx.LevelInfo:
		return logrus.InfoLevel
	case logx.LevelDebug:
		return logrus.DebugLevel
	}

	return logrus.FatalLevel
}

func (impl *logrusImpl) SetLevel(level logx.Level) {
	impl.level = level
}

func (impl *logrusImpl) WithFields(fields ...logx.Field) logx.Logger {
	fs := make(map[string]interface{})
	for _, field := range fields {
		fs[field.K] = field.V
	}

	return &logrusImpl{
		rl:    impl.rl.WithFields(fs),
		level: impl.level,
	}
}

func (impl *logrusImpl) Log(level logx.Level, a ...interface{}) {
	if level > impl.level {
		return
	}

	mLevel := impl.mapLevel(level)
	if mLevel == logrus.FatalLevel {
		impl.rl.Fatal(a...)
	} else {
		impl.rl.Log(mLevel, a...)
	}
}

func (impl *logrusImpl) Logf(level logx.Level, format string, a ...interface{}) {
	if level > impl.level {
		return
	}

	mLevel := impl.mapLevel(level)
	if mLevel == logrus.FatalLevel {
		impl.rl.Fatalf(format, a...)
	} else {
		impl.rl.Logf(impl.mapLevel(level), format, a...)
	}
}
