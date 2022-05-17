package main

type KeyMap map[rune]string

var globalKeyMap = KeyMap{
	'j': "next",
	'k': "prev",
}
