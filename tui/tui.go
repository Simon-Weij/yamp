package tui

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"yamp/tui/components"
)

type page int

const (
	homePage     page = iota
	playlistPage page = iota
)

type model struct {
	currentPage page
	choices     []string
	cursor      int
	chosen      string
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
			m.chosen = m.choices[m.cursor]
			return m, tea.Quit
		case "escape":
			m.currentPage = homePage
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
		s = "playlist page"
	}
	v := tea.NewView(s)
	v.AltScreen = true
	return v
}

func RunTUI() {
	p := tea.NewProgram(initialModel())
	m, err := p.Run()
	if err != nil {
		fmt.Printf("There's been an error: %v", err)
		os.Exit(1)
	}
	if m, ok := m.(model); ok && m.chosen != "" {
		fmt.Println(m.chosen)
	}
}
