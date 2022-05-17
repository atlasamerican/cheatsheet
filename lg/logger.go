package lg

type FileLogger interface {
	Log(string, ...interface{})
	Close()
}

var Logger FileLogger
