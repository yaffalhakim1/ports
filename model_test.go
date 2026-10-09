package main

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// Tests of the model: what the scan finds, how it is filtered and sorted,
// and how the choice follows it. These need no window.

func TestFilterNarrowsTheList(t *testing.T) {
	want := map[string]int{"3000": 3000, "node": 3000, "web": 3000, "5432": 5432, "hermes": 8644}
	for query, port := range want {
		a := newTestApp()
		a.query = query
		a.refresh()
		if len(a.shown) != 1 || a.shown[0].Port != port {
			t.Errorf("query %q showed %v, want port %d alone", query, portsOf(a.shown), port)
		}
	}
}

// TestSortOrdersRows checks every sortable column orders the list, and

func TestSortOrdersRows(t *testing.T) {
	cases := []struct {
		column string
		first  int // the port expected first, ascending
	}{
		{"port", 445},
		{"name", 3000},   // node, postgres, python, then System, case-insensitively
		{"folder", 445},  // no folder sorts first
		{"address", 445}, // 0.0.0.0 sorts before 127.0.0.1 and ::1
		{"pid", 445},     // PID 4 is the lowest
	}
	for _, c := range cases {
		a := newTestApp()
		a.order = ui.SortOrder{Column: c.column}
		a.refresh()
		if len(a.shown) == 0 || a.shown[0].Port != c.first {
			t.Errorf("sorting by %s put %v first, want port %d", c.column, portsOf(a.shown), c.first)
		}

		a.order.Descending = true
		a.refresh()
		if last := a.shown[len(a.shown)-1]; last.Port != c.first {
			t.Errorf("reversing %s put %d last, want port %d", c.column, last.Port, c.first)
		}
	}
}

// TestChoiceFollowsThePort checks a refresh keeps the chosen row on the
// same port even when the list around it changes, which is what makes a

func TestChoiceFollowsThePort(t *testing.T) {
	a := newTestApp()
	a.order = ui.SortOrder{Column: "port"}
	a.refresh()
	a.chosen = 2
	want, _ := a.selected()

	// A scan that finds one more port, which shifts the rows.
	a.apply(append(fakePorts(), Port{Proto: "tcp", Addr: "127.0.0.1", Port: 1000, PID: 99, Name: "new.exe"}))
	got, ok := a.selected()
	if !ok || got.Key() != want.Key() {
		t.Errorf("after a refresh the choice is %v (%v), want %v", got, ok, want)
	}
}

// TestChoiceClearedWhenPortGoes checks the choice is dropped, not moved to
// another process, when the chosen port stops listening: acting on the

func TestChoiceClearedWhenPortGoes(t *testing.T) {
	a := newTestApp()
	a.refresh()
	a.chosen = 1
	gone, _ := a.selected()

	kept := make([]Port, 0, len(a.ports))
	for _, p := range a.ports {
		if p.Key() != gone.Key() {
			kept = append(kept, p)
		}
	}
	a.apply(kept)
	if got, ok := a.selected(); ok {
		t.Errorf("the choice followed a port that stopped listening: %v", got)
	}
}

// TestExposedOnlyKeepsExposedPorts checks the filter keeps what another
// machine could reach and drops what only this machine can, which is the

func TestExposedOnlyKeepsExposedPorts(t *testing.T) {
	a := newTestApp()
	a.only = true
	a.refresh()
	for _, p := range a.shown {
		if p.Loopback() {
			t.Errorf("exposed-only kept the loopback port %d", p.Port)
		}
	}
	// 0.0.0.0:5432 and 0.0.0.0:445 are reachable; 127.0.0.1 and ::1 are not.
	if len(a.shown) != 2 {
		t.Errorf("exposed-only showed %v, want 2 rows", portsOf(a.shown))
	}
}

// TestSystemProcessCannotBeStopped checks the guard that keeps the app

func TestSystemProcessCannotBeStopped(t *testing.T) {
	for _, p := range fakePorts() {
		if got, want := p.systemPID(), p.PID == 4; got != want {
			t.Errorf("PID %d blocked=%v, want %v", p.PID, got, want)
		}
	}
}

// TestSplitCommandLine checks the command-line splitter against the
// quoting Windows actually uses, which decides whether a restart starts

func TestPortURL(t *testing.T) {
	cases := []struct {
		addr string
		port int
		want string
	}{
		{"0.0.0.0", 3000, "http://localhost:3000"},
		{"::", 3000, "http://localhost:3000"},
		{"127.0.0.1", 8644, "http://127.0.0.1:8644"},
		{"::1", 5173, "http://[::1]:5173"},
		{"192.168.1.4", 80, "http://192.168.1.4:80"},
	}
	for _, c := range cases {
		p := Port{Addr: c.addr, Port: c.port}
		if got := p.URL(); got != c.want {
			t.Errorf("%s:%d opened as %q, want %q", c.addr, c.port, got, c.want)
		}
	}
}

// TestFolderIsTheLastPart checks the folder column shows the project name

func TestFolderIsTheLastPart(t *testing.T) {
	cases := map[string]string{
		`C:\work\web`:         "web",
		`C:\work\web\`:        "web",
		`C:\Users\ada\hermes`: "hermes",
		`C:/Users/ada/hermes`: "hermes",
		``:                    "",
	}
	for dir, want := range cases {
		if got := (Port{Dir: dir}).Folder(); got != want {
			t.Errorf("folder of %q is %q, want %q", dir, got, want)
		}
	}
}

func portsOf(ports []Port) []int {
	out := make([]int, len(ports))
	for i, p := range ports {
		out[i] = p.Port
	}
	return out
}

// TestFooterActionsAppearOnce checks each action is labelled once. A
// duplicated label is invisible in the render but doubles the button's
