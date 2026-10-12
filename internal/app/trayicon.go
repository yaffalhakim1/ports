package app

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math"
)

// trayIconPNG draws the tray icon: the app's socket mark — a bold ring with
// a gap and the live dot in it — as a black template icon on transparent,
// so the system recolours it for a light or dark menu bar.
//
// It is drawn rather than embedded so the icon stays in the repository as
// code, with no binary asset to keep in step.
func trayIconPNG() []byte {
	const size = 32 // 16 points at 2x
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	ink := color.NRGBA{R: 0, G: 0, B: 0, A: 255}

	// The mark, in the same proportions as the master: a ring (outer r=84,
	// inner r=44) with a 70° gap at the upper right, and the dot (r=20)
	// centred in the gap. Scaled to fit the 32px tile with a small margin.
	const (
		outer = 11.0 // ring outer radius
		inner = 5.8  // ring inner radius
		dotR  = 2.6  // the live dot
		cx    = 15.0
		cy    = 15.5
	)
	// The gap runs from -80° to -10°; the dot sits on its centre line.
	gapFrom, gapTo := -80.0, -10.0
	dotX := cx + math.Cos(-45*math.Pi/180)*((outer+inner)/2)
	dotY := cy + math.Sin(-45*math.Pi/180)*((outer+inner)/2)

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			// Sample at the pixel centre, and supersample 3x3 for an edge
			// that reads cleanly at 16 points.
			if coverage(x, y, cx, cy, outer, inner, dotX, dotY, dotR, gapFrom, gapTo) > 0.5 {
				img.Set(x, y, ink)
			}
		}
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil
	}
	return buf.Bytes()
}

// coverage returns the fraction of a pixel covered by the mark, so the ring
// and dot have smooth edges.
func coverage(px, py int, cx, cy, outer, inner, dotX, dotY, dotR, gapFrom, gapTo float64) float64 {
	const samples = 3
	hits := 0
	for sy := 0; sy < samples; sy++ {
		for sx := 0; sx < samples; sx++ {
			x := float64(px) + (float64(sx)+0.5)/samples
			y := float64(py) + (float64(sy)+0.5)/samples
			if inMark(x, y, cx, cy, outer, inner, dotX, dotY, dotR, gapFrom, gapTo) {
				hits++
			}
		}
	}
	return float64(hits) / (samples * samples)
}

// inMark reports whether a point is inside the ring (minus its gap) or the dot.
func inMark(x, y, cx, cy, outer, inner, dotX, dotY, dotR, gapFrom, gapTo float64) bool {
	if (x-dotX)*(x-dotX)+(y-dotY)*(y-dotY) <= dotR*dotR {
		return true
	}
	d := math.Hypot(x-cx, y-cy)
	if d < inner || d > outer {
		return false
	}
	// Angles measured like the master: 0° is east, negative is up.
	ang := math.Atan2(y-cy, x-cx) * 180 / math.Pi
	return ang < gapFrom || ang > gapTo
}
