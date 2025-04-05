package list

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/o-kaisan/text-clipper/interface/bubbletea/constants"
)

const (
	adjustedListItemPosition = 2
	adjustedHeight           = -10
	minWidth                 = 80
	titleBarMinWidth         = 85
	choicesMinWidth          = 45
	previewMinWidth          = 35
	helpMinWidth             = 85
)

var (
	mainView            = lipgloss.NewStyle()
	listViewStyle       = lipgloss.NewStyle().PaddingTop(1)
	itemStyle           = lipgloss.NewStyle().PaddingLeft(4) // リストが揃うように
	noItemStyle         = lipgloss.NewStyle().PaddingLeft(adjustedListItemPosition)
	selectedItemStyle   = lipgloss.NewStyle().PaddingLeft(adjustedListItemPosition).Foreground(lipgloss.Color("170"))
	paginationStyle     = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	titleStyle          = lipgloss.NewStyle().Reverse(true).PaddingLeft(1).Italic(true).Width(18).Foreground(lipgloss.Color(constants.BgColor))
	titleAreaStyle      = lipgloss.NewStyle().PaddingLeft(1).PaddingBottom(1)
	previewStyle        = lipgloss.NewStyle().Margin(1, 1, 0, 4).Padding(1).Foreground(lipgloss.Color("#FAFAFA")).Background(lipgloss.Color(constants.BgColor)).Border(lipgloss.NormalBorder(), false, false, false, true)
	previewContentStyle = lipgloss.NewStyle().PaddingRight(4)
	helpStyle           = lipgloss.NewStyle().PaddingLeft(2).PaddingTop(1).PaddingBottom(1).Height(5)
)

func (m model) View() string {
	width := m.width
	helpWidth := width
	if width <= minWidth {
		width = titleBarMinWidth
		helpWidth = helpMinWidth
	}

	titleView := m.pageTitleView(width)
	listView := m.listView(width, m.height+adjustedHeight)
	preview := m.previewView(width-constants.AdjustedWidth, m.height+adjustedHeight)
	helpView := m.helpView(helpWidth)

	return mainView.Render(titleView + lipgloss.JoinHorizontal(lipgloss.Top, listView, preview) + "\n" + helpView)
}

func (m model) pageTitleView(width int) string {
	return titleAreaStyle.Width(width).Render(titleStyle.Render("# Active Items"))
}

func (m model) listView(width, height int) string {
	m.list.SetHeight(height)
	m.list.Styles.NoItems = noItemStyle.Width(constants.AdjustedWidth)
	m.list.Styles.PaginationStyle = paginationStyle
	m.list.Styles.StatusBarFilterCount = list.DefaultStyles().StatusBarFilterCount.Width(21)
	return listViewStyle.Render(m.list.View())
}

func (m model) previewView(width, height int) string {
	if m.list.SelectedItem() == nil {
		return previewStyle.Width(width).Height(height).Render("")
	}
	previewContentWidth := width - 10

	selectedItem := m.list.SelectedItem().(activeItem)
	titleBar := lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, true, false).Width(width - 2).Render(selectedItem.Title)
	previewContent := previewContentStyle.Width(previewContentWidth).Render(selectedItem.TruncateContent(height, previewContentWidth))

	return previewStyle.Width(width).Height(height).Render(titleBar + "\n" + previewContent)
}

func (m model) helpView(width int) string {
	return helpStyle.Width(width).Render(m.help.View(keys))
}
