package archive

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/o-kaisan/text-clipper/domain/service"
	"github.com/o-kaisan/text-clipper/interface/bubbletea/constants"
)

type model struct {
	cs        service.ClipService
	listModel tea.Model
	width     int
	help      help.Model
	height    int
	list      list.Model
	err       error
}

func NewArchive(cs service.ClipService, listModel tea.Model) (model, tea.Cmd) {
	var m model
	clips, err := cs.GetArchivedClips()
	if err != nil {
		m.err = fmt.Errorf("failed to initial archive model: %w", err)
		return m, nil
	}

	m.cs = cs
	// *item.Item のスライスを ListItem のスライスに変換
	m.list = convertArchivedClipsToListItems(clips)
	m.listModel = listModel
	m.help = help.New()
	m.help.ShowAll = true // フルヘルプをはじめから表示する
	return m, func() tea.Msg { return constants.ErrMsg(err) }
}

func (m model) Init() tea.Cmd {
	return nil
}
