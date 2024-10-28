package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	screen tea.Model
}

const (
	logPath = "log"
)

func main() {
	f, err := tea.LogToFile(logPath, "")
	if err != nil {
		log.Fatal(err)

	}
	defer f.Close()

	os.Truncate(logPath, 0)

	p := tea.NewProgram(initModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func initModel() model {
	return model{
		screen: newSearchScreen(),
	}
}

func (m model) Init() tea.Cmd {
	return m.screen.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	case error:
		log.Println(msg)
		return m, nil
	}

	m.screen, cmd = m.screen.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.screen.View()
}
