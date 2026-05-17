package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type view int

const (
	playlists view = iota
	song      view = iota
)

type Model struct {
	view   view
	width  int
	height int
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
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
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
	tabBar := m.renderTabBar()
	content := m.renderContent(tabBar)

	full := lipgloss.JoinVertical(lipgloss.Left, tabBar, content)
	return tea.View{
		Content:   full,
		AltScreen: m.currentView().AltScreen,
	}
}

func (m Model) renderTabBar() string {
	labels := []string{"Playlists", "Song"}
	styled := make([]string, len(labels))

	for i, t := range labels {
		if view(i) == m.view {
			styled[i] = tabActive().Render(t)
		} else {
			styled[i] = tabInactive().Render(t)
		}
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, styled...)
	gap := tabGap().Render(strings.Repeat(" ", max(0, m.width-lipgloss.Width(row))))
	return lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)
}

func (m Model) renderContent(tabBar string) string {
	style := lipgloss.NewStyle().
		Width(max(0, m.width)).
		Height(max(0, m.height-lipgloss.Height(tabBar)))

	return style.Render(m.currentView().Content)
}

func (m Model) currentView() tea.View {
	switch m.view {
	case playlists:
		return playlistModel{}.View()
	case song:
		return songModel{}.View()
	}
	return tea.View{}
}

func tabActive() lipgloss.Style {
	theme := currentTheme()
	return lipgloss.NewStyle().
		Border(lipgloss.Border{
			Top:         "─",
			Bottom:      " ",
			Left:        "│",
			Right:       "│",
			TopLeft:     "╭",
			TopRight:    "╮",
			BottomLeft:  "┘",
			BottomRight: "└",
		}, true).
		BorderForeground(lipgloss.Color(theme.borderColour)).
		Padding(0, 1).
		Bold(true)
}

func tabInactive() lipgloss.Style {
	theme := currentTheme()
	return lipgloss.NewStyle().
		Border(lipgloss.Border{
			Top:         "─",
			Bottom:      "─",
			Left:        "│",
			Right:       "│",
			TopLeft:     "╭",
			TopRight:    "╮",
			BottomLeft:  "┴",
			BottomRight: "┴",
		}, true).
		BorderForeground(lipgloss.Color(theme.borderColour)).
		Padding(0, 1)
}

func tabGap() lipgloss.Style {
	theme := currentTheme()
	return lipgloss.NewStyle().
		Border(lipgloss.Border{
			Top:         " ",
			Bottom:      "─",
			Left:        " ",
			Right:       " ",
			TopLeft:     " ",
			TopRight:    " ",
			BottomLeft:  "─",
			BottomRight: "─",
		}, true).
		BorderForeground(lipgloss.Color(theme.borderColour)).
		BorderTop(false).
		BorderBottom(true).
		BorderLeft(false).
		BorderRight(false)
}

func RunTUI() error {
	p := tea.NewProgram(initialModel())
	_, err := p.Run()
	return err
}
