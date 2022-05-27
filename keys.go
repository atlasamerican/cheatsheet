package main

import "github.com/gdamore/tcell/v2"

type KeyPress struct {
	key tcell.Key
	ch  rune
}

type KeyMap map[KeyPress]string

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
