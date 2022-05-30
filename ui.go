package main

import (
	"fmt"
	"strings"

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
	if !state {
		return
	}
	dp := w.root.(*DataPages)
	dp.SwitchToPage(w.name)
	dp.ui.footer.updatePage(dp)
}

func (w *PageWidget) handleCommand(cmd string) bool {
	if cmd == "nextSection" {
		cmd = "next"
	}
	if cmd == "prevSection" {
		cmd = "prev"
	}
	return w.BaseWidget.handleCommand(cmd)
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

func (w *DataPages) handleCommand(cmd string) bool {
	if cmd == "nextPage" || cmd == "nextSection" {
		cmd = "next"
	}
	if cmd == "prevPage" || cmd == "prevSection" {
		cmd = "prev"
	}
	return w.BaseWidget.handleCommand(cmd)
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
	view      *tview.TextView
	page      *TldrPage
	firstView bool
}

func (v *PageView) setPage(page *TldrPage) {
	v.page = page
	v.firstView = true
}

func (v *PageView) render(width int) bool {
	if v.page == nil {
		return false
	}
	text := string(markdown.Render(v.page.content, width, 0))
	v.view.SetText(tview.TranslateANSI(text))
	return true
}

func newPageView() *PageView {
	view := tview.NewTextView().SetWrap(true).SetDynamicColors(true)

	grid := tview.NewGrid().SetColumns(-2, -1, -1, -1, -1, -1, -1, -1, -2)
	grid.AddItem(view, 0, 1, 1, 7, 0, 0, true)
	grid.AddItem(view, 0, 3, 1, 4, 0, 125, true)

	frame := tview.NewFrame(grid).SetBorders(2, 0, 0, 0, 0, 0)

	v := &PageView{frame, view, nil, false}

	view.SetDrawFunc(func(_ tcell.Screen, x, y, w, h int) (int, int, int, int) {
		if v.render(w) && v.firstView {
			v.firstView = false
			view.ScrollToBeginning()
		}
		return view.GetInnerRect()
	})

	return v
}

type Footer struct {
	*tview.Flex
	page *tview.TextView
	info *tview.TextView
}

func (f *Footer) updatePage(dp *DataPages) {
	name, _ := dp.GetFrontPage()
	f.page.Clear()
	f.page.SetText(fmt.Sprintf("%s/%d", name, dp.GetPageCount()))
}

type UI struct {
	app       *tview.Application
	mainFlex  *tview.Flex
	mainPages *tview.Pages
	dataPages *DataPages
	pageView  *PageView
	dataset   *Dataset
	footer    *Footer
	viewing   bool
	hinting   bool
	errors    int
	keyMap    KeyMap
	hintKey   string
	hints     string
}

func (ui *UI) viewPage(w *CommandWidget) {
	page, err := ui.dataset.getPage(w.data)
	if err != nil {
		logger.Log("[error] %v", err)
		return
	}
	if page != nil {
		debugLogger.Log("[widget] view %s", page.name)
		ui.pageView.setPage(page)
		ui.mainPages.SwitchToPage("View")
		ui.viewing = true
	}
}

func (ui *UI) unviewPage() {
	ui.mainPages.SwitchToPage("Sections")
	ui.viewing = false
}

func (ui *UI) resetFooter(hintKey bool) {
	ui.footer.info.Clear()
	if hintKey {
		ui.footer.info.SetText(ui.hintKey)
	}
	ui.hinting = false
	ui.errors = 0
}

func (ui *UI) toggleHints() {
	if ui.hinting {
		ui.resetFooter(true)
		return
	}

	ui.resetFooter(false)

	ui.hinting = true
	ui.footer.info.SetText(ui.hints)
}

func (ui *UI) showError(msg string) {
	ui.hinting = false

	var text string
	if ui.errors > 0 {
		text = ui.footer.info.GetText(true) + "\n"
	}
	ui.footer.info.Clear()

	ui.footer.info.SetText("[red]" + text + tview.Escape(msg))
	ui.errors++
}

func (ui *UI) handleKey(ev *tcell.EventKey) *tcell.EventKey {
	cmd, ok := ui.keyMap.event2command(ev)
	if !ok {
		return ev
	}

	switch cmd {
	case "hint":
		ui.toggleHints()
		return ev
	case "clear":
		ui.resetFooter(true)
		return ev
	case "quit":
		ui.app.Stop()
		return ev
	}

	if ui.viewing {
		switch cmd {
		case "view":
			fallthrough
		case "unview":
			ui.unviewPage()
		}
		return ev
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

	return ev
}

func newFooter() *Footer {
	f := &Footer{
		tview.NewFlex().SetDirection(tview.FlexRow),
		// page
		tview.NewTextView().SetTextAlign(tview.AlignRight),
		// info
		tview.NewTextView().
			SetDynamicColors(true).SetWrap(true).SetWordWrap(true).
			SetTextAlign(tview.AlignLeft),
	}
	f.AddItem(f.page, 1, 1, false)
	f.AddItem(f.info, 0, 1, false)

	return f
}

func newUI(config Config) *UI {
	ui := &UI{
		app:       tview.NewApplication(),
		mainFlex:  tview.NewFlex().SetDirection(tview.FlexRow),
		mainPages: tview.NewPages(),
		pageView:  newPageView(),
		dataset:   newDataset(config.appDirs.UserData()),
		footer:    newFooter(),
		keyMap:    config.keyMap,
	}

	ui.hintKey = ui.keyMap.generateHintKey()
	ui.hints = strings.Join(ui.keyMap.generateHints(), " ")

	ui.mainPages.SetInputCapture(ui.handleKey)

	ui.mainFlex.AddItem(ui.mainPages, 0, 6, true)
	ui.mainFlex.AddItem(ui.footer, 0, 1, false)

	ui.dataPages = newDataPages(ui, config.sectionsPerPage)
	ui.mainPages.AddPage("Sections", ui.dataPages, true, true)
	ui.mainPages.AddPage("View", ui.pageView, true, false)

	ui.resetFooter(true)
	ui.app.SetRoot(ui.mainFlex, true).SetFocus(ui.mainPages)

	go func() {
		for msg := range logger.queue {
			ui.app.QueueUpdateDraw(func() {
				ui.showError(msg)
			})
		}
	}()

	return ui
}
