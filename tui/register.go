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

type Register struct {
	textId       uint
	title        textinput.Model
	err          error
	contents     textarea.Model
	focusedIndex int
	maxIndex     int
	help         help.Model
	validateErr  string
}

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	cursorStyle  = focusedStyle.Copy()

	blurredStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	helpStyle        = blurredStyle.Copy()
	validateErrStyle = focusedStyle.Copy()

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

func InitialRegister(text *text.Text) Register {
	log.Printf("id=%d, title=%s", text.ID, text.Title)

	ti := textinput.New()
	ti.Cursor.Style = cursorStyle // TODO あってもなくてもよい
	ti.Focus()
	ti.CharLimit = 32
	ti.Placeholder = "Title"
	ti.PromptStyle = focusedStyle
	ti.TextStyle = focusedStyle
	ti.SetValue(text.Title)

	ta := textarea.New()
	ta.SetValue(text.Contents)
	m := Register{
		textId:       text.ID, // 0の場合は新規保存
		title:        ti,
		contents:     ta,
		maxIndex:     2, // title, contents, submitの合計 -1
		focusedIndex: 0, // 初期表示はタイトルにフォーカス
		help:         help.New(),
	}
	return m
}

// 直接Updateから呼び出すので実際使われることはない
func (m Register) Init() tea.Cmd {
	return nil
}

type (
	errMsg error
)

func (m Register) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		constants.WindowSizeMsg = msg
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.Keymap.Back):
			choice, err := InitialChoice()
			if err != nil {
				log.Fatal(err)
			}
			return choice.Update(constants.WindowSizeMsg)
		case key.Matches(msg, constants.Keymap.Submit):
			// submitにフォーカスがある場合に登録処理を実行する
			if m.focusedIndex == m.maxIndex {
				var targetText *text.Text
				if m.textId != 0 {
					targetText = constants.Tr.FindByID(m.textId)
					targetText.Title = m.title.Value()
					targetText.Contents = m.contents.Value()
				} else {
					targetText = &text.Text{
						Title:    m.title.Value(),
						Contents: m.contents.Value(),
					}
				}
				if len(targetText.Title) > 0 && len(targetText.Contents) > 0 {
					err := saveOrUpdateText(constants.Tr, targetText)
					if err != nil {
						log.Fatal(err)
					}
					// 元の画面に戻る
					choice, _ := InitialChoice()
					return choice.Update(constants.WindowSizeMsg)
				} else {
					// 未入力項目があるためエラーメッセージを表示する
					m.validateErr = "Please fill in all required fields."
				}
			}

		case key.Matches(msg, constants.Keymap.Next):
			if m.focusedIndex < m.maxIndex {
				m.focusedIndex++
			}
			// 一旦フォーカスを解除
			m.title.Blur()
			m.contents.Blur()

			// テキストエリアにフォーカス
			if m.focusedIndex == 1 {
				cmd = m.contents.Focus()
				cmds = append(cmds, cmd)
			}

		case key.Matches(msg, constants.Keymap.Prev):
			if m.focusedIndex > 0 {
				m.focusedIndex--
			}
			// 一旦フォーカスを解除
			m.title.Blur()
			m.contents.Blur()

			// タイトルにフォーカス
			if m.focusedIndex == 0 {
				cmd = m.title.Focus()
				cmds = append(cmds, cmd)
			}
			// コンテンツにフォーカス
			if m.focusedIndex == 1 {
				cmd = m.contents.Focus()
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
		m.contents, cmd = m.contents.Update(msg)
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func saveOrUpdateText(tr *text.GormRepository, text *text.Text) error {
	var err error
	log.Printf("target text id=%d, title=%s", text.ID, text.Title)
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
	// // 元の画面に戻る
	// if m.backTo {
	// 	return InitialChoice().View()
	// }

	var b strings.Builder
	b.WriteString("Register your new texts...")
	b.WriteString("\n")
	b.WriteString(m.title.View())
	b.WriteString("\n\n")
	b.WriteString(m.contents.View())
	button := &blurredButton
	if m.focusedIndex == m.maxIndex {
		button = &focusedButton
	}
	if len(m.validateErr) > 0 {
		b.WriteString(validateErrStyle.Render(m.validateErr))
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)
	b.WriteString(helpStyle.Render(m.helpView()))

	return b.String()
}

func (m Register) helpView() string {
	help := m.help.ShortHelpView([]key.Binding{
		constants.Keymap.Submit,
		constants.Keymap.Back,
	})
	return help
}
