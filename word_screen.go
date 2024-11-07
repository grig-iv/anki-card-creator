package main

import (
	"fmt"
	"reflect"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/grig-iv/anki-card-creator/ld"
)

type wordScreen struct {
	page            ld.WordPage
	selectedEntry   int
	selectedSense   int
	selectedExample int
	windowHeigth    int
}

var (
	titleStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("4")).Bold(true)
	hyphenationStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true)
	pronunciationStyle   = lipgloss.NewStyle().Italic(true)
	partOfSpeachStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	grammarStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	crossRefStyle        = lipgloss.NewStyle().Bold(true)
	colloquialStyle      = lipgloss.NewStyle().Bold(true)
	senseDefinitionStyle = lipgloss.NewStyle().Bold(true)
	exampleTextStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Italic(true)
)

func newWordScreen(page ld.WordPage) wordScreen {
	screen := wordScreen{}
	screen.page = page

	if len(page.Entries) == 0 {
		screen.selectedEntry = -1
	}

	if len(page.Entries[0].Senses) == 0 {
		screen.selectedSense = -1
	}

	return screen
}

func (w wordScreen) Init() tea.Cmd {
	return tea.WindowSize()
}

func (w wordScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		w.windowHeigth = msg.Height
		return w, nil
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyDown:
			w.selectNext()
			return w, nil
		case tea.KeyUp:
			w.selectPrev()
			return w, nil
		}
	case pageMsg:
		switch page := msg.(type) {
		case ld.WordPage:
			w.page = page
			return w, nil
		}
	}

	return w, nil
}

func (w *wordScreen) selectNext() {
	currEntry := w.page.Entries[w.selectedEntry]

	w.selectedSense += 1
	if w.selectedSense >= len(currEntry.Senses) {
		w.selectedEntry += 1
		if w.selectedEntry >= len(w.page.Entries) {
			w.selectedEntry = 0
		}

		currEntry := w.page.Entries[w.selectedEntry]
		if len(currEntry.Senses) != 0 {
			w.selectedSense = 0
		}
	}
}

func (w *wordScreen) selectPrev() {
	w.selectedSense -= 1
	if w.selectedSense < 0 {
		w.selectedEntry -= 1
		if w.selectedEntry < 0 {
			w.selectedEntry = len(w.page.Entries) - 1
		}

		currEntry := w.page.Entries[w.selectedEntry]
		if len(currEntry.Senses) != 0 {
			w.selectedSense = len(currEntry.Senses) - 1
		}
	}
}

func (w *wordScreen) isSeleceted(sense ld.LdSense) bool {
	if w.selectedEntry < 0 || w.selectedEntry >= len(w.page.Entries) {
		return false
	}

	selectedEntry := w.page.Entries[w.selectedEntry]
	if w.selectedSense < 0 || w.selectedSense >= len(selectedEntry.Senses) {
		return false
	}

	selectedSense := selectedEntry.Senses[w.selectedSense]
	return reflect.DeepEqual(selectedSense, sense)
}

func (w wordScreen) View() string {
	builder := &strings.Builder{}

	for _, e := range w.page.Entries {
		builder.WriteString("\n")
		w.renderEtrie(builder, e)
	}

	return builder.String()
}

func (w wordScreen) renderEtrie(builder *strings.Builder, entry ld.DictEntry) {
	if entry.Hyphenation != "" {
		builder.WriteString(hyphenationStyle.Render(entry.Hyphenation))
		builder.WriteString(" ")
	}

	if entry.Pronunciation != "" {
		builder.WriteString(pronunciationStyle.Render(entry.Pronunciation))
		builder.WriteString(" ")
	}

	if entry.PartOfSpeach != "" {
		builder.WriteString(partOfSpeachStyle.Render(entry.PartOfSpeach))
	}

	builder.WriteString("\n")

	for i, s := range entry.Senses {
		isSelected := w.isSeleceted(s)
		builder.WriteString("  ")
		if len(entry.Senses) > 1 {
			builder.WriteString(fmt.Sprintf("%d. ", i+1))
		}
		switch s := s.(type) {
		case ld.Sense:
			renderSense(builder, s, isSelected)
		case ld.CrossRefSense:
			renderCrossRefSense(builder, s, isSelected)
		}
	}
}

func renderSense(builder *strings.Builder, sense ld.Sense, isSeleceted bool) {
	if sense.Grammar != "" {
		builder.WriteString("[")
		builder.WriteString(grammarStyle.Render(sense.Grammar))
		builder.WriteString("] ")
	}

	defStyle := senseDefinitionStyle
	if isSeleceted {
		defStyle = senseDefinitionStyle.Foreground(lipgloss.Color("4"))
	}
	builder.WriteString(defStyle.Render(sense.Definition))

	builder.WriteString("\n")

	lastColloquial := ""
	for _, e := range sense.Examples {
		builder.WriteString("    ")

		if e.Colloquial != "" && e.Colloquial != lastColloquial {
			builder.WriteString(colloquialStyle.Render(e.Colloquial))
			lastColloquial = e.Colloquial
			builder.WriteString("\n    ")
		}

		builder.WriteString("- \"")
		builder.WriteString(exampleTextStyle.Render(e.Text))
		builder.WriteString("\"\n")
	}
}

func renderCrossRefSense(builder *strings.Builder, sense ld.CrossRefSense, isSelecetd bool) {
	style := crossRefStyle
	if isSelecetd {
		style = crossRefStyle.Foreground(lipgloss.Color("4"))
	}
	builder.WriteString(style.Render("->", sense.Text))
	builder.WriteString("\n")
}
