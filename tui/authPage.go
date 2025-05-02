package tui

import (
	comps "echo/tui/components"
	"echo/tui/keymaps"
	"echo/tui/styles"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AuthModel struct {
	width     int
	height    int
	authForum comps.AuthForumModel
	help      help.Model
}

func InitialAuthModel() AuthModel {

	m := AuthModel{
		authForum: comps.InitialAuthForumModel(),
		help:      help.New(),
	}
	return m
}

func (m AuthModel) Init() tea.Cmd {

	return nil
}

func (m AuthModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// var cmds []tea.Cmd = nil
	var cmd tea.Cmd = nil
	switch msg := msg.(type) {
	// Handle Window Size Changes
	case tea.WindowSizeMsg:
		m.width = msg.Width   // Update the width
		m.height = msg.Height // Update the height
		m.help.Width = msg.Width
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keymaps.AuthKeyMaps.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		}
	}
	// add the forum update function here to handle text input, cursor blinking and authentication logic launching
	updatedAuthForum, cmd := m.authForum.Update(msg)

	m.authForum = updatedAuthForum.(comps.AuthForumModel)

	return m, cmd
}

// Assuming 'styles' is the package name from the theme files
func (m AuthModel) View() string {
	// build your auth page body as before
	header := lipgloss.JoinVertical(
		lipgloss.Center,
		styles.Title.Width(m.width).Align(lipgloss.Center).Render(" welcome to Echo "),
		styles.Subtitle.Width(m.width).Align(lipgloss.Center).Render(" a small chat application served over ssh \n"),
		styles.Subtitle.Render(m.authForum.AuthMode.String()+"\n"),
	)

	forum := m.authForum.View()

	help := m.help.View(keymaps.AuthKeyMaps)

	// Add some vertical spacing if needed (lipgloss.Height includes margin)
	// vspace := lipgloss.Height(header) +
	// 	lipgloss.Height(forum) +
	// 	lipgloss.Height(help)

	body := lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		forum,
		lipgloss.NewStyle().Height(2).Render(""), // Adjust height calculation as needed
		help,
	)

	// Use lipgloss.Place WITHOUT the background color option
	return styles.ClientRenderer.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		body,
		// lipgloss.WithWhitespaceForeground(styles.ColorBg), // REMOVED
	)
}
