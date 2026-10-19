package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"

	"ports/internal/host"
	"ports/internal/port"
)

// ask opens the confirmation for stopping or restarting a port's process.
// Ending a process cannot be undone, so nothing happens until the user
// says so.
func (a *App) ask(verb string, p port.Port) {
	a.pending = &pendingAction{
		verb:   verb,
		port:   p,
		detail: fmt.Sprintf("%s (PID %d), listening on %s:%d", p.Name, p.PID, p.Addr, p.Port),
	}
}

// run carries out a confirmed action off the main thread, then re-scans so
// the list reflects what happened. What it reports is queued as a toast as
// well as written to the status line: the eye is usually on the row that was
// acted on, not on the bar under it.
func (a *App) run(action pendingAction) {
	p := action.port
	go func() {
		var err error
		switch action.verb {
		case "Restart":
			err = host.Restart(p)
		default:
			err = host.Stop(p.PID)
		}
		a.win.Update(func() {
			if err != nil {
				a.status = failed(action.verb, p, err)
				a.say(action.verb+" failed", a.status, components.ToastError)
			} else {
				a.status = done(action.verb, p)
				a.say(a.status, p.Label(), components.ToastSuccess)
			}
		})
		// Let the system settle before reading it again, so the list does
		// not still show the process that was just ended.
		time.Sleep(300 * time.Millisecond)
		a.scanNow()
	}()
}

// open hands a port's address to the browser.
func (a *App) open(c *ui.Context, p port.Port) {
	u := p.URL()
	a.status = "Opening " + u
	c.OpenURL(u)
}

// copy puts a port's address on the clipboard. The toast is the visible
// confirmation: the clipboard gives no feedback of its own.
func (a *App) copy(c *ui.Context, p port.Port) {
	u := p.URL()
	c.WriteClipboard(u)
	c.Announce("Copied " + u)
	a.status = "Copied " + u
	a.say("Copied", u, components.ToastSuccess)
}

var _ = mygo.Shell // the shell module opens URLs; kept for the dependency note

// done is the past tense of an action, written out rather than built from
// the verb: appending "ed" gives "Stoped".
func done(verb string, p port.Port) string {
	switch verb {
	case "Restart":
		return "Restarted " + p.Name
	default:
		return "Stopped " + p.Name
	}
}

// failed says what went wrong and what the state of the process now is.
// A restart that could not start again is the one case where the user has
// to be told the process is gone, because the port is now closed.
func failed(verb string, p port.Port, err error) string {
	switch verb {
	case "Restart":
		if strings.Contains(err.Error(), "could not start it again") {
			return "Stopped " + p.Name + ", but it did not start again. " +
				"Start it yourself: " + firstCommand(p)
		}
		return "Could not restart " + p.Name + ". It is still running."
	default:
		return "Could not stop " + p.Name + ". It may belong to another user; " +
			"stopping it needs permission this app does not have."
	}
}

// firstCommand names the program to run, for the one message where the
// user has to start a process by hand.
func firstCommand(p port.Port) string {
	argv, err := host.CommandLine(p.PID)
	if err != nil || len(argv) == 0 {
		return p.Name
	}
	return argv[0]
}
