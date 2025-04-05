package command

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/o-kaisan/text-clipper/interface/bubbletea/constants"
)

// 画面間でWindowsSizeを共有するためのtea.Cmd
func SendWindowSizeCmd(msg tea.WindowSizeMsg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}

// list画面のリストを最新化するために別画面からlistに渡すためのtea.Cmd
func SendFetchDataCmd() tea.Cmd {
	return func() tea.Msg {
		return constants.FetchDataMsg("") // 文字列までは確認しないため暫定的に空文字
	}
}
