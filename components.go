package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ComponentWidget[T Section | Command] struct {
	*tview.Flex
	checkbox *tview.Checkbox
	data     T
	height   int
}

func (w *ComponentWidget[T]) focus(state bool) {
	w.checkbox.SetChecked(state)
}

func (c Command) widget() *ComponentWidget[Command] {
	var (
		flex     = tview.NewFlex().SetDirection(tview.FlexRow)
		checkbox = tview.NewCheckbox().
				SetFieldBackgroundColor(tcell.ColorBlack).
				SetFieldTextColor(tcell.ColorWhite)
		descText = tview.NewTextView().
				SetText("[yellow]" + c.Description).
				SetTextAlign(tview.AlignLeft).SetDynamicColors(true)
		examText = tview.NewTextView().
				SetText(c.getExample()).
				SetTextAlign(tview.AlignLeft)
		rows = make([]*tview.Flex, 2)
	)

	for i := range rows {
		rows[i] = tview.NewFlex().SetDirection(tview.FlexColumn)
	}

	rows[0].AddItem(checkbox, 2, 1, false)
	rows[0].AddItem(descText, 0, 2, false)
	rows[1].AddItem(nil, 4, 1, false)
	rows[1].AddItem(examText, 0, 2, false)

	for _, r := range rows {
		flex.AddItem(r, 0, 1, false)
	}

	return &ComponentWidget[Command]{
		flex,
		checkbox,
		c,
		flex.GetItemCount() + 1,
	}
}

func (s Section) widget() *ComponentWidget[Section] {
	var (
		flex     = tview.NewFlex().SetDirection(tview.FlexColumn)
		checkbox = tview.NewCheckbox().
				SetFieldBackgroundColor(tcell.ColorBlack).
				SetFieldTextColor(tcell.ColorWhite)
		name = tview.NewTextView().
			SetText(s.Name).
			SetTextAlign(tview.AlignLeft)
	)

	flex.AddItem(checkbox, 2, 1, false)
	flex.AddItem(name, 0, 2, false)

	return &ComponentWidget[Section]{
		flex,
		checkbox,
		s,
		flex.GetItemCount() + 1,
	}
}
