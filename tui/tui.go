package tui

import (
	"os"

	tea "charm.land/bubbletea/v2"
	"yamp/tui/components"
)

type page int

const (
	homePage          page = iota
	playlistPage      page = iota
	anotherOptionPage page = iota
)

type model struct {
	currentPage page
	choices     []string
	cursor      int
}

func initialModel() model {
	return model{
		choices: []string{"playlists", "another option"},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			if m.currentPage == homePage {
				switch m.choices[m.cursor] {
				case "playlists":
					m.currentPage = playlistPage
				case "another option":
					m.currentPage = anotherOptionPage
				}
				m.cursor = 0
			}
		case "esc":
			m.currentPage = homePage
			m.cursor = 0
		}
	}

	return m, nil
}

func (m model) View() tea.View {
	var s string
	switch m.currentPage {
	case homePage:
		s = components.HomeView(m.choices, m.cursor)
	case playlistPage:
		s = components.PlaylistView()
	case anotherOptionPage:
		s = components.AnotherOptionView()
	}
	v := tea.NewView(s)
	v.AltScreen = true
	return v
}

func RunTUI() {
	p := tea.NewProgram(initialModel())
	_, err := p.Run()
	if err != nil {
		os.Exit(1)
	}
}
