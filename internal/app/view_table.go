package app

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"

	"ports/internal/port"
)

// columns are the table's columns. The port is fixed at the left so it stays
// visible while the rest scroll sideways, and it is wide enough for the
// state badge beside the number; the others the user may reorder and resize.
// Each column's ID is the field it sorts by, so a header click needs no
// translation.
func columns() []ui.TableColumn {
	return []ui.TableColumn{
		{Title: "Port", ID: "port", Width: 136, MinWidth: 120, MaxWidth: 160, Fixed: true, Sortable: true},
		{Title: "Process", ID: "name", Width: 150, MinWidth: 90, Sortable: true},
		{Title: "Address", ID: "address", Width: 108, MinWidth: 80, Sortable: true},
		// The pid is a number, so it ends where the column does and its
		// digits line up down the column.
		{Title: "PID", ID: "pid", Width: 70, MinWidth: 56, Align: ui.End, Sortable: true},
		// The folder takes the room the others leave, and it goes last so
		// that room falls at the end of the row rather than between the
		// port and the address it is bound to. It is also the field that
		// can be longest: it is what tells two servers of the same program
		// apart.
		{Title: "Folder", ID: "folder", Width: 0, MinWidth: 110, Sortable: true},
	}
}

// table shows the ports, one row each. It is a real table so the user can
// sort by a header, resize the columns and choose a row with the keyboard,
// and so only the rows in view are built.
func (a *App) table(c *ui.Context) {
	a.wireTable()

	// The table is inset by the window's margin, so its columns start and
	// end on the same edges as the toolbar and the action bar.
	components.TableVirtualized(c, &a.tableState, columns(), len(a.shown), func(row, col int) {
		p := a.shown[row]
		switch col {
		case 0:
			a.portCell(c, p)
		case 1:
			a.processCell(c, p)
		case 2:
			a.addressCell(c, p)
		case 3:
			a.pidCell(c, p)
		case 4:
			a.folderCell(c, p)
		}
	}).Grow(1).MinHeight(0).PaddingX(pageMargin).Children(func() {
		if len(a.shown) == 0 {
			a.empty(c)
		}
	})
}

// wireTable hands the table the app's own state: which row is chosen, how to
// name a row, the sort order a header click writes, and what to do when it
// changes. It runs every frame because the row count and the rows themselves
// change with every scan.
func (a *App) wireTable() {
	a.tableState.List.Key = func(row int) any { return a.shown[row].Key() }
	a.tableState.List.Label = func(row int) string { return a.shown[row].Label() }
	a.tableState.List.Selected = &a.chosen
	// The header writes the chosen order into the app's own sort order and
	// reports the click, so there is one copy of the order rather than two
	// to keep in step.
	a.tableState.Sort = &a.order
	a.tableState.OnSort = a.refresh
}

// The cells take their theme from the frame rather than from the context:
// they are built in a later pass, where the context no longer carries it, and
// a cell that asked the context there would paint the library's default
// palette into the window.

// portCell shows the port number and its state. The state is a word in a
// tinted pill rather than a colour or a glyph on its own, so it is readable
// without seeing the tint, and the address column spells out the same thing.
//
// The state is also the one thing the old list left unsaid: an exposed port —
// the kind another machine can reach — had no mark of its own at all.
func (a *App) portCell(c *ui.Context, p port.Port) {
	t := a.themeOf()
	ui.Row(c).AlignItems(ui.Center).Gap(t.Space(2)).Children(func() {
		components.Text(c, fmt.Sprint(p.Port)).FontSize(t.FontSize - 1).
			Bold().FontFeatures("tnum").Shrink(0)
		a.stateChip(c, p)
	})
}

// stateChip is the port's state as a pill. It takes its tint from the state
// rather than from the badge's own neutral face, which is built for a page
// background and goes white-on-white against a dark row; the tint is the
// theme's tone over the canvas, so both appearances get a chip that holds its
// text.
func (a *App) stateChip(c *ui.Context, p port.Port) {
	t := a.themeOf()
	tint := t.TextMuted
	if p.Exposed() {
		tint = t.Warning
	}
	ui.Row(c).AlignItems(ui.Center).Shrink(0).
		Padding(1, 6).Radius(t.Radius-2).
		Background(t.Background.Mix(tint, 0.18)).
		Border(1, t.Background.Mix(tint, 0.45)).
		Children(func() {
			components.Text(c, stateWord(p)).FontSize(t.FontSize - 3).TextColor(tint)
		})
}

// stateWord is the port's state as the pill spells it: what the network can
// reach, and what only this machine can.
func stateWord(p port.Port) string {
	if p.Exposed() {
		return "exposed"
	}
	return "local"
}

// processCell shows the process name, with a warning icon when the process
// cannot be stopped, so the disabled action is explained where the user looks
// for it.
func (a *App) processCell(c *ui.Context, p port.Port) {
	t := a.themeOf()
	ui.Row(c).AlignItems(ui.Center).Gap(t.Space(1.5)).Children(func() {
		if p.SystemPID() {
			warn := components.Icon(c, lucide.TriangleAlert).Size(12, 12).
				TextColor(t.Warning).Shrink(0)
			components.Tooltip(c, warn, "A system process, which cannot be stopped")
		}
		components.Text(c, p.Name).FontSize(t.FontSize - 1).SingleLine()
	})
}

// folderCell shows the last part of the working directory, which is what
// tells two servers of the same program apart, with the whole path on the
// tooltip because two projects can share a folder name.
func (a *App) folderCell(c *ui.Context, p port.Port) {
	t := a.themeOf()
	name := p.Folder()
	if name == "" {
		components.Text(c, "unknown").FontSize(t.FontSize - 1).TextColor(t.TextMuted)
		return
	}
	ui.Row(c).AlignItems(ui.Center).Gap(t.Space(1.5)).Children(func() {
		components.Icon(c, lucide.Folder).Size(12, 12).TextColor(t.TextMuted).Shrink(0)
		label := components.Text(c, name).FontSize(t.FontSize - 1).SingleLine().Ellipsis("…")
		components.Tooltip(c, label, p.Dir)
	})
}

// addressCell shows the bound address, which is the source of truth the
// state badge beside the port number agrees with.
func (a *App) addressCell(c *ui.Context, p port.Port) {
	t := a.themeOf()
	components.Text(c, p.Addr).FontSize(t.FontSize - 1).TextColor(t.TextMuted).
		SingleLine()
}

// pidCell shows the process id in tabular figures, so the digits line up
// down the column.
func (a *App) pidCell(c *ui.Context, p port.Port) {
	t := a.themeOf()
	components.Text(c, fmt.Sprint(p.PID)).FontSize(t.FontSize - 1).
		TextColor(t.TextMuted).FontFeatures("tnum")
}

// empty is what the table shows when nothing matches, which is also where a
// failed scan is reported. It names what was narrowed the list away and
// offers the way out, rather than only reporting that there is nothing.
func (a *App) empty(c *ui.Context) {
	t := c.Theme()
	message, hint := "No listening ports", "Nothing is listening on this machine."
	switch {
	case a.lastErr != nil:
		message = "Could not read the system"
		hint = a.lastErr.Error() + ". Scanning again may help."
	case a.query != "":
		message = "No ports match “" + a.query + "”"
		hint = "Clear the search to see every port."
	case a.folder != "":
		message = "Nothing is listening in " + folderName(a.folder)
		hint = "“" + ModeAll.Label() + "” shows every folder again."
	case a.mode != ModeAll:
		message, hint = a.mode.Empty()
	}

	ui.Column(c).Fill().Center().Gap(8).Padding(40).Children(func() {
		components.Text(c, message).FontSize(t.FontSize + 1).Bold().TextColor(t.TextMuted)
		components.Text(c, hint).FontSize(t.FontSize - 1).TextColor(t.TextMuted).MaxWidth(420)
		if a.query != "" {
			clear := components.Button(c, "", components.ButtonOutline)
			clear.Children(func() { components.Text(c, "Clear search") })
			clear.OnClick(func() {
				a.query = ""
				a.refresh()
			})
		}
	})
}

// folderName is the last part of a working directory, for a message that has
// to name it without printing the whole path.
func folderName(dir string) string {
	if i := strings.LastIndexAny(dir, `\/`); i >= 0 {
		return dir[i+1:]
	}
	return dir
}
