package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type CommandWidget struct {
	*BaseWidget
	*tview.Flex
	checkbox *tview.Checkbox
	data     Command
}

func (w *CommandWidget) handleFocus(state bool) {
	w.checkbox.SetChecked(state)
	w.BaseWidget.handleFocus(state)
}

func (c Command) widget() *CommandWidget {
	var (
		flex     = tview.NewFlex().SetDirection(tview.FlexColumn)
		checkbox = tview.NewCheckbox().
				SetFieldBackgroundColor(tcell.ColorBlack).
				SetFieldTextColor(tcell.ColorWhite)
		descText = tview.NewTextView().
				SetText(c.Description).
				SetTextAlign(tview.AlignLeft)
		examText = tview.NewTextView().
				SetText(c.GetExample()).
				SetTextAlign(tview.AlignRight)
	)

	flex.AddItem(checkbox, 2, 1, false)
	flex.AddItem(descText, 0, 2, false)
	flex.AddItem(examText, 0, 1, false)

	return &CommandWidget{
		newBaseWidget(c.GetExample(), "command"),
		flex,
		checkbox,
		c,
	}
}

type SectionWidget struct {
	*BaseWidget
	*tview.Frame
	data Section
	size int
}

func (section Section) widget() *SectionWidget {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	frame := tview.NewFrame(flex).
		AddText(section.Name, true, tview.AlignCenter, tcell.ColorWhite)

	size := 0
	base := newBaseWidget(section.Name, "section")

	for _, v := range section.Commands {
		w := v.widget()
		w.element = base.widgets.PushBack(w)
		flex.AddItem(w, 1, 1, false)
		size += 1
	}

	return &SectionWidget{
		base,
		frame,
		section,
		size,
	}
}

type PageWidget struct {
	*BaseWidget
	*tview.Grid
	content *tview.Flex
	pages   *DataPages
}

func (w *PageWidget) handleFocus(state bool) {
	w.BaseWidget.handleFocus(state)
	if state {
		w.pages.SwitchToPage(w.name)
	}
}

func newPageWidget(pages *DataPages, num int) *PageWidget {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	grid := tview.NewGrid().SetColumns(-4, -1, -1, -1, -2, -1, -1, -1, -4)
	grid.AddItem(flex, 0, 1, 1, 7, 0, 0, true)
	grid.AddItem(flex, 0, 2, 1, 5, 0, 125, true)
	grid.AddItem(flex, 0, 3, 1, 3, 0, 175, true)

	return &PageWidget{
		newBaseWidget(fmt.Sprintf("Page %d", num), "page"),
		grid,
		flex,
		pages,
	}
}

type DataPages struct {
	*BaseWidget
	*tview.Pages
}

func newDataPages(ds *Dataset, perPage int) *DataPages {
	var (
		dp = &DataPages{
			newBaseWidget("main", "dataPages"),
			tview.NewPages(),
		}
		pageNum = 1
		page    = newPageWidget(dp, pageNum)
	)

	for i, section := range ds.sections {
		if i > 0 && i%perPage == 0 {
			page.element = dp.widgets.PushBack(page)
			dp.AddPage(page.name, page, true, true)
			pageNum++
			page = newPageWidget(dp, pageNum)
		}

		w := section.widget()
		w.element = page.widgets.PushBack(w)
		page.content.AddItem(w, w.size+3, 1, false)
	}

	page.element = dp.widgets.PushBack(page)
	dp.AddPage(page.name, page, true, true)

	dp.init()
	dp.setFocus(dp.widgets.Front())

	return dp
}

type UI struct {
	app       *tview.Application
	mainPages *tview.Pages
	dataPages *DataPages
	dataset   *Dataset
	modal     bool
	keyMap    KeyMap
}

func (ui *UI) handleKey(key rune) {
	cmd, ok := ui.keyMap[key]
	if !ok {
		return
	}

	if ui.modal {
		return
	}

	var (
		page    = ui.dataPages.focus.(*PageWidget)
		section = page.focus.(*SectionWidget)
	)

	if !section.handleCommand(cmd) {
		if !page.handleCommand(cmd) {
			ui.dataPages.handleCommand(cmd)
		}
	}
}

func newUI(config Config) *UI {
	ui := &UI{
		app:       tview.NewApplication(),
		mainPages: tview.NewPages(),
		dataset:   newDataset(config.datasetPath),
		keyMap:    config.keyMap,
	}

	ui.mainPages.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {
		key := ev.Rune()
		ui.handleKey(key)
		return ev
	})

	ui.dataPages = newDataPages(ui.dataset, config.sectionsPerPage)
	ui.mainPages.AddPage("Sections", ui.dataPages, true, true)

	ui.app.SetRoot(ui.mainPages, true).SetFocus(ui.mainPages)

	return ui
}
