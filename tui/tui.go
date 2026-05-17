package tui

import (
	tea "charm.land/bubbletea/v2"
)

type view int

const (
	playlists view = iota
	song      view = iota
)

type Model struct {
	view view
}

func (m Model) Init() tea.Cmd {
	return nil
}

func initialModel() Model {
	return Model{
		view: playlists,
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			m.view = (m.view + 1) % 2
		case "shift+tab":
			m.view = (m.view - 1 + 2) % 2
		}
	}
	return m, nil
}

func (m Model) View() tea.View {
	switch m.view {
	case playlists:
		pm := playlistModel{}
		return pm.View()
	case song:
		sm := songModel{}
		return sm.View()
	}
	return playlistModel{}.View()
}

func RunTUI() error {
	p := tea.NewProgram(initialModel())
	_, err := p.Run()
	return err
}
