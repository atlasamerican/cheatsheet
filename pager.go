package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	PAGER_SECTION = "SectionPager"
	PAGER_COMMAND = "CommandPager"
)

type Pager[T Section | Command] struct {
	*tview.Frame
	name      string
	ui        *UI
	columns   *tview.Flex
	columnsN  int
	width     int
	height    int
	start     int
	end       int
	last      int
	focus     int
	page      int
	pageStart []int
	widgets   []ComponentWidget[T]
}

func (r *Pager[T]) handleCommand(cmd string) {
	debugLogger.Log("[pager] %s handling command: %s", r.name, cmd)

	switch cmd {
	case "next":
		if r.focus == r.last {
			break
		}
		if r.focus == r.end {
			if len(r.pageStart) <= r.page {
				r.pageStart = append(r.pageStart, r.start)
			}
			r.page++
			r.start = r.end + 1
			r.update()
		}
		r.setFocus(r.focus + 1)
	case "prev":
		if r.focus == 0 {
			break
		}
		if r.focus == r.start {
			r.page--
			r.start = r.pageStart[r.page]
			r.update()
		}
		r.setFocus(r.focus - 1)
	case "nextPage":
		if r.end == r.last {
			break
		}
		if len(r.pageStart) <= r.page {
			r.pageStart = append(r.pageStart, r.start)
		}
		r.page++
		r.start = r.end + 1
		r.update()
		r.setFocus(r.start)
	case "prevPage":
		if r.start == 0 {
			break
		}
		r.page--
		r.start = r.pageStart[r.page]
		r.update()
		r.setFocus(r.start)
	case "view":
		switch r.name {
		case PAGER_SECTION:
			s := any(r.widgets[r.focus].getData()).(Section)
			r.ui.switchToCommandPager(s)
		case PAGER_COMMAND:
			c := any(r.widgets[r.focus].getData()).(Command)
			r.ui.viewTldr(c)
		}
	case "back":
		if r.name == PAGER_COMMAND {
			r.ui.switchToSectionPager()
		}
	}
}

func (r *Pager[T]) setFocus(i int) {
	debugLogger.Log("[pager] %s setting focus: %d", r.name, i)

	if r.focus != -1 {
		r.widgets[r.focus].focus(false)
	}
	r.widgets[i].focus(true)
	r.focus = i
}

func (r *Pager[T]) draw(w, h int) {
	debugLogger.Log("[pager] %s draw: %d, %d", r.name, w, h)

	if r.width == w && r.height == h {
		return
	}

	r.width = w
	r.height = h

	r.focus = -1
	r.start = 0
	r.page = 0
	r.pageStart = make([]int, 0)

	r.columnsN = 0
	for _, i := range breakpoints {
		if w > i {
			r.columnsN++
		}
	}

	debugLogger.Log("[pager] %s columns: %d", r.name, r.columnsN)

	r.update()
}

func (r *Pager[T]) clear() {
	for i := 0; i < r.columns.GetItemCount(); i++ {
		r.columns.GetItem(i).(*tview.Flex).Clear()
	}
	r.columns.Clear()
}

func (r *Pager[T]) update() {
	debugLogger.Log("[pager] %s updating...", r.name)

	r.clear()

	maxHeight := r.height - r.ui.footer.height
	height := maxHeight
	columnsN := r.columnsN - 1
	column := tview.NewFlex().SetDirection(tview.FlexRow)
	r.columns.AddItem(column, 0, 1, false)

	var (
		i int
		w ComponentWidget[T]
	)

	for i, w = range r.widgets[r.start:] {
		height -= w.getHeight()

		if height < 1 ||
			(r.ui.maxPerColumn > 0 && column.GetItemCount() == r.ui.maxPerColumn) {
			height = maxHeight
			columnsN--
			if columnsN < 0 {
				i--
				break
			}
			column = tview.NewFlex().SetDirection(tview.FlexRow)
			r.columns.AddItem(column, 0, 1, false)
		}

		column.AddItem(w, w.getHeight(), 1, false)
	}

	r.end = r.start + i
	debugLogger.Log("[pager] %s showing: %d, %d (%d)", r.name, r.start, r.end, r.last)

	if r.focus == -1 {
		r.setFocus(r.start)
	}

	r.updatePage()
}

func (r *Pager[T]) updatePage() {
	r.ui.footer.updatePage(r.page, r.end < r.last)
}

func (r *Pager[T]) init() {
	r.last = len(r.widgets) - 1
	r.end = r.last

	r.SetDrawFunc(func(_ tcell.Screen, _, _, _, _ int) (int, int, int, int) {
		x, y, w, h := r.GetInnerRect()
		r.draw(w, h)
		return x, y, w, h
	})
}

func newSectionPager(ui *UI) *Pager[Section] {
	cols := tview.NewFlex().SetDirection(tview.FlexColumn)
	r := &Pager[Section]{
		Frame:   tview.NewFrame(cols).AddText("Sections", true, tview.AlignCenter, tcell.ColorWhite),
		name:    PAGER_SECTION,
		ui:      ui,
		columns: cols,
		widgets: make([]ComponentWidget[Section], len(ui.dataset.sections)),
	}

	i := 0
	for _, d := range ui.dataset.sections {
		r.widgets[i] = d.widget()
		i++
	}

	r.init()
	return r
}

func newCommandPager(ui *UI, s Section) *Pager[Command] {
	cols := tview.NewFlex().SetDirection(tview.FlexColumn)
	r := &Pager[Command]{
		Frame:   tview.NewFrame(cols).AddText(s.Name, true, tview.AlignCenter, tcell.ColorWhite),
		name:    PAGER_COMMAND,
		ui:      ui,
		columns: cols,
		widgets: make([]ComponentWidget[Command], len(s.Commands)),
	}

	for i, d := range s.Commands {
		r.widgets[i] = d.widget()
	}

	r.init()
	return r
}
