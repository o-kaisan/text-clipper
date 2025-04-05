package list

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Deactivate key.Binding
	Archive    key.Binding
	Up         key.Binding
	Down       key.Binding
	Select     key.Binding
	Add        key.Binding
	Quit       key.Binding
	Copy       key.Binding
	Paste      key.Binding
	Next       key.Binding
	Prev       key.Binding
	Edit       key.Binding
	Help       key.Binding
	Home       key.Binding
	End        key.Binding
}

var keys = keyMap{
	Archive: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "move to archive view"),
	),
	Deactivate: key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("ctrl+d", "archive item"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select item"),
	),
	Add: key.NewBinding(
		key.WithKeys("ctrl+a"),
		key.WithHelp("ctrl+a", "add new item"),
	),
	Edit: key.NewBinding(
		key.WithKeys("ctrl+e"),
		key.WithHelp("ctrl+e", "edit item"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "esc"),
		key.WithHelp("ctrl+c/esc", "quit"),
	),
	Next: key.NewBinding(
		key.WithKeys("→", "l", "pgdown"),
		key.WithHelp("→/l/pgdown", "next page"),
	),
	Prev: key.NewBinding(
		key.WithKeys("←", "l", "pgup"),
		key.WithHelp("←/h/pgup", "next page"),
	),
	Copy: key.NewBinding(
		key.WithKeys("ctrl+y"),
		key.WithHelp("ctrl+y", "copy item."),
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
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Home, k.End},
		{k.Add, k.Edit, k.Copy, k.Deactivate},
		{k.Select, k.Archive, k.Quit, k.Help},
	}
}
