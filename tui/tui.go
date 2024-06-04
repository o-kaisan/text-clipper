package tui

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/o-kaisan/text-clipper/common"
	"github.com/o-kaisan/text-clipper/text"
	"github.com/o-kaisan/text-clipper/tui/constants"
)

// StartTea the entry point for the UI. Initializes the model.
func StartTea(tr text.GormRepository) error {
	// デフォルトのデータベースパス
	defaultLogFilePath, err := common.GetPathFromDefaultPath("debug.log")
	if err != nil {
		log.Fatal("failed in getting log file path: %w", err)
		os.Exit(1)
	}
	// 環境変数で指定されたパスがあればそれを使用
	logFilePath := os.Getenv("TEXT_CLIPPER__LOG_FILE_PATH")
	if defaultLogFilePath == "" {
		defaultLogFilePath = logFilePath
	}

	f, err := tea.LogToFile(defaultLogFilePath, "debug")
	if err != nil {
		log.Fatal("failed in setting log file: %w", err)
		os.Exit(1)
	}
	defer f.Close()

	constants.Tr = &tr

	// TODO: can we acknowledge this error
	// エラーがtea.Cmdなのでアプリがスタートする前にキャッチできない
	m, _ := InitialList()
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		log.Fatal("Error while running program:", err)
		os.Exit(1)
	}
	return nil
}
