package main

import (
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/grig-iv/anki-card-creator/ankiConnect"
)

type startupScreen struct {
	state startupScreenState
}

type startupScreenState int

const (
	checkingAnki startupScreenState = iota
	askingOpenAnki
)

type askOpenAnkiMsg struct{}
type openSearchMsg struct{}

func (s startupScreen) Init() tea.Cmd {
	return func() tea.Msg {
		isRunning, err := ankiConnect.IsRunning()

		if err != nil {
			return askOpenAnkiMsg{}
		}

		if isRunning {
			return openSearchMsg{}
		} else {
			return askOpenAnkiMsg{}
		}
	}
}

func (s startupScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case askOpenAnkiMsg:
		s.state = askingOpenAnki
	case tea.KeyMsg:
		if s.state == askingOpenAnki {
			switch msg.String() {
			case "y":
				exec.Command("anki").Start()
				return s, func() tea.Msg { return openSearchMsg{} }
			case "n":
				return s, func() tea.Msg { return openSearchMsg{} }
			}
		}
	}

	return s, nil
}

func (s startupScreen) View() string {
	switch s.state {
	case checkingAnki:
		return "Checking if anki connect is running..."
	case askingOpenAnki:
		return "Anki connect is not runnig, do you want to start anki?\n[y/n]"
	}

	return ""
}
