//go:build windows

package main

import (
	"fmt"
	"net"
	"sort"
	"syscall"
	"unsafe"
)

// The tables come from iphlpapi rather than from parsing netstat: netstat
// output is localized, so its columns and address text differ per system
// language, while these tables do not.
var (
	iphlpapi                = syscall.NewLazyDLL("iphlpapi.dll")
	procGetExtendedTcpTable = iphlpapi.NewProc("GetExtendedTcpTable")
)

const (
	afINet  = 2
	afINet6 = 23
	// tcpTableOwnerPIDListener asks for the listeners alone, each with the
	// PID that owns it.
	tcpTableOwnerPIDListener = 3
	// errorInsufficientBuffer is what the sizing call reports, and what
	// tells that the table grew between two calls.
	errorInsufficientBuffer = 122
)

// listener is one table entry, before its name and directory are resolved.
type listener struct {
	proto string
	addr  string
	port  int
	pid   int
}

// mibTCPRowOwnerPID mirrors MIB_TCPROW_OWNER_PID, whose addresses and
// ports are in network byte order.
type mibTCPRowOwnerPID struct {
	State      uint32
	LocalAddr  uint32
	LocalPort  uint32
	RemoteAddr uint32
	RemotePort uint32
	OwningPID  uint32
}

// mibTCP6RowOwnerPID mirrors MIB_TCP6ROW_OWNER_PID.
type mibTCP6RowOwnerPID struct {
	LocalAddr     [16]byte
	LocalScopeID  uint32
	LocalPort     uint32
	RemoteAddr    [16]byte
	RemoteScopeID uint32
	RemotePort    uint32
	State         uint32
	OwningPID     uint32
}

// listeners returns every TCP port this machine listens on, by port.
func listeners() ([]listener, error) {
	v4, err := listen4()
	if err != nil {
		return nil, err
	}
	v6, err := listen6()
	if err != nil {
		return nil, err
	}
	all := append(v4, v6...)
	sort.Slice(all, func(i, j int) bool {
		if all[i].port != all[j].port {
			return all[i].port < all[j].port
		}
		return all[i].proto < all[j].proto
	})
	return all, nil
}

func listen4() ([]listener, error) {
	rows, err := tcpTable(afINet, unsafe.Sizeof(mibTCPRowOwnerPID{}))
	if err != nil {
		return nil, err
	}
	out := make([]listener, 0, len(rows))
	for _, raw := range rows {
		r := (*mibTCPRowOwnerPID)(unsafe.Pointer(&raw[0]))
		out = append(out, listener{
			proto: "tcp",
			addr: net.IPv4(byte(r.LocalAddr), byte(r.LocalAddr>>8),
				byte(r.LocalAddr>>16), byte(r.LocalAddr>>24)).String(),
			port: portOf(r.LocalPort),
			pid:  int(r.OwningPID),
		})
	}
	return out, nil
}

func listen6() ([]listener, error) {
	rows, err := tcpTable(afINet6, unsafe.Sizeof(mibTCP6RowOwnerPID{}))
	if err != nil {
		return nil, err
	}
	out := make([]listener, 0, len(rows))
	for _, raw := range rows {
		r := (*mibTCP6RowOwnerPID)(unsafe.Pointer(&raw[0]))
		out = append(out, listener{
			proto: "tcp6",
			addr:  net.IP(r.LocalAddr[:]).String(),
			port:  portOf(r.LocalPort),
			pid:   int(r.OwningPID),
		})
	}
	return out, nil
}

// tcpTable returns the rows of a listener table, each as raw bytes of
// rowSize that the caller casts to the row type it asked for. The table
// can grow between sizing it and reading it, which the call reports as an
// insufficient buffer; retry a few times rather than drop the scan.
func tcpTable(family uint32, rowSize uintptr) ([][]byte, error) {
	var size uint32
	for attempt := 0; attempt < 5; attempt++ {
		r, _, _ := procGetExtendedTcpTable.Call(0,
			uintptr(unsafe.Pointer(&size)), 0, uintptr(family), tcpTableOwnerPIDListener, 0)
		if r != 0 && r != errorInsufficientBuffer {
			return nil, fmt.Errorf("size the tcp table of family %d: code %d", family, r)
		}
		if size == 0 {
			return nil, fmt.Errorf("size the tcp table of family %d", family)
		}

		buf := make([]byte, size)
		r, _, _ = procGetExtendedTcpTable.Call(uintptr(unsafe.Pointer(&buf[0])),
			uintptr(unsafe.Pointer(&size)), 0, uintptr(family), tcpTableOwnerPIDListener, 0)
		switch r {
		case 0:
			return splitRows(buf, rowSize), nil
		case errorInsufficientBuffer:
			// The table grew; size carries the new length, so try again.
			continue
		default:
			return nil, fmt.Errorf("read the tcp table of family %d: code %d", family, r)
		}
	}
	return nil, fmt.Errorf("read the tcp table of family %d: it kept changing", family)
}

// splitRows cuts a table into rows: a count of rows, then the rows.
func splitRows(buf []byte, rowSize uintptr) [][]byte {
	if len(buf) < 4 || rowSize == 0 {
		return nil
	}
	n := int(*(*uint32)(unsafe.Pointer(&buf[0])))
	rows := make([][]byte, 0, n)
	for off := 4; off+int(rowSize) <= len(buf) && len(rows) < n; off += int(rowSize) {
		rows = append(rows, buf[off:off+int(rowSize)])
	}
	return rows
}

// portOf reads a port the tables store in network byte order, which
// Windows keeps in the low half of the word.
func portOf(v uint32) int {
	lo := uint16(v)
	return int(lo>>8 | lo<<8)
}
