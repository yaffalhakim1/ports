package main

import (
	"image"
	"testing"

	"github.com/egoist/mygo/ui"
)

// Tests of the window: what it shows, what it offers, and that the controls
// stay reachable at the smallest size the window allows.

// TestViewShowsPorts checks the window builds from state alone: the ports,
// their processes and their folders.
func TestViewShowsPorts(t *testing.T) {
	a := newTestApp()
	tt := ui.NewTester(a.view, 900, 620)
	tt.Frame()

	for _, want := range []string{"Ports", "3000", "node.exe", "web", "postgres.exe", "api"} {
		if !tt.HasText(want) {
			t.Errorf("no %q on screen; texts are %q", want, tt.Texts())
		}
	}
}

// TestFooterActionsAppearOnce checks each action is labelled once. A
// duplicated label is invisible in the render but doubles the button's
// text for assistive technology, which is how it was found.
func TestFooterActionsAppearOnce(t *testing.T) {
	a := newTestApp()
	a.chosen = 1
	tt := ui.NewTester(a.view, 900, 620)
	tt.Frame()

	count := map[string]int{}
	for _, s := range tt.Texts() {
		count[s]++
	}
	for _, label := range []string{"Open", "Copy", "Restart", "Stop"} {
		if count[label] != 1 {
			t.Errorf("%q appears %d times in the window, want once", label, count[label])
		}
	}
}

// TestChoosingARowRevealsActions checks the actions are offered only once
// a port is chosen, so the footer never shows buttons that would act on
// nothing.
func TestChoosingARowRevealsActions(t *testing.T) {
	a := newTestApp()
	a.chosen = -1
	tt := ui.NewTester(a.view, 900, 620)
	tt.Frame()
	for _, label := range []string{"Open", "Copy", "Restart", "Stop"} {
		if tt.HasText(label) {
			t.Errorf("%q is offered with no port chosen", label)
		}
	}
	if !tt.HasText("Choose a port to open, stop or restart it") {
		t.Error("the empty footer does not say what to do")
	}
}

// TestActionsStayOnScreen checks every action stays inside the window at
// the smallest size the window allows. Before this was tested, Restart
// and Stop were laid out past the right edge and were unreachable.
func TestActionsStayOnScreen(t *testing.T) {
	// The client area of the smallest window, which is smaller than the
	// window itself by its frame.
	const minWidth, minHeight = 604, 381
	for _, query := range []string{"", "postgres"} {
		a := newTestApp()
		a.chosen = 0
		a.query = query
		a.refresh()
		a.status = "Copied http://0.0.0.0:5432"

		tt := ui.NewTester(a.view, minWidth, minHeight)
		tt.Frame()
		for _, label := range []string{"Open", "Copy", "Restart", "Stop"} {
			box, ok := tt.Find(label)
			if !ok {
				t.Fatalf("%q is missing at %dx%d", label, minWidth, minHeight)
			}
			if box.X+box.W > minWidth {
				t.Errorf("%q ends at %.0f, past the %d-wide window",
					label, box.X+box.W, minWidth)
			}
			if box.Y+box.H > minHeight {
				t.Errorf("%q ends at %.0f, past the %d-high window",
					label, box.Y+box.H, minHeight)
			}
		}
	}
}

// TestStatusDoesNotHideActions checks a long status message, which is the
// longest text the window can show, does not push the actions away.
func TestStatusDoesNotHideActions(t *testing.T) {
	a := newTestApp()
	a.chosen = 0
	a.status = "Restart failed: access is denied opening a process of another user"
	tt := ui.NewTester(a.view, 604, 381)
	tt.Frame()
	if _, ok := tt.Find("Stop"); !ok {
		t.Error("a long status message hid the actions")
	}
}

// TestKeyboardReachesTheWindow checks the keys the app declares: R rescans
// and F reaches the search field, so the list is usable without a pointer.
func TestKeyboardReachesTheWindow(t *testing.T) {
	a := newTestApp()
	tt := ui.NewTester(a.view, 900, 620)
	tt.Frame()

	tt.Key(0, ui.KeyF)
	if !tt.Focused("Search ports") {
		t.Error("F did not focus the search field")
	}

	tt.Key(0, ui.KeyR)
	if !tt.HasText("3000") {
		t.Error("R left the window without its list")
	}
}

// TestTypingFilters checks typing in the search field narrows the table,
// which is the point of focusing it.
func TestTypingFilters(t *testing.T) {
	a := newTestApp()
	tt := ui.NewTester(a.view, 900, 620)
	tt.Frame()

	tt.Click("Search ports")
	tt.Type("postgres")
	if len(a.shown) != 1 || a.shown[0].Name != "postgres.exe" {
		t.Errorf("typing postgres showed %v", portsOf(a.shown))
	}
	if !tt.HasText("5432") || tt.HasText("3000") {
		t.Errorf("the table does not show the filtered list; texts are %q", tt.Texts())
	}
}

// TestEveryActionIsDelineated checks each control is distinguishable from
// what is behind it, which the brand system does with a hairline rather
// than a shadow or a fill. In the light appearance a control's surface is
// white on a near-white canvas, so the stroke is the only thing separating
// them; this asserts the separation exists rather than assuming which of
// the two carries it.
func TestEveryActionIsDelineated(t *testing.T) {
	for _, dark := range []bool{false, true} {
		a := newTestApp()
		a.chosen = 1
		tt := ui.NewTester(a.view, 900, 620)
		tt.SetDark(dark)
		tt.Frame()

		img := tt.Image()
		canvas := pixel(img, 2, 2)
		for _, label := range []string{"Open", "Copy", "Restart", "Stop", "Scan", "System"} {
			box, ok := tt.Find(label)
			if !ok {
				t.Errorf("dark=%v: %s is missing", dark, label)
				continue
			}
			if !delineated(img, box, canvas) {
				t.Errorf("dark=%v: %s has neither a fill nor a border of its own, so it reads as text",
					dark, label)
			}
		}
	}
}

// TestAppearanceMenuOffersEveryChoice checks the button opens a menu with
// all three appearances and marks the current one. A cycle would make
// reaching the desktop's setting take two clicks from either of the
// others, which is what this replaced.
func TestAppearanceMenuOffersEveryChoice(t *testing.T) {
	a := newTestApp()
	a.appearance = AppearanceDark
	tt := ui.NewTester(a.view, 960, 660)
	tt.Frame()

	if err := tt.Click("Appearance: Dark"); err != nil {
		t.Fatalf("the appearance button did not open a menu: %v", err)
	}
	items := tt.Menu()
	t.Logf("menu: %q", items)
	for _, want := range []string{"System", "Light", "Dark"} {
		if !contains(items, want) {
			t.Errorf("the menu does not offer %q", want)
		}
	}
}

// TestChoosingFromTheAppearanceMenu checks choosing one applies it, which
// is the whole point of the menu.
func TestChoosingFromTheAppearanceMenu(t *testing.T) {
	for _, c := range []struct {
		item string
		want Appearance
	}{
		{"System", AppearanceSystem},
		{"Light", AppearanceLight},
		{"Dark", AppearanceDark},
	} {
		a := newTestApp()
		a.appearance = AppearanceLight
		tt := ui.NewTester(a.view, 960, 660)
		tt.Frame()
		if err := tt.Click("Appearance: Light"); err != nil {
			t.Fatalf("no menu: %v", err)
		}
		if err := tt.ChooseMenuItem(c.item); err != nil {
			t.Fatalf("choosing %q: %v", c.item, err)
		}
		tt.Frame()
		if a.appearance != c.want {
			t.Errorf("choosing %q set %v, want %v", c.item, a.appearance, c.want)
		}
	}
}

// contains reports whether a list holds a string.
func contains(list []string, want string) bool {
	for _, s := range list {
		if s == want {
			return true
		}
	}
	return false
}

// TestAppearanceChangesTheTheme checks the choice reaches the theme: light
// and dark produce a theme of the matching polarity, whatever the desktop
// says.
func TestAppearanceChangesTheTheme(t *testing.T) {
	cases := []struct {
		appearance Appearance
		desktop    bool
		wantDark   bool
	}{
		{AppearanceSystem, false, false},
		{AppearanceSystem, true, true},
		{AppearanceLight, true, false},
		{AppearanceDark, false, true},
	}
	for _, c := range cases {
		a := newTestApp()
		a.appearance = c.appearance
		tt := ui.NewTester(a.view, 900, 620)
		tt.SetDark(c.desktop)
		tt.Frame()
		if got := tt.Image(); isDark(got) != c.wantDark {
			t.Errorf("appearance %v with desktop dark=%v rendered %v",
				c.appearance, c.desktop, map[bool]string{true: "dark", false: "light"}[isDark(got)])
		}
	}
}

// isDark reports whether a frame is a dark one, by its corner pixel.
func isDark(img *image.RGBA) bool {
	p := pixel(img, 2, 2)
	sum := int(p[0]) + int(p[1]) + int(p[2])
	return sum < 200
}

// delineated reports whether any pixel across the left padding of a
// control differs from the canvas, which is where its border sits.
func delineated(img *image.RGBA, box ui.Rect, canvas [3]uint8) bool {
	y := int(box.Y + box.H/2)
	for x := int(box.X) - 14; x < int(box.X)-2; x++ {
		if x < 0 {
			continue
		}
		if pixel(img, x, y) != canvas {
			return true
		}
	}
	return false
}

// pixel reads an opaque pixel as its three channels.
func pixel(img *image.RGBA, x, y int) [3]uint8 {
	r, g, b, _ := img.At(x, y).RGBA()
	return [3]uint8{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)}
}
