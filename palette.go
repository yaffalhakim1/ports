package main

import "github.com/egoist/mygo/ui"

// The app's colour system, from the waku brand palette: "neutral graphite"
// surfaces with hue reserved for meaning. The desktop app is canonical
// here; the values are the same ones the web and mobile clients use, so
// the three stay one product.
//
// Surfaces are near-greyscale, so colour always means something. The brand
// green is the accent and carries no structure beyond the chosen row and
// the focus ring, which is where a desktop app needs one.
var (
	// Dark: graphite canvas, ink text.
	darkCanvas   = ui.Hex("#1A1A1A")
	darkRaised   = ui.Hex("#232323")
	darkComposer = ui.Hex("#212121")
	darkInset    = ui.Hex("#151515")
	// The hairlines: white 7% and 14% over the canvas, resolved.
	darkBorder       = ui.Hex("#2A2A2A")
	darkBorderStrong = ui.Hex("#3A3A3A")

	darkText          = ui.Hex("#E2E2E2")
	darkTextSecondary = ui.Hex("#A3A3A3")
	darkTextTertiary  = ui.Hex("#7D7D7D")
	darkTextGhost     = ui.Hex("#575757")

	// Light: paper canvas, ink text.
	lightCanvas   = ui.Hex("#F6F5F6")
	lightRaised   = ui.Hex("#ECECEC")
	lightComposer = ui.Hex("#FFFFFF")
	lightInset    = ui.Hex("#E6E6E6")
	// hsla(220,10%,12%,.08) and .15 over the canvas, resolved.
	lightBorder       = ui.Hex("#E5E4E5")
	lightBorderStrong = ui.Hex("#D5D5D7")

	lightText          = ui.Hex("#242424")
	lightTextSecondary = ui.Hex("#666666")
	lightTextTertiary  = ui.Hex("#858585")
	lightTextGhost     = ui.Hex("#A4A4A4")

	// The brand green, and the near-black that sits on it. The green is
	// bright in both appearances, so its text is dark in both: white on it
	// measures under 2:1 and would be unreadable.
	brandGreen      = ui.Hex("#26E085")
	brandGreenLight = ui.Hex("#00DA7D")
	onBrandGreen    = ui.Hex("#17181C")

	// Status colours, the only other saturated hues.
	statusSuccessDark  = ui.Hex("#62C987")
	statusSuccessLight = ui.Hex("#2F8F52")
	statusWarningDark  = ui.Hex("#E0B36A")
	statusWarningLight = ui.Hex("#A66B20")
	statusDangerDark   = ui.Hex("#E2726A")
	statusDangerLight  = ui.Hex("#C64A42")

	// Primary buttons are the inverse surface, not the accent, as the
	// brand system has it: the accent stays non-structural.
	inverseDark    = ui.Hex("#E7E9EC")
	onInverseDark  = ui.Hex("#17181C")
	inverseLight   = ui.Hex("#202227")
	onInverseLight = ui.Hex("#F8F8F9")
)

// Derived values, each one measured rather than chosen by eye. The tests
// in brand_contrast_test.go assert every pair these take part in.

var (
	// brandGreenLightDark is the brand's own inline-code green, reused as
	// the focus ring on the light canvas: the light brand green measures
	// 1.71:1 there, and a focus indicator needs 3:1.
	brandGreenLightDark = ui.Hex("#00814A")

	// The light status colours as text, darkened from the brand's fills
	// until each passes 4.5:1 on the light canvas at 13px.
	statusSuccessLightText = ui.Hex("#1E7A40")
	statusWarningLightText = ui.Hex("#8A5714")
	statusDangerLightText  = ui.Hex("#B03A33")
)
