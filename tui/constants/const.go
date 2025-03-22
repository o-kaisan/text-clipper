package constants

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/o-kaisan/text-clipper/item"
)

var (
	Ir            *item.ItemRepository
	WindowSizeMsg tea.WindowSizeMsg
	True          = boolPtr(true)
	False         = boolPtr(false)
)

// bool のポインタを作成するヘルパー関数
func boolPtr(v bool) *bool {
	return &v
}
