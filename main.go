package main

import (
	"cheatsheet/lg"
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/ProtonMail/go-appdir"
	"github.com/invopop/jsonschema"
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
	schema := flag.Bool("schema", false, "Print JSON schema")
	flag.Parse()

	if *schema {
		s := jsonschema.Reflect(&Section{})
		s.Definitions["Command"].Title = "Command definition"
		s.Definitions["Section"].Title = "Section definition"

		b, err := json.MarshalIndent(s, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(b))
		return
	}

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
