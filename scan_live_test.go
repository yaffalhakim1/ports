package main

// TestScanLive reads the machine this test runs on. It is the only test
// that touches the system, so it is skipped in short mode.
import "testing"

func TestScanLive(t *testing.T) {
	if testing.Short() {
		t.Skip("reads the system")
	}
	ports, err := scan()
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	if len(ports) == 0 {
		t.Fatal("no listening ports found")
	}
	named, withDir := 0, 0
	for _, p := range ports {
		if p.Name != "" && p.Name != "pid "+itoa(p.PID) {
			named++
		}
		if p.Dir != "" {
			withDir++
		}
	}
	// Ports owned by another user or a protected process cannot be read,
	// so this asserts only that the scan resolves what it is allowed to.
	if named != len(ports) {
		t.Errorf("%d of %d ports had no process name", len(ports)-named, len(ports))
	}
	t.Logf("%d listening ports, %d with a working directory", len(ports), withDir)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b []byte
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	return string(b)
}
