package tui

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/o-kaisan/text-clipper/text"
	"github.com/o-kaisan/text-clipper/tui/constants"
)

type model struct {
	width   int
	height  int
	cursor  int
	choices []*text.Text
	help    help.Model
}

func InitialChoice() (model, tea.Cmd) {
	texts, err := getTextList(constants.Tr)
	m := model{
		choices: texts,
		help:    help.New(),
	}

	return m, func() tea.Msg { return errMsg(err) }
}

func getTextList(tr *text.GormRepository) ([]*text.Text, error) {
	texts, err := tr.List()
	if err != nil {
		return nil, fmt.Errorf("cannot get all texts: %w", err)
	}
	return texts, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

// --------------------------------------------------------------------------------
// update
// --------------------------------------------------------------------------------
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		constants.WindowSizeMsg = msg
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.Keymap.Quit):
			return m, tea.Quit
		case key.Matches(msg, constants.Keymap.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, constants.Keymap.Down):
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case key.Matches(msg, constants.Keymap.Edit):
			var cmds []interface{}
			cmds = append(cmds, constants.WindowSizeMsg)
			cmds = append(cmds, textinput.Blink) // TODO 遷移直後のBlinkが効かない

			var targetText *text.Text
			for i, choice := range m.choices {
				if m.cursor == i {
					targetText = choice
				}
			}
			register := InitialRegister(targetText)
			return register.Update(cmds)
		case key.Matches(msg, constants.Keymap.Add):
			var cmds []interface{}
			cmds = append(cmds, constants.WindowSizeMsg)
			cmds = append(cmds, textinput.Blink) // TODO 遷移直後のBlinkが効かない
			// 登録画面に遷移
			initialText := &text.Text{
				Title:    "",
				Contents: "",
			}
			register := InitialRegister(initialText)
			return register.Update(cmds)
		case key.Matches(msg, constants.Keymap.Delete):
			// 対象のitemを削除する
			for i, choice := range m.choices {
				if m.cursor == i {
					deleteText(constants.Tr, choice)
				}
				texts, err := getTextList(constants.Tr)
				if err != nil {
					return nil, func() tea.Msg { return errMsg(err) }
				}
				m.choices = texts
			}

		case key.Matches(msg, constants.Keymap.Select):
			for i, choice := range m.choices {
				if m.cursor == i {
					// 選択したテキストをクリップボードに登録する
					err := clipboard.WriteAll(choice.Contents)
					if err != nil {
						fmt.Println(fmt.Errorf("failed to clip the text to clipboard: text=%s, err=%w", choice.Contents, err))
					}
				}
			}
			return m, tea.Quit
		}
	}
	return m, nil
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
	choiceViewWidth := m.width / 2  // 幅を2で割って左右のパネルの幅を決定
	previewViewWidth := m.width / 2 // 幅を2で割って左右のパネルの幅を決定
	return lipgloss.JoinHorizontal(lipgloss.Top, m.choicesView(choiceViewWidth), m.previewView(previewViewWidth))
}

func (m model) choicesView(halfWidth int) string {
	// Iterate over our choices
	choicesView := "Select the item you want to call\n\n"
	for i, choice := range m.choices {
		choicesView += checkbox(choice.Title, m.cursor == i)
	}

	return lipgloss.NewStyle().Width(halfWidth).Render(choicesView + "\n" + m.helpView())
}

func (m model) previewView(halfWidth int) string {
	content := ""
	for i, choice := range m.choices {
		if m.cursor == i {
			content = choice.Contents
		}
	}

	previewStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true).
		BorderForeground(lipgloss.Color("#6272A4"))
	return previewStyle.Width(halfWidth).Render(content)
}

func (m model) helpView() string {
	help := m.help.ShortHelpView([]key.Binding{
		constants.Keymap.Up,
		constants.Keymap.Down,
		constants.Keymap.Select,
		constants.Keymap.Add,
		constants.Keymap.Delete,
		constants.Keymap.Quit,
	})
	return help
}

// --------------------------------------------------------------------------------
// utils
// --------------------------------------------------------------------------------
func checkbox(label string, checked bool) string {
	if checked {
		return colorFg("> "+label+"\n", "212")
	}
	return fmt.Sprintf("  %s\n", label)
}

// Color a string's foreground with the given value.
func colorFg(val, color string) string {
	return termenv.String(val).Foreground(term.Color(color)).String()
}

// General stuff for styling the view
var (
	term = termenv.EnvColorProfile()
)
