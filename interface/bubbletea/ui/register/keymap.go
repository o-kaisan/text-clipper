package register

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Submit key.Binding
	Back   key.Binding
	Next   key.Binding
	Prev   key.Binding
	Help   key.Binding
}

var keys = keyMap{
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
}

// ショートヘルプは使用しないため、空のkey.Bindingを返す
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{}
}
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Next, k.Prev, k.Submit, k.Back},
	}
}
