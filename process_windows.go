//go:build windows

package main

import (
	"errors"
	"fmt"
	"strings"
	"syscall"
	"unsafe"
)

// Windows exposes a process's working directory and command line only
// through its PEB, by NtQueryInformationProcess: the toolhelp snapshot and
// WMI report neither. Everything read here is read-only, and a process
// that refuses to be read yields an error the caller treats as "unknown",
// which is ordinary for services owned by another user.
var (
	kernel32              = syscall.NewLazyDLL("kernel32.dll")
	ntdll                 = syscall.NewLazyDLL("ntdll.dll")
	procOpenProcess       = kernel32.NewProc("OpenProcess")
	procReadProcessMemory = kernel32.NewProc("ReadProcessMemory")
	procNtQueryInfo       = ntdll.NewProc("NtQueryInformationProcess")
)

const (
	processQueryLimitedInformation = 0x1000
	processVMRead                  = 0x0010

	// The PEB holds a pointer to its RTL_USER_PROCESS_PARAMETERS. Within
	// those parameters, CURDIR.DosPath is the working directory and
	// CommandLine the command that started the process, in a 64-bit
	// process.
	pebProcessParameters = 0x20
	paramsCurrentDir     = 0x38
	paramsCommandLine    = 0x70

	maxPathBytes = 4096
	maxImageLen  = 32768
)

// processBasicInformation mirrors PROCESS_BASIC_INFORMATION.
type processBasicInformation struct {
	ExitStatus                   uintptr
	PebBaseAddress               uintptr
	AffinityMask                 uintptr
	BasePriority                 int32
	UniqueProcessID              uintptr
	InheritedFromUniqueProcessID uintptr
}

// unicodeString mirrors UNICODE_STRING in a 64-bit process.
type unicodeString struct {
	Length        uint16
	MaximumLength uint16
	_             uint32
	Buffer        uintptr
}

// processDir reads the working directory a process was started in.
func processDir(pid int) (string, error) { return processString(pid, paramsCurrentDir) }

// processCommandLine reads the command that started a process, as the
// system recorded it.
func processCommandLine(pid int) (string, error) { return processString(pid, paramsCommandLine) }

// processString reads a UNICODE_STRING field of a process's parameters.
func processString(pid int, field uintptr) (string, error) {
	h, err := openProcess(pid, processQueryLimitedInformation|processVMRead)
	if err != nil {
		return "", err
	}
	defer syscall.CloseHandle(h)

	var pbi processBasicInformation
	var written uint32
	status, _, _ := procNtQueryInfo.Call(uintptr(h), 0,
		uintptr(unsafe.Pointer(&pbi)), unsafe.Sizeof(pbi),
		uintptr(unsafe.Pointer(&written)))
	if status != 0 {
		return "", fmt.Errorf("query process %d: status %#x", pid, status)
	}
	if pbi.PebBaseAddress == 0 {
		return "", fmt.Errorf("process %d has no PEB", pid)
	}

	// The PEB holds a pointer to the parameters, not the parameters
	// themselves, so follow it before reading a field of them.
	var params uintptr
	if err := readProcess(h, pbi.PebBaseAddress+pebProcessParameters,
		unsafe.Pointer(&params), unsafe.Sizeof(params)); err != nil {
		return "", err
	}
	if params == 0 {
		return "", fmt.Errorf("process %d has no parameters", pid)
	}

	var us unicodeString
	if err := readProcess(h, params+field, unsafe.Pointer(&us), unsafe.Sizeof(us)); err != nil {
		return "", err
	}
	if us.Length == 0 || us.Buffer == 0 || us.Length%2 != 0 || us.Length > maxPathBytes {
		return "", errors.New("not recorded")
	}
	buf := make([]uint16, us.Length/2)
	if err := readProcess(h, us.Buffer, unsafe.Pointer(&buf[0]), uintptr(us.Length)); err != nil {
		return "", err
	}
	return syscall.UTF16ToString(buf), nil
}

// commandLine returns the command that started a process, split into the
// program and its arguments, falling back to the executable's path when
// the command line cannot be read. It is what a restart repeats.
func commandLine(pid int) ([]string, error) {
	line, err := processCommandLine(pid)
	if err != nil || strings.TrimSpace(line) == "" {
		image, ierr := processImage(pid)
		if ierr != nil {
			return nil, fmt.Errorf("command line of process %d: %w", pid, ierr)
		}
		return []string{image}, nil
	}
	argv := splitCommandLine(line)
	if len(argv) == 0 {
		return nil, fmt.Errorf("parse the command line of process %d", pid)
	}
	return argv, nil
}

// processImage returns the full path of a process's executable.
func processImage(pid int) (string, error) {
	h, err := openProcess(pid, processQueryLimitedInformation)
	if err != nil {
		return "", err
	}
	defer syscall.CloseHandle(h)

	buf := make([]uint16, maxImageLen)
	size := uint32(len(buf))
	r, _, e := procQueryImage.Call(uintptr(h), 0,
		uintptr(unsafe.Pointer(&buf[0])), uintptr(unsafe.Pointer(&size)))
	if r == 0 {
		return "", fmt.Errorf("image of process %d: %w", pid, e)
	}
	return syscall.UTF16ToString(buf[:size]), nil
}

// openProcess opens a process with the access a caller needs.
func openProcess(pid int, access uint32) (syscall.Handle, error) {
	if pid <= 0 {
		return 0, fmt.Errorf("invalid pid %d", pid)
	}
	h, _, err := procOpenProcess.Call(uintptr(access), 0, uintptr(pid))
	if h == 0 {
		return 0, fmt.Errorf("open process %d: %w", pid, err)
	}
	return syscall.Handle(h), nil
}

// readProcess copies size bytes out of a process at addr.
func readProcess(h syscall.Handle, addr uintptr, dst unsafe.Pointer, size uintptr) error {
	var n uintptr
	r, _, err := procReadProcessMemory.Call(uintptr(h), addr, uintptr(dst), size,
		uintptr(unsafe.Pointer(&n)))
	if r == 0 || n != size {
		return fmt.Errorf("read process memory at %#x: %w", addr, err)
	}
	return nil
}
