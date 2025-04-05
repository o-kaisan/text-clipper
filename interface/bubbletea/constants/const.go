package constants

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/o-kaisan/text-clipper/common"
)

type (
	ErrMsg       error
	FetchDataMsg string
)

var (
	// 画面間で共通して使う画面サイズ
	WindowSizeMsg tea.WindowSizeMsg
	// リスト画面のプレビューとページタイトルの背景色
	BgColor = getBgColor()
	// アーカイブ画面のプレビューとページタイトルの背景色
	ArchiveBgColor = getArchiveBgColor()
	// クリップのタイトルに付けられる最大文字数
	TitleMaxLength = 30
	ListVewWidth   = 33
	// 一覧画面の幅を調整するためのwidth
	AdjustedWidth = 41
)

var (
	FilterPrompt = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#ECFD65"})
	FilterCursor = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"})
)

func getBgColor() string {
	bgColor := common.Env("TEXT_CLIPPER_BG_COLOR", "#696969")
	return bgColor
}

func getArchiveBgColor() string {
	bgColor := common.Env("TEXT_CLIPPER_ARCHIVE_BG_COLOR", "#af00af")
	return bgColor
}
