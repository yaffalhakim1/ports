//go:build windows

package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

// Process names come from a toolhelp snapshot rather than tasklist, whose
// CSV output is localized and would not parse on a non-English system.
var (
	procCreateToolhelp32Snapshot = kernel32.NewProc("CreateToolhelp32Snapshot")
	procProcess32FirstW          = kernel32.NewProc("Process32FirstW")
	procProcess32NextW           = kernel32.NewProc("Process32NextW")
	procQueryImage               = kernel32.NewProc("QueryFullProcessImageNameW")
)

const (
	th32csSnapProcess = 0x00000002
	invalidHandle     = ^uintptr(0)
)

// processEntry32W mirrors PROCESSENTRY32W.
type processEntry32W struct {
	Size            uint32
	CntUsage        uint32
	ProcessID       uint32
	DefaultHeapID   uintptr
	ModuleID        uint32
	CntThreads      uint32
	ParentProcessID uint32
	PriClassBase    int32
	Flags           uint32
	ExeFile         [260]uint16
}

// processNames returns the executable name of every process, by PID. A
// name that cannot be read is left out, and the caller falls back to the
// PID, so a protected process still appears.
func processNames() (map[int]string, error) {
	snap, _, err := procCreateToolhelp32Snapshot.Call(th32csSnapProcess, 0)
	if snap == invalidHandle {
		return nil, fmt.Errorf("snapshot processes: %w", err)
	}
	defer syscall.CloseHandle(syscall.Handle(snap))

	names := map[int]string{}
	var entry processEntry32W
	entry.Size = uint32(unsafe.Sizeof(entry))
	ok, _, _ := procProcess32FirstW.Call(snap, uintptr(unsafe.Pointer(&entry)))
	for ok != 0 {
		names[int(entry.ProcessID)] = syscall.UTF16ToString(entry.ExeFile[:])
		ok, _, _ = procProcess32NextW.Call(snap, uintptr(unsafe.Pointer(&entry)))
	}
	return names, nil
}

// nameFor is the process name for a PID, or a placeholder that still
// identifies it when the name could not be read.
func nameFor(names map[int]string, pid int) string {
	if n := names[pid]; n != "" {
		return n
	}
	return "pid " + fmt.Sprint(pid)
}
