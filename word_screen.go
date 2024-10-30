package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/grig-iv/anki-card-creator/ld"
)

type wordScreen struct {
	page ld.WordPage
}

func newWordScreen() wordScreen {
	screen := wordScreen{}
	return screen
}

func (w wordScreen) Init() tea.Cmd {
	return nil
}

func (w wordScreen) Update(cmd tea.Msg) (tea.Model, tea.Cmd) {
	return w, nil
}

func (w wordScreen) View() string {
	return ""
}
