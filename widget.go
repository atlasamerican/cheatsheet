package main

import "container/list"

type Widget interface {
	init()
	handleCommand(string) bool
	handleFocus(bool)
	setFocus(*list.Element) bool
	getElement() *list.Element
}

type BaseWidget struct {
	name    string
	kind    string
	widgets *list.List
	element *list.Element
	focus   Widget
	root    Widget
}

func (w *BaseWidget) getElement() *list.Element {
	return w.element
}

func (w *BaseWidget) handleFocus(state bool) {
	logger.Log("[widget] focus %t: %s", state, w.name)
	if w.focus != nil {
		w.focus.handleFocus(state)
	}
}

func (w *BaseWidget) setFocus(el *list.Element) bool {
	if el == nil {
		return false
	}
	f := el.Value.(Widget)
	if w.focus != nil && w.focus != f {
		w.focus.handleFocus(false)
	}
	w.focus = f
	f.handleFocus(true)
	return true
}

func (w *BaseWidget) handleCommand(cmd string) bool {
	switch cmd {
	case "next":
		if w.focus == nil {
			break
		}
		if el := w.focus.getElement(); el != nil {
			return w.setFocus(el.Next())
		}
	case "prev":
		if w.focus == nil {
			break
		}
		if el := w.focus.getElement(); el != nil {
			return w.setFocus(el.Prev())
		}
	}
	return false
}

func (w *BaseWidget) init() {
	logger.Log("[widget] init %s: %s", w.kind, w.name)
	front := w.widgets.Front()
	for e := front; e != nil; e = e.Next() {
		f := e.Value.(Widget)
		if e == front {
			w.focus = f
		}
		f.init()
	}
}

func newBaseWidget(root Widget, name string, kind string) *BaseWidget {
	return &BaseWidget{
		name:    name,
		kind:    kind,
		widgets: list.New(),
		root:    root,
	}
}
