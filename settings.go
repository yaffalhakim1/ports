package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/egoist/mygo"
)

// settings is what the app remembers between runs: the appearance the user
// chose. The window's own bounds are remembered by MyGo, through the
// window's StateKey.
type settings struct {
	Appearance string `json:"appearance"`
}

// settingsFile is the settings file's name inside the app's data
// directory, which is per user and created on first use.
const settingsFile = "settings.json"

// loadSettings reads the remembered appearance. A missing or unreadable
// file is not an error: it means the app follows the desktop, which is
// what a first run should do.
func loadSettings() settings {
	var s settings
	path, err := settingsPath()
	if err != nil {
		return s
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return s
	}
	_ = json.Unmarshal(data, &s)
	return s
}

// saveSettings writes the appearance. A failure is ignored on purpose: the
// choice still applies to this run, and losing it on the next start is not
// worth interrupting the user over.
func (a *app) saveSettings() {
	path, err := settingsPath()
	if err != nil {
		return
	}
	data, err := json.Marshal(settings{Appearance: a.appearance.String()})
	if err != nil {
		return
	}
	_ = os.MkdirAll(filepath.Dir(path), 0o755)
	_ = os.WriteFile(path, data, 0o644)
}

// settingsPath is where the settings live: the app's per-user data
// directory, so it sits beside the window state MyGo keeps.
func settingsPath() (string, error) {
	dir, err := mygo.App.Path(mygo.PathUserData)
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, settingsFile), nil
}
