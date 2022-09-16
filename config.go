package main

import "github.com/gdamore/tcell/v2"

var globalKeyMap = KeyMap{
	KeyPress{tcell.KeyRune, 'q'}:       "quit",
	KeyPress{tcell.KeyRune, 'c'}:       "clear",
	KeyPress{tcell.KeyRune, '?'}:       "hint",
	KeyPress{tcell.KeyRune, 'j'}:       "next",
	KeyPress{tcell.KeyDown, ' '}:       "next",
	KeyPress{tcell.KeyRune, 'k'}:       "prev",
	KeyPress{tcell.KeyUp, ' '}:         "prev",
	KeyPress{tcell.KeyEnter, ' '}:      "view",
	KeyPress{tcell.KeyBackspace, ' '}:  "back",
	KeyPress{tcell.KeyBackspace2, ' '}: "back",
	KeyPress{tcell.KeyPgDn, ' '}:       "nextPage",
	KeyPress{tcell.KeyRight, ' '}:      "nextPage",
	KeyPress{tcell.KeyRune, 'l'}:       "nextPage",
	KeyPress{tcell.KeyPgUp, ' '}:       "prevPage",
	KeyPress{tcell.KeyLeft, ' '}:       "prevPage",
	KeyPress{tcell.KeyRune, 'h'}:       "prevPage",
}
