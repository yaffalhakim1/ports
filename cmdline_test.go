package main

import "testing"

// The command-line splitter decides whether a restart starts the same
// command it stopped, so its quoting rules are tested directly.

func TestSplitCommandLine(t *testing.T) {
	cases := []struct {
		line string
		want []string
	}{
		{``, nil},
		{`   `, nil},
		{`node server.js`, []string{"node", "server.js"}},
		{`"C:\Program Files\node\node.exe" server.js --port 3000`,
			[]string{`C:\Program Files\node\node.exe`, "server.js", "--port", "3000"}},
		{`python -m hermes_cli.main gateway run`,
			[]string{"python", "-m", "hermes_cli.main", "gateway", "run"}},
		{`node "a b.js"`, []string{"node", "a b.js"}},
		{`node "say ""hi"""`, []string{"node", `say "hi"`}},
		{`C:\tools\node.exe x`, []string{`C:\tools\node.exe`, "x"}},
		{`node "trailing\\"`, []string{"node", `trailing\`}},
	}
	for _, c := range cases {
		got := splitCommandLine(c.line)
		if len(got) != len(c.want) {
			t.Errorf("%q split into %q, want %q", c.line, got, c.want)
			continue
		}
		for i := range got {
			if got[i] != c.want[i] {
				t.Errorf("%q split into %q, want %q", c.line, got, c.want)
				break
			}
		}
	}
}

// TestPortURL checks the address a browser is given for each way a port
