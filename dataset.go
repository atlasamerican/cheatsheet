package main

import (
	"embed"
	"encoding/json"
	"log"
	"path"
)

//go:embed data
var data embed.FS

type Dataset struct {
	sections []Section
	archive  *TldrArchive
}

type Command struct {
	Name        string `json:"name" jsonschema:"description=Command name--corresponds to tldr page"`
	Description string `json:"description" jsonschema:"description=Command description"`
	Example     string `json:"example,omitempty" jsonschema:"description=Command example--defaults to 'name' if omitted"`
}

type Section struct {
	Name     string    `json:"section" jsonschema:"description=Section name"`
	Commands []Command `json:"commands" jsonschema:"description=List of commands"`
}

func (c Command) GetExample() string {
	if c.Example != "" {
		return c.Example
	}
	return c.Name
}

func newDataset(archivePath string) *Dataset {
	ds := &Dataset{
		sections: make([]Section, 0),
		archive:  newTldrArchive(archivePath),
	}

	files, err := data.ReadDir("data")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		f, err := data.ReadFile(path.Join("data", file.Name()))
		if err != nil {
			log.Fatal(err)
		}

		var section Section
		err = json.Unmarshal(f, &section)
		if err != nil {
			log.Fatal(err)
		}
		ds.sections = append(ds.sections, section)
	}

	return ds
}

func (ds *Dataset) getPage(c Command) (*TldrPage, error) {
	return ds.archive.getPage(c.Name)
}
