package main

import (
	"cheatsheet/lg"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/ProtonMail/go-appdir"
)

var debugLogger = lg.Logger

var logger *Logger

type Config struct {
	sectionsPerPage int
	keyMap          KeyMap
	appDirs         appdir.Dirs
	dataPath        string
}

func validate(path string) (bool, string) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return false, fmt.Sprintf("Failed to read file: %s\n", path)
	}

	cmds := make([]Command, 0)
	fs := make(map[string]Filter)

	if filepath.Base(path) == filtersFile {
		_, err = readFiltersBuf(f, fs)
	} else {
		_, err = readCommandsBuf(f, cmds)
	}

	if err != nil {
		return false, fmt.Sprintf("FAIL %s: %v", path, err)
	}
	return true, fmt.Sprintf("PASS %s", path)
}

func main() {
	defer debugLogger.Close()

	update := flag.Bool("update", false, "Update archive")
	val := flag.String("validate", "", "File to validate")
	path := flag.String("path", "", "Path to local data")
	flag.Parse()

	if *val != "" {
		ok, msg := validate(*val)
		fmt.Println(msg)
		if !ok {
			os.Exit(1)
		}
		return
	}

	config := Config{
		sectionsPerPage: 3,
		keyMap:          globalKeyMap,
		appDirs:         appdir.New("cheatsheet"),
		dataPath:        *path,
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
