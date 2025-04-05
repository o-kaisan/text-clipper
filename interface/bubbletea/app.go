package bubbletea

import (
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/o-kaisan/text-clipper/common"
	"github.com/o-kaisan/text-clipper/di"
	"github.com/o-kaisan/text-clipper/interface/bubbletea/ui/list"
)

// StartTea the entry point for the UI. Initializes the model.
func StartTea(c di.Container) error {

	// ログ出力設定
	debug_env := os.Getenv("TEXT_CLIPPER_DEBUG")
	debug := strings.ToLower(debug_env) == "true"

	if debug {
		f := getBubbleTeaLogger()
		defer f.Close()
	}

	m, _ := list.NewListModel(c.Cs)
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		log.Fatal("Error while running program:", err)
		os.Exit(1)
	}
	return nil
}

func getBubbleTeaLogger() *os.File {
	logFilePath, err := common.GetPathFromPath("debug.log")
	if err != nil {
		log.Fatal("failed in getting log file path: %w", err)
		os.Exit(1)
	}
	f, err := tea.LogToFile(logFilePath, "debug")
	if err != nil {
		log.Fatal("failed in setting log file: %w", err)
		os.Exit(1)
	}
	return f
}
