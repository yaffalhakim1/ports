package main

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"
)

// view builds the window from the app's state. It runs on the main thread
// for every frame, so it only reads state and asks for actions; the scan
// that fills that state runs elsewhere.
func (a *app) view(c *ui.Context) {
	a.now = c.Now()
	a.applyTheme(c)
	t := c.Theme()

	a.watchSort()
	a.shortcuts(c)

	ui.Column(c).Fill().Background(t.Background).Children(func() {
		a.header(c)
		ui.Divider(c).Background(t.Border)
		a.table(c)
		ui.Divider(c).Background(t.Border)
		a.footer(c)
	})

	a.confirm(c)
}

// watchSort notices a header click, which the table records in the sort
// order, and re-sorts the list to match. The table shows rows in the order
// the view gives them, so sorting is the view's work.
func (a *app) watchSort() {
	if a.order == a.appliedOrder {
		return
	}
	a.appliedOrder = a.order
	a.refresh()
}

// header is the window's top bar: what the list shows, the search, and the
// actions that apply to the whole list.
func (a *app) header(c *ui.Context) {
	t := c.Theme()
	ui.Column(c).Padding(pageMargin, 18).Gap(14).Children(func() {
		ui.Row(c).AlignItems(ui.Center).Gap(12).Children(func() {
			ui.Text(c, "Ports").FontSize(22).FontWeight(650).Shrink(0)
			ui.Text(c, a.summary()).TextColor(t.TextMuted).FontSize(13).
				SingleLine().Ellipsis("…").Shrink(1).MinWidth(0)
			ui.Spacer(c)
			a.appearanceButton(c)
			a.scanButton(c)
		})
		ui.Row(c).AlignItems(ui.Center).Gap(10).Children(func() {
			a.search(c)
			ui.Spacer(c)
			ui.Switch(c, &a.only).OnChange(a.refresh)
			ui.Text(c, "Exposed only").FontSize(13).TextColor(t.TextMuted).Tooltip("Show only ports another machine can reach")
		})
	})
}

// scanButton re-reads the system on demand and shows that a scan is
// running, so a manual refresh has a visible result even when nothing
// changed.
func (a *app) scanButton(c *ui.Context) {
	t := c.Theme()
	b := ui.ButtonBase(c).Padding(buttonPadV, buttonPadH).Radius(8).
		Focusable().Shrink(0).
		Tooltip("Scan now (R)").Disabled(a.scanning).
		OnClick(a.scanOnce).
		Background(t.Surface).Border(1, t.Border)
	b.Children(func() {
		ui.Row(c).AlignItems(ui.Center).Gap(iconLabelGap).Children(func() {
			icon := ui.Icon(c, iconRefresh).Size(15, 15).TextColor(t.TextMuted)
			if a.scanning {
				icon.Rotate(icon.Loop("spin", 900*time.Millisecond, ui.Linear) * 360)
			}
			ui.Text(c, "Scan").FontSize(13).TextColor(t.TextMuted)
		})
	})
	if b.Hovered() && !a.scanning {
		b.Background(t.SurfaceHover)
	}
}

// scanOnce starts a scan off the main thread, so the window stays
// responsive while it runs.
func (a *app) scanOnce() { go a.scanNow() }

// search is the query field, with the icon inside its leading edge.
func (a *app) search(c *ui.Context) {
	t := c.Theme()
	row := ui.Row(c).AlignItems(ui.Center).Gap(iconLabelGap).Width(300).
		Padding(buttonPadV, buttonPadH).Radius(8).Border(1, controlBorder(t.Dark)).Background(t.Surface)
	row.Children(func() {
		ui.Icon(c, iconSearch).Size(14, 14).TextColor(t.TextMuted)
		ui.TextInputBase(c, &a.query).
			Label("Search ports").Placeholder("Search ports, processes, folders").
			Bind(&a.searchBox).OnChange(a.refresh)
	})
	if row.Hovered() {
		row.Border(1, controlBorder(t.Dark))
	}
}

// statusText reports what just happened, which is where the app confirms
// an action or explains why one failed. It sits in the header rather than
// the footer so it can never push the footer's actions off screen.
func (a *app) statusText(c *ui.Context) {
	if a.status == "" {
		return
	}
	ui.Text(c, a.status).FontSize(12).TextColor(c.Theme().TextMuted).
		SingleLine().Ellipsis("…").MaxWidth(280).Shrink(1)
}

// summary says how many ports are listening and when they were read, which
// is the window's only claim about freshness.
func (a *app) summary() string {
	switch {
	case a.lastErr != nil:
		return "could not read the system"
	case a.lastScan.IsZero():
		return "reading…"
	}
	shown, total := len(a.shown), len(a.ports)
	count := fmt.Sprintf("%d listening", total)
	if shown != total {
		count = fmt.Sprintf("%d of %d shown", shown, total)
	}
	return count + " · updated " + ago(a.now.Sub(a.lastScan))
}

// ago writes a short elapsed time, as a status line reads it.
func ago(d time.Duration) string {
	switch {
	case d < time.Second:
		return "now"
	case d < time.Minute:
		return fmt.Sprintf("%ds ago", int(d.Seconds()))
	default:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	}
}

// shortcuts declares the window's keys. They are declared in the root
// scope, so they work wherever the focus is, and the field takes the focus
// when R is pressed so the next keystroke reaches the table.
func (a *app) shortcuts(c *ui.Context) {
	c.OnShortcut(0, ui.KeyR, func() {
		a.scanOnce()
		a.list.Handle.Focus(c)
	})
	c.OnShortcut(0, ui.KeyF, func() { a.focusSearch(c) })
}

// focusSearch moves the keyboard focus to the filter field. The handle is
// bound to the field in the view, and Focus waits until the field is
// built, so the key works from any frame.
func (a *app) focusSearch(c *ui.Context) { a.searchBox.Focus(c) }

// appearanceButton opens the appearance menu. It shows the icon and name of
// what the app is doing now, and the menu offers all three choices with the
// current one checked, so changing to the desktop's setting is one click
// rather than a cycle back through the others.
func (a *app) appearanceButton(c *ui.Context) {
	icon, label := a.appearanceIcon(), a.appearance.String()
	b := ui.ButtonBase(c).Padding(buttonPadV, buttonPadH).Radius(8).
		Focusable().Shrink(0).
		Background(c.Theme().Surface).Border(1, controlBorder(c.Theme().Dark)).
		Label("Appearance: " + label).
		Tooltip("Appearance: " + label).
		Menu(func(m *ui.Menu) {
			for _, choice := range []Appearance{AppearanceSystem, AppearanceLight, AppearanceDark} {
				if m.Item(choice.MenuLabel()).Checked(a.appearance == choice).Chosen() {
					a.setAppearance(choice)
				}
			}
		})
	b.Children(func() {
		ui.Row(c).AlignItems(ui.Center).Gap(iconLabelGap).Children(func() {
			ui.Icon(c, icon).Size(14, 14)
			ui.Text(c, label).FontSize(13).Shrink(0)
			// A small arrow, so the button reads as one that opens a menu
			// rather than one that acts.
			ui.Icon(c, iconCaret).Size(11, 11)
		})
	})
	if b.Hovered() {
		b.Background(c.Theme().SurfaceHover)
	}
}

// appearanceIcon is the icon of what the app is doing now: the desktop's
// screen when it follows the system, else the sun or the moon.
func (a *app) appearanceIcon() *ui.SVG {
	switch a.appearance {
	case AppearanceLight:
		return iconSun
	case AppearanceDark:
		return iconMoon
	default:
		return iconSystem
	}
}

// setAppearance chooses an appearance and remembers it, so the window opens
// the way it was left.
func (a *app) setAppearance(want Appearance) {
	if a.appearance == want {
		return
	}
	a.appearance = want
	a.saveSettings()
}
