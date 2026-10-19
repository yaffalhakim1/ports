package app

import (
	"strconv"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
	"github.com/oneo-me/onevra/source/shells/sidebar"
)

// side is the shell's left sidebar: which view the list shows, then which
// project. It is navigation, not settings — every row is a place to go, and
// exactly one of them is where the list is.
//
// It runs behind the shell's floating header, so it makes room for that
// header itself: the shell does not inset a region's builder.
func (a *App) side(c *ui.Context) {
	top := sidebar.HeaderHeight(c)
	ui.Column(c).Fill().MinWidth(0).Padding(top, 0, 0, 0).Children(func() {
		components.ScrollView(c, components.ScrollViewOptions{}).Fill().Children(func() {
			components.Sidebar(c).Children(func() {
				components.SidebarSection(c, "Ports", func() {
					a.portViews(c)
				})
				if fs := folders(a.ports); len(fs) > 0 {
					components.SidebarSection(c, "Folders", func() {
						for _, f := range fs {
							a.folderRow(c, f)
						}
					})
				}
			})
		})
	})
}

// portViews are the three questions the list answers, with their counts.
func (a *App) portViews(c *ui.Context) {
	for _, m := range modes {
		count := 0
		for _, p := range a.ports {
			if m.Keeps(p) {
				count++
			}
		}
		a.navRow(c, navItem{
			label:    m.Label(),
			icon:     modeIcon(m),
			selected: a.folder == "" && a.mode == m,
			count:    count,
			open:     func() { a.setView(m, "") },
		})
	}
}

// folderRow is one project: the ports it holds, under the folder it was
// started in. The whole path is the tooltip, because two projects can share
// a folder name and the row only shows the last part of it.
func (a *App) folderRow(c *ui.Context, f Folder) {
	a.navRow(c, navItem{
		label:    f.Name,
		icon:     lucide.Folder,
		selected: a.folder == f.Dir,
		count:    f.Ports,
		tooltip:  f.Dir,
		open:     func() { a.setView(ModeAll, f.Dir) },
	})
}

// navItem describes one row of the sidebar.
type navItem struct {
	label    string
	icon     func() *ui.SVG
	selected bool
	count    int
	tooltip  string
	open     func()
}

// navRow builds one row. The count is a trailing slot rather than part of
// the label, so the row keeps its identity — and its accessible name —
// while the number under it changes every two seconds.
//
// The leading icon takes the accent only while its row is the chosen one:
// the brand spends colour on the chosen row and the focus ring, and nowhere
// else.
func (a *App) navRow(c *ui.Context, row navItem) {
	t := c.Theme()
	item := components.SidebarItemWithSlots(c, row.label, row.selected,
		func(c *ui.Context, size float32) {
			color := t.TextMuted
			if row.selected {
				color = t.Accent
			}
			components.Icon(c, row.icon).Size(size, size).TextColor(color)
		},
		func(c *ui.Context, size float32) {
			components.Text(c, strconv.Itoa(row.count)).
				FontSize(t.FontSize - 1).TextColor(t.TextMuted).Shrink(0)
		})
	if row.tooltip != "" {
		components.Tooltip(c, item, row.tooltip)
	}
	if item.Clicked() {
		row.open()
	}
}

// modeIcon is the glyph of each view: everything, what another machine can
// reach, and what only this one can.
func modeIcon(m Mode) func() *ui.SVG {
	switch m {
	case ModeExposed:
		return lucide.Globe
	case ModeLocal:
		return lucide.Lock
	default:
		return lucide.Layers
	}
}
