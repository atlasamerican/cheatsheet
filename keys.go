package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
)

type KeyPress struct {
	key tcell.Key
	ch  rune
}

type KeyMap map[KeyPress]string

func generateHint(cmd string, kps ...KeyPress) string {
	keys := make([]string, 0, len(kps))
	for _, kp := range kps {
		if kp.key == tcell.KeyRune {
			keys = append(keys, string(kp.ch))
		} else {
			keys = append(keys, tcell.KeyNames[kp.key])
		}
	}
	sort.Strings(keys)

	return fmt.Sprintf("[black:white]%s[-:-] %s", strings.Join(keys, ","), cmd)
}

func (km KeyMap) generateHintKey() string {
	for kp, cmd := range km {
		if cmd == "hint" {
			return generateHint(cmd, kp)
		}
	}
	return ""
}

func (km KeyMap) generateHints() []string {
	order := make([]string, 0)
	cmds := make(map[string][]KeyPress)

	for kp, cmd := range km {
		if kps, ok := cmds[cmd]; ok {
			cmds[cmd] = append(kps, kp)
		} else {
			kps = make([]KeyPress, 1)
			kps[0] = kp
			cmds[cmd] = kps
			order = append(order, cmd)
		}
	}

	sort.Strings(order)

	hs := make([]string, 0, len(cmds))
	for _, cmd := range order {
		hs = append(hs, generateHint(cmd, cmds[cmd]...))
	}
	return hs
}

func (km KeyMap) event2command(ev *tcell.EventKey) (string, bool) {
	key := ev.Key()
	kp := KeyPress{key, ' '}
	if key == tcell.KeyRune {
		kp.ch = ev.Rune()
	}
	if cmd, ok := km[kp]; ok {
		return cmd, true
	}
	return "", false
}
