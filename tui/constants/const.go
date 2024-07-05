package constants

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/o-kaisan/text-clipper/text"
)

const DefaultWidth = 20
const DefaultHeight = 14

var (
	Tr            *text.GormRepository
	WindowSizeMsg tea.WindowSizeMsg
)
