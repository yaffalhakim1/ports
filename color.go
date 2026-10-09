package main

import (
	"math"

	"github.com/egoist/mygo/ui"
)

// Contrast measurement, shared by the theme and the tests so both agree on
// what a passing pair is. Contrast is one of the few interface concerns
// with an exact answer, so it is computed rather than judged.

// lum is WCAG relative luminance.
func lum(c ui.Color) float64 {
	channel := func(v uint8) float64 {
		s := float64(v) / 255
		if s <= 0.04045 {
			return s / 12.92
		}
		return math.Pow((s+0.055)/1.055, 2.4)
	}
	return 0.2126*channel(c.R) + 0.7152*channel(c.G) + 0.0722*channel(c.B)
}

// ratio is the WCAG contrast ratio between two opaque colours.
func ratio(a, b ui.Color) float64 {
	la, lb := lum(a), lum(b)
	if la < lb {
		la, lb = lb, la
	}
	return (la + 0.05) / (lb + 0.05)
}
