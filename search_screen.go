package main

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/grig-iv/anki-card-creator/ld"
)

type searchScreen struct {
	state          searcScreenState
	searchBox      textinput.Model
	suggestionList list.Model
}

type pageMsg ld.Page

type suggestionItem string

type suggestionItemDelegate struct{}

type searcScreenState uint8

const (
	typing searcScreenState = iota
	loading
	suggestion
)

var (
	suggestionTitleStyle = lipgloss.NewStyle().PaddingLeft(0)
	selectedItemStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("170"))
)

func (i suggestionItem) FilterValue() string { return "" }

func (d suggestionItemDelegate) Height() int                             { return 1 }
func (d suggestionItemDelegate) Spacing() int                            { return 0 }
func (d suggestionItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d suggestionItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(suggestionItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := func(s string) string { return s }
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render(s)
		}
	}

	fmt.Fprint(w, fn(str))
}

func newSearchScreen() searchScreen {
	screen := searchScreen{}

	screen.state = typing

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
	switch s.state {
	case typing:
		return s.updateTyping(msg)
	case suggestion:
		return s.updateSuggestions(msg)
	case loading:
		return s.updateLoading(msg)
	default:
		return s, nil
	}
}

func (s searchScreen) updateTyping(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			text := s.searchBox.Value()
			text = strings.TrimSpace(text)
			if text == "" {
				return s, nil
			}
			s.state = loading
			return s, search(text)
		case tea.KeyEsc:
			return s, tea.Quit
		}
	}

	var cmd tea.Cmd
	s.searchBox, cmd = s.searchBox.Update(msg)
	return s, cmd
}

func (s searchScreen) updateSuggestions(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			index := s.suggestionList.Index()
			items := s.suggestionList.Items()
			selected, _ := items[index].(suggestionItem)
			s.state = loading
			return s, search(string(selected))
		case tea.KeyEsc:
			s.state = typing
			return s, nil
		}
	}

	var cmd tea.Cmd
	s.suggestionList, cmd = s.suggestionList.Update(msg)
	return s, cmd
}

func (s searchScreen) updateLoading(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case pageMsg:
		switch page := msg.(type) {
		case ld.SearchPage:
			items := make([]list.Item, 0, len(page.Results))
			for _, r := range page.Results {
				items = append(items, suggestionItem(r))
			}

			s.suggestionList = list.New(items, suggestionItemDelegate{}, 20, 14)
			s.suggestionList.Title = "Did you mean:"
			s.suggestionList.SetShowStatusBar(false)
			s.suggestionList.SetFilteringEnabled(false)
			s.suggestionList.SetShowHelp(false)
			s.suggestionList.Styles.Title = suggestionTitleStyle

			s.state = suggestion

			var cmd tea.Cmd
			s.suggestionList, cmd = s.suggestionList.Update(msg)
			return s, cmd
		}
	}

	return s, nil
}

func (s searchScreen) View() string {
	switch s.state {
	case typing:
		return s.searchBox.View()
	case loading:
		return "loading..."
	case suggestion:
		return s.suggestionList.View()
	default:
		return "something wrong"
	}
}

func search(text string) tea.Cmd {
	return func() tea.Msg {
		log.Println("Searching:", text)
		page, err := ld.Search(text)
		if err != nil {
			return err
		}

		return pageMsg(page)
	}
}
