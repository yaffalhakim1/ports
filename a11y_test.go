package main

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// TestCopyAnnounces checks copying reaches assistive technology, since the
// clipboard change is otherwise invisible to a screen reader.
func TestCopyAnnounces(t *testing.T) {
	a := newTestApp()
	a.chosen = 1
	tt := ui.NewTester(a.view, 900, 620)
	tt.Frame()

	if err := tt.Click("Copy"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()

	want := "http://127.0.0.1:3000"
	if got := tt.Clipboard(); got != want {
		t.Errorf("clipboard is %q, want %q", got, want)
	}
	told := tt.Announcements()
	t.Logf("announcements: %q", told)
	if len(told) == 0 {
		t.Error("copying announced nothing to assistive technology")
	}
	if a.status == "" {
		t.Error("copying left no visible confirmation")
	}
}

// TestEveryControlIsNamed checks every interactive element carries a name
// for assistive technology, which is the escalation trigger that matters
// most for a window made of icon buttons.
func TestEveryControlIsNamed(t *testing.T) {
	a := newTestApp()
	a.chosen = 1
	tt := ui.NewTester(a.view, 900, 620)
	tt.Frame()

	for _, name := range []string{"Scan", "Search ports", "Exposed only", "Open", "Copy", "Restart", "Stop"} {
		if _, ok := tt.Find(name); !ok {
			t.Errorf("%q has no accessible name", name)
		}
	}
	// The table names its rows, so a screen reader reads the port, the
	// process and the address together rather than three loose cells.
	if !tt.HasText("node.exe 127.0.0.1:3000") {
		t.Error("the table rows carry no combined name")
	}
}
