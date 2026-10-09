package main

import "github.com/egoist/mygo/ui"

// theme builds the theme for an appearance: the brand's neutral graphite
// surfaces with the brand green as the one accent.
//
// Three values depart from the brand hexes because the roles here are not
// the roles they were chosen for, and the brand's own value fails the
// requirement in this one:
//
//   - The focus ring is the brand green in the dark appearance, but the
//     light brand green measures 1.71:1 on the light canvas, far under the
//     3:1 a focus indicator needs. The light ring is the darker green the
//     brand already uses for inline code.
//   - Status colours are text here rather than diff fills, so the light
//     ones are darkened until they pass 4.5:1 as 13px text.
//   - The brand's tertiary and ghost greys measure 4.23:1 and 2.41:1 on
//     the canvas, so this app uses the brand's secondary grey for all
//     muted text and leaves the quieter steps for decoration.
//
// The desktop's own accent is deliberately not adopted. The brand system
// makes the green the one accent, and no text colour reaches 4.5:1 on
// every accent a desktop can supply, so following it would mean shipping
// an unreadable chosen row on some machines.
func theme(dark bool) *ui.Theme {
	t := &ui.Theme{Dark: dark}
	if dark {
		t.Background = darkCanvas
		t.Surface = darkRaised
		t.SurfaceHover = darkComposer
		t.SurfacePressed = darkInset
		t.Border = darkBorder
		t.Text = darkText
		t.TextMuted = darkTextSecondary
		t.Selection = brandGreen.Alpha(0.30)
		t.Focus = brandGreen
		t.Inverse = inverseDark
		t.InverseText = onInverseDark
		t.Scrollbar = darkBorderStrong
		t.Warning = statusWarningDark
		t.Success = statusSuccessDark
		t.Danger = statusDangerDark

		t.Accent = brandGreen
		t.AccentHover = brandGreen
		t.AccentPressed = brandGreen
	} else {
		t.Background = lightCanvas
		t.Surface = lightComposer
		t.SurfaceHover = lightRaised
		t.SurfacePressed = lightInset
		t.Border = lightBorder
		t.Text = lightText
		t.TextMuted = lightTextSecondary
		t.Selection = brandGreenLight.Alpha(0.30)
		t.Focus = brandGreenLightDark
		t.Inverse = inverseLight
		t.InverseText = onInverseLight
		t.Scrollbar = lightBorderStrong
		t.Warning = statusWarningLightText
		t.Success = statusSuccessLightText
		t.Danger = statusDangerLightText

		t.Accent = brandGreenLight
		t.AccentHover = brandGreenLight
		t.AccentPressed = brandGreenLightDark
	}
	// The brand green is bright in both appearances, so the text on it is
	// the brand's near-black in both: white on it measures under 2:1.
	t.AccentText = onBrandGreen
	t.Radius = 8
	return t
}

// controlBorder is the stroke around a control, which is stronger than a
// divider: a control has to read as a control, while a divider only has to
// separate. The brand system draws hairlines rather than shadows, so this
// stroke is the only thing that makes a button distinct from the surface
// behind it.
func controlBorder(dark bool) ui.Color {
	if dark {
		return darkBorderStrong
	}
	return lightBorderStrong
}
