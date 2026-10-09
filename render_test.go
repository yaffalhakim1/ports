package main

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
	writePNG(t, "testdata/window.png", tt.Image())
	writePNG(t, "testdata/window-dark.png", renderDark(a))
}

// TestRenderStates writes the states that are hard to reach by hand: a
// system process chosen, where Stop and Restart are unavailable.
func TestRenderStates(t *testing.T) {
	a := newTestApp()
	a.chosen = 3 // the System row
	tt := ui.NewTester(a.view, 900, 620)
	tt.Frame()
	writePNG(t, "testdata/window-disabled.png", tt.Image())
}

func renderDark(a *app) *image.RGBA {
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
		writePNG(t, "testdata/appearance-"+c.name+".png", tt.Image())
	}
}
