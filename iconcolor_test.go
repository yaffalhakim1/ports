package main

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// TestDestructiveIconMatchesItsLabel checks the Stop icon takes the danger
// colour along with its label, rather than staying muted: a half-coloured
// destructive control reads as an accident.
func TestDestructiveIconMatchesItsLabel(t *testing.T) {
	for _, dark := range []bool{false, true} {
		a := newTestApp()
		a.chosen = 1
		tt := ui.NewTester(a.view, 900, 620)
		tt.SetDark(dark)
		tt.Frame()
		img := tt.Image()

		stop, ok := tt.Find("Stop")
		if !ok {
			t.Fatal("Stop is missing")
		}
		danger := theme(dark).Danger
		// The icon sits just left of the label, inside the button.
		iconX := int(stop.X) - 20
		y := int(stop.Y + stop.H/2)
		found := false
		for x := iconX - 6; x < int(stop.X); x++ {
			if x < 0 {
				continue
			}
			r, g, b, _ := img.At(x, y).RGBA()
			c := ui.RGB(uint8(r>>8), uint8(g>>8), uint8(b>>8))
			if near(c, danger, 40) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("dark=%v: the Stop icon is not drawn in the danger colour %v", dark, danger)
		}
	}
}

// near reports whether two colours are within tolerance on every channel.
func near(a, b ui.Color, tol int) bool {
	diff := func(x, y uint8) int {
		if x > y {
			return int(x - y)
		}
		return int(y - x)
	}
	return diff(a.R, b.R) <= tol && diff(a.G, b.G) <= tol && diff(a.B, b.B) <= tol
}
