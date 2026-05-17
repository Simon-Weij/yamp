package tui

import tea "charm.land/bubbletea/v2"

type songModel struct{}

func (m songModel) Init() tea.Cmd {
	return nil
}

func (m songModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m songModel) View() tea.View {
	return tea.View{
		Content:   "song!\n",
		AltScreen: true,
	}
}
