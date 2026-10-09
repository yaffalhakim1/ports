package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

// ask opens the confirmation for stopping or restarting a port's process.
// Ending a process cannot be undone, so nothing happens until the user
// says so.
func (a *app) ask(verb string, p Port) {
	a.pending = &pendingAction{
		verb:   verb,
		port:   p,
		detail: fmt.Sprintf("%s (PID %d), listening on %s:%d", p.Name, p.PID, p.Addr, p.Port),
	}
}

// confirm shows the pending action, if any, and runs it when the user
// agrees. The alert names the process, because a port number alone does
// not say what is about to end.
func (a *app) confirm(c *ui.Context) {
	if a.pending == nil {
		return
	}
	open := true
	message := a.pending.detail
	if a.pending.verb == "Restart" {
		message += "\nIt is started again with the same command, in the same folder."
	} else {
		message += "\nAnything it is serving stops."
	}

	switch ui.AlertDialog(c, &open, a.pending.verb+" "+a.pending.port.Name+"?", message,
		"Cancel", a.pending.verb) {
	case 1:
		action := *a.pending
		a.pending = nil
		a.run(action)
	case 0:
		a.pending = nil
		a.status = ""
	}
}

// run carries out a confirmed action off the main thread, then re-scans so
// the list reflects what happened.
func (a *app) run(action pendingAction) {
	p := action.port
	go func() {
		var err error
		switch action.verb {
		case "Restart":
			err = restart(p)
		default:
			err = stop(p.PID)
		}
		a.win.Update(func() {
			if err != nil {
				a.status = failed(action.verb, p, err)
			} else {
				a.status = done(action.verb, p)
			}
		})
		// Let the system settle before reading it again, so the list does
		// not still show the process that was just ended.
		time.Sleep(300 * time.Millisecond)
		a.scanNow()
	}()
}

// restart ends a process and starts it again the way it was started. The
// command is read before the process ends, because it is gone with it:
// reading it afterwards is what makes a restart silently lose its
// arguments.
func restart(p Port) error {
	argv, err := commandLine(p.PID)
	if err != nil {
		return fmt.Errorf("could not read how it was started: %w", err)
	}
	if err := stop(p.PID); err != nil {
		return err
	}
	time.Sleep(400 * time.Millisecond)
	if err := relaunch(p, argv); err != nil {
		return fmt.Errorf("stopped, but could not start it again: %w", err)
	}
	return nil
}

// open hands a port's address to the browser.
func (a *app) open(c *ui.Context, p Port) {
	u := p.URL()
	a.status = "Opening " + u
	c.OpenURL(u)
}

// copy puts a port's address on the clipboard.
func (a *app) copy(c *ui.Context, p Port) {
	u := p.URL()
	c.WriteClipboard(u)
	c.Announce("Copied " + u)
	a.status = "Copied " + u
}

var _ = mygo.Shell // the shell module opens URLs; kept for the dependency note

// done is the past tense of an action, written out rather than built from
// the verb: appending "ed" gives "Stoped".
func done(verb string, p Port) string {
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
func failed(verb string, p Port, err error) string {
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
func firstCommand(p Port) string {
	argv, err := commandLine(p.PID)
	if err != nil || len(argv) == 0 {
		return p.Name
	}
	return argv[0]
}
