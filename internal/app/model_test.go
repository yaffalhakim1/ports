package app

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"ports/internal/port"
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
	a.apply(append(fakePorts(), port.Port{Proto: "tcp", Addr: "127.0.0.1", Port: 1000, PID: 99, Name: "new.exe"}))
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

	kept := make([]port.Port, 0, len(a.ports))
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
	a.mode = ModeExposed
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

// TestLocalOnlyKeepsLoopbackPorts is the other half of the same filter, and
// the view the sidebar made reachable: the dev servers a browser on this
// machine can use.

func TestLocalOnlyKeepsLoopbackPorts(t *testing.T) {
	a := newTestApp()
	a.mode = ModeLocal
	a.refresh()
	for _, p := range a.shown {
		if !p.Loopback() {
			t.Errorf("local-only kept the reachable port %d", p.Port)
		}
	}
	// 127.0.0.1:3000 and ::1:8644 answer on this machine alone.
	if len(a.shown) != 2 {
		t.Errorf("local-only showed %v, want 2 rows", portsOf(a.shown))
	}
}

// TestFolderViewKeepsItsFolder checks a folder view shows that project's
// ports and nothing else, which is what makes the sidebar's Folders section
// worth having.

func TestFolderViewKeepsItsFolder(t *testing.T) {
	a := newTestApp()
	a.setView(ModeAll, `C:\work\web`)
	if len(a.shown) != 1 || a.shown[0].Port != 3000 {
		t.Fatalf("the web folder showed %v, want port 3000 alone", portsOf(a.shown))
	}
	// The filter is the whole path, not its last part: another project could
	// be called "web" too.
	a.setView(ModeAll, `C:\work\api`)
	if len(a.shown) != 1 || a.shown[0].Port != 5432 {
		t.Errorf("the api folder showed %v, want port 5432 alone", portsOf(a.shown))
	}
}

// TestFoldersAreCountedAndOrdered checks the sidebar's Folders section: one
// row per working directory, busiest first, and no row for a port whose
// directory could not be read.

func TestFoldersAreCountedAndOrdered(t *testing.T) {
	got := folders(fakePorts())
	if len(got) != 3 {
		t.Fatalf("counted %d folders, want 3: %v", len(got), got)
	}
	for _, f := range got {
		if f.Ports != 1 {
			t.Errorf("folder %q counted %d ports, want 1", f.Name, f.Ports)
		}
		if f.Dir == "" {
			t.Error("a folder row has no path to filter by")
		}
	}
	if names := []string{got[0].Name, got[1].Name, got[2].Name}; names[0] > names[1] || names[1] > names[2] {
		t.Errorf("equal folders are not ordered by name: %v", names)
	}
}

// TestSystemProcessCannotBeStopped checks the guard that keeps the app

func TestSystemProcessCannotBeStopped(t *testing.T) {
	for _, p := range fakePorts() {
		if got, want := p.SystemPID(), p.PID == 4; got != want {
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
		p := port.Port{Addr: c.addr, Port: c.port}
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
		if got := (port.Port{Dir: dir}).Folder(); got != want {
			t.Errorf("folder of %q is %q, want %q", dir, got, want)
		}
	}
}

func portsOf(ports []port.Port) []int {
	out := make([]int, len(ports))
	for i, p := range ports {
		out[i] = p.Port
	}
	return out
}

// TestFooterActionsAppearOnce checks each action is labelled once. A
// duplicated label is invisible in the render but doubles the button's
