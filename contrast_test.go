package main

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// The pairs are asserted rather than eyeballed. The ratios are WCAG 2:
// 4.5:1 for text under 24px, 3:1 for graphical objects and focus rings.

// TestContrastOfEveryPair checks each foreground against the background it
// actually renders on, in both appearances, because a pair that passes in
// the light theme can fail in the dark one.
func TestContrastOfEveryPair(t *testing.T) {
	for _, dark := range []bool{false, true} {
		got := theme(dark)

		pairs := []struct {
			name string
			fg   ui.Color
			bg   ui.Color
			need float64
		}{
			{"body text on the window", got.Text, got.Background, 4.5},
			{"muted text on the window", got.TextMuted, got.Background, 4.5},
			{"muted text on a control", got.TextMuted, got.Surface, 4.5},
			{"muted text on a hovered control", got.TextMuted, got.SurfaceHover, 4.5},
			{"text in the search field", got.Text, got.Surface, 4.5},
			{"the chosen row's text", got.AccentText, got.Accent, 4.5},
			{"the danger label", got.Danger, got.Background, 4.5},
			{"the warning icon", got.Warning, got.Background, 3.0},
			{"the success colour", got.Success, got.Background, 3.0},
			{"the focus ring on the window", got.Focus, got.Background, 3.0},
			// The scrollbar and the hairlines only have to be perceptible.
			{"a hairline border", got.Border, got.Background, 1.05},
			{"the scrollbar", got.Scrollbar, got.Background, 1.05},
		}
		for _, p := range pairs {
			if r := ratio(p.fg, p.bg); r < p.need {
				t.Errorf("dark=%v: %s is %.2f:1, needs %.1f:1",
					dark, p.name, r, p.need)
			}
		}
	}
}

// TestEveryTextRoleIsLegibleOnEverySurface checks the text ramp against
// each surface it can land on, not only the window background: a control's
// label sits on the control, and a chosen row sits on the accent.
func TestEveryTextRoleIsLegibleOnEverySurface(t *testing.T) {
	for _, dark := range []bool{false, true} {
		got := theme(dark)
		surfaces := []struct {
			name string
			c    ui.Color
		}{
			{"the window", got.Background},
			{"a control", got.Surface},
			{"a hovered control", got.SurfaceHover},
		}
		for _, s := range surfaces {
			if r := ratio(got.Text, s.c); r < 4.5 {
				t.Errorf("dark=%v: text on %s is %.2f:1", dark, s.name, r)
			}
			if r := ratio(got.TextMuted, s.c); r < 4.5 {
				t.Errorf("dark=%v: muted text on %s is %.2f:1", dark, s.name, r)
			}
		}
	}
}
