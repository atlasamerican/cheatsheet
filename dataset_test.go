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
		cmds = readCommandsBuf(f, cmds)
		assert.Equal(t, len(commands), len(cmds))
		for i, c := range commands {
			assert.Equal(t, c, cmds[i])
		}
	}
}
