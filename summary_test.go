package main

import "testing"

// TestSummaryReportsAFailedScan checks the header's status line does not
// claim the list is fresh when the last scan failed.
func TestSummaryReportsAFailedScan(t *testing.T) {
	a := newTestApp()
	a.lastErr = errScan{}
	got := a.summary()
	t.Logf("summary after a failed scan: %q", got)
	if got != "could not read the system" {
		t.Errorf("summary is %q, want the failure named", got)
	}

	// And with a good scan it reports the count and the age.
	b := newTestApp()
	b.lastScan = b.now
	if got := b.summary(); got == "" || got == "could not read the system" {
		t.Errorf("summary after a good scan is %q", got)
	}
}
