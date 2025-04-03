package logx

import "context"

func NewWrapper(l Logger) Wrapper {
	return &wrapperImpl{
		l: l,
	}
}

type wrapperImpl struct {
	l Logger
}

func (wrapper *wrapperImpl) GetLogger() Logger {
	return wrapper.l
}

func (wrapper *wrapperImpl) WithFields(fields ...Field) Wrapper {
	l := wrapper.l.WithFields(fields...)

	return NewWrapper(l)
}

func (wrapper *wrapperImpl) Fatal(a ...interface{}) {
	wrapper.l.Log(LevelFatal, a...)
	panic(a)
}

func (wrapper *wrapperImpl) Fatalf(format string, a ...interface{}) {
	wrapper.l.Logf(LevelFatal, format, a...)
	panic(a)
}

func (wrapper *wrapperImpl) Error(a ...interface{}) {
	wrapper.l.Log(LevelError, a...)
}

func (wrapper *wrapperImpl) Errorf(format string, a ...interface{}) {
	wrapper.l.Logf(LevelError, format, a...)
}

func (wrapper *wrapperImpl) Warn(a ...interface{}) {
	wrapper.l.Log(LevelWarn, a...)
}

func (wrapper *wrapperImpl) Warnf(format string, a ...interface{}) {
	wrapper.l.Logf(LevelWarn, format, a...)
}

func (wrapper *wrapperImpl) Info(a ...interface{}) {
	wrapper.l.Log(LevelInfo, a...)
}

func (wrapper *wrapperImpl) Infof(format string, a ...interface{}) {
	wrapper.l.Logf(LevelInfo, format, a...)
}

func (wrapper *wrapperImpl) Debug(a ...interface{}) {
	wrapper.l.Log(LevelDebug, a...)
}

func (wrapper *wrapperImpl) Debugf(format string, a ...interface{}) {
	wrapper.l.Logf(LevelDebug, format, a...)
}

func (wrapper *wrapperImpl) GetWrapperWithContext() WrapperWithContext {
	lc, _ := wrapper.l.(LoggerWithContext)

	return &contextWrapperImpl{
		l:  wrapper.l,
		lc: lc,
	}
}

type contextWrapperImpl struct {
	l  Logger
	lc LoggerWithContext
}

func (wrapper *contextWrapperImpl) GetLogger() Logger {
	return wrapper.l
}

func (wrapper *contextWrapperImpl) WithFields(fields ...Field) WrapperWithContext {
	l := wrapper.l.WithFields(fields...)

	return NewWrapper(l).GetWrapperWithContext()
}

func (wrapper *contextWrapperImpl) logWithContext(ctx context.Context, level Level, a ...interface{}) {
	if wrapper.lc != nil {
		wrapper.lc.LogWithContext(ctx, level, a...)
	} else {
		wrapper.l.Log(level, a...)
	}
}

func (wrapper *contextWrapperImpl) logWithContextf(ctx context.Context, level Level, format string, a ...interface{}) {
	if wrapper.lc != nil {
		wrapper.lc.LogWithContextf(ctx, level, format, a...)
	} else {
		wrapper.l.Logf(level, format, a...)
	}
}

func (wrapper *contextWrapperImpl) Fatal(ctx context.Context, a ...interface{}) {
	wrapper.logWithContext(ctx, LevelFatal, a...)
	panic(a)
}

func (wrapper *contextWrapperImpl) Fatalf(ctx context.Context, format string, a ...interface{}) {
	wrapper.logWithContextf(ctx, LevelFatal, format, a...)
	panic(a)
}

func (wrapper *contextWrapperImpl) Error(ctx context.Context, a ...interface{}) {
	wrapper.logWithContext(ctx, LevelError, a...)
}

func (wrapper *contextWrapperImpl) Errorf(ctx context.Context, format string, a ...interface{}) {
	wrapper.logWithContextf(ctx, LevelError, format, a...)
}

func (wrapper *contextWrapperImpl) Warn(ctx context.Context, a ...interface{}) {
	wrapper.logWithContext(ctx, LevelWarn, a...)
}

func (wrapper *contextWrapperImpl) Warnf(ctx context.Context, format string, a ...interface{}) {
	wrapper.logWithContextf(ctx, LevelWarn, format, a...)
}

func (wrapper *contextWrapperImpl) Info(ctx context.Context, a ...interface{}) {
	wrapper.logWithContext(ctx, LevelInfo, a...)
}

func (wrapper *contextWrapperImpl) Infof(ctx context.Context, format string, a ...interface{}) {
	wrapper.logWithContextf(ctx, LevelInfo, format, a...)
}

func (wrapper *contextWrapperImpl) Debug(ctx context.Context, a ...interface{}) {
	wrapper.logWithContext(ctx, LevelDebug, a...)
}

func (wrapper *contextWrapperImpl) Debugf(ctx context.Context, format string, a ...interface{}) {
	wrapper.logWithContextf(ctx, LevelDebug, format, a...)
}
