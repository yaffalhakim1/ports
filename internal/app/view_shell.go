package app

import (
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/shells/sidebar"
)

// defaultWidth and defaultHeight are the size the window opens at, which is
// also the size the README's capture shows.
const (
	defaultWidth  = 1120
	defaultHeight = 720
)

// sidebarWidth is the sidebar's width: how much room a folder name and its
// count need. The tests lay their own copy of the frame out from it.
const sidebarWidth = 240

// ShellOptions describes the window this app wants: its name and size, the
// region builders the shell calls every frame, and the brand theme it should
// install before it builds anything.
//
// The window is frameless, so the shell draws the titlebar and the window
// controls. That is why the controls that used to live in the app's own
// header row now sit in the titlebar (the summary, beside the title) and in
// the toolbar above the list (the search and the two buttons), rather than
// being laid out by the app from the top of the window down.
func (a *App) ShellOptions() sidebar.ShellOptions {
	return sidebar.ShellOptions{
		Title:     "Ports",
		Width:     defaultWidth,
		Height:    defaultHeight,
		MinWidth:  620,
		MinHeight: 420,
		StateKey:  "main",

		// Holds a folder name and its count without the name ellipsising,
		// which is the narrowest the sidebar can be and still be worth
		// opening.
		LeftWidth:    sidebarWidth,
		LeftMinWidth: 200,

		ContentTitle: a.titlebar,
		Left:         a.side,
		Content:      a.content,

		Toasts: components.ToastOptions{Placement: components.ToastBottomEnd},
		Theme:  a.beginFrame,
	}
}
