package clipboard

import (
	"fmt"

	"github.com/atotto/clipboard"
)

// クリップボードにテキストをコピーする関数
func CopyToClipboard(text string) error {
	err := clipboard.WriteAll(text)
	if err != nil {
		return fmt.Errorf("failed to write to clipboard: %v", err)
	}
	return nil
}
