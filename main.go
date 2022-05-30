package main

import (
	"cheatsheet/lg"
	"flag"
	"fmt"
	"log"

	"github.com/ProtonMail/go-appdir"
)

var debugLogger = lg.Logger

var logger *Logger

type Config struct {
	sectionsPerPage int
	keyMap          KeyMap
	appDirs         appdir.Dirs
}

func main() {
	defer debugLogger.Close()

	config := Config{
		sectionsPerPage: 3,
		keyMap:          globalKeyMap,
		appDirs:         appdir.New("cheatsheet"),
	}

	update := flag.Bool("update", false, "Update archive")
	flag.Parse()

	if *update {
		debugLogger.Overload(func(f string, v ...interface{}) {
			fmt.Printf(f+"\n", v...)
		})
		newTldrArchive(config.appDirs.UserData()).waitForUpdate()
		return
	}

	logger = newLogger(config.appDirs.UserLogs())
	defer logger.Close()

	ui := newUI(config)

	if err := ui.app.Run(); err != nil {
		log.Fatal(err)
	}
}
