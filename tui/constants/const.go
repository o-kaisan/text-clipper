package constants

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/o-kaisan/text-clipper/common"
	"github.com/o-kaisan/text-clipper/item"
)

var (
	Ir             *item.ItemRepository
	WindowSizeMsg  tea.WindowSizeMsg
	True           = boolPtr(true)
	False          = boolPtr(false)
	BgColor        = getBgColor()
	ArchiveBgColor = getArchiveBgColor()
)

// bool のポインタを作成するヘルパー関数
func boolPtr(v bool) *bool {
	return &v
}

func getBgColor() string {
	bgColor := common.Env("TEXT_CLIPPER_BG_COLOR", "#696969")
	return bgColor
}

func getArchiveBgColor() string {
	bgColor := common.Env("TEXT_CLIPPER_ARCHIVE_BG_COLOR", "#af00af")
	return bgColor
}
