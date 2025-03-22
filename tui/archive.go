package tui

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/o-kaisan/text-clipper/common"
	"github.com/o-kaisan/text-clipper/item"
	"github.com/o-kaisan/text-clipper/tui/constants"
)

// --------------------------------------------------------------------------------
// Help
// --------------------------------------------------------------------------------
type archiveKeyMap struct {
	restore key.Binding
	Delete  key.Binding
	Up      key.Binding
	Down    key.Binding
	Next    key.Binding
	Prev    key.Binding
	Home    key.Binding
	End     key.Binding
	Help    key.Binding
	Back    key.Binding
}

var archiveKeys = archiveKeyMap{
	restore: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "restore item"),
	),
	Back: key.NewBinding(
		key.WithKeys("ctrl+c", "esc"),
		key.WithHelp("ctrl+c/esc", "back to list view"),
	),
	Next: key.NewBinding(
		key.WithKeys("→", "l", "pgdown"),
		key.WithHelp("→/l/pgdown", "next page"),
	),
	Prev: key.NewBinding(
		key.WithKeys("←", "l", "pgup"),
		key.WithHelp("←/h/pgup", "next page"),
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
		key.WithHelp("g", "top"),
	),
	End: key.NewBinding(
		key.WithKeys("G"),
		key.WithHelp("G", "end"),
	),
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k archiveKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Back}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k archiveKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Home, k.End},
		{k.restore, k.Delete},
		{k.Back, k.Help},
	}
}

// ---------------------------------------------------------------
// Style
// ---------------------------------------------------------------
var (
	archiveItemStyle          = lipgloss.NewStyle().PaddingLeft(4)
	archiveNoItemStyle        = lipgloss.NewStyle().PaddingLeft(2).Width(38)
	archiveSelectedItemStyle  = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("127"))
	archivePaginationStyle    = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	archivePreViewPadding     = 1
	archiveListTitleStyle     = lipgloss.NewStyle().Reverse(true).PaddingLeft(1).Italic(true).Width(18).Foreground(lipgloss.Color("#af00af"))
	archiveListTitleViewStyle = lipgloss.NewStyle().PaddingLeft(1).PaddingBottom(1)
	archivePreviewStyle       = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, false, true).MarginLeft(4).MarginRight(1).Foreground(lipgloss.Color("#FAFAFA")).Background(lipgloss.Color("#af00af")).MarginTop(1).Padding(archivePreViewPadding)
	archiveListHelpStyle      = lipgloss.NewStyle().PaddingLeft(2).PaddingTop(1).PaddingBottom(1).Height(5)
)

// ---------------------------------------------------------------
// Delegate
// ---------------------------------------------------------------

// ArchiveItemはlist.ArchiveItemインターフェースを実装するためのラッパー
type ArchiveItem struct {
	*item.Item
}

func (l ArchiveItem) FilterValue() string {
	return l.Item.Title
}

type ArchiveItemDelegate struct{}

func (d ArchiveItemDelegate) Height() int                             { return 1 }
func (d ArchiveItemDelegate) Spacing() int                            { return 0 }
func (d ArchiveItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d ArchiveItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {

	i, ok := listItem.(ArchiveItem)
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

	fn := archiveItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return archiveSelectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

// ---------------------------------------------------------------
// Model
// ---------------------------------------------------------------
type archive struct {
	width  int
	help   help.Model
	height int
	list   list.Model
	err    error
}

func InitialArchive() (archive, tea.Cmd) {
	var m archive
	items, err := getInActiveItemList(constants.Ir)
	if err != nil {
		m.err = fmt.Errorf("failed to initial choice: %w", err)
		return m, nil
	}

	// *item.Item のスライスを ListItem のスライスに変換
	l := convertArchiveItemsToListItems(items)
	m.list = l
	m.help = help.New()
	m.help.ShowAll = true // フルヘルプをはじめから表示する
	return m, func() tea.Msg { return errMsg(err) }
}

func convertArchiveItemsToListItems(items []*item.Item) list.Model {
	listItems := make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = ArchiveItem{Item: item}
	}

	// リストモデルの初期化 //
	// リストの高さと幅は Update と View で決定するため0で初期化する
	l := list.New(listItems, ArchiveItemDelegate{}, 0, 0)
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.NoItems = archiveNoItemStyle
	l.Styles.PaginationStyle = archivePaginationStyle
	l.SetShowHelp(false)  // helpは独自で定義するためここで明示的に無効化する
	l.SetShowTitle(false) // Titleはここで明示的に無効化する

	return l
}

func (m archive) Init() tea.Cmd {
	return nil
}

// --------------------------------------------------------------------------------
// update
// --------------------------------------------------------------------------------
func (m archive) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case key.Matches(msg, archiveKeys.Back):
			choice, _ := InitialList()
			return choice.Update(constants.WindowSizeMsg)

		case key.Matches(msg, archiveKeys.Help):
			m.help.ShowAll = !m.help.ShowAll

		case key.Matches(msg, archiveKeys.Delete):
			// アイテムが無ければなにもしない
			items := m.list.Items()
			if len(items) <= 0 {
				choice, _ := InitialList()
				return choice.Update(constants.WindowSizeMsg)
			}

			// 対象のitemを削除する
			selectedItem := m.list.SelectedItem().(ArchiveItem)
			choice := constants.Ir.FindByID(selectedItem.ID)
			deleteItem(constants.Ir, choice)

			inActiveItems, err := getInActiveItemList(constants.Ir)
			if err != nil {
				return nil, func() tea.Msg { return errMsg(err) }
			}
			m.list = convertArchiveItemsToListItems(inActiveItems)

		case key.Matches(msg, archiveKeys.restore):
			// アイテムが無ければなにもしない
			items := m.list.Items()
			if len(items) <= 0 {
				choice, _ := InitialList()
				return choice.Update(constants.WindowSizeMsg)
			}
			// 選択したアイテムを取得
			choice := m.list.SelectedItem().(ArchiveItem)
			err := clipboard.WriteAll(choice.Content)
			if err != nil {
				fmt.Println(fmt.Errorf("failed to clip the item to clipboard: item=%s, err=%w", choice.Content, err))
			}
			// アイテムを有効化する
			target := constants.Ir.FindByID(choice.ID)
			target.IsActive = constants.True
			constants.Ir.Update(target)

			updatedItems, err := getInActiveItemList(constants.Ir)
			if err != nil {
				return nil, func() tea.Msg { return errMsg(err) }
			}
			m.list = convertArchiveItemsToListItems(updatedItems)
		}
	}

	// This will also call our delegate's update function.
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func getInActiveItemList(ir *item.ItemRepository) ([]*item.Item, error) {
	order := common.Env("TEXT_CLIPPER_SORT", "createdAtDesc")
	items, err := ir.ListOfInactive(order)
	if err != nil {
		return nil, fmt.Errorf("cannot get all items: %w", err)
	}
	return items, nil
}

// --------------------------------------------------------------------------------
// View
// --------------------------------------------------------------------------------
func (m archive) View() string {
	mainView := lipgloss.NewStyle()

	var titleWidth int
	var choicesWidth int
	var previewWidth int
	var helpWidth int
	if m.width <= 80 {
		// 最低限の幅を設定する
		titleWidth = 85
		choicesWidth = 45
		previewWidth = 35
		helpWidth = 85
	} else {
		titleWidth = m.width
		// choicesViewが画面の2/3の幅を使用
		// choicesWidth = m.width * 1 / 3
		choicesWidth = 45
		// previewViewが画面の1/3の幅を使用
		previewWidth = m.width - choicesWidth
		helpWidth = m.width
	}

	title := m.titleView(titleWidth)
	choices := m.choicesView(choicesWidth, m.height-6)
	preview := m.previewView(previewWidth, m.height-6)
	help := m.helpView(helpWidth)

	return mainView.Render(title + lipgloss.JoinHorizontal(lipgloss.Top, choices, preview) + "\n" + help)
}

func (m archive) titleView(width int) string {
	title := archiveListTitleStyle.Render("# Archived Items")
	return archiveListTitleViewStyle.Width(width - 3).Render(title)
}

func (m archive) choicesView(width int, height int) string {
	// リストの画面サイズ
	m.list.SetWidth(width)
	m.list.SetHeight(height)
	return lipgloss.NewStyle().Render(m.list.View())
}

func (m archive) previewView(width int, height int) string {

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
			selectedItem := m.list.SelectedItem().(ArchiveItem)
			// タイトルとコンテンツの区切り線
			line := strings.Repeat(" ", width-(archivePreViewPadding*2)-adjustedPadding)
			titleStyle := lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, true, false)

			preview = selectedItem.Title + titleStyle.Render(line) + "\n" + truncateString(selectedItem.Content, maxLength, maxLines)
		}
	}

	return archivePreviewStyle.Width(width).Height(height).Render(preview)
}

func (m archive) helpView(width int) string {
	return archiveListHelpStyle.Width(width).Render(m.help.View(archiveKeys))
}
