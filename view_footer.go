package main

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"
)

// buttonPadV and buttonPadH size an action so its hit area is 32 DIPs
// tall, which clears the 24x24 floor of WCAG 2.5.8 with room for the
// focus ring. The label's own box is only as tall as its text.
const (
	buttonPadV = 7
	buttonPadH = 12
)

// footer is the action bar: what is chosen, and what can be done with it.
//
// The left side grows and shrinks, so its text truncates when the window
// is narrow, while the actions keep their size and stay on screen. The
// status message lives in the header for the same reason: it is the
// longest text and would otherwise push the actions out of the window.
func (a *app) footer(c *ui.Context) {
	t := c.Theme()
	ui.Row(c).AlignItems(ui.Center).Gap(footerGap).Padding(pageMargin, 14).Children(func() {
		p, ok := a.selected()
		if !ok {
			ui.Text(c, "Choose a port to open, stop or restart it").
				FontSize(13).TextColor(t.TextMuted).SingleLine().Ellipsis("…").
				Grow(1).Shrink(1).MinWidth(0)
			return
		}

		// What just happened, or what is chosen, in the space the actions
		// leave. An action's result belongs beside the action, and the
		// header already carries how fresh the list is.
		if a.status != "" {
			ui.Text(c, a.status).FontSize(12).TextColor(t.TextMuted).
				SingleLine().Ellipsis("…").Grow(1).Shrink(1).MinWidth(0)
		} else {
			ui.Row(c).Grow(1).Shrink(1).MinWidth(0).
				AlignItems(ui.Center).Gap(8).Children(func() {
				ui.Text(c, p.Name).FontSize(13).FontWeight(600).
					SingleLine().Ellipsis("…").Shrink(1).MinWidth(0)
				ui.Text(c, "·").FontSize(13).TextColor(t.TextMuted).Shrink(0)
				ui.Text(c, p.URL()).FontSize(13).TextColor(t.TextMuted).
					SingleLine().Ellipsis("…").Shrink(1).MinWidth(0)
			})
		}

		a.actions(c, p)
	})
}

// actions are what a chosen port offers. They never shrink, so a narrow
// window truncates the text beside them rather than hiding them. Stop and
// restart are disabled for a process that cannot be ended, with the
// reason on the button's tooltip.
func (a *app) actions(c *ui.Context, p Port) {
	blocked := p.systemPID()

	iconButton(c, iconOpen, "Open", disabledReason("Open", p, false), func() { a.open(c, p) })
	iconButton(c, iconCopy, "Copy", disabledReason("Copy", p, false), func() { a.copy(c, p) })

	restart := iconButton(c, iconRestart, "Restart", disabledReason("Restart", p, blocked),
		func() { a.ask("Restart", p) })
	restart.Disabled(blocked)

	stop := iconButton(c, iconStop, "Stop", disabledReason("Stop", p, blocked),
		func() { a.ask("Stop", p) })
	stop.Disabled(blocked)
	// A destructive action is marked in its own colour and by its wording,
	// so the meaning is never carried by the colour alone.
	stop.TextColor(c.Theme().Danger)
	if stop.Hovered() && !blocked {
		stop.Background(c.Theme().Danger.Alpha(0.12))
	}
}

// iconButton is one action in the footer: an icon and a label in a button
// that takes the keyboard focus.
//
// Every action carries a border and a surface, so it reads as a control
// rather than as the static text beside it; only the primary action gets
// a fill. The padding is what gives the button a hit area, so it is set
// here rather than left to the label's own box.
func iconButton(c *ui.Context, icon *ui.SVG, label, tip string, onClick func()) ui.Element {
	t := c.Theme()
	b := ui.ButtonBase(c).Padding(buttonPadV, buttonPadH).Radius(8).
		Focusable().Shrink(0).
		Background(t.Surface).Border(1, controlBorder(t.Dark)).
		Tooltip(tip).OnClick(onClick).
		// The icon inherits this, so a button that overrides its colour
		// (the destructive one) recolours its icon with its label.
		TextColor(t.TextMuted)
	// A row with a gap, so the icon does not touch its label.
	b.Children(func() {
		ui.Row(c).AlignItems(ui.Center).Gap(iconLabelGap).Children(func() {
			ui.Icon(c, icon).Size(14, 14)
			ui.Text(c, label).FontSize(13).Shrink(0)
		})
	})
	if b.Hovered() {
		b.Background(t.SurfaceHover)
	}
	return b
}

// disabledReason is the tooltip of an action: why it is unavailable, or
// exactly what it will do.
func disabledReason(verb string, p Port, blocked bool) string {
	if blocked {
		return "A system process cannot be " + strings.ToLower(verb) + "ed"
	}
	switch verb {
	case "Open":
		return "Open " + p.URL() + " in the browser"
	case "Copy":
		return "Copy " + p.URL()
	}
	return verb + " " + p.Name + " (PID " + fmt.Sprint(p.PID) + ")"
}

// Spacing constants for the footer. pageMargin is the window's horizontal
// margin, which the header, the table and the footer all share, so the
// right edge of the last action lines up with the table's left edge.
const (
	pageMargin   = 20
	footerGap    = 12
	iconLabelGap = 8
)
