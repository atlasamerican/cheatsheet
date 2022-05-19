package main

import (
	"fmt"
	"log"

	markdown "github.com/MichaelMure/go-term-markdown"
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

func (w *CommandWidget) handleCommand(cmd string) bool {
	if cmd == "view" {
		w.root.(*DataPages).ui.viewPage(w)
		return true
	}
	return false
}

func (c Command) widget(root Widget) *CommandWidget {
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
		newBaseWidget(root, c.GetExample(), "command"),
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

func (section Section) widget(root Widget) *SectionWidget {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	frame := tview.NewFrame(flex).
		AddText(section.Name, true, tview.AlignCenter, tcell.ColorWhite)

	size := 0
	base := newBaseWidget(root, section.Name, "section")

	for _, v := range section.Commands {
		w := v.widget(root)
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
}

func (w *PageWidget) handleFocus(state bool) {
	w.BaseWidget.handleFocus(state)
	if state {
		w.root.(*DataPages).SwitchToPage(w.name)
	}
}

func newPageWidget(pages *DataPages, num int) *PageWidget {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	grid := tview.NewGrid().SetColumns(-4, -1, -1, -1, -2, -1, -1, -1, -4)
	grid.AddItem(flex, 0, 1, 1, 7, 0, 0, true)
	grid.AddItem(flex, 0, 2, 1, 5, 0, 125, true)
	grid.AddItem(flex, 0, 3, 1, 3, 0, 175, true)

	return &PageWidget{
		newBaseWidget(pages, fmt.Sprintf("Page %d", num), "page"),
		grid,
		flex,
	}
}

type DataPages struct {
	*BaseWidget
	*tview.Pages
	ui *UI
}

func newDataPages(ui *UI, perPage int) *DataPages {
	var (
		dp = &DataPages{
			newBaseWidget(nil, "main", "dataPages"),
			tview.NewPages(),
			ui,
		}
		pageNum = 1
		page    = newPageWidget(dp, pageNum)
	)

	for i, section := range ui.dataset.sections {
		if i > 0 && i%perPage == 0 {
			page.element = dp.widgets.PushBack(page)
			dp.AddPage(page.name, page, true, true)
			pageNum++
			page = newPageWidget(dp, pageNum)
		}

		w := section.widget(dp)
		w.element = page.widgets.PushBack(w)
		page.content.AddItem(w, w.size+3, 1, false)
	}

	page.element = dp.widgets.PushBack(page)
	dp.AddPage(page.name, page, true, true)

	dp.init()
	dp.setFocus(dp.widgets.Front())

	return dp
}

type PageView struct {
	*tview.Frame
	view *tview.TextView
}

func (v *PageView) render(page *TldrPage, width, leftPad int) {
	v.Clear()
	v.AddText(page.name, true, tview.AlignCenter, tcell.ColorWhite)
	text := string(markdown.Render(page.content, width, leftPad))
	v.view.SetText(tview.TranslateANSI(text))
}

func newPageView() *PageView {
	view := tview.NewTextView().SetWrap(true).SetDynamicColors(true)

	grid := tview.NewGrid().SetColumns(-2, -1, -1, -1, -1, -1, -1, -1, -2)
	grid.AddItem(view, 0, 1, 1, 7, 0, 0, true)
	grid.AddItem(view, 0, 3, 1, 7, 0, 125, true)

	return &PageView{tview.NewFrame(grid), view}
}

type UI struct {
	app       *tview.Application
	mainPages *tview.Pages
	dataPages *DataPages
	pageView  *PageView
	dataset   *Dataset
	viewing   bool
	keyMap    KeyMap
}

func (ui *UI) viewPage(w *CommandWidget) {
	page, err := ui.dataset.getPage(w.data)
	if err != nil {
		// TODO: Log this to a file in UserLogs
		log.Println(err)
		return
	}
	if page != nil {
		logger.Log("[widget] view %s", page.name)
		ui.pageView.render(page, 80, 0)
		ui.mainPages.SwitchToPage("View")
		ui.viewing = true
	}
}

func (ui *UI) unviewPage() {
	ui.mainPages.SwitchToPage("Sections")
	ui.viewing = false
}

func (ui *UI) handleKey(key rune) {
	cmd, ok := ui.keyMap[key]
	if !ok {
		return
	}

	if ui.viewing {
		switch cmd {
		case "view":
			ui.unviewPage()
		}
		return
	}

	var (
		page    = ui.dataPages.focus.(*PageWidget)
		section = page.focus.(*SectionWidget)
		command = section.focus.(*CommandWidget)
	)

	if !command.handleCommand(cmd) {
		if !section.handleCommand(cmd) {
			if !page.handleCommand(cmd) {
				ui.dataPages.handleCommand(cmd)
			}
		}
	}
}

func newUI(config Config) *UI {
	ui := &UI{
		app:       tview.NewApplication(),
		mainPages: tview.NewPages(),
		pageView:  newPageView(),
		dataset:   newDataset(config.appDirs.UserData()),
		keyMap:    config.keyMap,
	}

	ui.mainPages.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {
		key := ev.Rune()
		ui.handleKey(key)
		return ev
	})

	ui.dataPages = newDataPages(ui, config.sectionsPerPage)
	ui.mainPages.AddPage("Sections", ui.dataPages, true, true)
	ui.mainPages.AddPage("View", ui.pageView, true, false)

	ui.app.SetRoot(ui.mainPages, true).SetFocus(ui.mainPages)

	return ui
}
