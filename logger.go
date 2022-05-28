package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Logger struct {
	file       *os.File
	fileLogger *log.Logger
	queue      chan string
}

func (lg *Logger) Log(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	lg.fileLogger.Println(msg)
	go func() {
		lg.queue <- msg
	}()
}

func (lg *Logger) Close() {
	lg.file.Close()
}

func newLogger(path string) *Logger {
	if err := os.MkdirAll(path, 0700); err != nil {
		log.Fatal(err)
	}
	file, err := os.Create(filepath.Join(path, "cheatsheet.log"))
	if err != nil {
		log.Fatal(err)
	}
	return &Logger{
		file:       file,
		fileLogger: log.New(file, "", 0),
		queue:      make(chan string),
	}
}
