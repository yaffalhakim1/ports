package main

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
)

// trayIconPNG draws the tray icon: a small stack of horizontal bars, which
// reads as a list of ports at 16 points and stays legible when macOS tints
// it. It is black on transparent because a template icon is what lets the
// system recolour it for a light or dark menu bar.
//
// It is drawn rather than embedded so the icon stays in the repository as
// code, with no binary asset to keep in step.
func trayIconPNG() []byte {
	const size = 32 // 16 points at 2x
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	ink := color.NRGBA{R: 0, G: 0, B: 0, A: 255}

	// Three bars, the middle one shorter, like rows of a list. Each is
	// rounded by drawing a slightly inset rectangle with its ends capped.
	bars := []struct{ y, h, x0, x1 int }{
		{6, 4, 6, 26},
		{14, 4, 6, 20},
		{22, 4, 6, 23},
	}
	for _, b := range bars {
		fillRoundRect(img, b.x0, b.y, b.x1, b.y+b.h, b.h/2, ink)
	}
	// A dot on the right, which reads as a live listener.
	fillCircle(img, 26, 24, 3, ink)

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil
	}
	return buf.Bytes()
}

// fillRoundRect fills a rectangle with rounded ends.
func fillRoundRect(img *image.RGBA, x0, y0, x1, y1, r int, c color.NRGBA) {
	fillRect(img, x0+r, y0, x1-r, y1, c)
	fillCircle(img, x0+r, (y0+y1)/2, r, c)
	fillCircle(img, x1-r, (y0+y1)/2, r, c)
}

func fillRect(img *image.RGBA, x0, y0, x1, y1 int, c color.NRGBA) {
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			if image.Pt(x, y).In(img.Rect) {
				img.Set(x, y, c)
			}
		}
	}
}

func fillCircle(img *image.RGBA, cx, cy, r int, c color.NRGBA) {
	for y := cy - r; y <= cy+r; y++ {
		for x := cx - r; x <= cx+r; x++ {
			dx, dy := x-cx, y-cy
			if dx*dx+dy*dy <= r*r && image.Pt(x, y).In(img.Rect) {
				img.Set(x, y, c)
			}
		}
	}
}
