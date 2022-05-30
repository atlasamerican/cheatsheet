package lg

type LogFunc func(string, ...interface{})

type FileLogger interface {
	Log(string, ...interface{})
	Overload(LogFunc)
	Close()
}

type BaseLogger struct {
	overloaded LogFunc
}

func (bl *BaseLogger) Overload(lf LogFunc) {
	bl.overloaded = lf
}

var Logger FileLogger
