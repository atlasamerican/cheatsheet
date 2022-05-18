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
}

type Command struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Example     string `json:"example"`
}

type Section struct {
	Name     string    `json:"section"`
	Commands []Command `json:"commands"`
}

func (c Command) GetExample() string {
	if c.Example != "" {
		return c.Example
	}
	return c.Name
}

func newDataset() *Dataset {
	ds := &Dataset{
		sections: make([]Section, 0),
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
