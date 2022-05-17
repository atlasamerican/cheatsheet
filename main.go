package main

import (
	"cheatsheet/lg"
	"log"
)

var logger = lg.Logger

type Config struct {
	datasetPath     string
	sectionsPerPage int
	keyMap          KeyMap
}

func main() {
	config := Config{
		datasetPath:     "data",
		sectionsPerPage: 3,
		keyMap:          globalKeyMap,
	}

	ui := newUI(config)

	if err := ui.app.Run(); err != nil {
		log.Fatal(err)
	}
}
