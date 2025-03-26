package register

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	descriptionStyle  = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("#FAFAFA")).Background(lipgloss.Color("#696969"))
	registerViewStyle = lipgloss.NewStyle().PaddingLeft(1).PaddingBottom(1)
	blurredStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	registerHelpStyle = blurredStyle.Copy()
	validateErrStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	focusedButton     = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render("[ Submit ]")
	blurredButton     = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

func (m model) View() string {
	descriptionMessage := "# Register a new item. "
	if m.clipId != 0 {
		descriptionMessage = "# Edit the selected item. "
	}

	registerViewWidth := m.width - 5
	registerViewHeight := 18
	helpWidth := m.width - 5

	// 最低限のwidthを定義
	if m.width <= 65 {
		helpWidth = 65
		registerViewWidth = 65
	}

	m.content.SetWidth(registerViewWidth - 3)

	var b strings.Builder
	b.WriteString(descriptionStyle.MarginLeft(-1).Render(descriptionMessage))
	b.WriteString("\n\n")
	b.WriteString(m.title.View())
	b.WriteString("\n\n")
	b.WriteString(m.content.View())
	button := &blurredButton
	if m.focusedIndex == m.maxIndex {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	if len(m.validateErr) > 0 {
		b.WriteString(validateErrStyle.Render(m.validateErr))
		b.WriteString("\n\n")
	}
	b.WriteString(m.helpView(helpWidth))

	return registerViewStyle.Width(registerViewWidth).Height(registerViewHeight).Render(b.String())
}

func (m model) helpView(width int) string {
	return registerHelpStyle.Width(width).Render(m.help.View(keys))
}
