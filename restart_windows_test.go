//go:build windows

package main

import (
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"testing"
	"time"
)

// TestStopAndRestartARealProcess starts a throwaway listener, finds it the
// way the window does, stops it and starts it again. It is the only test
// that changes the system, so it is skipped in short mode.
//
// The listener is a PowerShell TCP socket rather than a stub: it has a
// real command line and a real working directory, which is exactly what a
// restart has to reproduce.
func TestStopAndRestartARealProcess(t *testing.T) {
	if testing.Short() {
		t.Skip("starts and stops a process")
	}
	port := freePort(t)
	cmd := startListener(t, port)
	defer func() { _ = cmd.Process.Kill() }()

	// Find it the way the window does.
	found, err := waitForPort(port, 5*time.Second)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("found %s (PID %d) in %q", found.Name, found.PID, found.Dir)

	if found.Dir == "" {
		t.Error("the scan did not read the working directory, which a restart needs")
	}

	// Read how it was started before stopping it, as the app does: the
	// command line is gone with the process.
	argv, err := commandLine(found.PID)
	if err != nil {
		t.Fatalf("command line: %v", err)
	}
	t.Logf("started as %q", argv)

	// Stop it, and check the port is really gone.
	if err := stop(found.PID); err != nil {
		t.Fatalf("stop: %v", err)
	}
	if err := waitForPortGone(port, 5*time.Second); err != nil {
		t.Fatalf("the port is still listening after a stop: %v", err)
	}
	t.Log("stopped")

	// Start it again the way it was started, and check it is back.
	if err := relaunch(found, argv); err != nil {
		t.Fatalf("relaunch: %v", err)
	}
	back, err := waitForPort(port, 10*time.Second)
	if err != nil {
		t.Fatalf("the port did not come back: %v", err)
	}
	t.Logf("restarted as PID %d in %q", back.PID, back.Dir)

	if back.Dir != found.Dir {
		t.Errorf("restarted in %q, want the original %q", back.Dir, found.Dir)
	}
	if back.PID == found.PID {
		t.Error("the restarted process has the same PID, so it was never stopped")
	}
}

// TestStopRefusesASystemProcess checks the guard, since ending one of the
// kernel's processes would take the machine down.
func TestStopRefusesASystemProcess(t *testing.T) {
	for _, pid := range []int{0, 4} {
		if err := stop(pid); err == nil {
			t.Errorf("stop(%d) was allowed", pid)
		}
	}
}

// freePort returns a port nothing is listening on, for the test listener.
func freePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

// startListener starts a process that listens on a port and waits.
func startListener(t *testing.T, port int) *exec.Cmd {
	t.Helper()
	script := fmt.Sprintf(
		"$l=[System.Net.Sockets.TcpListener]::new([System.Net.IPAddress]::Loopback,%d);"+
			"$l.Start();Start-Sleep 300", port)
	cmd := exec.Command("powershell.exe", "-NoProfile", "-Command", script)
	if err := cmd.Start(); err != nil {
		t.Fatalf("start the test listener: %v", err)
	}
	return cmd
}

// waitForPort scans until the port appears, so the test does not depend on
// how long a process takes to start.
func waitForPort(port int, within time.Duration) (Port, error) {
	deadline := time.Now().Add(within)
	for time.Now().Before(deadline) {
		ports, err := scan()
		if err != nil {
			return Port{}, err
		}
		for _, p := range ports {
			if p.Port == port {
				return p, nil
			}
		}
		time.Sleep(200 * time.Millisecond)
	}
	return Port{}, fmt.Errorf("port %d never started listening", port)
}

// waitForPortGone scans until the port disappears.
func waitForPortGone(port int, within time.Duration) error {
	deadline := time.Now().Add(within)
	for time.Now().Before(deadline) {
		ports, err := scan()
		if err != nil {
			return err
		}
		gone := true
		for _, p := range ports {
			if p.Port == port {
				gone = false
				break
			}
		}
		if gone {
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	return fmt.Errorf("port %d still listens after %s", port, within)
}

var _ = strconv.Itoa
