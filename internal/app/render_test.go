package app

import (
	"image"
	"image/png"
	"os"
	"testing"

	"github.com/egoist/mygo/ui"
)

// TestRender writes PNGs of the window, so a change to the interface can
// be looked at without launching the app.
func TestRender(t *testing.T) {
	a := newTestApp()
	a.chosen = 1
	a.status = "Copied http://127.0.0.1:3000"
	tt := ui.NewTester(a.view, 900, 620)
	tt.Frame()
	writePNG(t, "../../testdata/window.png", tt.Image())
	writePNG(t, "../../testdata/window-dark.png", renderDark(a))
}

// TestRenderStates writes the states that are hard to reach by hand: a
// system process chosen, where Stop and Restart are unavailable.
func TestRenderStates(t *testing.T) {
	a := newTestApp()
	row := -1
	for i, p := range a.shown {
		if p.SystemPID() {
			row = i
			break
		}
	}
	if row < 0 {
		t.Fatal("the fixture has no system process to render")
	}
	a.chosen = row
	tt := ui.NewTester(a.view, 900, 620)
	tt.Frame()
	writePNG(t, "../../testdata/window-disabled.png", tt.Image())
}

// TestRenderLive writes the window at the size it opens at, which is the
// capture the README and the landing page show.
func TestRenderLive(t *testing.T) {
	a := newTestApp()
	a.chosen = 1
	a.status = "Copied http://127.0.0.1:3000"
	tt := ui.NewTester(a.view, defaultWidth, defaultHeight)
	tt.Frame()
	writePNG(t, "../../testdata/live.png", tt.Image())
}

// TestRenderDialogs writes the two surfaces that are only on screen while
// something is being decided: the confirmation before a process is ended,
// and the command palette. Both are the library's dialogs wearing the app's
// theme, so they are worth looking at as a whole rather than one control at
// a time.
func TestRenderDialogs(t *testing.T) {
	a := newTestApp()
	a.chosen = 1
	a.ask("Stop", a.shown[1])
	tt := ui.NewTester(a.view, 900, 620)
	tt.Frame()
	tt.Frame()
	writePNG(t, "../../testdata/confirm.png", tt.Image())

	b := newTestApp()
	b.paletteOpen = true
	tt2 := ui.NewTester(b.view, 900, 620)
	tt2.Frame()
	tt2.Frame()
	writePNG(t, "../../testdata/palette.png", tt2.Image())
}

func renderDark(a *App) *image.RGBA {
	tt := ui.NewTester(a.view, 900, 620)
	tt.SetDark(true)
	tt.Frame()
	return tt.Image()
}

// writePNG writes a rendered frame where a person can look at it.
func writePNG(t *testing.T, path string, img *image.RGBA) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		t.Fatal(err)
	}
	t.Logf("wrote %s", path)
}

// TestRenderAppearances writes both appearances, so the brand palette can
// be looked at in each rather than only measured.
func TestRenderAppearances(t *testing.T) {
	for _, c := range []struct {
		name       string
		appearance Appearance
	}{
		{"light", AppearanceLight},
		{"dark", AppearanceDark},
	} {
		a := newTestApp()
		a.appearance = c.appearance
		a.chosen = 1
		a.status = "Copied http://127.0.0.1:3000"
		tt := ui.NewTester(a.view, 960, 660)
		tt.Frame()
		writePNG(t, "../../testdata/appearance-"+c.name+".png", tt.Image())
	}
}
