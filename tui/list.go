package tui

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/o-kaisan/text-clipper/common"
	"github.com/o-kaisan/text-clipper/item"
	"github.com/o-kaisan/text-clipper/tui/constants"
)

// --------------------------------------------------------------------------------
// Help
// --------------------------------------------------------------------------------
type listKeyMap struct {
	Deactivate key.Binding
	Archive    key.Binding
	Up         key.Binding
	Down       key.Binding
	Select     key.Binding
	Add        key.Binding
	Quit       key.Binding
	Copy       key.Binding
	Paste      key.Binding
	Next       key.Binding
	Prev       key.Binding
	Edit       key.Binding
	Help       key.Binding
	Home       key.Binding
	End        key.Binding
}

var listKeys = listKeyMap{
	Archive: key.NewBinding(
		key.WithKeys("ctrl+l"),
		key.WithHelp("ctrl+l", "archive item"),
	),
	Deactivate: key.NewBinding(
		key.WithKeys("delete"),
		key.WithHelp("delete", "move to archive item list"),
	),
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
	Next: key.NewBinding(
		key.WithKeys("→", "l", "pgdown"),
		key.WithHelp("→/l/pgdown", "next page"),
	),
	Prev: key.NewBinding(
		key.WithKeys("←", "l", "pgup"),
		key.WithHelp("←/h/pgup", "next page"),
	),
	Copy: key.NewBinding(
		key.WithKeys("ctrl+v"),
		key.WithHelp("ctrl+v", "copy item."),
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
		key.WithHelp("g", "top"),
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
		{k.Up, k.Down, k.Home, k.End},
		{k.Add, k.Edit, k.Copy, k.Deactivate},
		{k.Select, k.Help, k.Quit},
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
	*item.Item
}

func (l Item) FilterValue() string {
	return l.Item.Title
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
	items, err := getActiveItemList(constants.Ir)
	if err != nil {
		m.err = fmt.Errorf("failed to initial choice: %w", err)
		return m, nil
	}

	// *item.Item のスライスを ListItem のスライスに変換
	l := convertItemsToListItems(items)
	m.list = l
	m.help = help.New()
	m.help.ShowAll = true // フルヘルプをはじめから表示する
	return m, func() tea.Msg { return errMsg(err) }
}

func convertItemsToListItems(items []*item.Item) list.Model {
	listItems := make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = Item{Item: item}
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
			// アイテムが無ければなにもしない
			items := m.list.Items()
			if len(items) <= 0 {
				choice, _ := InitialList()
				return choice.Update(constants.WindowSizeMsg)
			}
			cmds = append(cmds, constants.WindowSizeMsg)
			cmds = append(cmds, textinput.Blink)
			selectedItem := m.list.SelectedItem().(Item)
			targetItem := constants.Ir.FindByID(selectedItem.ID)
			register := InitialRegister(targetItem)
			return register.Update(cmds)

		case key.Matches(msg, listKeys.Archive):
			choice, _ := InitialArchive()
			return choice.Update(constants.WindowSizeMsg)

		case key.Matches(msg, listKeys.Add):
			var cmds []interface{}
			cmds = append(cmds, constants.WindowSizeMsg)
			cmds = append(cmds, textinput.Blink)
			var unsetTime time.Time // 登録時に時刻を更新するのでここではゼロ値を設定する
			initialItem := item.NewItem("", "", constants.True, unsetTime, unsetTime, unsetTime)
			// 登録画面に遷移
			register := InitialRegister(initialItem)
			return register.Update(cmds)

		case key.Matches(msg, listKeys.Deactivate):
			// アイテムが無ければなにもしない
			items := m.list.Items()
			if len(items) <= 0 {
				choice, _ := InitialList()
				return choice.Update(constants.WindowSizeMsg)
			}

			// 選択したアイテムを取得
			selectedItem := m.list.SelectedItem().(Item)
			choice := constants.Ir.FindByID(selectedItem.ID)
			archiveItem(constants.Ir, choice)

			ActiveItems, err := getActiveItemList(constants.Ir)
			if err != nil {
				return nil, func() tea.Msg { return errMsg(err) }
			}
			m.list = convertItemsToListItems(ActiveItems)

		case key.Matches(msg, listKeys.Copy):
			// アイテムが無ければなにもしない
			items := m.list.Items()
			if len(items) <= 0 {
				choice, _ := InitialList()
				return choice.Update(constants.WindowSizeMsg)
			}
			// 選択したアイテムを取得
			selectedItem := m.list.SelectedItem().(Item)
			constants.Ir.Copy(selectedItem.ID) // 複製機能の呼び出し

			ActiveItems, err := getActiveItemList(constants.Ir)
			if err != nil {
				return nil, func() tea.Msg { return errMsg(err) }
			}
			m.list = convertItemsToListItems(ActiveItems) // リストを更新

		case key.Matches(msg, listKeys.Select):
			// アイテムが無ければなにもしない
			items := m.list.Items()
			if len(items) <= 0 {
				return m, tea.Quit
			}
			// 選択したアイテムを取得
			choice := m.list.SelectedItem().(Item)
			err := clipboard.WriteAll(choice.Content)
			if err != nil {
				fmt.Println(fmt.Errorf("failed to clip the item to clipboard: item=%s, err=%w", choice.Content, err))
			}
			// 最終利用日時を更新する
			target := constants.Ir.FindByID(choice.ID)
			target.LastUsedAt = time.Now()
			constants.Ir.Update(target)

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

func archiveItem(ir *item.ItemRepository, item *item.Item) error {
	item.IsActive = constants.False // 無効化する
	err := ir.Update(item)
	if err != nil {
		return fmt.Errorf("cannot delete item: title=%s id=%d err=%w", item.Title, item.ID, err)
	}
	return nil
}

func deleteItem(ir *item.ItemRepository, item *item.Item) error {
	err := ir.Delete(item)
	if err != nil {
		return fmt.Errorf("cannot delete item: title=%s id=%d err=%w", item.Title, item.ID, err)
	}
	return nil
}

func getActiveItemList(ir *item.ItemRepository) ([]*item.Item, error) {
	order := common.Env("TEXT_CLIPPER_SORT", "createdAtDesc")
	items, err := ir.ListOfActive(order)
	if err != nil {
		return nil, fmt.Errorf("cannot get all items: %w", err)
	}
	return items, nil
}

// --------------------------------------------------------------------------------
// View
// --------------------------------------------------------------------------------
func (m model) View() string {
	mainView := lipgloss.NewStyle()

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
		// 一部のターミナルで変な改行が入ってしまうので調整用のpadding
		adjustedPadding := 8
		if m.list.SelectedItem() != nil {
			selectedItem := m.list.SelectedItem().(Item)
			// タイトルとコンテンツの区切り線
			line := strings.Repeat(" ", width-(preViewPadding*2)-adjustedPadding)
			titleStyle := lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, true, false)

			preview = selectedItem.Title + titleStyle.Render(line) + "\n" + truncateString(selectedItem.Content, maxLength, maxLines)
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
