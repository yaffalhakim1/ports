package app

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/shells/sidebar"
)

// view builds the whole window the way the shell lays it out: the sidebar
// beside the content, both full height, under a header that floats over them
// and is aligned to the content area.
//
// Tests drive the app through this rather than through the shell itself,
// because the shell's window controls, sidebar animation and resize handle
// are the library's, and testing them would test Onevra. What it does copy is
// the geometry that the tests are really asking about: which region is inset
// for the header, where the titlebar starts, and that a window too narrow for
// the sidebar gets the whole width for the list.
func (a *App) view(c *ui.Context) {
	t := a.beginFrame(c)
	c.SetTheme(t)

	width, _ := c.Size()
	a.narrow = width < narrowWidth
	left := float32(0)
	if !a.narrow {
		left = sidebarWidth
	}

	ui.Row(c).Fill().Children(func() {
		if !a.narrow {
			ui.Column(c).Width(left).Shrink(0).FillHeight().
				Background(t.Surface).Children(func() { a.side(c) })
		}
		ui.Column(c).Grow(1).MinWidth(0).FillHeight().
			Background(t.Background).Children(func() { a.content(c) })
	})
	ui.Row(c).Absolute().Top(0).Left(0).Right(0).
		Height(sidebar.HeaderHeight(c)).AlignItems(ui.Center).Children(func() {
		ui.Box(c).Width(left).Shrink(0)
		a.titlebar(c)
	})
}
