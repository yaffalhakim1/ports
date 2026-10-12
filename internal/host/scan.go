package host

import (
	"sync"

	"ports/internal/port"
)

// Scan reads every listening TCP port with the process that owns it, its
// name and its working directory. It is the only place that touches the
// system, and it is safe to call from any goroutine: the view calls it
// through a background scan, never on the main thread.
func Scan() ([]port.Port, error) {
	raw, err := listeners()
	if err != nil {
		return nil, err
	}
	names, _ := processNames() // names are optional; a row falls back to its PID

	ports := make([]port.Port, 0, len(raw))
	for _, l := range raw {
		if l.pid == 0 {
			continue
		}
		ports = append(ports, port.Port{
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
func directories(ports []port.Port) map[int]string {
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
