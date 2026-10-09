package main

import (
	"fmt"
	"net"
	"path/filepath"
	"strconv"
	"strings"
)

// Port is one TCP port a process is listening on.
type Port struct {
	Proto string // "tcp" or "tcp6"
	Addr  string // the local address it binds, as the system reports it
	Port  int
	PID   int
	Name  string // the process's executable name, "pid N" when unknown
	Dir   string // the process's working directory, "" when unreadable
}

// Key identifies the listener across scans, so the choice and the sort
// follow it while the list is rebuilt.
func (p Port) Key() string {
	return p.Proto + "|" + p.Addr + "|" + strconv.Itoa(p.Port)
}

// Label is what a row shows and what assistive technology reads.
func (p Port) Label() string {
	return fmt.Sprintf("%s %s:%d", p.Name, p.Addr, p.Port)
}

// URL is where the port serves, for the addresses a browser can reach. A
// port bound to every interface is opened as localhost, which is what a
// browser on this machine should use.
func (p Port) URL() string {
	host := strings.Trim(p.Addr, "[]")
	switch host {
	case "0.0.0.0", "::", "*", "":
		host = "localhost"
	}
	// JoinHostPort brackets an IPv6 host itself, so the host is passed
	// bare here.
	return "http://" + net.JoinHostPort(host, strconv.Itoa(p.Port))
}

// Loopback reports whether the port serves this machine alone. It is what
// separates a development server from something another machine can
// reach, and the window uses it for the "Exposed only" filter.
func (p Port) Loopback() bool {
	switch p.Addr {
	case "127.0.0.1", "::1", "localhost":
		return true
	}
	return strings.HasPrefix(p.Addr, "127.")
}

// Exposed reports whether another machine can reach the port, which is
// what the "Exposed only" filter keeps.
func (p Port) Exposed() bool { return !p.Loopback() }

// Folder is the last part of the working directory, to show beside the
// process name, with the whole path in its tooltip.
func (p Port) Folder() string {
	if p.Dir == "" {
		return ""
	}
	trimmed := strings.TrimRight(p.Dir, `\/`)
	base := filepath.Base(trimmed)
	if base == "." || base == string(filepath.Separator) || base == trimmed && trimmed == "" {
		return ""
	}
	return base
}

// systemPID reports whether the port belongs to a process the app must
// not end: the kernel's own processes are the lowest PIDs, and they cannot
// be stopped at all.
func (p Port) systemPID() bool { return p.PID <= 4 }
