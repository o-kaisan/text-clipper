package tui

import (
	"fmt"
	"io"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/o-kaisan/text-clipper/text"
	"github.com/o-kaisan/text-clipper/tui/constants"
)

// ---------------------------------------------------------------
var (
	titleStyle          = lipgloss.NewStyle().MarginLeft(2)
	itemStyle           = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle   = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle     = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	choiceViewHelpStyle = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

// ---------------------------------------------------------------
// Itemはlist.Itemインターフェースを実装するためのラッパー
// ---------------------------------------------------------------
type Item struct {
	*text.Text
}

func (l Item) FilterValue() string {
	return l.Text.Title
}

// ---------------------------------------------------------------
// リスト表示
// ---------------------------------------------------------------
type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Item)
	if !ok {
		return
	}
	str := fmt.Sprintf("%d. %s", index+1, i.Title)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func getTextList(tr *text.GormRepository) ([]*text.Text, error) {
	texts, err := tr.List()
	if err != nil {
		return nil, fmt.Errorf("cannot get all texts: %w", err)
	}
	return texts, nil
}

// ---------------------------------------------------------------
type model struct {
	width  int
	height int
	list   list.Model
	err    error
}

func InitialList() (model, tea.Cmd) {
	var m model

	texts, err := getTextList(constants.Tr)
	if err != nil {
		m.err = fmt.Errorf("failed to initial choice: %w", err)
		return m, nil
	}

	// *text.Text のスライスを ListItem のスライスに変換
	// 検索機能有効化
	l := convertTextsToListItems(texts)

	m.list = l

	return m, func() tea.Msg { return errMsg(err) }
}

func convertTextsToListItems(texts []*text.Text) list.Model {
	listItems := make([]list.Item, len(texts))
	for i, txt := range texts {
		listItems[i] = Item{Text: txt}
	}
	l := list.New(listItems, itemDelegate{}, 0, 0) // リストのサイズはUpdateで決定する
	l.Title = "Select the item you want to copy to clipboard"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = choiceViewHelpStyle
	l.SetItems(listItems)
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			constants.Keymap.Select,
			constants.Keymap.Add,
			constants.Keymap.Delete,
		}
	}
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			constants.Keymap.Select,
			constants.Keymap.Add,
			constants.Keymap.Delete,
		}
	}
	return l
}

func (m model) Init() tea.Cmd {
	return nil
}

// --------------------------------------------------------------------------------
// update
// --------------------------------------------------------------------------------
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		constants.WindowSizeMsg = msg
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, constants.Keymap.Quit):
			return m, tea.Quit
		case key.Matches(msg, constants.Keymap.Edit):
			var cmds []interface{}
			cmds = append(cmds, constants.WindowSizeMsg)
			cmds = append(cmds, textinput.Blink)
			selectedItem := m.list.SelectedItem().(Item)
			targetText := constants.Tr.FindByID(selectedItem.ID)
			register := InitialRegister(targetText)
			return register.Update(cmds)

		case key.Matches(msg, constants.Keymap.Add):
			var cmds []interface{}
			cmds = append(cmds, constants.WindowSizeMsg)
			cmds = append(cmds, textinput.Blink)
			// 登録画面に遷移
			initialText := &text.Text{
				Title:    "",
				Contents: "",
			}
			register := InitialRegister(initialText)
			return register.Update(cmds)

		case key.Matches(msg, constants.Keymap.Delete):
			// 対象のitemを削除する
			selectedItem := m.list.SelectedItem().(Item)
			choice := constants.Tr.FindByID(selectedItem.ID)
			deleteText(constants.Tr, choice)

			texts, err := getTextList(constants.Tr)
			if err != nil {
				return nil, func() tea.Msg { return errMsg(err) }
			}
			m.list = convertTextsToListItems(texts)

		case key.Matches(msg, constants.Keymap.Select):

			choice := m.list.SelectedItem().(Item)
			err := clipboard.WriteAll(choice.Contents)
			if err != nil {
				fmt.Println(fmt.Errorf("failed to clip the text to clipboard: text=%s, err=%w", choice.Contents, err))
			}
			return m, tea.Quit
		}
	}

	// This will also call our delegate's update function.
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func deleteText(tr *text.GormRepository, text *text.Text) error {
	err := tr.Delete(text)
	if err != nil {
		return fmt.Errorf("cannot delete item: title=%s id=%d err=%w", text.Title, text.ID, err)
	}
	return nil
}

// --------------------------------------------------------------------------------
// view
// --------------------------------------------------------------------------------
func (m model) View() string {
	m.list.SetSize(m.width/3, m.height) // リストの画面サイズ
	return lipgloss.JoinHorizontal(lipgloss.Top, m.choicesView(), m.previewView(m.width/2, m.height-4))
}

func (m model) previewView(width int, height int) string {
	contents := ""

	// フィルタリングでItemがあれば選択しているItemの内容をプレビューに表示する
	if m.list.FilterState() != list.Filtering {
		selectedItem := m.list.SelectedItem().(Item)
		contents = selectedItem.Contents
	}

	previewStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true).
		BorderForeground(lipgloss.Color("#6272A4"))
	return previewStyle.Width(width).Height(height).Render(contents)
}

func (m model) choicesView() string {
	return m.list.View()
}