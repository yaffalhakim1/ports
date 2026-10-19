package app

import (
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/components/command"
	"github.com/oneo-me/onevra/source/shells/sidebar"

	"ports/internal/host"
	"ports/internal/port"
)

// refreshEvery is how often the list re-reads the system on its own. A
// scan touches every process, so it runs on a background goroutine and
// never on the UI thread.
const refreshEvery = 2 * time.Second

// App is the window's state. The view builds the interface from it on the
// main thread; a background scan changes it through Window.Update.
type App struct {
	win *mygo.Window

	// ports is the last scan; shown is that list filtered and sorted the
	// way the window currently shows it.
	ports []port.Port
	shown []port.Port
	query string
	order ui.SortOrder

	// mode is which view the sidebar has chosen, and folder narrows it to
	// one working directory; "" means every folder.
	mode   Mode
	folder string

	status string

	// scanning is true while a scan is in flight, so the view can say so.
	scanning bool
	lastErr  error
	lastScan time.Time

	// chosen is the selected row, as an index into shown. table is the
	// table's own state: the keyboard cursor, the column layout, the sort
	// order and the scrolling. chosen and table.List.Selected point at the
	// same number, one for the app to read and one for the table to move.
	chosen     int
	tableState components.TableVirtualizedState

	// searchBox is the filter field, so the keyboard can reach it.
	searchBox ui.Handle

	// toast is an outcome waiting to be shown. A toast can only be added
	// while a frame is being built, and the actions that raise one run off
	// the main thread, so they queue it here and the frame emits it.
	toast *ui.Toast

	// appearance is which look to show, and quitting is set when the app
	// is really going away rather than closing to the tray.
	appearance Appearance
	quitting   bool
	tray       *mygo.Tray

	// shell is the window's frame, which owns the sidebar; the app asks it
	// to collapse when the window is too narrow to hold both.
	shell *sidebar.Shell
	// narrow records whether the sidebar was collapsed for want of room,
	// so widening the window gives it back rather than fighting the user's
	// own toggle.
	narrow bool

	// palette is the command palette's query and active row, and open is
	// whether it is showing.
	paletteState command.State
	paletteOpen  bool

	// pending names an action waiting for the user to confirm it, since
	// ending a process cannot be undone.
	pending *pendingAction

	// now is the frame's time, read once per frame.
	now time.Time

	// frame is the theme the frame's appearance resolves to. beginFrame
	// installs it in the context and keeps it here, because the table's
	// cells are built in a later pass than the rest of the window and the
	// context no longer carries the app's theme by then: a cell that asked
	// the context there would paint the library's default palette, which is
	// light, into a dark window.
	frame *ui.Theme
}

// themeOf is the frame's theme, for the code that cannot read it from the
// context: the table's cells. Everything else asks the context, which is
// where beginFrame put the same value.
func (a *App) themeOf() *ui.Theme {
	if a.frame == nil {
		a.frame = theme(false)
	}
	return a.frame
}

// New returns an app with the remembered appearance, ready to be shown.
func New(appearance Appearance) *App {
	return &App{chosen: -1, appearance: appearance}
}

// SetWindow records the window the app draws into.
func (a *App) SetWindow(w *mygo.Window) { a.win = w }

// SetShell records the window's frame, which owns the sidebar. The app asks
// it to collapse when the window is too narrow for both it and the list.
func (a *App) SetShell(s *sidebar.Shell) { a.shell = s }

// setView switches the sidebar's view, keeping the chosen port on the same
// listener when the new view still holds it, for the same reason a scan
// does: switching to "Local only" should not lose the port being read.
func (a *App) setView(m Mode, folder string) {
	was, had := a.selected()
	a.mode, a.folder = m, folder
	a.refresh()
	if !had {
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

// pendingAction is a stop or restart waiting for confirmation.
type pendingAction struct {
	verb   string // "Stop" or "Restart"
	port   port.Port
	detail string
}

// selected returns the port the chosen row shows, and whether there is one.
func (a *App) selected() (port.Port, bool) {
	if a.chosen < 0 || a.chosen >= len(a.shown) {
		return port.Port{}, false
	}
	return a.shown[a.chosen], true
}

// Start scans at once and then on every tick, so the window is never
// showing a stale list for long. It runs on the main thread, which is where
// the theme source has to be set.
func (a *App) Start() {
	mygo.Theme.SetSource(a.appearance.Source())
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
func (a *App) scanNow() {
	if a.win == nil {
		return
	}
	a.win.Update(func() { a.scanning = true })
	ports, err := host.Scan()
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
func (a *App) apply(ports []port.Port) {
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

// refresh filters and sorts the last scan into what the window shows. The
// order is the sidebar's view, then its folder, then the search: from the
// broadest narrowing to the narrowest, so a port hidden by the mode is
// never reported as a missing search hit.
func (a *App) refresh() {
	ports := keepMode(a.ports, a.mode)
	if a.folder != "" {
		ports = keepFolder(ports, a.folder)
	}
	a.shown = filter(ports, a.query)
	sortPorts(a.shown, a.order)
	if a.chosen >= len(a.shown) {
		a.chosen = -1
	}
}

// appearanceOf resolves the appearance to show: the user's choice, or what
// the desktop asks for when the choice is to follow it. The frame's theme is
// MyGo's own until the shell installs the brand's, so reading its dark flag
// here reads the desktop's setting rather than this app's answer.
func (a *App) appearanceOf(c *ui.Context) bool {
	switch a.appearance {
	case AppearanceLight:
		return false
	case AppearanceDark:
		return true
	default:
		return c.Theme().Dark
	}
}
