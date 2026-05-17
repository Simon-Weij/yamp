package tui

import (
	tea "charm.land/bubbletea/v2"
)

type Model struct {
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() tea.View {
	return tea.View{
		Content:   "Hello, World!\n",
		AltScreen: true,
	}
}

func RunTUI() error {
	p := tea.NewProgram(Model{})
	_, err := p.Run()
	return err
}
