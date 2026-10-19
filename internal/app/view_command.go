package app

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/components/command"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

// palettePrefix namespaces a command's ID, so one prefix decides what
// activating it does.
const (
	portCommand = "port:"
	viewCommand = "view:"
	appCommand  = "app:"
)

// palette is the command palette: every listening port the last scan found,
// then the app's own commands. It is built into the frame, so it is called
// unconditionally and shows itself only while open.
//
// The ports come from the whole scan rather than the list on screen: the
// reason to reach for the keyboard is usually that the port being looked for
// is hidden behind the search or the sidebar.
func (a *App) palette(c *ui.Context) {
	var items []command.Item
	if a.paletteOpen {
		items = a.paletteItems()
	}
	command.CommandDialog(c, &a.paletteOpen, &a.paletteState, items,
		command.Options{
			SearchLabel: "Search ports and commands",
			Placeholder: "Ports or commands",
			OnActivate:  a.runCommand,
		},
		components.DialogOptions{Label: "Commands"})
}

// paletteItems builds the palette's list from the current scan. Titles lead
// with the port number because that is what a person remembers; the process
// name and the folder are what they search by.
func (a *App) paletteItems() []command.Item {
	items := make([]command.Item, 0, len(a.ports)+7)
	for _, p := range a.ports {
		items = append(items, command.Item{
			ID:          portCommand + p.Key(),
			Title:       fmt.Sprintf("%d · %s", p.Port, p.Name),
			Description: p.URL(),
			Group:       "Ports",
			Icon:        lucide.Server(),
			Keywords:    []string{p.Addr, p.Folder(), p.Dir, fmt.Sprint(p.PID), p.Proto},
		})
	}
	for _, m := range modes {
		items = append(items, command.Item{
			ID:          viewCommand + m.String(),
			Title:       m.Label(),
			Description: "Show " + strings.ToLower(m.Label()),
			Group:       "View",
			Icon:        modeIcon(m)(),
		})
	}
	for _, choice := range []Appearance{AppearanceSystem, AppearanceLight, AppearanceDark} {
		items = append(items, command.Item{
			ID:          appCommand + choice.String(),
			Title:       choice.String(),
			Description: "Appearance",
			Group:       "Appearance",
			Icon:        appearanceIcon(choice)(),
		})
	}
	if a.folder != "" {
		items = append(items, command.Item{
			ID:          viewCommand + "folders",
			Title:       "Show every folder",
			Description: "Clear the folder filter",
			Group:       "View",
			Icon:        lucide.Folder(),
		})
	}
	return append(items, command.Item{
		ID:          "scan",
		Title:       "Scan now",
		Description: "Re-read listening ports",
		Group:       "App",
		Icon:        lucide.RefreshCw(),
		Shortcut:    "R",
	})
}

// runCommand carries out the palette's chosen row.
func (a *App) runCommand(id string) {
	switch {
	case strings.HasPrefix(id, portCommand):
		a.jumpTo(strings.TrimPrefix(id, portCommand))
	case id == viewCommand+"folders":
		a.setView(ModeAll, "")
	case strings.HasPrefix(id, viewCommand):
		for _, m := range modes {
			if id == viewCommand+m.String() {
				a.setView(m, "")
				return
			}
		}
	case strings.HasPrefix(id, appCommand):
		for _, choice := range []Appearance{AppearanceSystem, AppearanceLight, AppearanceDark} {
			if id == appCommand+choice.String() {
				a.setAppearance(choice)
				return
			}
		}
	case id == "scan":
		a.scanOnce()
	}
}

// jumpTo puts the choice on a port found in the palette. Everything that
// hides a port is cleared first — the view, the folder, the search — so the
// row is in the list whatever the window was narrowed to, and the list
// scrolls to it rather than leaving the choice off screen.
func (a *App) jumpTo(key string) {
	a.mode, a.folder, a.query = ModeAll, "", ""
	a.refresh()
	for i, p := range a.shown {
		if p.Key() == key {
			a.chosen = i
			a.tableState.List.ScrollIntoView(i)
			a.status = "Chose " + p.Label()
			return
		}
	}
	a.chosen = -1
}

// openPalette shows the palette at the whole scan, rather than whatever was
// typed last time it was open, and gives the search field the keyboard.
func (a *App) openPalette(c *ui.Context) {
	a.paletteState.Query = ""
	a.paletteState.ActiveID = ""
	a.paletteOpen = true
	a.paletteState.Focus()
}
