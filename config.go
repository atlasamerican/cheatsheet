package main

import "github.com/gdamore/tcell/v2"

var globalKeyMap = KeyMap{
	KeyPress{tcell.KeyRune, 'q', 0}:       "quit",
	KeyPress{tcell.KeyRune, 'c', 0}:       "clear",
	KeyPress{tcell.KeyRune, '?', 0}:       "hint",
	KeyPress{tcell.KeyRune, 'j', 0}:       "next",
	KeyPress{tcell.KeyDown, ' ', 1}:       "next",
	KeyPress{tcell.KeyRune, 'k', 0}:       "prev",
	KeyPress{tcell.KeyUp, ' ', 1}:         "prev",
	KeyPress{tcell.KeyEnter, ' ', 0}:      "view",
	KeyPress{tcell.KeyBackspace, ' ', 0}:  "back",
	KeyPress{tcell.KeyBackspace2, ' ', 1}: "back",
	KeyPress{tcell.KeyRune, 'l', 0}:       "nextPage",
	KeyPress{tcell.KeyRight, ' ', 1}:      "nextPage",
	KeyPress{tcell.KeyPgDn, ' ', 2}:       "nextPage",
	KeyPress{tcell.KeyRune, 'h', 0}:       "prevPage",
	KeyPress{tcell.KeyLeft, ' ', 1}:       "prevPage",
	KeyPress{tcell.KeyPgUp, ' ', 2}:       "prevPage",
}
