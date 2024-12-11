package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/grig-iv/anki-card-creator/ld"
)

type model struct {
	screen tea.Model
}

func newModel() model {
	return model{
		screen: startupScreen{},
	}
}

func (m model) Init() tea.Cmd {
	return m.screen.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case openSearchMsg:
		m.screen = newSearchScreen()
		return m, m.screen.Init()
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	case error:
		panic(msg)
	case pageMsg:
		switch page := msg.page.(type) {
		case ld.WordPage:
			if _, ok := m.screen.(wordScreen); !ok {
				m.screen = newWordScreen(page)
				return m, m.screen.Init()
			}
		case ld.SearchPage:
			if _, ok := m.screen.(searchScreen); !ok {
				m.screen = newSearchScreen()
				return m, m.screen.Init()
			}
		}
	}

	m.screen, cmd = m.screen.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.screen.View()
}
