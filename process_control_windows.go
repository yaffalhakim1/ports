//go:build windows

package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"syscall"
)

// Ending and starting processes. A stop cannot be undone, so the view
// always confirms one first; these functions only carry it out.
var (
	procTerminateProcess = kernel32.NewProc("TerminateProcess")
)

const (
	processTerminate = 0x0001
	// createNewProcessGroup keeps a relaunched process from sharing this
	// app's console and signal handling, so it outlives the window.
	createNewProcessGroup = 0x00000200
)

// stop ends a process. A process this user may not end, such as a service
// of another user, returns an error the view reports rather than hides.
func stop(pid int) error {
	if pid <= 4 {
		return fmt.Errorf("process %d is a system process", pid)
	}
	h, err := openProcess(pid, processTerminate)
	if err != nil {
		return err
	}
	defer syscall.CloseHandle(h)
	if r, _, e := procTerminateProcess.Call(uintptr(h), 1); r == 0 {
		return fmt.Errorf("stop process %d: %w", pid, e)
	}
	return nil
}

// relaunch starts a process again as it was started: the same command, in
// the same working directory, detached so that it outlives this app. The
// command is passed in because a stopped process no longer has one.
func relaunch(p Port, argv []string) error {
	if p.systemPID() {
		return fmt.Errorf("process %d is a system process", p.PID)
	}
	if len(argv) == 0 {
		return fmt.Errorf("no command to start process %d", p.PID)
	}
	dir := p.Dir
	if dir == "" {
		dir = filepath.Dir(argv[0])
	}
	cmd := exec.Command(argv[0], argv[1:]...)
	cmd.Dir = dir
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: createNewProcessGroup}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start %s: %w", argv[0], err)
	}
	// Nothing waits for the process, as a server started this way is meant
	// to run on its own.
	return cmd.Process.Release()
}
