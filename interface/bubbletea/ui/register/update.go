package register

import (
	"log"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/o-kaisan/text-clipper/interface/bubbletea/command"
	"github.com/o-kaisan/text-clipper/interface/bubbletea/constants"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		constants.WindowSizeMsg = msg
		m.width = msg.Width - 4
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Back):
			return moveToListView(m)

			// return m.listModel.Update(constants.WindowSizeMsg)
		case key.Matches(msg, keys.Submit):
			// submitにフォーカスがある場合に登録処理を実行する
			if m.focusedIndex == m.maxIndex {
				return handleSubmitKey(m)
			}

		case key.Matches(msg, keys.Next):
			return handleNextKey(m)

		case key.Matches(msg, keys.Prev):
			return handlePrevKey(m)
		}
	}

	var cmd tea.Cmd
	switch m.focusedIndex {
	case 0:
		m.title, cmd = m.title.Update(msg)
	case 1:
		m.content, cmd = m.content.Update(msg)
	}

	return m, cmd
}

func handlePrevKey(m model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if m.focusedIndex > 0 {
		m.focusedIndex--
	}
	// 一旦フォーカスを解除
	m.title.Blur()
	m.content.Blur()

	// タイトルにフォーカス
	if m.focusedIndex == 0 {
		cmd = m.title.Focus()
	}
	// コンテンツにフォーカス
	if m.focusedIndex == 1 {
		cmd = m.content.Focus()
	}
	return m, cmd
}

func handleNextKey(m model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if m.focusedIndex < m.maxIndex {
		m.focusedIndex++
	}
	// 一旦フォーカスを解除
	m.title.Blur()
	m.content.Blur()

	// テキストエリアにフォーカス
	if m.focusedIndex == 1 {
		cmd = m.content.Focus()
	}
	return m, cmd
}

func handleSubmitKey(m model) (tea.Model, tea.Cmd) {
	err := m.cs.RegisterClip(m.clipId, m.title.Value(), m.content.Value())
	if err != nil {
		log.Fatal(err)
	}
	// 元の画面に戻る
	return moveToListView(m)
}

func moveToListView(m model) (tea.Model, tea.Cmd) {
	return m.listModel, tea.Batch(command.SendWindowSizeCmd(constants.WindowSizeMsg), command.SendFetchDataCmd())
}
