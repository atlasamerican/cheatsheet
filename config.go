package main

import "github.com/gdamore/tcell/v2"

var globalKeyMap = KeyMap{
	KeyPress{tcell.KeyRune, 'q'}:  "quit",
	KeyPress{tcell.KeyRune, 'j'}:  "next",
	KeyPress{tcell.KeyDown, ' '}:  "next",
	KeyPress{tcell.KeyRune, 'k'}:  "prev",
	KeyPress{tcell.KeyUp, ' '}:    "prev",
	KeyPress{tcell.KeyEnter, ' '}: "view",
	KeyPress{tcell.KeyRune, 'J'}:  "nextSection",
	KeyPress{tcell.KeyRune, 'n'}:  "nextSection",
	KeyPress{tcell.KeyRune, 'K'}:  "prevSection",
	KeyPress{tcell.KeyRune, 'p'}:  "prevSection",
	KeyPress{tcell.KeyPgDn, ' '}:  "nextPage",
	KeyPress{tcell.KeyRight, ' '}: "nextPage",
	KeyPress{tcell.KeyRune, 'l'}:  "nextPage",
	KeyPress{tcell.KeyPgUp, ' '}:  "prevPage",
	KeyPress{tcell.KeyLeft, ' '}:  "prevPage",
	KeyPress{tcell.KeyRune, 'h'}:  "prevPage",
}
