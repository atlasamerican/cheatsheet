package main

import (
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ComponentWidget[T Command | Section] interface {
	tview.Primitive
	focus(bool)
	getData() T
	getHeight() int
	setWidth(int)
}

type BaseWidget struct {
	*tview.Flex
	checkbox *tview.Checkbox
	height   int
}

func (w *BaseWidget) focus(state bool) {
	w.checkbox.SetChecked(state)
}

func (w *BaseWidget) getHeight() int {
	return w.height
}

func (w *BaseWidget) setWidth(width int) {
}

type CommandWidget struct {
	*BaseWidget
	data        Command
	description *tview.TextView
	example     *tview.TextView
}

type SectionWidget struct {
	*BaseWidget
	data Section
}

func (w *CommandWidget) getData() Command {
	return w.data
}

func shortenText(str string, max int) string {
	lastSpaceIx := -1
	len := 0
	for i, r := range str {
		if unicode.IsSpace(r) {
			lastSpaceIx = i
		}
		len++
		if len >= max {
			if lastSpaceIx != -1 {
				return str[:lastSpaceIx] + "..."
			}
			return str[:max]
		}
	}
	return str
}

func (w *CommandWidget) setWidth(width int) {
	w.description.SetText("[yellow]" + shortenText(w.data.Description, width))
	w.example.SetText(shortenText(w.data.Example, width))
}

func (c Command) widget() *CommandWidget {
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

	return &CommandWidget{
		&BaseWidget{
			flex,
			checkbox,
			flex.GetItemCount() + 1,
		},
		c,
		descText,
		examText,
	}
}

func (w *SectionWidget) getData() Section {
	return w.data
}

func (s Section) widget() *SectionWidget {
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

	return &SectionWidget{
		&BaseWidget{
			flex,
			checkbox,
			flex.GetItemCount() + 1,
		},
		s,
	}
}
