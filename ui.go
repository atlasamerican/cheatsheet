package main

import (
	"fmt"
	"strings"

	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	PAGE_SECTIONS = "Sections"
	PAGE_COMMANDS = "Commands"
	PAGE_VIEWER   = "Viewer"
)

var breakpoints = [...]int{0, 115, 165}

type UI struct {
	app *tview.Application
	// root contains router and footer.
	root *tview.Flex
	// router displays the appropriate pager or viewer,
	// one of: section, command, tldr.
	// Pagers handle displaying their content.
	router        *tview.Pages
	page          string
	footer        *Footer
	sectionPager  *Pager[Section]
	commandPager  *Pager[Command]
	commandPagers map[string]*Pager[Command]
	maxPerColumn  int
	viewer        *Viewer
	dataset       *Dataset
	logLength     int
	keyMap        KeyMap
	hintKey       string
	hints         string
	fullHints     string
	fullHinting   bool
}

type Viewer struct {
	*tview.Frame
	view      *tview.TextView
	page      *TldrPage
	firstView bool
}

func (v *Viewer) setPage(page *TldrPage) {
	v.page = page
	v.firstView = true
}

func (v *Viewer) render(width int) bool {
	if v.page == nil {
		return false
	}
	text := string(markdown.Render(v.page.content, width, 0))
	v.view.SetText(tview.TranslateANSI(text))
	return true
}

func newViewer() *Viewer {
	view := tview.NewTextView().SetWrap(true).SetDynamicColors(true)
	frame := tview.NewFrame(view).SetBorders(2, 0, 0, 0, 0, 0)

	v := &Viewer{frame, view, nil, false}

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
	ui     *UI
	page   *tview.TextView
	hints  *tview.TextView
	log    *tview.TextView
	height int
}

func newFooter(ui *UI) *Footer {
	f := &Footer{
		tview.NewFlex().SetDirection(tview.FlexRow),
		ui,
		// page
		tview.NewTextView().SetTextAlign(tview.AlignRight),
		// hints
		tview.NewTextView().
			SetDynamicColors(true).SetWrap(true).SetWordWrap(true).
			SetTextAlign(tview.AlignLeft),
		// log
		tview.NewTextView().
			SetDynamicColors(true).SetWrap(true).SetWordWrap(true).
			SetTextAlign(tview.AlignLeft),
		3,
	}

	f.AddItem(f.page, 1, 1, false)
	f.AddItem(f.hints, 0, 1, false)
	f.AddItem(f.log, 0, 1, false)

	// Force redraw when footer page changes.
	f.page.SetChangedFunc(func() {
		ui.app.Draw()
	})

	return f
}

func (f *Footer) clearPage() {
	f.page.Clear()
}

func (f *Footer) updatePage(i int, more bool) {
	debugLogger.Log("[footer] update page: %d, %t", i, more)

	var s string
	if i > 0 {
		s += "< "
	} else {
		s += "  "
	}
	s += fmt.Sprintf("Page %d", i+1)
	if more {
		s += " >"
	} else {
		s += "  "
	}
	f.page.SetText(s)
}

func (ui *UI) clearLog() {
	ui.footer.log.Clear()
	ui.logLength = 0
}

func (ui *UI) resetHints(full bool) {
	ui.footer.hints.Clear()
	if full {
		ui.footer.hints.SetText(ui.fullHints)
	} else {
		ui.footer.hints.SetText(ui.hints)
	}
	ui.fullHinting = full
}

func (ui *UI) viewTldr(c Command) {
	page, err := ui.dataset.getPage(c)
	if err != nil {
		logger.Log("[error] %v", err)
		return
	}
	if page != nil {
		debugLogger.Log("[widget] view %s", page.name)
		ui.viewer.setPage(page)
		ui.router.SwitchToPage(PAGE_VIEWER)
	}
}

func (ui *UI) unviewTldr() {
	ui.router.SwitchToPage(PAGE_COMMANDS)
}

func (ui *UI) toggleHints() {
	ui.resetHints(!ui.fullHinting)
}

func (ui *UI) showMessage(msg string, color string) {
	var text string
	if ui.logLength > 0 {
		text = ui.footer.log.GetText(true) + "\n"
	}
	ui.footer.log.Clear()

	ui.footer.log.SetText(color + text + tview.Escape(msg))
	ui.logLength++
}

func (ui *UI) handleKey(ev *tcell.EventKey) *tcell.EventKey {
	cmd, ok := ui.keyMap.event2command(ev)
	if !ok {
		return ev
	}

	debugLogger.Log("[key] command: %s", cmd)

	switch cmd {
	case "hint":
		ui.toggleHints()
		return ev
	case "clear":
		ui.resetHints(false)
		ui.clearLog()
		return ev
	case "quit":
		ui.app.Stop()
		return ev
	}

	switch ui.page {
	case PAGE_VIEWER:
		switch cmd {
		case "view":
			fallthrough
		case "back":
			ui.unviewTldr()
		}
	case PAGE_COMMANDS:
		ui.commandPager.handleCommand(cmd)
	case PAGE_SECTIONS:
		ui.sectionPager.handleCommand(cmd)
	}

	return ev
}

func (ui *UI) switchToSectionPager() {
	ui.router.SwitchToPage(PAGE_SECTIONS)
}

func (ui *UI) switchToCommandPager(s Section) {
	p, ok := ui.commandPagers[s.Name]
	if !ok {
		debugLogger.Log("[ui] creating new CommandPager: %s", s.Name)
		p = newCommandPager(ui, s)
		ui.commandPagers[s.Name] = p
	}
	ui.commandPager = p
	ui.router.AddAndSwitchToPage(PAGE_COMMANDS, p, true)
}

func (ui *UI) routerChanged() {
	n, _ := ui.router.GetFrontPage()
	if n == ui.page {
		return
	}
	ui.page = n
	switch ui.page {
	case PAGE_VIEWER:
		ui.footer.clearPage()
	case PAGE_COMMANDS:
		ui.commandPager.updatePage()
	case PAGE_SECTIONS:
		ui.sectionPager.updatePage()
	}
}

func newUI(config Config) *UI {

	dataset, updated := newDataset(config)

	ui := &UI{
		app:           tview.NewApplication(),
		root:          tview.NewFlex().SetDirection(tview.FlexRow),
		router:        tview.NewPages(),
		dataset:       dataset,
		commandPagers: make(map[string]*Pager[Command]),
		keyMap:        config.keyMap,
		maxPerColumn:  -1,
	}
	ui.footer = newFooter(ui)

	go func() {
		status, ok := <-updated
		if ok && status {
			ui.app.QueueUpdateDraw(func() {
				ui.showMessage("Updated! Restart required", "[green]")
			})
		}
	}()

	ui.router.SetChangedFunc(ui.routerChanged)

	ui.hintKey = ui.keyMap.generateHintKey()
	ui.hints = strings.Join(ui.keyMap.generateHints(0), " ")
	ui.fullHints = strings.Join(ui.keyMap.generateHints(2), " ")

	ui.root.SetInputCapture(ui.handleKey)

	ui.root.AddItem(ui.router, 0, 6, true)
	ui.root.AddItem(ui.footer, 0, 1, false)

	ui.sectionPager = newSectionPager(ui)
	ui.viewer = newViewer()

	ui.router.AddPage(PAGE_SECTIONS, ui.sectionPager, true, true)
	ui.router.AddPage(PAGE_VIEWER, ui.viewer, true, false)

	ui.resetHints(false)
	ui.clearLog()
	ui.app.SetRoot(ui.root, true)

	go func() {
		for msg := range logger.queue {
			ui.app.QueueUpdateDraw(func() {
				ui.showMessage(msg, "[red]")
			})
		}
	}()

	return ui
}
