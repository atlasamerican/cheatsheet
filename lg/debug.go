//go:build debug

package lg

import (
	"log"
	"os"
)

type DebugLogger struct {
	file       *os.File
	fileLogger *log.Logger
}

func (lg *DebugLogger) Log(format string, v ...interface{}) {
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
		file:       file,
		fileLogger: log.New(file, "", 0),
	}
}
