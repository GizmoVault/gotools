package logx

import (
	"fmt"
	"strings"
	"time"
)

type Recorder interface {
	Log(level Level, a ...interface{})
	Logf(level Level, format string, a ...interface{})
}

type CommLogger struct {
	level      Level
	fields     []Field
	recordTime bool
	recorder   []Recorder
}

func NewCommLogger(recorder ...Recorder) Logger {
	return &CommLogger{
		level:      LevelInfo,
		fields:     nil,
		recordTime: true,
		recorder:   recorder,
	}
}

func NewCommLoggerEx(recordTime bool, recorder ...Recorder) Logger {
	return &CommLogger{
		level:      LevelInfo,
		fields:     nil,
		recordTime: recordTime,
		recorder:   recorder,
	}
}

func (l *CommLogger) AddRecorder(recorder ...Recorder) Logger {
	l.recorder = append(l.recorder, recorder...)

	return l
}

func (l *CommLogger) clone() *CommLogger {
	newLogger := &CommLogger{}
	newLogger.level = l.level
	newLogger.fields = append(newLogger.fields, l.fields...)
	newLogger.recorder = l.recorder
	newLogger.recordTime = l.recordTime

	return newLogger
}

func (l *CommLogger) SetLevel(level Level) {
	l.level = level
}

func (l *CommLogger) WithFields(fields ...Field) Logger {
	if len(fields) == 0 {
		return l
	}

	newLogger := l.clone()
	newLogger.fields = append(newLogger.fields, fields...)

	return newLogger
}

func (l *CommLogger) levelShouldRecord(level Level) bool {
	return level <= l.level
}

func (*CommLogger) mapFields(fields ...Field) string {
	ss := &strings.Builder{}
	ss.WriteString("[")

	var f bool

	for _, field := range fields {
		switch field.T {
		case FieldTypeString:
			if f {
				ss.WriteString(" ")
			}

			f = true

			ss.WriteString(field.K)
			ss.WriteString(":")

			if s, ok := field.V.(string); ok {
				ss.WriteString(s)
			}
		case FieldTypeError:
			if f {
				ss.WriteString(" ")
			}

			f = true

			ss.WriteString(field.K)
			ss.WriteString(":")

			if s, ok := field.V.(error); ok {
				ss.WriteString(s.Error())
			}
		default:
			if f {
				ss.WriteString(" ")
			}

			f = true

			ss.WriteString(field.K)
			ss.WriteString(":")
			_, _ = fmt.Fprintf(ss, "%v", field.V)
		}
	}

	ss.WriteString("]")

	return ss.String()
}

func (l *CommLogger) Log(level Level, a ...interface{}) {
	if !l.levelShouldRecord(level) {
		return
	}

	vs := make([]interface{}, 0, len(l.fields)+len(a)+1)

	if l.recordTime {
		vs = append(vs, time.Now().Format(time.RFC3339))
	}

	vs = append(vs, l.mapFields(l.fields...))
	vs = append(vs, a...)

	for _, recorder := range l.recorder {
		recorder.Log(level, vs...)
	}

	if level == LevelFatal {
		panic(a)
	}
}

func (l *CommLogger) Logf(level Level, format string, a ...interface{}) {
	if !l.levelShouldRecord(level) {
		return
	}

	tag := l.mapFields(l.fields...)

	if l.recordTime {
		tag = time.Now().Format(time.RFC3339) + " " + tag
	}

	for _, recorder := range l.recorder {
		recorder.Logf(level, tag+" "+format, a...)
	}

	if level == LevelFatal {
		panic(a)
	}
}
