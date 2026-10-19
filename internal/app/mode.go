package app

import (
	"sort"

	"ports/internal/port"
)

// Mode is which view of the last scan the list shows. It is the sidebar's
// Ports section: the three questions a person actually asks of this list.
//
// It replaces the single "Exposed only" switch, which could only ever name
// one of the three and left "local only" — the *dev servers I can reach*
// case — unreachable.
type Mode int

const (
	// ModeAll is every listening port.
	ModeAll Mode = iota
	// ModeExposed is only the ports another machine can reach.
	ModeExposed
	// ModeLocal is only the ports this machine alone can reach, which is
	// the dev servers a browser here can use.
	ModeLocal
)

// modes is the order the sidebar shows them in, which is also the order of
// widening scope: everything, then each half.
var modes = []Mode{ModeAll, ModeExposed, ModeLocal}

// Label is the mode's name in the sidebar.
func (m Mode) Label() string {
	switch m {
	case ModeExposed:
		return "Exposed only"
	case ModeLocal:
		return "Local only"
	default:
		return "All ports"
	}
}

// String is the view's stable name, for the command IDs the palette carries
// and for anything that has to persist which view is showing.
func (m Mode) String() string {
	switch m {
	case ModeExposed:
		return "exposed"
	case ModeLocal:
		return "local"
	default:
		return "all"
	}
}

// Empty is what a view says when a filter it owns has emptied the list, so
// the window explains itself rather than looking like a machine with nothing
// running on it.
func (m Mode) Empty() (message, hint string) {
	switch m {
	case ModeExposed:
		return "No exposed ports", "Nothing here answers to another machine."
	case ModeLocal:
		return "No local-only ports", "Everything listening here is reachable from the network."
	default:
		return "No listening ports", "Nothing is listening on this machine."
	}
}

// Keeps reports whether a port belongs in this view.
func (m Mode) Keeps(p port.Port) bool {
	switch m {
	case ModeExposed:
		return p.Exposed()
	case ModeLocal:
		return p.Loopback()
	default:
		return true
	}
}

// keepMode filters the scan to one view's ports.
func keepMode(ports []port.Port, m Mode) []port.Port {
	if m == ModeAll {
		return ports
	}
	out := make([]port.Port, 0, len(ports))
	for _, p := range ports {
		if m.Keeps(p) {
			out = append(out, p)
		}
	}
	return out
}

// keepFolder filters to the ports started in one working directory. The
// directory is the whole path rather than its last part, because two
// projects can share a folder name and the sidebar must be able to tell
// them apart.
func keepFolder(ports []port.Port, dir string) []port.Port {
	out := make([]port.Port, 0, len(ports))
	for _, p := range ports {
		if p.Dir == dir {
			out = append(out, p)
		}
	}
	return out
}

// Folder is one working directory in the scan, with the ports started in
// it. It is a sidebar row: the project's name, how many ports it holds,
// and the whole path behind them.
type Folder struct {
	Dir   string
	Name  string
	Ports int
}

// folders lists the working directories the scan holds, most ports first
// and then by name, so the busiest project is where the eye lands. Ports
// whose directory could not be read are left out: a group named "unknown"
// would sort first and lead nowhere.
func folders(ports []port.Port) []Folder {
	counts := make(map[string]*Folder)
	for _, p := range ports {
		if p.Dir == "" {
			continue
		}
		f, ok := counts[p.Dir]
		if !ok {
			f = &Folder{Dir: p.Dir, Name: p.Folder()}
			counts[p.Dir] = f
		}
		f.Ports++
	}
	out := make([]Folder, 0, len(counts))
	for _, f := range counts {
		out = append(out, *f)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Ports != out[j].Ports {
			return out[i].Ports > out[j].Ports
		}
		return lessFold(out[i].Name, out[j].Name)
	})
	return out
}
