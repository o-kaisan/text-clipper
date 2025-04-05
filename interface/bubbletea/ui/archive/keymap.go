package archive

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	restore key.Binding
	Delete  key.Binding
	Up      key.Binding
	Down    key.Binding
	Next    key.Binding
	Prev    key.Binding
	Home    key.Binding
	End     key.Binding
	Help    key.Binding
	Back    key.Binding
}

var keys = keyMap{
	restore: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "restore item"),
	),
	Back: key.NewBinding(
		key.WithKeys("ctrl+c", "esc"),
		key.WithHelp("ctrl+c/esc", "back to list view"),
	),
	Next: key.NewBinding(
		key.WithKeys("→", "l", "pgdown"),
		key.WithHelp("→/l/pgdown", "next page"),
	),
	Prev: key.NewBinding(
		key.WithKeys("←", "l", "pgup"),
		key.WithHelp("←/h/pgup", "next page"),
	),
	Delete: key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("ctrl+d", "delete item"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Home: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("g", "top"),
	),
	End: key.NewBinding(
		key.WithKeys("G"),
		key.WithHelp("G", "end"),
	),
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Back}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Home, k.End},
		{k.restore, k.Delete},
		{k.Back, k.Help},
	}
}
