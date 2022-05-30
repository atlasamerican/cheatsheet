//go:build !debug

package lg

type DummyLogger struct {
	*BaseLogger
}

func (lg DummyLogger) Log(format string, v ...interface{}) {
}

func (lg DummyLogger) Close() {
}

func init() {
	Logger = DummyLogger{&BaseLogger{}}
}
