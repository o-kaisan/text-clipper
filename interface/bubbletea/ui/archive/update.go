package archive

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/o-kaisan/text-clipper/interface/bubbletea/command"
	"github.com/o-kaisan/text-clipper/interface/bubbletea/constants"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		constants.WindowSizeMsg = msg // 別の画面にサイズを渡すため
		m.width = msg.Width
		m.height = msg.Height - 3

	case tea.KeyMsg:
		// フィルタリング入力中
		if m.list.FilterState() == list.Filtering {
			if key.Matches(msg, keys.Back) {
				m.list.ResetFilter() // フィルタリングを解除
				return m, nil        // そのまま継続
			}
			break
		}

		// フィルタリング適用中
		if m.list.FilterState() == list.FilterApplied {
			if key.Matches(msg, keys.Back) {
				m.list.ResetFilter() // フィルタリングを解除
				return m, nil        // そのまま継続
			}
		}

		switch {
		case key.Matches(msg, keys.Back):
			return moveToListView(m)

		case key.Matches(msg, keys.Help):
			m.help.ShowAll = !m.help.ShowAll

		case key.Matches(msg, keys.Delete):
			return handleDeleteKey(m)

		case key.Matches(msg, keys.restore):
			return handleRestoreKey(m)
		}
	}

	// bubble tea標準のupdate処理でフィルタリング機能を有効化する
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func moveToListView(m model) (tea.Model, tea.Cmd) {
	return m.listModel, tea.Batch(command.SendWindowSizeCmd(constants.WindowSizeMsg), command.SendFetchDataCmd())
}

func handleDeleteKey(m model) (tea.Model, tea.Cmd) {
	// アイテムが無ければなにもしない
	clips := m.list.Items()
	if len(clips) <= 0 {
		return m, nil
	}

	// 対象のitemを削除する
	selectedClip := m.list.SelectedItem().(archivedClip)
	m.cs.DeleteClip(selectedClip.ID)

	inActiveClips, err := m.cs.GetArchivedClips()
	if err != nil {
		return nil, func() tea.Msg { return constants.ErrMsg(err) }
	}

	m.list = convertArchivedClipsToListItems(inActiveClips)

	return m, nil
}

func handleRestoreKey(m model) (tea.Model, tea.Cmd) {
	// アイテムが無ければなにもしない
	clips := m.list.Items()
	if len(clips) <= 0 {
		return moveToListView(m)
	}
	// 選択したアイテムを取得
	selectedClip := m.list.SelectedItem().(archivedClip)
	err := clipboard.WriteAll(selectedClip.Content)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to clip the item to clipboard: item=%s, err=%w", selectedClip.Content, err))
	}
	// アイテムを有効化する
	m.cs.ActivateClip(selectedClip.ID)

	newClips, err := m.cs.GetArchivedClips()
	if err != nil {
		return nil, func() tea.Msg { return constants.ErrMsg(err) }
	}
	m.list = convertArchivedClipsToListItems(newClips)
	return m, nil

}
