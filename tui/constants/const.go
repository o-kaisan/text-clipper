package constants

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/o-kaisan/text-clipper/text"
)

const DefaultWidth = 20
const DefaultHeight = 14

var (
	Tr            *text.GormRepository
	WindowSizeMsg tea.WindowSizeMsg
)

type keymap struct {
	// Up     key.Binding
	// Down   key.Binding
	Select key.Binding
	Submit key.Binding
	Add    key.Binding
	Quit   key.Binding
	Back   key.Binding
	Next   key.Binding
	Prev   key.Binding
	Paste  key.Binding
	Delete key.Binding
	Edit   key.Binding
}

var Keymap = keymap{
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
	Submit: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "over the submit button to register the item"),
	),
	Back: key.NewBinding(
		key.WithKeys("ctrl+c", "esc"),
		key.WithHelp("ctrl+c/esc", "back to list view"),
	),
	Next: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next input"),
	),
	Prev: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "previous input"),
	),
	Delete: key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("ctrl+d", "delete item"),
	),
}
