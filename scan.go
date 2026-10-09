package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/egoist/mygo/ui"
)

// scan reads every listening TCP port with the process that owns it, its
// name and its working directory. It is the only place that touches the
// system, and it is safe to call from any goroutine: the view calls it
// through a background scan, never on the main thread.
func scan() ([]Port, error) {
	raw, err := listeners()
	if err != nil {
		return nil, err
	}
	names, _ := processNames() // names are optional; a row falls back to its PID

	ports := make([]Port, 0, len(raw))
	for _, l := range raw {
		if l.pid == 0 {
			continue
		}
		ports = append(ports, Port{
			Proto: l.proto,
			Addr:  l.addr,
			Port:  l.port,
			PID:   l.pid,
			Name:  nameFor(names, l.pid),
		})
	}

	// The working directory costs a separate handle and read per process,
	// so read each PID once even when it holds several ports.
	dirs := directories(ports)
	for i := range ports {
		ports[i].Dir = dirs[ports[i].PID]
	}
	return ports, nil
}

// directories reads the working directory of each distinct PID, a few at
// a time, since each is a process handle and a read.
func directories(ports []Port) map[int]string {
	seen := make(map[int]bool, len(ports))
	pids := make([]int, 0, len(ports))
	for _, p := range ports {
		if !seen[p.PID] {
			seen[p.PID] = true
			pids = append(pids, p.PID)
		}
	}

	const workers = 8
	out := make(map[int]string, len(pids))
	var mu sync.Mutex
	var wg sync.WaitGroup
	work := make(chan int)

	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for pid := range work {
				dir, err := processDir(pid)
				if err != nil {
					continue // unreadable: shown as unknown, not as an error
				}
				mu.Lock()
				out[pid] = dir
				mu.Unlock()
			}
		}()
	}
	for _, pid := range pids {
		work <- pid
	}
	close(work)
	wg.Wait()
	return out
}

// sortPorts orders ports by the column a table header asked for. The
// column IDs are those of the table's columns, so a header click needs no
// translation. Ties fall back to the port and then the address, so rows
// never shuffle between scans.
func sortPorts(ports []Port, order ui.SortOrder) {
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
func filter(ports []Port, query string) []Port {
	query = strings.TrimSpace(query)
	if query == "" {
		return ports
	}
	q := strings.ToLower(query)
	out := make([]Port, 0, len(ports))
	for _, p := range ports {
		if matches(p, q) {
			out = append(out, p)
		}
	}
	return out
}

func matches(p Port, q string) bool {
	return strings.Contains(strconv.Itoa(p.Port), q) ||
		strings.Contains(strings.ToLower(p.Name), q) ||
		strings.Contains(strings.ToLower(p.Addr), q) ||
		strings.Contains(strings.ToLower(p.Dir), q)
}
