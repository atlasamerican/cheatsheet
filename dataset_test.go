package main

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

const dataPath = "test"

func TestCommands(t *testing.T) {
	tests := map[string][]Command{
		"networking.yml": {
			{"ip", "Show IP address", "ip address show", "Networking", []string(nil)},
			{"ip", "Show routing table", "ip route show", "Networking", []string(nil)},
			{"ip", "Show IPv6 routes", "ip -6 route", "Networking", []string(nil)},
		},
		"single-command.yml": {
			{"ip", "Display the routing table", "ip -brief link", "Networking", []string(nil)},
		},
	}

	for file, commands := range tests {
		f, err := ioutil.ReadFile(path.Join(dataPath, file))
		if err != nil {
			t.Fatal(err)
		}
		cmds := make([]Command, 0)
		cmds, err = readCommandsBuf(f, cmds)
		assert.Nil(t, err)
		assert.Equal(t, len(commands), len(cmds))
		for i, c := range commands {
			assert.Equal(t, c, cmds[i])
		}
	}
}

func TestFilters(t *testing.T) {
	tests := map[string]Filter{
		"rpm": {
			Os:      "linux",
			Distros: []string{"fedora", "almalinux"},
		},
		"pacman": {
			Os:      "linux",
			Distros: []string{"arch"},
		},
	}

	for filter, attrs := range tests {
		f, err := ioutil.ReadFile(path.Join(dataPath, "__filters__.yml"))
		if err != nil {
			t.Fatal(err)
		}
		fs := make(map[string]Filter)
		fs, err = readFiltersBuf(f, fs)
		assert.Nil(t, err)
		assert.Equal(t, attrs, fs[filter])
	}
}
