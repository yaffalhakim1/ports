package app

import (
	"sort"
	"strconv"
	"strings"

	"github.com/egoist/mygo/ui"

	"ports/internal/port"
)

// sortPorts orders ports by the column a table header asked for. The
// column IDs are those of the table's columns, so a header click needs no
// translation. Ties fall back to the port and then the address, so rows
// never shuffle between scans.
func sortPorts(ports []port.Port, order ui.SortOrder) {
	less := func(i, j int) bool {
		a, b := ports[i], ports[j]
		switch order.Column {
		case "name":
			return lessFold(a.Name, b.Name)
		case "folder":
			return lessFold(a.Folder(), b.Folder())
		case "address":
			return a.Addr < b.Addr
		case "pid":
			return a.PID < b.PID
		default: // "port"
			return a.Port < b.Port
		}
	}
	sort.SliceStable(ports, func(i, j int) bool {
		// Ties are broken by the port and then the address, in the same
		// direction, so reversing the order reverses the whole list.
		flip := func(less func(i, j int) bool) bool {
			if order.Descending {
				return less(j, i)
			}
			return less(i, j)
		}
		if less(i, j) || less(j, i) {
			return flip(less)
		}
		if ports[i].Port != ports[j].Port {
			return flip(func(i, j int) bool { return ports[i].Port < ports[j].Port })
		}
		return flip(func(i, j int) bool { return ports[i].Key() < ports[j].Key() })
	})
}

// lessFold compares two names without regard to case, and by their raw
// text when they are equal that way, so the order is the same on every
// run rather than depending on the sort's stability.
func lessFold(a, b string) bool {
	la, lb := strings.ToLower(a), strings.ToLower(b)
	if la != lb {
		return la < lb
	}
	return a < b
}

// filter keeps the ports matching a query, which matches the port, the
// process name, the address and the folder, case-insensitively.
func filter(ports []port.Port, query string) []port.Port {
	query = strings.TrimSpace(query)
	if query == "" {
		return ports
	}
	q := strings.ToLower(query)
	out := make([]port.Port, 0, len(ports))
	for _, p := range ports {
		if matches(p, q) {
			out = append(out, p)
		}
	}
	return out
}

func matches(p port.Port, q string) bool {
	return strings.Contains(strconv.Itoa(p.Port), q) ||
		strings.Contains(strings.ToLower(p.Name), q) ||
		strings.Contains(strings.ToLower(p.Addr), q) ||
		strings.Contains(strings.ToLower(p.Dir), q)
}
