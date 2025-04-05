package model

import (
	"log"
	"strings"
	"time"
	"unicode"

	"gorm.io/gorm"
)

type Clip struct {
	gorm.Model
	Title      string
	Content    string
	IsActive   *bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
	LastUsedAt time.Time
}

func NewClip(title, content string, isActive *bool, createdAt time.Time, updatedAt time.Time, lastUsedAt time.Time) *Clip {
	return &Clip{
		Title:      title,
		Content:    content,
		IsActive:   isActive,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
		LastUsedAt: lastUsedAt,
	}
}

func (c *Clip) TruncateContent(height int, width int) string {
	// コンテンツを指定の幅で折り返して分割
	// log.Print(width)
	// log.Print(height)
	wrappedLines := wrapText(c.Content, width)

	if height <= 0 || len(wrappedLines) <= height {
		return c.Content
	}

	// 指定の高さにトランケート
	return strings.Join(wrappedLines[:height], "\n") + "\n…"
}

// 指定の幅でテキストを折り返す
func wrapText(text string, width int) []string {
	lines := strings.Split(text, "\n")
	var wrapped []string

	log.Print(len(lines))
	for _, line := range lines {
		wrapped = append(wrapped, SplitByDisplayWidth(line, width-8)...)
	}
	log.Print(len(wrapped))
	return wrapped
}

// 表示幅を考慮して文字列を分割する関数
func SplitByDisplayWidth(s string, maxWidth int) []string {
	var result []string
	var line []rune
	width := 0

	for _, r := range s {
		w := RuneWidth(r)
		if width+w > maxWidth {
			result = append(result, string(line))
			line = []rune{r}
			width = w
		} else {
			line = append(line, r)
			width += w
		}
	}
	if len(line) > 0 {
		result = append(result, string(line))
	}

	return result
}

// 日本語などの全角文字は幅2、英数字などは幅1とする
func RuneWidth(r rune) int {
	// CJK (中日韓) の文字か、全角記号・ひらがな・カタカナ・漢字など
	if unicode.In(r,
		unicode.Han,      // 漢字
		unicode.Hiragana, // ひらがな
		unicode.Katakana, // カタカナ
		unicode.Hangul,   // 韓国語
	) {
		return 2
	}
	// その他は半角とみなす
	return 1
}
