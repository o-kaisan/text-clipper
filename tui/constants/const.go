package constants

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/o-kaisan/text-clipper/common"
)

var (
	WindowSizeMsg  tea.WindowSizeMsg
	BgColor        = getBgColor()
	ArchiveBgColor = getArchiveBgColor()
)

func getBgColor() string {
	bgColor := common.Env("TEXT_CLIPPER_BG_COLOR", "#696969")
	return bgColor
}

func getArchiveBgColor() string {
	bgColor := common.Env("TEXT_CLIPPER_ARCHIVE_BG_COLOR", "#af00af")
	return bgColor
}
