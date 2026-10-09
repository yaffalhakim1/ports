package main

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// The brand palette is a web and mobile system, where the colours above the
// second step are fills rather than text. Three of its values do not hold
// up in the roles this app uses them for, so the theme derives from them.
// These tests pin that reasoning: if someone later "corrects" the theme
// back to the brand hexes, this says why that fails.

// TestBrandGreenNeedsDarkText checks the pairing the brand specifies. White
// on the brand green measures under 2:1, so a theme that pairs them is
// unreadable rather than merely off-brand.
func TestBrandGreenNeedsDarkText(t *testing.T) {
	white := ui.Hex("#ffffff")
	for _, green := range []ui.Color{brandGreen, brandGreenLight} {
		dark := ratio(onBrandGreen, green)
		light := ratio(white, green)
		t.Logf("on %v: brand near-black %.2f:1, white %.2f:1", green, dark, light)
		if dark < 4.5 {
			t.Errorf("the brand text colour on %v is %.2f:1", green, dark)
		}
		if light >= 4.5 {
			t.Errorf("white unexpectedly passes on %v at %.2f:1", green, light)
		}
	}
}

// TestLightBrandGreenCannotBeAFocusRing checks the derivation: the light
// brand green is invisible on the light canvas, so the ring uses the
// darker green the brand already defines for inline code.
func TestLightBrandGreenCannotBeAFocusRing(t *testing.T) {
	canvas := ui.Hex("#F6F5F6")
	raw := ratio(brandGreenLight, canvas)
	used := ratio(brandGreenLightDark, canvas)
	t.Logf("light brand green on the light canvas: %.2f:1", raw)
	t.Logf("the derived ring on the light canvas:    %.2f:1", used)
	if raw >= 3.0 {
		t.Errorf("the raw green unexpectedly passes as a ring at %.2f:1", raw)
	}
	if used < 3.0 {
		t.Errorf("the derived ring fails at %.2f:1", used)
	}
	if theme(false).Focus != brandGreenLightDark {
		t.Error("the light theme does not use the derived ring")
	}
}

// TestLightStatusIsDarkenedForText checks the derivation: the brand's
// status colours are diff fills, and as 13px text on the light canvas they
// fall short, so the theme uses darkened steps.
func TestLightStatusIsDarkenedForText(t *testing.T) {
	canvas := ui.Hex("#F6F5F6")
	for _, c := range []struct {
		name string
		raw  ui.Color
		used ui.Color
	}{
		{"success", statusSuccessLight, statusSuccessLightText},
		{"warning", statusWarningLight, statusWarningLightText},
		{"danger", statusDangerLight, statusDangerLightText},
	} {
		rawRatio, usedRatio := ratio(c.raw, canvas), ratio(c.used, canvas)
		t.Logf("%-8s raw %.2f:1, used %.2f:1", c.name, rawRatio, usedRatio)
		if usedRatio < 4.5 {
			t.Errorf("%s: the derived colour is %.2f:1, under 4.5", c.name, usedRatio)
		}
	}

	got := theme(false)
	if got.Danger != statusDangerLightText {
		t.Error("the light theme does not use the derived danger colour")
	}
}

// TestQuietGreysAreDecorationOnly checks the third derivation: the brand's
// quieter greys fall short as text, so the theme uses the brand's
// secondary grey for every muted role and leaves the rest for decoration.
func TestQuietGreysAreDecorationOnly(t *testing.T) {
	for _, c := range []struct {
		name   string
		canvas ui.Color
		quiet  []ui.Color
		used   ui.Color
	}{
		{"dark", darkCanvas, []ui.Color{darkTextTertiary, darkTextGhost}, darkTextSecondary},
		{"light", lightCanvas, []ui.Color{lightTextTertiary, lightTextGhost}, lightTextSecondary},
	} {
		for _, q := range c.quiet {
			t.Logf("%s quiet grey %v on the canvas: %.2f:1", c.name, q, ratio(q, c.canvas))
		}
		if r := ratio(c.used, c.canvas); r < 4.5 {
			t.Errorf("%s: the muted text colour is %.2f:1", c.name, r)
		}
	}

	if theme(true).TextMuted != darkTextSecondary {
		t.Error("the dark theme does not use the secondary grey for muted text")
	}
	if theme(false).TextMuted != lightTextSecondary {
		t.Error("the light theme does not use the secondary grey for muted text")
	}
}
