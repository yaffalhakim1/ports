package main

// fakePorts is a fixed list, so tests do not depend on what this machine
// happens to be listening on, and newTestApp is an app showing it as if a
// scan had just finished.

func fakePorts() []Port {
	return []Port{
		{Proto: "tcp", Addr: "127.0.0.1", Port: 3000, PID: 4242, Name: "node.exe", Dir: `C:\work\web`},
		{Proto: "tcp", Addr: "0.0.0.0", Port: 5432, PID: 5150, Name: "postgres.exe", Dir: `C:\work\api`},
		{Proto: "tcp6", Addr: "::1", Port: 8644, PID: 21740, Name: "python.exe", Dir: `C:\Users\ada\hermes`},
		{Proto: "tcp", Addr: "0.0.0.0", Port: 445, PID: 4, Name: "System"},
	}
}

// newTestApp is an app showing a fixed list, as if a scan had finished.
func newTestApp() *app {
	a := &app{chosen: -1}
	a.ports = fakePorts()
	a.refresh()
	return a
}

// TestViewShowsPorts checks the window builds from state alone: the ports,
