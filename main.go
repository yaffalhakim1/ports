package main

import (
	"log"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/plugins/updater/native"
	"github.com/oneo-me/onevra/source/shells/sidebar"

	"ports/internal/app"
)

func main() {
	// The native update window: the app draws its own UI, so it needs no
	// webview. It checks for a newer release in the background.
	mygo.Use(native.Plugin)

	s := app.LoadSettings()
	a := app.New(app.ParseAppearance(s.Appearance))

	mygo.App.WhenReady(func() {
		// The shell owns the window: its titlebar, its window controls and
		// its sidebar. The app owns what goes in them.
		shell := sidebar.NewShell(a.ShellOptions())
		a.SetShell(shell)
		a.SetWindow(shell.Window())
		a.StartTray()
		a.WatchClose()
		a.Start()
	})

	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
