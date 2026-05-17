package tui

import tea "charm.land/bubbletea/v2"

type playlistModel struct{}

func (m playlistModel) Init() tea.Cmd {
	return nil
}

func (m playlistModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m playlistModel) View() tea.View {
	return tea.View{
		Content:   "playlist!\n",
		AltScreen: true,
	}
}
