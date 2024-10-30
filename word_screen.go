package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/grig-iv/anki-card-creator/ld"
)

type wordScreen struct {
	page          ld.WordPage
	selectedEntry uint
	selectedSense uint
}

var (
	titleStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("4")).Bold(true)
	hyphenationStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true)
	pronunciationStyle = lipgloss.NewStyle().Italic(true)
	partOfSpeachStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	grammarStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	crossRefStyle      = lipgloss.NewStyle().Bold(true)
	colloquialStyle    = lipgloss.NewStyle().Bold(true)
)

func newWordScreen(page ld.WordPage) wordScreen {
	return wordScreen{page, 1, 1}
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
	builder := &strings.Builder{}

	builder.WriteString("\n")

	for _, e := range w.page.Entries {
		renderEtrie(builder, e)
	}

	return builder.String()
}

func renderEtrie(builder *strings.Builder, entry ld.DictEntry) {
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

	for _, s := range entry.Senses {
		builder.WriteString("  ")
		switch s := s.(type) {
		case ld.Sense:
			renderSense(builder, s)
		case ld.CrossRefSense:
			renderCrossRefSense(builder, s)
		}
	}
}

func renderSense(builder *strings.Builder, sense ld.Sense) {
	if sense.Grammar != "" {
		builder.WriteString("[")
		builder.WriteString(grammarStyle.Render(sense.Grammar))
		builder.WriteString("] ")
	}

	builder.WriteString(sense.Definition)
	builder.WriteString("\n")

	lastColloquial := ""
	for _, e := range sense.Examples {
		builder.WriteString("    ")

		if e.Colloquial != "" && e.Colloquial != lastColloquial {
			builder.WriteString(colloquialStyle.Render(e.Colloquial))
			lastColloquial = e.Colloquial
			builder.WriteString("\n")
		}

		builder.WriteString("- \"")
		builder.WriteString(e.Text)
		builder.WriteString("\"\n")
	}
}

func renderCrossRefSense(builder *strings.Builder, sense ld.CrossRefSense) {
	builder.WriteString(crossRefStyle.Render("-> ", sense.Text))
	builder.WriteString("\n")
}
