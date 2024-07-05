package tui

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/o-kaisan/text-clipper/text"
	"github.com/o-kaisan/text-clipper/tui/constants"
)

// --------------------------------------------------------------------------------
// Help
// --------------------------------------------------------------------------------
type listKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Add    key.Binding
	Quit   key.Binding
	Paste  key.Binding
	Delete key.Binding
	Edit   key.Binding
	Help   key.Binding
	Home   key.Binding
	End    key.Binding
}

var listKeys = listKeyMap{
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select item"),
	),
	Add: key.NewBinding(
		key.WithKeys("ctrl+a"),
		key.WithHelp("ctrl+a", "add new item"),
	),
	Edit: key.NewBinding(
		key.WithKeys("ctrl+e"),
		key.WithHelp("ctrl+e", "edit item"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "esc"),
		key.WithHelp("ctrl+c/esc", "quit"),
	),
	Delete: key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("ctrl+d", "delete item"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Home: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("g", "home"),
	),
	End: key.NewBinding(
		key.WithKeys("G"),
		key.WithHelp("G", "end"),
	),
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k listKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k listKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Home, k.End},
		{k.Select, k.Help},
		{k.Delete, k.Edit, k.Quit},
	}
}

// ---------------------------------------------------------------
// Style
// ---------------------------------------------------------------
var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	noItemStyle       = lipgloss.NewStyle().PaddingLeft(2).Width(38)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	preViewPadding    = 1
	previewStyle      = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, false, true).MarginLeft(4).MarginRight(1).Foreground(lipgloss.Color("#FAFAFA")).Background(lipgloss.Color("#696969")).MarginTop(1).Padding(preViewPadding)
	listHelpStyle     = lipgloss.NewStyle().PaddingLeft(2).PaddingTop(1).PaddingBottom(1).Height(5)
)

// ---------------------------------------------------------------
// Delegate
// ---------------------------------------------------------------

// Itemはlist.Itemインターフェースを実装するためのラッパー
type Item struct {
	*text.Text
}

func (l Item) FilterValue() string {
	return l.Text.Title
}

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {

	i, ok := listItem.(Item)
	if !ok {
		return
	}

	// 画面全体を使うようにフォーマットを変更する
	choicesWidth := m.Width() * 2 / 3
	itemNum := index + 1

	// リストのタイトルが揃うように10個未満は空白を2つ入れる
	var str string
	if itemNum >= 10 {
		str = fmt.Sprintf("%d. %-"+strconv.Itoa(choicesWidth)+"s", itemNum, i.Title)
	} else {
		str = fmt.Sprintf("%d.  %-"+strconv.Itoa(choicesWidth)+"s", itemNum, i.Title)
	}

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

// ---------------------------------------------------------------
// Model
// ---------------------------------------------------------------
type model struct {
	width  int
	help   help.Model
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
	m.help = help.New()
	return m, func() tea.Msg { return errMsg(err) }
}

func convertTextsToListItems(texts []*text.Text) list.Model {
	listItems := make([]list.Item, len(texts))
	for i, text := range texts {
		listItems[i] = Item{Text: text}
	}

	// リストモデルの初期化 //
	// リストの高さと幅は Update と View で決定するため0で初期化する
	l := list.New(listItems, itemDelegate{}, 0, 0)
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.NoItems = noItemStyle
	l.Styles.PaginationStyle = paginationStyle
	l.SetShowHelp(false)  // helpは独自で定義するためここで明示的に無効化する
	l.SetShowTitle(false) // Titleはここで明示的に無効化する

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
		constants.WindowSizeMsg = msg // 別の画面にサイズを渡すため
		m.width = msg.Width - 5
		m.height = msg.Height - 3
	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, listKeys.Quit):
			return m, tea.Quit
		case key.Matches(msg, listKeys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, listKeys.Edit):
			var cmds []interface{}
			cmds = append(cmds, constants.WindowSizeMsg)
			cmds = append(cmds, textinput.Blink)
			selectedItem := m.list.SelectedItem().(Item)
			targetText := constants.Tr.FindByID(selectedItem.ID)
			register := InitialRegister(targetText)
			return register.Update(cmds)

		case key.Matches(msg, listKeys.Add):
			var cmds []interface{}
			cmds = append(cmds, constants.WindowSizeMsg)
			cmds = append(cmds, textinput.Blink)
			// 登録画面に遷移
			initialText := &text.Text{
				Title:   "",
				Content: "",
			}
			register := InitialRegister(initialText)
			return register.Update(cmds)

		case key.Matches(msg, listKeys.Delete):
			// 対象のitemを削除する
			selectedItem := m.list.SelectedItem().(Item)
			choice := constants.Tr.FindByID(selectedItem.ID)
			deleteText(constants.Tr, choice)

			texts, err := getTextList(constants.Tr)
			if err != nil {
				return nil, func() tea.Msg { return errMsg(err) }
			}
			m.list = convertTextsToListItems(texts)

		case key.Matches(msg, listKeys.Select):

			choice := m.list.SelectedItem().(Item)
			err := clipboard.WriteAll(choice.Content)
			if err != nil {
				fmt.Println(fmt.Errorf("failed to clip the text to clipboard: text=%s, err=%w", choice.Content, err))
			}
			// アイテムを選択したらアプリを閉じる
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

func getTextList(tr *text.GormRepository) ([]*text.Text, error) {
	texts, err := tr.List()
	if err != nil {
		return nil, fmt.Errorf("cannot get all texts: %w", err)
	}
	return texts, nil
}

// --------------------------------------------------------------------------------
// View
// --------------------------------------------------------------------------------
func (m model) View() string {
	mainView := lipgloss.NewStyle().Border(lipgloss.RoundedBorder())

	var choicesWidth int
	var previewWidth int
	var helpWidth int
	if m.width <= 80 {
		// 最低限の幅を設定する
		choicesWidth = 45
		previewWidth = 35
		helpWidth = 80
	} else {
		// choicesViewが画面の2/3の幅を使用
		// choicesWidth = m.width * 1 / 3
		choicesWidth = 45
		// previewViewが画面の1/3の幅を使用
		previewWidth = m.width - choicesWidth
		helpWidth = m.width
	}

	choices := m.choicesView(choicesWidth, m.height-6)
	preview := m.previewView(previewWidth, m.height-6)
	help := m.helpView(helpWidth)

	return mainView.Render(lipgloss.JoinHorizontal(lipgloss.Top, choices, preview) + "\n" + help)
}

func (m model) choicesView(width int, height int) string {
	// リストの画面サイズ
	m.list.SetWidth(width)
	m.list.SetHeight(height) // m.list.Styles.Title.Width(30)

	return lipgloss.NewStyle().Render(m.list.View())
}

func (m model) previewView(width int, height int) string {
	preview := ""
	maxLength := 135
	maxLines := 2
	if height > 0 {
		maxLength = calculateMaxLength(height)
		maxLines = calculateMaxLines(height)
	}

	// フィルタリングでItemがあれば選択しているItemの内容をプレビューに表示する
	if m.list.FilterState() != list.Filtering {
		if m.list.SelectedItem() != nil {
			selectedItem := m.list.SelectedItem().(Item)
			line := strings.Repeat("-", width-(preViewPadding*2)) // タイトルとコンテンツの区切り線
			preview = selectedItem.Title + "\n" + line + "\n" + truncateString(selectedItem.Content, maxLength, maxLines)
		}
	}

	return previewStyle.Width(width).Height(height).Render(preview)
}

func calculateMaxLines(height int) int {
	// 8の近辺で1または2になる指数関数的な式
	if height <= 4 {
		return 1
	} else if height <= 8 {
		return int(math.Pow(float64(height), 0.2))
	} else if height <= 10 {
		return int(math.Pow(float64(height), 0.5)) + 1
	} else if height <= 15 {
		return int(math.Pow(float64(height), 0.8)) + 1
	}
	return int(math.Pow(float64(height), 0.89)) + 1
}

// 高さに応じて指数関数的にmaxLengthを計算する関数
func calculateMaxLength(height int) int {
	return int(math.Pow(float64(height), 1.3)) * 10
}

func truncateString(s string, maxLength int, maxLines int) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")

	truncatedLines := lines
	if len(lines) > maxLines {
		truncatedLines = lines[:maxLines]
		truncated := strings.Join(truncatedLines, "\n")
		return truncated + "\n" + "..."
	}

	if len(lines) == 1 && len(s) > maxLength {
		truncated := s[:maxLength-3]
		return truncated + "..."
	}

	return s
}

func (m model) helpView(width int) string {
	return listHelpStyle.Width(width).Render(m.help.View(listKeys))
}
