//go:build debug

package lg

import (
	"log"
	"os"
)

type DebugLogger struct {
	*BaseLogger
	file       *os.File
	fileLogger *log.Logger
}

func (lg *DebugLogger) Log(format string, v ...interface{}) {
	if lg.overloaded != nil {
		lg.overloaded(format, v...)
		return
	}
	lg.fileLogger.Printf(format, v...)
}

func (lg *DebugLogger) Close() {
	lg.file.Close()
}

func init() {
	file, err := os.Create("./cheatsheet.log")
	if err != nil {
		log.Fatal(err)
	}
	Logger = &DebugLogger{
		&BaseLogger{},
		file,
		log.New(file, "", 0),
	}
}
