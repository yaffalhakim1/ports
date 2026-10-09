package main

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

// Tests of the actions: the confirmation before a destructive one, and
// the wording of what it reports afterwards.

func TestStopAsksBeforeActing(t *testing.T) {
	a := newTestApp()
	a.chosen = 1
	want, _ := a.selected()

	tt := ui.NewTester(a.view, 900, 620)
	if err := tt.Click("Stop"); err != nil {
		t.Fatal(err)
	}
	if a.pending == nil {
		t.Fatal("clicking Stop did not ask for confirmation")
	}
	if a.pending.port.Key() != want.Key() {
		t.Errorf("the confirmation is about %v, want %v", a.pending.port, want)
	}
	if !tt.HasText("Stop " + want.Name + "?") {
		t.Errorf("the confirmation does not name the process; texts are %q", tt.Texts())
	}
}

// TestCancelLeavesTheProcessAlone checks cancelling the confirmation

func TestCancelLeavesTheProcessAlone(t *testing.T) {
	a := newTestApp()
	a.chosen = 1
	a.ask("Stop", a.shown[1])

	tt := ui.NewTester(a.view, 900, 620)
	if err := tt.Click("Cancel"); err != nil {
		t.Fatal(err)
	}
	if a.pending != nil {
		t.Error("cancelling left the action pending")
	}
}

// TestActionsStayOnScreen checks every action stays inside the window at
// the smallest size the window allows. Before this was tested, Restart

func TestActionMessagesAreWords(t *testing.T) {
	p := Port{Name: "node.exe", PID: 42, Port: 3000}
	for _, c := range []struct {
		verb string
		want string
	}{
		{"Stop", "Stopped node.exe"},
		{"Restart", "Restarted node.exe"},
	} {
		if got := done(c.verb, p); got != c.want {
			t.Errorf("after %s the message is %q, want %q", c.verb, got, c.want)
		}
	}
}

// TestStopFailureSaysWhatToDo checks a failure names a way forward rather

func TestStopFailureSaysWhatToDo(t *testing.T) {
	p := Port{Name: "mysqld.exe", PID: 6416, Port: 3306}
	msg := failed("Stop", p, errAccessDenied{})
	if !strings.Contains(msg, "mysqld.exe") {
		t.Errorf("the failure does not name the process: %q", msg)
	}
	if !strings.Contains(msg, "permission") {
		t.Errorf("the failure does not say what to do: %q", msg)
	}
}

// TestRestartFailureWarnsTheProcessIsGone checks the one case where a
// failure leaves the port closed: the user has to be told, because the

func TestRestartFailureWarnsTheProcessIsGone(t *testing.T) {
	p := Port{Name: "node.exe", PID: 42, Port: 3000}
	msg := failed("Restart", p, errStart{})
	if !strings.Contains(msg, "did not start") {
		t.Errorf("the failure does not warn that it is gone: %q", msg)
	}
}

type errAccessDenied struct{}

func (errAccessDenied) Error() string { return "access is denied" }

type errStart struct{}

func (errStart) Error() string { return "could not start it again: not found" }

// TestEveryActionHasAShape checks the controls are distinguishable from
// the text beside them. A control styled like static text does not read
