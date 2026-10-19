package app

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"

	"ports/internal/port"
)

// footerGap is the space between the actions, and between the text and the
// first of them: one step of the eight-pixel grid the window is laid out on.
const footerGap = 8

// footer is the action bar: what is chosen, and what can be done with it.
//
// The left side grows and shrinks, so its text truncates when the window is
// narrow, while the actions keep their size and stay on screen. The summary
// of the list lives in the titlebar for the same reason: it is the longest
// text and would otherwise push the actions out of the window.
func (a *App) footer(c *ui.Context) {
	t := c.Theme()
	ui.Column(c).Shrink(0).Children(func() {
		ui.Box(c).FillWidth().Height(1).Background(t.Border)
		ui.Row(c).AlignItems(ui.Center).Gap(footerGap).Padding(footerPadV, pageMargin).Children(func() {
			p, ok := a.selected()
			if !ok {
				components.Text(c, "Choose a port to open, stop or restart it").
					FontSize(t.FontSize - 1).TextColor(t.TextMuted).
					SingleLine().Ellipsis("…").Grow(1).Shrink(1).MinWidth(0)
				return
			}
			a.footerText(c, p)
			a.actions(c, p)
		})
	})
}

// footerText is what just happened, or what is chosen, in the space the
// actions leave. An action's result belongs beside the action, and the
// titlebar already carries how fresh the list is.
func (a *App) footerText(c *ui.Context, p port.Port) {
	t := c.Theme()
	if a.status != "" {
		components.Text(c, a.status).FontSize(t.FontSize - 1).TextColor(t.TextMuted).
			SingleLine().Ellipsis("…").Grow(1).Shrink(1).MinWidth(0)
		return
	}
	ui.Row(c).Grow(1).Shrink(1).MinWidth(0).AlignItems(ui.Center).Gap(t.Space(2)).Children(func() {
		components.Text(c, p.Name).FontSize(t.FontSize - 1).Bold().
			SingleLine().Ellipsis("…").Shrink(1).MinWidth(0)
		components.Text(c, "·").FontSize(t.FontSize - 1).TextColor(t.TextMuted).Shrink(0)
		components.Text(c, p.URL()).FontSize(t.FontSize - 1).TextColor(t.TextMuted).
			SingleLine().Ellipsis("…").Shrink(1).MinWidth(0)
	})
}

// action is one thing the action bar can do with the chosen port: its icon,
// its word, what its tooltip says, and whether it is the destructive one.
type action struct {
	icon     func() *ui.SVG
	label    string
	tip      string
	onClick  func()
	disabled bool
	danger   bool
}

// actions are what a chosen port offers. They never shrink, so a narrow
// window truncates the text beside them rather than hiding them. Stop and
// restart are unavailable for a process that cannot be ended, with the
// reason on the button's tooltip.
func (a *App) actions(c *ui.Context, p port.Port) {
	blocked := p.SystemPID()

	a.actionButton(c, action{
		icon:    lucide.ExternalLink,
		label:   "Open",
		tip:     disabledReason("Open", p, false),
		onClick: func() { a.open(c, p) },
	})
	a.actionButton(c, action{
		icon:    lucide.Copy,
		label:   "Copy",
		tip:     disabledReason("Copy", p, false),
		onClick: func() { a.copy(c, p) },
	})
	a.actionButton(c, action{
		icon:     lucide.RotateCcw,
		label:    "Restart",
		tip:      disabledReason("Restart", p, blocked),
		disabled: blocked,
		onClick:  func() { a.ask("Restart", p) },
	})
	a.actionButton(c, action{
		icon:     lucide.Square,
		label:    "Stop",
		tip:      disabledReason("Stop", p, blocked),
		disabled: blocked,
		danger:   true,
		onClick:  func() { a.ask("Stop", p) },
	})
}

// actionButton builds one action from Onevra's outline variant, so its hit
// area, its radius, its hover and pressed faces all come from the theme
// rather than from numbers written here.
//
// The label is left to the button's own text rather than being set as an
// accessible name as well: MyGo reads a control's name from its text, so
// setting both would make a screen reader say the word twice.
func (a *App) actionButton(c *ui.Context, act action) {
	t := c.Theme()
	b := components.Button(c, "", components.ButtonOutline).
		Shrink(0).OnClick(act.onClick)
	b.Children(func() {
		ui.Row(c).AlignItems(ui.Center).Gap(t.Space(1.5)).Children(func() {
			// The icon inherits the button's colour, so recolouring the
			// button below recolours its icon with its label.
			components.Icon(c, act.icon).Size(14, 14)
			components.Text(c, act.label).FontSize(t.FontSize - 1).Shrink(0)
		})
	})
	components.Tooltip(c, b, act.tip)

	if act.disabled {
		// Onevra dims a disabled control while styling it, which is before
		// this action knows it is disabled. The dimming is applied here
		// instead, at the same weight.
		if act.danger {
			b.TextColor(t.Danger)
		}
		b.Disabled(true).Opacity(disabledOpacity)
		return
	}
	if act.danger {
		dangerFace(c, b)
		return
	}
	b.TextColor(t.TextMuted)
}

// dangerFace marks a destructive control in the danger colour, with the
// pressed and hover tints of its own, so pressing it looks like pressing it
// rather than like touching a neutral button. The wording carries the same
// meaning, so the colour is never the only signal.
func dangerFace(c *ui.Context, b ui.Element) {
	t := c.Theme()
	b.TextColor(t.Danger)
	switch {
	case b.Pressed():
		b.Background(t.Danger.Alpha(0.20))
	case b.Hovered():
		b.Background(t.Danger.Alpha(0.12))
	}
}

// disabledOpacity is how far a control fades when it cannot be used, which
// is the weight Onevra's own controls fade to.
const disabledOpacity = 0.45

// disabledReason is the tooltip of an action: why it is unavailable, or
// exactly what it will do.
func disabledReason(verb string, p port.Port, blocked bool) string {
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
