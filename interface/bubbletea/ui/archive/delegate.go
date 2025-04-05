package archive

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	dmodel "github.com/o-kaisan/text-clipper/domain/model"
	"github.com/o-kaisan/text-clipper/interface/bubbletea/constants"
)

// archivedClipはlist.archivedClipインターフェースを実装するためのラッパー
type archivedClip struct {
	*dmodel.Clip
}

func (l archivedClip) FilterValue() string {
	return l.Clip.Title
}

type delegate struct{}

func (d delegate) Height() int                             { return 1 }
func (d delegate) Spacing() int                            { return 0 }
func (d delegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d delegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {

	itemNum := index + 1
	i, ok := listItem.(archivedClip)
	if !ok {
		return
	}

	// リストのタイトルが揃うように10個未満は空白を2つ入れる
	// ListViewの画面の長さはconstants.ListVewWidthに依存する
	var str string
	if itemNum >= 10 {
		str = fmt.Sprintf("%d. %-"+strconv.Itoa(constants.ListVewWidth)+"s", itemNum, i.Title)
	} else {
		str = fmt.Sprintf("%d.  %-"+strconv.Itoa(constants.ListVewWidth)+"s", itemNum, i.Title)
	}

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func convertArchivedClipsToListItems(clips []*dmodel.Clip) list.Model {
	listItems := make([]list.Item, len(clips))
	for i, item := range clips {
		listItems[i] = archivedClip{Clip: item}
	}

	// リストモデルの初期化 //
	// リストの高さと幅はdelegate.Render()で決まるため0で初期化する
	l := list.New(listItems, delegate{}, 0, 0)
	l.FilterInput.CharLimit = constants.TitleMaxLength
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)  // helpは独自で定義するためここで明示的に無効化する
	l.SetShowTitle(false) // Titleはここで明示的に無効化する

	return l
}
