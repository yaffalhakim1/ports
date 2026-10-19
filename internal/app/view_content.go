package app

import (
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
	"github.com/oneo-me/onevra/source/shells/sidebar"
)

// pageMargin is the window's horizontal margin, which the toolbar, the list
// and the action bar all share, so their two edges line up. It and the two
// vertical paddings are multiples of eight, so the whole window sits on one
// grid.
const (
	pageMargin  = 24
	toolbarPadV = 16
	footerPadV  = 16
	contentGapV = 16
)

// content builds what the app owns inside the frame: the toolbar over the
// list, the list, and the action bar under it. The frame around it — the
// titlebar, the sidebar and the window controls — belongs to the shell, and
// is built by titlebar and by side.
//
// The shell's header floats over the full-height regions and the shell does
// not inset a region's builder, so the content makes room for it here. Without
// this the toolbar would sit under the titlebar and no control in it could be
// clicked.
//
// It runs on the main thread for every frame, so it only reads state and
// asks for actions; the scan that fills that state runs elsewhere.
func (a *App) content(c *ui.Context) {
	a.flushToast(c)

	ui.Column(c).Fill().MinWidth(0).
		Padding(sidebar.HeaderHeight(c), 0, 0, 0).Children(func() {
		a.toolbar(c)
		a.table(c)
		a.footer(c)
		a.confirm(c)
		a.palette(c)
	})
}

// toolbar is the row above the list: the search, and the actions that apply
// to the whole list rather than to a chosen row. What the list is showing —
// the count and its freshness — is the titlebar's job, next to the title.
func (a *App) toolbar(c *ui.Context) {
	ui.Column(c).Shrink(0).Padding(toolbarPadV, pageMargin).Gap(contentGapV).Children(func() {
		ui.Row(c).AlignItems(ui.Center).Gap(12).Children(func() {
			components.TextBox(c, &a.query, components.TextBoxOptions{
				// The magnifier is the field's own affordance, so the
				// placeholder does not have to say "search" twice.
				LeadingIcon: lucide.Search(),
			}).Label("Search ports").
				Placeholder("Ports, processes, folders").
				Bind(&a.searchBox).OnChange(a.refresh).
				Width(320).Shrink(0)
			ui.Spacer(c)
			a.scanButton(c)
			a.appearanceButton(c)
		})
		// Onevra leaves dividers out on purpose: a hairline is a box the
		// theme colours, and a widget for it would be a widget for nothing.
		ui.Box(c).FillWidth().Height(1).Background(c.Theme().Border)
	})
}

// scanButton re-reads the system on demand and shows that a scan is
// running, so a manual refresh has a visible result even when nothing
// changed.
func (a *App) scanButton(c *ui.Context) {
	t := c.Theme()
	b := components.Button(c, "", components.ButtonOutline).Shrink(0).
		Disabled(a.scanning).OnClick(a.scanOnce)
	b.Children(func() {
		ui.Row(c).AlignItems(ui.Center).Gap(t.Space(1.5)).Children(func() {
			icon := components.Icon(c, lucide.RefreshCw).Size(14, 14)
			if a.scanning {
				icon.Rotate(icon.Loop("spin", 900*time.Millisecond, ui.Linear) * 360)
			}
			components.Text(c, "Scan").FontSize(t.FontSize - 1).Shrink(0)
		})
	})
	components.Tooltip(c, b, "Scan now (R)")
}

// scanOnce starts a scan off the main thread, so the window stays
// responsive while it runs.
func (a *App) scanOnce() { go a.scanNow() }

// appearanceButton opens the appearance menu. It shows the icon and name of
// what the app is doing now, and the menu offers all three choices with the
// current one checked, so changing to the desktop's setting is one click
// rather than a cycle back through the others.
func (a *App) appearanceButton(c *ui.Context) {
	t := c.Theme()
	label := a.appearance.String()
	b := components.Button(c, "", components.ButtonOutline).Shrink(0).
		// The accessible name says what the control is as well as what it
		// shows, because the visible word alone ("Dark") does not say what
		// it would change.
		Label("Appearance: " + label).
		Menu(func(m *ui.Menu) {
			for _, choice := range []Appearance{AppearanceSystem, AppearanceLight, AppearanceDark} {
				if m.Item(choice.MenuLabel()).Checked(a.appearance == choice).Chosen() {
					a.setAppearance(choice)
				}
			}
		})
	b.Children(func() {
		ui.Row(c).AlignItems(ui.Center).Gap(t.Space(1.5)).Children(func() {
			components.Icon(c, appearanceIcon(a.appearance)).Size(14, 14)
			components.Text(c, label).FontSize(t.FontSize - 1).SingleLine().Shrink(0)
			// A caret, so the button reads as one that opens a menu rather
			// than one that acts.
			components.Icon(c, lucide.ChevronDown).Size(12, 12).TextColor(t.TextMuted)
		})
	})
	components.Tooltip(c, b, "Appearance: "+label)
}

// appearanceIcon is what the app is doing now: the desktop's screen when it
// follows the system, else the sun or the moon. It is returned as an
// accessor, which is what the icon component takes.
func appearanceIcon(a Appearance) func() *ui.SVG {
	switch a {
	case AppearanceLight:
		return lucide.Sun
	case AppearanceDark:
		return lucide.Moon
	default:
		return lucide.Monitor
	}
}

// setAppearance chooses an appearance and remembers it, so the window opens
// the way it was left. The desktop's own chrome follows through the theme
// source, which also covers the menus the window opens.
func (a *App) setAppearance(want Appearance) {
	if a.appearance == want {
		return
	}
	a.appearance = want
	a.saveSettings()
	a.applySource()
}

// applySource points MyGo's theme source at the chosen appearance. SetSource
// waits for the main loop, so it cannot be called from a frame build; the
// request goes through the window, which runs it on that loop. A test builds
// frames without a loop or a window, so there is nothing to point.
func (a *App) applySource() {
	if a.win == nil {
		return
	}
	a.win.Update(func() { mygo.Theme.SetSource(a.appearance.Source()) })
}

// flushToast shows a toast the app asked for. A toast can only be added
// while a frame is being built, and the actions that raise one run off the
// main thread, so they queue it and this emits it on the next frame.
func (a *App) flushToast(c *ui.Context) {
	if a.toast == nil {
		return
	}
	c.AddToast(*a.toast)
	a.toast = nil
}

// say queues an outcome to show as a toast. The status line in the action
// bar already reports the same thing, so this is the notification half of
// it: the toast is for when the eye is on the list, not the bar.
func (a *App) say(title, detail, kind string) {
	a.toast = &ui.Toast{Title: title, Description: detail, Type: kind}
}

// focusSearch moves the keyboard focus to the filter field. The handle is
// bound to the field in the view, and Focus waits until the field is
// built, so the key works from any frame.
func (a *App) focusSearch(c *ui.Context) { a.searchBox.Focus(c) }
