package main

import "github.com/gdamore/tcell/v2"

var globalKeyMap = KeyMap{
	KeyPress{tcell.KeyRune, 'j'}:  "next",
	KeyPress{tcell.KeyDown, ' '}:  "next",
	KeyPress{tcell.KeyRune, 'k'}:  "prev",
	KeyPress{tcell.KeyUp, ' '}:    "prev",
	KeyPress{tcell.KeyRune, 'l'}:  "view",
	KeyPress{tcell.KeyEnter, ' '}: "view",
}
