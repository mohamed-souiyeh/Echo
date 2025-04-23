package tui

import (
	"echo/tui/styles"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)


type AuthModel struct {
	width  int
	height int
	count  int
}

func InitialAuthModel() AuthModel {
	return AuthModel{}
}

func (m AuthModel) Init() tea.Cmd {

	return nil
}

func (m AuthModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Handle Window Size Changes
	case tea.WindowSizeMsg:
		m.width = msg.Width   // Update the width
		m.height = msg.Height // Update the height
	case tea.KeyMsg:
		switch msg.String() {
		case "k":
			m.count++
			return m, nil
		case "j":
			m.count--
			return m, nil
		}
	}
	return m, nil
}

// Assuming 'styles' is the package name from the theme files
func (m AuthModel) View() string {
	// build your auth page body as before
	header := lipgloss.JoinVertical(
		lipgloss.Center,
		styles.Title.Render(" welcome to Echo "),
		styles.Subtitle.Render(" a small chat application served over ssh "),
	)

	forum := lipgloss.JoinVertical(lipgloss.Center,
		styles.Input.Render(" Username "),
		styles.Input.Render(" Password "),
		styles.Button.Render(" lemme in "),
	)

	help := styles.HelpBar.Render("?:Toggle Help • ⏎: submit • ctrl+c: quit")

	// Add some vertical spacing if needed (lipgloss.Height includes margin)
	// vspace := lipgloss.Height(header) +
	// 	lipgloss.Height(userInput) +
	// 	lipgloss.Height(passInput) +
	// 	lipgloss.Height(submitBtn) +
	// 	lipgloss.Height(help)

	body := lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		forum,
		// Add flexible space above help if desired
		// lipgloss.NewStyle().Height(m.height-vspace-2).Render(""), // Adjust height calculation as needed
		help,
	)
	debug.print(func() {
		fmt.Println("we are in auth returning the current view")
	})
	// return fmt.Sprintf("auth count is => %d", m.count) + "\n" + help
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
