package app

import "github.com/egoist/mygo"

// Appearance is which of the two looks the app is showing. The desktop's
// own setting is the third choice, and the default: the app follows the
// system until the user says otherwise.
type Appearance int

const (
	AppearanceSystem Appearance = iota
	AppearanceLight
	AppearanceDark
)

// String is the appearance's name, as the settings file stores it and as
// the button shows it.
func (a Appearance) String() string {
	switch a {
	case AppearanceLight:
		return "Light"
	case AppearanceDark:
		return "Dark"
	default:
		return "System"
	}
}

// MenuLabel is what the appearance menu calls each choice. "System" rather
// than "Follow the desktop": it is what the setting is called everywhere
// else, and the label has to fit beside an icon.
func (a Appearance) MenuLabel() string { return a.String() }

// Source is the choice as MyGo's theme source, which the native chrome —
// menus, tooltips, the platforms' own controls — follows. The app draws its
// own surfaces from the same choice through theme.
func (a Appearance) Source() mygo.ThemeSource {
	switch a {
	case AppearanceLight:
		return mygo.ThemeLight
	case AppearanceDark:
		return mygo.ThemeDark
	default:
		return mygo.ThemeSystem
	}
}

// ParseAppearance reads a stored appearance, falling back to following the
// system for anything unrecognised.
func ParseAppearance(s string) Appearance {
	switch s {
	case "light", "Light":
		return AppearanceLight
	case "dark", "Dark":
		return AppearanceDark
	default:
		return AppearanceSystem
	}
}
