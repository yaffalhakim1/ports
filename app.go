package main

import (
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

// refreshEvery is how often the list re-reads the system on its own. A
// scan touches every process, so it runs on a background goroutine and
// never on the UI thread.
const refreshEvery = 2 * time.Second

// app is the window's state. The view builds the interface from it on the
// main thread; a background scan changes it through Window.Update.
type app struct {
	win *mygo.Window

	// ports is the last scan; shown is that list filtered and sorted the
	// way the window currently shows it.
	ports []Port
	shown []Port
	query string
	order ui.SortOrder
	only  bool // show only the ports another machine can reach

	status string

	// scanning is true while a scan is in flight, so the view can say so.
	scanning bool
	lastErr  error
	lastScan time.Time

	// chosen is the selected row, as an index into shown.
	chosen int
	list   ui.ListState

	// appliedOrder is the sort order the shown list was last sorted by, so a
	// header click, which the table records in order, is applied once.
	appliedOrder ui.SortOrder

	// searchBox is the filter field, so the keyboard can reach it.
	searchBox ui.Handle

	// appearance is which look to show, and quitting is set when the app
	// is really going away rather than closing to the tray.
	appearance Appearance
	quitting   bool
	tray       *mygo.Tray

	// pending names an action waiting for the user to confirm it, since
	// ending a process cannot be undone.
	pending *pendingAction

	// now is the frame's time, read once per frame.
	now time.Time
}

// pendingAction is a stop or restart waiting for confirmation.
type pendingAction struct {
	verb   string // "Stop" or "Restart"
	port   Port
	detail string
}

// selected returns the port the chosen row shows, and whether there is one.
func (a *app) selected() (Port, bool) {
	if a.chosen < 0 || a.chosen >= len(a.shown) {
		return Port{}, false
	}
	return a.shown[a.chosen], true
}

// start scans at once and then on every tick, so the window is never
// showing a stale list for long.
func (a *app) start() {
	go func() {
		a.scanNow()
		ticker := time.NewTicker(refreshEvery)
		defer ticker.Stop()
		for range ticker.C {
			a.scanNow()
		}
	}()
}

// scanNow reads the system and hands the result to the window's state on
// the main thread. A scan that fails keeps the last good list and reports
// the error, rather than emptying the window.
func (a *app) scanNow() {
	if a.win == nil {
		return
	}
	a.win.Update(func() { a.scanning = true })
	ports, err := scan()
	a.win.Update(func() {
		a.scanning = false
		a.lastScan = time.Now()
		if err != nil {
			a.lastErr = err
			a.status = "Scan failed: " + err.Error()
			return
		}
		a.lastErr = nil
		a.apply(ports)
	})
}

// apply replaces the list and rebuilds what the window shows, keeping the
// chosen row on the same port while it is still listening.
func (a *app) apply(ports []Port) {
	was, hadChoice := a.selected()
	a.ports = ports
	a.refresh()
	if !hadChoice {
		return
	}
	for i, p := range a.shown {
		if p.Key() == was.Key() {
			a.chosen = i
			return
		}
	}
	a.chosen = -1
}

// refresh filters and sorts the last scan into what the window shows.
func (a *app) refresh() {
	ports := a.ports
	if a.only {
		exposed := make([]Port, 0, len(ports))
		for _, p := range ports {
			if p.Exposed() {
				exposed = append(exposed, p)
			}
		}
		ports = exposed
	}
	a.shown = filter(ports, a.query)
	sortPorts(a.shown, a.order)
	if a.chosen >= len(a.shown) {
		a.chosen = -1
	}
}

// appearanceOf resolves the appearance to show: the user's choice, or what
// the desktop asks for when the choice is to follow it.
func (a *app) appearanceOf(c *ui.Context) bool {
	switch a.appearance {
	case AppearanceLight:
		return false
	case AppearanceDark:
		return true
	default:
		return c.Theme().Dark
	}
}

// applyTheme installs the theme for the current appearance. It is called
// at the top of the view, so a change of appearance takes effect on the
// next frame.
func (a *app) applyTheme(c *ui.Context) {
	c.SetTheme(theme(a.appearanceOf(c)))
}
