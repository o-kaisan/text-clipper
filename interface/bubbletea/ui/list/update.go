package list

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/o-kaisan/text-clipper/interface/bubbletea/constants"
	"github.com/o-kaisan/text-clipper/interface/bubbletea/ui/archive"
	"github.com/o-kaisan/text-clipper/interface/bubbletea/ui/register"
)

// TODO 関数分割
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case constants.FetchDataMsg:
		return handleFetchData(m)

	case tea.WindowSizeMsg:
		constants.WindowSizeMsg = msg // 別の画面にサイズを渡すため
		m.width = msg.Width
		m.height = msg.Height - 3

	case tea.KeyMsg:
		// フィルタリング入力中
		if m.list.FilterState() == list.Filtering {
			if key.Matches(msg, keys.Quit) {
				m.list.ResetFilter() // フィルタリングを解除
				return m, nil        // そのまま継続
			}
			break
		}

		// フィルタリング適用中
		if m.list.FilterState() == list.FilterApplied {
			if key.Matches(msg, keys.Quit) {
				m.list.ResetFilter() // フィルタリングを解除
				return m, nil        // そのまま継続
			}
		}

		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, keys.Help):
			m.help.ShowAll = !m.help.ShowAll

		case key.Matches(msg, keys.Edit):
			return handleEditKey(m)

		case key.Matches(msg, keys.Archive):
			return handleArchiveKey(m)

		case key.Matches(msg, keys.Add):
			return handleAddKey(m)

		case key.Matches(msg, keys.Deactivate):
			return handleDeactivateKey(m)

		case key.Matches(msg, keys.Copy):
			return handleCopyKey(m)

		case key.Matches(msg, keys.Select):
			// アイテムが無ければなにもしない
			return handleSelectKey(m)
		}
	}

	// bubble tea標準のupdate処理でフィルタリング機能を有効化する
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func handleDeactivateKey(m model) (tea.Model, tea.Cmd) {
	// アイテムが無ければなにもしない
	items := m.list.Items()
	if len(items) <= 0 {
		return m, nil
	}

	selectedItem := m.list.SelectedItem().(activeItem)
	m.cs.DeactivateClip(selectedItem.ID)

	newItems, err := m.cs.GetActiveClips()
	if err != nil {
		return nil, func() tea.Msg { return constants.ErrMsg(err) }
	}
	m.list = convertActiveClipsToListItems(newItems)
	return m, nil
}

func handleAddKey(m model) (tea.Model, tea.Cmd) {
	var cmds []interface{}
	cmds = append(cmds, constants.WindowSizeMsg)
	cmds = append(cmds, textinput.Blink)
	// 登録画面に遷移
	registerModel := register.NewRegister(m.cs, m, 0, "", "")
	return registerModel.Update(cmds)
}

func handleArchiveKey(m model) (tea.Model, tea.Cmd) {
	archiveModel, _ := archive.NewArchive(m.cs, m)
	return archiveModel.Update(constants.WindowSizeMsg)
}

func handleEditKey(m model) (tea.Model, tea.Cmd) {
	var cmds []interface{}
	// アイテムが無ければなにもしない
	items := m.list.Items()
	if len(items) <= 0 {
		return m, nil
	}
	cmds = append(cmds, constants.WindowSizeMsg)
	cmds = append(cmds, textinput.Blink)
	selectedItem := m.list.SelectedItem().(activeItem)
	registerModel := register.NewRegister(m.cs, m, selectedItem.ID, selectedItem.Title, selectedItem.Content)
	return registerModel, nil
}
func handleFetchData(m model) (tea.Model, tea.Cmd) {
	ActiveClips, err := m.cs.GetActiveClips()
	if err != nil {
		return nil, func() tea.Msg { return constants.ErrMsg(err) }
	}
	m.list = convertActiveClipsToListItems(ActiveClips) // リストを更新
	return m, nil
}

func handleSelectKey(m model) (tea.Model, tea.Cmd) {
	items := m.list.Items()
	if len(items) <= 0 {
		return m, nil
	}
	// 選択したアイテムを取得
	selectedClip := m.list.SelectedItem().(activeItem)
	err := clipboard.WriteAll(selectedClip.Content)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to clip the item to clipboard: item=%s, err=%w", selectedClip.Content, err))
	}
	// 最終利用日時を更新する

	// アイテムを選択したらアプリを閉じる
	return m, tea.Quit
}

func handleCopyKey(m model) (tea.Model, tea.Cmd) {
	// アイテムがなければなにもしない
	items := m.list.Items()
	if len(items) <= 0 {
		return m, nil
	}

	// 選択したアイテムを取得
	selectedItem := m.list.SelectedItem().(activeItem)
	m.cs.CopyClip(selectedItem.ID)

	ActiveClips, err := m.cs.GetActiveClips()
	if err != nil {
		return nil, func() tea.Msg { return constants.ErrMsg(err) }
	}

	m.list = convertActiveClipsToListItems(ActiveClips)
	return m, nil
}
