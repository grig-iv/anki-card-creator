package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

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

	p := tea.NewProgram(newModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
