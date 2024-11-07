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
	page  ld.WordPage
	state wordScreenState

	currEntry   int
	currSense   int
	currExample int

	windowHeigth int
}

type wordScreenState uint8

const (
	selectingSense wordScreenState = iota
	selectingExample
)

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
		screen.currEntry = -1
	}

	if len(page.Entries[0].Senses) == 0 {
		screen.currSense = -1
	}

	screen.currExample = -1

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
			switch w.state {
			case selectingExample:
				w.nextExample()
			case selectingSense:
				w.nextSense()
			}
			return w, nil
		case tea.KeyUp:
			switch w.state {
			case selectingExample:
				w.prevExample()
			case selectingSense:
				w.prevSense()
			}
			return w, nil
		case tea.KeyEnter:
			if w.state == selectingSense {
				w.state = selectingExample
				w.currExample = -1
				w.nextExample()
			}
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

func (w *wordScreen) nextSense() {
	currEntry := w.page.Entries[w.currEntry]

	w.currSense += 1
	if w.currSense >= len(currEntry.Senses) {
		w.currEntry += 1
		if w.currEntry >= len(w.page.Entries) {
			w.currEntry = 0
		}

		currEntry := w.page.Entries[w.currEntry]
		if len(currEntry.Senses) != 0 {
			w.currSense = 0
		}
	}
}

func (w *wordScreen) prevSense() {
	w.currSense -= 1
	if w.currSense < 0 {
		w.currEntry -= 1
		if w.currEntry < 0 {
			w.currEntry = len(w.page.Entries) - 1
		}

		currEntry := w.page.Entries[w.currEntry]
		if len(currEntry.Senses) != 0 {
			w.currSense = len(currEntry.Senses) - 1
		}
	}
}

func (w *wordScreen) nextExample() {
	currLdSense, _ := w.getCurrSense()
	currSense, ok := currLdSense.(ld.Sense)
	if !ok {
		return
	}

	if len(currSense.Examples) == 0 {
		return
	}

	w.currExample += 1
	if w.currExample >= len(currSense.Examples) {
		w.currExample = 0
	}
}

func (w *wordScreen) prevExample() {
	currLdSense, _ := w.getCurrSense()
	currSense, ok := currLdSense.(ld.Sense)
	if !ok {
		return
	}

	if len(currSense.Examples) == 0 {
		return
	}

	w.currExample -= 1
	if w.currExample < 0 {
		w.currExample = len(currSense.Examples) - 1
	}
}

func (w *wordScreen) isSelecetedSense(sense ld.LdSense) bool {
	currSense, ok := w.getCurrSense()
	if !ok {
		return false
	}

	return reflect.DeepEqual(currSense, sense)
}

func (w *wordScreen) isSelecetedExample(sense ld.Sense, example ld.Example) bool {
	currSense, ok := w.getCurrSense()
	if !ok || !reflect.DeepEqual(currSense, sense) {
		return false
	}

	if w.currExample == -1 {
		return false
	}

	currExample := sense.Examples[w.currExample]
	return reflect.DeepEqual(currExample, example)
}

func (w *wordScreen) getCurrEntry() (ld.DictEntry, bool) {
	if w.currEntry < 0 || w.currEntry >= len(w.page.Entries) {
		return ld.DictEntry{}, false
	}

	return w.page.Entries[w.currEntry], true
}

func (w *wordScreen) getCurrSense() (ld.LdSense, bool) {
	currEntry, ok := w.getCurrEntry()
	if !ok {
		return nil, false
	}

	if w.currSense < 0 || w.currSense >= len(currEntry.Senses) {
		return nil, false
	}

	return currEntry.Senses[w.currSense], true
}

func (w wordScreen) View() string {
	builder := &strings.Builder{}

	switch w.state {
	case selectingSense:
		for _, e := range w.page.Entries {
			builder.WriteString("\n")
			w.renderEntry(builder, e)
		}
	case selectingExample:
		w.renderExampleSelection(builder)
	}

	return builder.String()
}

func (w wordScreen) renderEntry(builder *strings.Builder, entry ld.DictEntry) {
	w.renderEntryHeader(builder, entry)

	for i, s := range entry.Senses {
		builder.WriteString("  ")
		if len(entry.Senses) > 1 {
			builder.WriteString(fmt.Sprintf("%d. ", i+1))
		}
		switch s := s.(type) {
		case ld.Sense:
			w.renderSense(builder, s)
		case ld.CrossRefSense:
			w.renderCrossRefSense(builder, s)
		}
	}
}

func (w wordScreen) renderEntryHeader(builder *strings.Builder, entry ld.DictEntry) {
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
}

func (w wordScreen) renderSense(builder *strings.Builder, sense ld.Sense) {
	w.renderSenseHeader(builder, sense)
	w.renderSenseExamples(builder, sense)
}

func (w wordScreen) renderSenseHeader(builder *strings.Builder, sense ld.Sense) {
	if sense.Grammar != "" {
		builder.WriteString("[")
		builder.WriteString(grammarStyle.Render(sense.Grammar))
		builder.WriteString("] ")
	}

	defStyle := senseDefinitionStyle
	if w.isSelecetedSense(sense) {
		defStyle = senseDefinitionStyle.Foreground(lipgloss.Color("4"))
	}
	builder.WriteString(defStyle.Render(sense.Definition))

	builder.WriteString("\n")
}

func (w wordScreen) renderSenseExamples(builder *strings.Builder, sense ld.Sense) {
	lastColloquial := ""
	for _, e := range sense.Examples {
		builder.WriteString("    ")

		if e.Colloquial != "" && e.Colloquial != lastColloquial {
			builder.WriteString(colloquialStyle.Render(e.Colloquial))
			lastColloquial = e.Colloquial
			builder.WriteString("\n    ")
		}

		style := exampleTextStyle
		if w.isSelecetedExample(sense, e) {
			style = exampleTextStyle.Foreground(lipgloss.Color("4"))
		}

		builder.WriteString("- ")
		builder.WriteString(style.Render("\""))
		builder.WriteString(style.Render(e.Text))
		builder.WriteString(style.Render("\""))
		builder.WriteString("\n")
	}
}

func (w wordScreen) renderCrossRefSense(builder *strings.Builder, sense ld.CrossRefSense) {
	style := crossRefStyle
	if w.isSelecetedSense(sense) {
		style = crossRefStyle.Foreground(lipgloss.Color("4"))
	}
	builder.WriteString(style.Render("->", sense.Text))
	builder.WriteString("\n")
}

func (w wordScreen) renderExampleSelection(builder *strings.Builder) {
	entry, ok := w.getCurrEntry()
	if !ok {
		return
	}

	w.renderEntryHeader(builder, entry)

	ldSense, ok := w.getCurrSense()
	if !ok {
		return
	}

	sense, ok := ldSense.(ld.Sense)
	if !ok {
		return
	}

	w.renderSenseHeader(builder, sense)
	w.renderSenseExamples(builder, sense)
}
