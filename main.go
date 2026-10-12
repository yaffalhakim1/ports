package main

import (
	"log"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/plugins/updater/native"
	"github.com/egoist/mygo/ui"

	"ports/internal/app"
)

func main() {
	// The native update window: the app draws its own UI, so it needs no
	// webview. It checks for a newer release in the background.
	mygo.Use(native.Plugin)

	s := app.LoadSettings()
	a := app.New(app.ParseAppearance(s.Appearance))

	mygo.App.WhenReady(func() {
		win := mygo.NewWindow(mygo.WindowOptions{
			Title:     "Ports",
			Width:     960,
			Height:    660,
			MinWidth:  620,
			MinHeight: 420,
			// Opens where the user left it last time.
			StateKey: "main",
			Content:  ui.View(a.View),
		})
		a.SetWindow(win)
		a.StartTray()
		a.WatchClose()
		a.Start()
	})

	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
