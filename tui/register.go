package tui

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/o-kaisan/text-clipper/text"
	"github.com/o-kaisan/text-clipper/tui/constants"
)

type (
	errMsg error
)

// ---------------------------------------------------------------
// Style
// ---------------------------------------------------------------
var (
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	// TODO widthの調整が必要
	// TODO Help独自実装してもよいかも
	registerHelpStyle = blurredStyle.Copy()
	validateErrStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	focusedButton     = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render("[ Submit ]")
	blurredButton     = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

// ---------------------------------------------------------------
// Keys
// ---------------------------------------------------------------
type registerKeyMap struct {
	Submit key.Binding
	Back   key.Binding
	Next   key.Binding
	Prev   key.Binding
}

var registerKeys = registerKeyMap{
	Submit: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "over the submit button to register the item"),
	),
	Back: key.NewBinding(
		key.WithKeys("ctrl+c", "esc"),
		key.WithHelp("ctrl+c/esc", "back to list view"),
	),
	Next: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next input"),
	),
	Prev: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "previous input"),
	),
}

// ---------------------------------------------------------------
// Model
// ---------------------------------------------------------------
type Register struct {
	width        int
	textId       uint
	title        textinput.Model
	err          error
	content      textarea.Model
	focusedIndex int
	maxIndex     int
	help         help.Model
	validateErr  string
}

func InitialRegister(text *text.Text) Register {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 30
	ti.Placeholder = "Title"
	ti.SetValue(text.Title)

	ta := textarea.New()
	ta.ShowLineNumbers = false
	ta.SetWidth(constants.WindowSizeMsg.Width)
	ta.Placeholder = "Enter your Content here..."
	ta.SetValue(text.Content)
	m := Register{
		textId:       text.ID, // 0の場合は新規保存
		title:        ti,
		content:      ta,
		maxIndex:     2, // title, content, submitの合計 -1
		focusedIndex: 0, // 初期表示はタイトルにフォーカス
		help:         help.New(),
	}
	return m
}

// 直接Updateから呼び出すので実際使われることはない
func (m Register) Init() tea.Cmd {
	return nil
}

// ---------------------------------------------------------------
// Update
// ---------------------------------------------------------------
func (m Register) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		constants.WindowSizeMsg = msg
		m.width = msg.Width
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, registerKeys.Back):
			// TODO: can we acknowledge this error
			// エラーがtea.Cmdなのでエラーがキャッチできない
			choice, _ := InitialList()
			return choice.Update(constants.WindowSizeMsg)
		case key.Matches(msg, registerKeys.Submit):
			// submitにフォーカスがある場合に登録処理を実行する
			if m.focusedIndex == m.maxIndex {
				var targetText *text.Text
				if m.textId != 0 {
					targetText = constants.Tr.FindByID(m.textId)
					targetText.Title = m.title.Value()
					targetText.Content = m.content.Value()
				} else {
					targetText = &text.Text{
						Title:   m.title.Value(),
						Content: m.content.Value(),
					}
				}
				if len(targetText.Title) > 0 && len(targetText.Content) > 0 {
					err := saveOrUpdateText(constants.Tr, targetText)
					if err != nil {
						log.Fatal(err)
					}
					// 元の画面に戻る
					choice, _ := InitialList()
					return choice.Update(constants.WindowSizeMsg)
				} else {
					// 未入力項目があるためエラーメッセージを表示する
					m.validateErr = "Please fill in all fields."
				}
			}

		case key.Matches(msg, registerKeys.Next):
			if m.focusedIndex < m.maxIndex {
				m.focusedIndex++
			}
			// 一旦フォーカスを解除
			m.title.Blur()
			m.content.Blur()

			// テキストエリアにフォーカス
			if m.focusedIndex == 1 {
				cmd = m.content.Focus()
				cmds = append(cmds, cmd)
			}

		case key.Matches(msg, registerKeys.Prev):
			if m.focusedIndex > 0 {
				m.focusedIndex--
			}
			// 一旦フォーカスを解除
			m.title.Blur()
			m.content.Blur()

			// タイトルにフォーカス
			if m.focusedIndex == 0 {
				cmd = m.title.Focus()
				cmds = append(cmds, cmd)
			}
			// コンテンツにフォーカス
			if m.focusedIndex == 1 {
				cmd = m.content.Focus()
				cmds = append(cmds, cmd)
			}

		}
	case errMsg:
		m.err = msg
		return m, nil
	}

	switch m.focusedIndex {
	case 0:
		m.title, cmd = m.title.Update(msg)
	case 1:
		m.content, cmd = m.content.Update(msg)
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func saveOrUpdateText(tr *text.GormRepository, text *text.Text) error {
	var err error
	if text.ID == 0 {
		err = tr.Crete(text)
	} else {
		err = tr.Update(text)
	}
	if err != nil {
		return fmt.Errorf("can not save new text: %w", err)
	}
	return nil
}

// --------------------------------------------------------------------------------
// View
// --------------------------------------------------------------------------------
func (m Register) View() string {
	if m.width > 0 {
		// m.title.Width = m.width
		m.content.SetWidth(m.width)
	}

	var b strings.Builder
	b.WriteString("Register your new item...")
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
	b.WriteString(m.helpView())

	return b.String()
}

func (m Register) helpView() string {
	help := m.help.ShortHelpView([]key.Binding{
		registerKeys.Submit,
		registerKeys.Back,
		registerKeys.Next,
		registerKeys.Prev,
	})
	return registerHelpStyle.Width(m.width).Render(help)
}
