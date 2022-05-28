package main

import (
	"cheatsheet/lg"
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
	config := Config{
		sectionsPerPage: 3,
		keyMap:          globalKeyMap,
		appDirs:         appdir.New("cheatsheet"),
	}

	logger = newLogger(config.appDirs.UserLogs())
	defer logger.Close()

	ui := newUI(config)

	if err := ui.app.Run(); err != nil {
		log.Fatal(err)
	}
}
