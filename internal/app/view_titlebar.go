package app

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

// narrowWidth is the window width below which the sidebar is put away for
// want of room. The action bar needs about 600 DIPs, the sidebar is 240
// beside it, and the window's own frame takes the rest.
const narrowWidth = 880

// beginFrame is the shell's Theme option: everything the app does once per
// frame, at the frame's root, before a single widget is built. It installs
// the brand palette for the current appearance, so the shell's own chrome and
// every Onevra control draw the brand rather than Onevra's default palette,
// and it declares the window's keys.
//
// Both belong here rather than deeper in the tree. A key declared from inside
// a region only reaches that region, and the first press of it lands a frame
// late; declared at the root, the palette covers the whole window and the
// first press works. The frame's clock is taken here too, because the summary
// says how fresh the list is and must read the same clock as the rest of the
// frame.
func (a *App) beginFrame(c *ui.Context) *ui.Theme {
	a.now = c.Now()
	a.shortcuts(c)
	a.frame = theme(a.appearanceOf(c))
	return a.frame
}

// shortcuts declares the window's keys. They are declared in the root scope,
// so they work wherever the focus is. R also hands the keyboard to the list,
// so the arrow keys reach it without a click first.
func (a *App) shortcuts(c *ui.Context) {
	c.OnShortcut(0, ui.KeyR, func() {
		a.scanOnce()
		a.tableState.List.Focus(c)
	})
	c.OnShortcut(0, ui.KeyF, func() { a.focusSearch(c) })
	c.OnShortcut(ui.Cmd, ui.KeyK, func() { a.openPalette(c) })
}

// titlebar is the shell's leading titlebar area: what the window is, and
// what the list is showing. It sits in the draggable strip, so it holds no
// control that a pointer has to drag.
func (a *App) titlebar(c *ui.Context) {
	a.watchWidth(c)

	t := c.Theme()
	ui.Row(c).AlignItems(ui.Center).Gap(t.Space(3)).Children(func() {
		components.Text(c, "Ports").FontSize(t.FontSize + 1).Bold().Shrink(0)
		components.Text(c, a.summary()).FontSize(t.FontSize - 1).TextColor(t.TextMuted).
			SingleLine().Ellipsis("…").Shrink(1).MinWidth(0)
	})
}

// watchWidth puts the sidebar away when the window is too narrow to hold it
// and the list together, and gives it back when the window grows again.
//
// While there is no room the sidebar is closed on every frame, so opening it
// by hand cannot squeeze the list: the width wins. Past the threshold it is
// opened once, when the window crosses back over, and after that it is the
// user's to toggle without the frame undoing it.
func (a *App) watchWidth(c *ui.Context) {
	if a.shell == nil {
		return
	}
	width, _ := c.Size()
	narrow := width < narrowWidth
	switch {
	case narrow:
		if a.shell.LeftExpanded() {
			a.shell.SetLeftExpanded(false)
		}
	case a.narrow:
		a.shell.SetLeftExpanded(true)
	}
	a.narrow = narrow
}

// summary says how many ports are listening and when they were read, which
// is the window's only claim about freshness.
func (a *App) summary() string {
	switch {
	case a.lastErr != nil:
		return "could not read the system"
	case a.lastScan.IsZero():
		return "reading…"
	}
	shown, total := len(a.shown), len(a.ports)
	count := fmt.Sprintf("%d listening", total)
	if shown != total {
		count = fmt.Sprintf("%d of %d shown", shown, total)
	}
	return count + " · updated " + ago(a.now.Sub(a.lastScan))
}

// ago writes a short elapsed time, as a status line reads it.
func ago(d time.Duration) string {
	switch {
	case d < time.Second:
		return "now"
	case d < time.Minute:
		return fmt.Sprintf("%ds ago", int(d.Seconds()))
	default:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	}
}
