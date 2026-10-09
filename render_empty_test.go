package main

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// TestRenderEmptyState writes the two empty states, which are otherwise
// hard to reach by hand and are where the copy matters most.
func TestRenderEmptyState(t *testing.T) {
	// Nothing matches the search: the state names the query and offers a
	// way out.
	a := newTestApp()
	a.query = "nothing-matches-this"
	a.refresh()
	tt := ui.NewTester(a.view, 900, 620)
	tt.Frame()
	writePNG(t, "testdata/empty-search.png", tt.Image())

	// The scan failed: the state says so and suggests scanning again.
	b := newTestApp()
	b.lastErr = errScan{}
	b.ports, b.shown = nil, nil
	tt2 := ui.NewTester(b.view, 900, 620)
	tt2.Frame()
	writePNG(t, "testdata/empty-error.png", tt2.Image())
}

type errScan struct{}

func (errScan) Error() string { return "access is denied reading the tcp table" }
