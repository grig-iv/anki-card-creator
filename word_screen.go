package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/grig-iv/anki-card-creator/ld"
)

type wordScreen struct {
	page ld.WordPage
}

func newWordScreen(page ld.WordPage) wordScreen {
	return wordScreen{page}
}

func (w wordScreen) Init() tea.Cmd {
	return nil
}

func (w wordScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case pageMsg:
		switch page := msg.(type) {
		case ld.WordPage:
			w.page = page
			return w, nil
		}
	}

	return w, nil
}

func (w wordScreen) View() string {
	return w.page.Title
}
