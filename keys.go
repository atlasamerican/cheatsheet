package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
)

type KeyPress struct {
	key    tcell.Key
	ch     rune
	weight int
}

type KeyMap map[KeyPress]string

func generateHint(cmd string, kps ...KeyPress) string {
	keys := make([]string, len(kps))
	for _, kp := range kps {
		var s string
		if kp.key == tcell.KeyRune {
			s = string(kp.ch)
		} else {
			s = tcell.KeyNames[kp.key]
		}
		keys[kp.weight] = s
	}

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

func (km KeyMap) generateHints(maxWeight int) []string {
	order := make([]string, 0)
	cmds := make(map[string][]KeyPress)

	for kp, cmd := range km {
		if kp.weight > maxWeight {
			continue
		}
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

func (km KeyMap) getCommand(kp KeyPress) (string, bool) {
	for k, cmd := range km {
		if k.key == kp.key && k.ch == kp.ch {
			return cmd, true
		}
	}
	return "", false
}

func (km KeyMap) event2command(ev *tcell.EventKey) (string, bool) {
	key := ev.Key()
	kp := KeyPress{key, ' ', -1}
	if key == tcell.KeyRune {
		kp.ch = ev.Rune()
	}
	if cmd, ok := km.getCommand(kp); ok {
		return cmd, true
	}
	return "", false
}
