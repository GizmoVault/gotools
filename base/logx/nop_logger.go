package logx

import "fmt"

type NopLogger struct{}

func (*NopLogger) SetLevel(_ Level) {

}

func (l *NopLogger) WithFields(_ ...Field) Logger {
	return l
}

func (*NopLogger) Log(level Level, a ...interface{}) {
	if level == LevelFatal {
		panic(fmt.Sprint(a...))
	}
}

func (*NopLogger) Logf(level Level, format string, a ...interface{}) {
	if level == LevelFatal {
		panic(fmt.Sprintf(format, a...))
	}
}
