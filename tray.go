package main

import (
	"log"

	"github.com/egoist/mygo"
)

// tray puts an icon in the notification area and keeps the app alive there
// when the window is closed, which is what a tool that watches the machine
// wants: it keeps scanning without occupying the screen.
//
// The window is hidden, never destroyed, so the scan loop, the chosen row
// and the scroll position all survive a close.
func (a *app) startTray() {
	tray, err := mygo.NewTray(mygo.TrayOptions{
		Icon:           trayIconPNG(),
		IconIsTemplate: true,
		ToolTip:        "Ports — listening TCP ports",
		Menu: mygo.NewMenu([]*mygo.MenuItem{
			{Label: "Show Ports", Click: func(*mygo.MenuItem, *mygo.Window) { a.showWindow() }},
			mygo.Separator(),
			{Label: "Quit", Click: func(*mygo.MenuItem, *mygo.Window) { a.quit() }},
		}),
	})
	if err != nil {
		// A tray that cannot be created is not fatal: the app falls back
		// to quitting when its window closes, which still works.
		log.Println("ports: no tray icon:", err)
		return
	}
	a.tray = tray
}

// watchClose keeps the app running when the window is closed, unless the
// user asked to quit. Closing the window is how the app goes to the tray.
func (a *app) watchClose() {
	a.win.OnClose(func(e *mygo.CloseEvent) {
		if a.quitting || a.tray == nil {
			return
		}
		e.PreventDefault()
		a.win.Hide()
	})
}

// showWindow brings the window back from the tray.
func (a *app) showWindow() {
	a.win.Show()
	a.win.Focus()
}

// quit ends the app for real, rather than closing to the tray.
func (a *app) quit() {
	a.quitting = true
	mygo.App.Quit()
}
