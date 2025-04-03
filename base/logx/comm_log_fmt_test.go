package logx

import "testing"

func TestCommLogger_Log(_ *testing.T) {
	r := &ConsoleRecorder{}
	l := NewCommLogger(r)
	l.WithFields(FieldString("key1", "val1"), FieldString("key2", "val2")).Log(LevelInfo, "hello, world")
	l.WithFields(FieldString("key1", "val1"), FieldString("key2", "val2")).Logf(LevelInfo, "hello, world。 I'm %s", "zjz")
}
