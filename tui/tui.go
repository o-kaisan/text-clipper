package tui

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/o-kaisan/text-clipper/text"
	"github.com/o-kaisan/text-clipper/tui/constants"
)

// StartTea the entry point for the UI. Initializes the model.
func StartTea(tr text.GormRepository) error {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal("err: %w", err)
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
