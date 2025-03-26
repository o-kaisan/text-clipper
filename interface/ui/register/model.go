package register

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/o-kaisan/text-clipper/domain/service"
	"github.com/o-kaisan/text-clipper/interface/constants"
)

type model struct {
	cs           service.ClipService
	listModel    tea.Model
	width        int
	clipId       uint
	title        textinput.Model
	err          error
	content      textarea.Model
	focusedIndex int
	maxIndex     int
	help         help.Model
	validateErr  string
}

func NewRegister(cs service.ClipService, listModel tea.Model, cid uint, title string, content string) model {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = constants.TitleMaxLength
	ti.Placeholder = "Title"
	ti.SetValue(title)

	ta := textarea.New()
	ta.ShowLineNumbers = true
	ta.CharLimit = 5000
	ta.SetWidth(constants.WindowSizeMsg.Width)
	ta.Placeholder = "Enter your Content here..."
	ta.SetValue(content)
	m := model{
		cs:           cs,
		listModel:    listModel,
		clipId:       cid, // 0の場合は新規保存
		title:        ti,
		content:      ta,
		maxIndex:     2, // 0からなのでtitle, content, submitの合計 -1
		focusedIndex: 0, // 初期表示はタイトルにフォーカス
		help:         help.New(),
		width:        constants.WindowSizeMsg.Width,
		validateErr:  "",
	}
	m.help.ShowAll = true // フルヘルプをはじめから表示する
	return m
}

// 直接Updateから呼び出すので実際使われることはない
func (m model) Init() tea.Cmd {
	return nil
}
