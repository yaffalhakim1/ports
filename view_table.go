package main

import (
	"fmt"

	"github.com/egoist/mygo/ui"
)

// columns are the table's columns. The port is fixed at the left so it
// stays visible while the rest scroll sideways; the others the user may
// reorder and resize. Each column's ID is the field it sorts by, so a
// header click needs no translation.
func columns() []ui.TableColumn {
	return []ui.TableColumn{
		{Title: "Port", ID: "port", Width: 84, MinWidth: 64, MaxWidth: 120, Fixed: true, Sortable: true},
		{Title: "Process", ID: "name", Width: 140, MinWidth: 90, Sortable: true},
		// The folder takes the room the others leave: it is what tells two
		// servers of the same program apart, so it deserves the space.
		{Title: "Folder", ID: "folder", Width: 0, MinWidth: 110, Sortable: true},
		{Title: "Address", ID: "address", Width: 108, MinWidth: 80, Sortable: true},
		{Title: "PID", ID: "pid", Width: 70, MinWidth: 56, Align: ui.End, Sortable: true},
	}
}

// table shows the ports, one row each. It is a real table so the user can
// sort by a header, resize the columns and choose a row with the keyboard,
// and so only the rows in view are built.
func (a *app) table(c *ui.Context) {
	a.list.Key = func(row int) any { return a.shown[row].Key() }
	a.list.Label = func(row int) string { return a.shown[row].Label() }
	a.list.Selected = &a.chosen
	a.list.Sort = &a.order

	// The table is inset by the window's margin, so its columns start and
	// end on the same edges as the header and the footer.
	ui.Table(c, &a.list, columns(), len(a.shown), func(row, col int) {
		p := a.shown[row]
		switch col {
		case 0:
			portCell(c, p)
		case 1:
			processCell(c, p)
		case 2:
			folderCell(c, p)
		case 3:
			addressCell(c, p)
		case 4:
			ui.Text(c, fmt.Sprint(p.PID)).FontSize(13).TextColor(c.Theme().TextMuted).
				FontFeatures("tnum")
		}
	}).Grow(1).PaddingX(pageMargin).Children(func() {
		if len(a.shown) == 0 {
			a.empty(c)
		}
	}).Submitted()
}

// portCell shows the port number, and marks whether the port serves this
// machine alone. The lock is not the only signal: the address column
// spells the same thing out, so the meaning is never carried by an icon
// alone.
func portCell(c *ui.Context, p Port) {
	t := c.Theme()
	ui.Row(c).AlignItems(ui.Center).Gap(6).Children(func() {
		ui.Text(c, fmt.Sprint(p.Port)).FontWeight(600).FontSize(13).
			FontFeatures("tnum")
		if p.Loopback() {
			ui.Icon(c, iconLock).Size(11, 11).TextColor(t.TextMuted).
				Tooltip("Serves this machine only")
		}
	})
}

// processCell shows the process name, with a warning icon when the process
// cannot be stopped, so the disabled action is explained where the user
// looks for it.
func processCell(c *ui.Context, p Port) {
	t := c.Theme()
	ui.Row(c).AlignItems(ui.Center).Gap(6).Children(func() {
		if p.systemPID() {
			ui.Icon(c, iconWarn).Size(12, 12).TextColor(t.Warning).
				Tooltip("A system process, which cannot be stopped")
		}
		ui.Text(c, p.Name).FontSize(13).SingleLine()
	})
}

// folderCell shows the last part of the working directory, which is what
// tells two servers of the same program apart.
func folderCell(c *ui.Context, p Port) {
	t := c.Theme()
	name := p.Folder()
	if name == "" {
		ui.Text(c, "unknown").FontSize(13).TextColor(t.TextMuted)
		return
	}
	ui.Row(c).AlignItems(ui.Center).Gap(6).Children(func() {
		ui.Icon(c, iconFolder).Size(12, 12).TextColor(t.TextMuted)
		ui.Text(c, name).FontSize(13).SingleLine().Ellipsis("…").Tooltip(p.Dir)
	})
}

// addressCell shows the bound address. A wildcard address is spelled out
// as well as shown, because "0.0.0.0" does not tell a reader that the port
// answers on every interface.
func addressCell(c *ui.Context, p Port) {
	ui.Text(c, p.Addr).FontSize(13).TextColor(c.Theme().TextMuted).SingleLine()
}

// empty is what the table shows when nothing matches, which is also where
// a failed scan is reported. It names what was searched for and offers a
// way out, rather than only reporting that there is nothing.
func (a *app) empty(c *ui.Context) {
	t := c.Theme()
	message, hint := "No listening ports", "Nothing is listening on this machine."
	switch {
	case a.lastErr != nil:
		message = "Could not read the system"
		hint = a.lastErr.Error() + ". Scanning again may help."
	case a.query != "":
		message = "No ports match “" + a.query + "”"
		hint = "Clear the search to see every port."
	}

	ui.Column(c).Fill().Center().Gap(8).Padding(40).Children(func() {
		ui.Text(c, message).FontSize(15).FontWeight(600).TextColor(t.TextMuted)
		ui.Text(c, hint).FontSize(13).TextColor(t.TextMuted).MaxWidth(420)
		if a.query != "" {
			ui.Button(c, "Clear search").OnClick(func() {
				a.query = ""
				a.refresh()
			})
		}
	})
}
