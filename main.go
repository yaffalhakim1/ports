package main

import (
	"log"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

func main() {
	s := loadSettings()
	a := &app{chosen: -1, appearance: parseAppearance(s.Appearance)}

	mygo.App.WhenReady(func() {
		win := mygo.NewWindow(mygo.WindowOptions{
			Title:     "Ports",
			Width:     960,
			Height:    660,
			MinWidth:  620,
			MinHeight: 420,
			// Opens where the user left it last time.
			StateKey: "main",
			Content:  ui.View(a.view),
		})
		a.win = win
		a.startTray()
		a.watchClose()
		a.start()
	})

	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
