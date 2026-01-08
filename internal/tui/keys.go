package tui

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines the key bindings for the application
type KeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Enter    key.Binding
	Quit     key.Binding
	Clear    key.Binding
	ClearAll key.Binding
}

// DefaultKeyMap returns the default key bindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "ctrl+p"),
			key.WithHelp("↑/ctrl+p", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "ctrl+n"),
			key.WithHelp("↓/ctrl+n", "move down"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Quit: key.NewBinding(
			key.WithKeys("esc", "ctrl+c", "q"),
			key.WithHelp("esc/q", "quit"),
		),
		Clear: key.NewBinding(
			key.WithKeys("backspace", "ctrl+h"),
			key.WithHelp("backspace", "delete char"),
		),
		ClearAll: key.NewBinding(
			key.WithKeys("ctrl+u"),
			key.WithHelp("ctrl+u", "clear input"),
		),
	}
}
