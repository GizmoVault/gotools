package logx

import "context"

type Level int

const (
	LevelFatal Level = 10
	LevelError Level = 11
	LevelWarn  Level = 12
	LevelInfo  Level = 13
	LevelDebug Level = 14
)

type Logger interface {
	SetLevel(level Level)
	WithFields(fields ...Field) Logger
	Log(level Level, a ...interface{})
	Logf(level Level, format string, a ...interface{})
}

type LoggerWithContext interface {
	Logger

	LogWithContext(ctx context.Context, level Level, a ...interface{})
	LogWithContextf(ctx context.Context, level Level, format string, a ...interface{})
}

type WrapperWithContext interface {
	GetLogger() Logger

	WithFields(fields ...Field) WrapperWithContext

	Fatal(ctx context.Context, a ...interface{})
	Fatalf(ctx context.Context, format string, a ...interface{})
	Error(ctx context.Context, a ...interface{})
	Errorf(ctx context.Context, format string, a ...interface{})
	Warn(ctx context.Context, a ...interface{})
	Warnf(ctx context.Context, format string, a ...interface{})
	Info(ctx context.Context, a ...interface{})
	Infof(ctx context.Context, format string, a ...interface{})
	Debug(ctx context.Context, a ...interface{})
	Debugf(ctx context.Context, format string, a ...interface{})
}

type Wrapper interface {
	GetLogger() Logger

	WithFields(fields ...Field) Wrapper

	Fatal(a ...interface{})
	Fatalf(format string, a ...interface{})
	Error(a ...interface{})
	Errorf(format string, a ...interface{})
	Warn(a ...interface{})
	Warnf(format string, a ...interface{})
	Info(a ...interface{})
	Infof(format string, a ...interface{})
	Debug(a ...interface{})
	Debugf(format string, a ...interface{})

	GetWrapperWithContext() WrapperWithContext
}
