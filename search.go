package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type searchScreen struct {
	searchBox textinput.Model
}

func newSearchScreen() searchScreen {
	screen := searchScreen{}
	screen.searchBox = textinput.New()
	screen.searchBox.Placeholder = "search"
	screen.searchBox.Focus()
	screen.searchBox.CharLimit = 50
	screen.searchBox.Width = 20

	return screen
}

func (s searchScreen) Init() tea.Cmd {
	return textinput.Blink
}

func (s searchScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	s.searchBox, cmd = s.searchBox.Update(msg)

	return s, cmd
}

func (s searchScreen) View() string {
	return s.searchBox.View()
}
