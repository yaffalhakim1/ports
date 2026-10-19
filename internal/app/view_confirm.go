package app

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

// confirm shows the pending action, if any, and runs it when the user
// agrees. The alert names the process, because a port number alone does not
// say what is about to end.
//
// It is Onevra's dialog rather than its alert preset: the preset's two
// buttons are the accent and the quiet secondary, and this app spends its
// accent on the chosen row and the focus ring rather than on a fill. The
// button that ends a process wears the same danger face the action bar gives
// it, so pressing Stop leads to a Stop that still looks like one.
func (a *App) confirm(c *ui.Context) {
	if a.pending == nil {
		return
	}
	action := *a.pending
	open := true
	result := components.DialogNone

	components.Dialog(c, &open, components.DialogOptions{
		Label: action.verb + " " + action.port.Name + "?",
	}, func(panel ui.Element) {
		panel.Role(ui.RoleAlertDialog)
		t := c.Theme()
		components.Text(c, action.verb+" "+action.port.Name+"?").Bold()
		a.confirmMessage(c, action)
		ui.Row(c).FillWidth().Wrap().Justify(ui.End).Gap(t.Space(2)).Children(func() {
			cancel := components.Button(c, "", components.ButtonSecondary)
			cancel.Children(func() { components.Text(c, "Cancel") })
			if cancel.Clicked() {
				result, open = components.DialogCancelled, false
			}
			confirm := components.Button(c, "", components.ButtonOutline)
			confirm.Children(func() { components.Text(c, action.verb) })
			dangerFace(c, confirm)
			if confirm.Clicked() {
				result, open = components.DialogConfirmed, false
			}
		})
	})

	if open {
		return
	}
	a.pending = nil
	if result == components.DialogConfirmed {
		a.run(action)
		return
	}
	a.status = ""
}

// confirmMessage says what the action will do to the process, which is the
// part a port number does not tell the reader.
func (a *App) confirmMessage(c *ui.Context, action pendingAction) {
	t := c.Theme()
	components.Text(c, action.detail).TextColor(t.TextMuted).Wrap()
	if action.verb == "Restart" {
		components.Text(c, "It is started again with the same command, in the same folder.").
			TextColor(t.TextMuted).Wrap()
		return
	}
	components.Text(c, "Anything it is serving stops.").TextColor(t.TextMuted).Wrap()
}
