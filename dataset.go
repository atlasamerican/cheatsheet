package main

import (
	"io/ioutil"
	"log"
	"path"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Dataset struct {
	sections map[string]Section
	tldr     *Archive[TldrPage]
}

type Command struct {
	Name        string
	Description string
	Example     string
	Section     string
}

type Section struct {
	Name     string
	Commands []Command
}

type Datafile struct {
	Section  *string
	Commands []Command `yaml:",flow"`
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

func readDataBuf(buf []byte, cmds *[]Command) {
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
			*cmds = append(*cmds, c)
		}
	} else if cmd.isValid(true) {
		*cmds = append(*cmds, cmd)
	}
}

func newDataset(dataPath string, archivePath string) *Dataset {
	ds := &Dataset{
		sections: make(map[string]Section),
		tldr:     newTldrArchive(archivePath),
	}
	cmds := newDataArchive(archivePath).getCommands()

	files, err := ioutil.ReadDir(dataPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".yml" || file.IsDir() {
			continue
		}

		f, err := ioutil.ReadFile(path.Join(dataPath, file.Name()))
		if err != nil {
			log.Fatal(err)
		}

		readDataBuf(f, &cmds)
	}

	for _, cmd := range cmds {
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
