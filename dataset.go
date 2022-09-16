package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

const osReleaseFile = "/etc/os-release"

type Filter struct {
	Os      string
	Distros []string
}

type Dataset struct {
	sections map[string]Section
	tldr     *Archive[TldrPage]
	filters  map[string]Filter
}

type Command struct {
	Name        string
	Description string
	Example     string
	Section     string
	Filters     []string
}

type Section struct {
	Name     string
	Commands []Command
}

type Datafile struct {
	Section  *string
	Commands []Command
}

func readDistroId() (string, bool) {
	file, err := os.Open(osReleaseFile)
	if err != nil {
		return "", false
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var id string
	for scanner.Scan() {
		line := scanner.Text()
		if _, err := fmt.Sscanf(line, "ID=%s", &id); err == nil {
			break
		}
	}
	return id, id != ""
}

func (f Filter) check() bool {
	if f.Os != "" && f.Os != runtime.GOOS {
		return false
	}
	if len(f.Distros) == 0 {
		return true
	}
	distro, ok := readDistroId()
	if !ok {
		logger.Log("[system] Failed to get distro ID")
		return false
	}
	for _, d := range f.Distros {
		if d == distro {
			return true
		}
	}
	return false
}

func (c Command) getExample() string {
	if c.Example != "" {
		return c.Example
	}
	return c.Name
}

func (c Command) isValid(section bool) bool {
	if section && c.Section == "" {
		section = false
	} else {
		section = true
	}
	return section && c.Name != "" && c.Description != ""
}

func (c Command) checkFilters(fs map[string]Filter) bool {
	for _, filter := range c.Filters {
		if f, ok := fs[filter]; ok && !f.check() {
			return false
		}
	}
	return true
}

func readCommandsBuf(buf []byte, cmds []Command) []Command {
	data := Datafile{
		Commands: make([]Command, 0),
	}
	var cmd Command
	var cmdErr error

	dataErr := yaml.Unmarshal(buf, &data)
	if dataErr != nil {
		if cmdErr = yaml.Unmarshal(buf, &cmd); cmdErr != nil {
			log.Fatal(cmdErr)
		}
	}

	if len(data.Commands) > 0 {
		for _, c := range data.Commands {
			if !c.isValid(false) {
				continue
			}
			if data.Section != nil {
				c.Section = *data.Section
			}
			cmds = append(cmds, c)
		}
	} else if cmd.isValid(true) {
		cmds = append(cmds, cmd)
	}

	return cmds
}

func readFiltersBuf(buf []byte, fs map[string]Filter) map[string]Filter {
	err := yaml.Unmarshal(buf, fs)
	if err != nil {
		log.Fatal(err)
	}
	return fs
}

func newDataset(dataPath string, archivePath string) *Dataset {
	cmds, filters := newDataArchive(archivePath).readData()

	ds := &Dataset{
		sections: make(map[string]Section),
		tldr:     newTldrArchive(archivePath),
		filters:  filters,
	}

	files, err := ioutil.ReadDir(dataPath)
	if err != nil {
		debugLogger.Log("[data] Local data directory not found; skipping")
	} else {
		for _, file := range files {
			if filepath.Ext(file.Name()) != ".yml" || file.IsDir() {
				continue
			}

			f, err := ioutil.ReadFile(path.Join(dataPath, file.Name()))
			if err != nil {
				log.Fatal(err)
			}

			readCommandsBuf(f, cmds)
		}
	}

	for _, cmd := range cmds {
		if !cmd.checkFilters(ds.filters) {
			continue
		}
		s, ok := ds.sections[cmd.Section]
		if !ok {
			s = Section{
				Name:     cmd.Section,
				Commands: make([]Command, 0),
			}
		}
		s.Commands = append(s.Commands, cmd)
		ds.sections[cmd.Section] = s
	}

	return ds
}

func (ds *Dataset) getPage(c Command) (*TldrPage, error) {
	return ds.tldr.getPage(c.Name)
}
