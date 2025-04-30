package tui

import (
	comps "echo/tui/components"
	"echo/tui/styles"
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AuthModel struct {
	width     int
	height    int
	isLoading bool
	authForum comps.AuthForumModel
	spinner   spinner.Model
	help      help.Model
}

func InitialAuthModel() AuthModel {
	spin := spinner.New()

	spin.Spinner = styles.EchoSpinner

	m := AuthModel{
		isLoading: false,
		authForum: comps.InitialAuthForumModel(),
		spinner:   spin,
		help:      help.New(),
	}
	return m
}

func (m AuthModel) Init() tea.Cmd {

	return nil
}

func (m AuthModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = nil
	var cmd tea.Cmd = nil
	switch msg := msg.(type) {
	// Handle Window Size Changes
	case tea.WindowSizeMsg:
		m.width = msg.Width   // Update the width
		m.height = msg.Height // Update the height
		m.help.Width = msg.Width
	case spinner.TickMsg:
		if m.isLoading {
			m.spinner, cmd = m.spinner.Update(msg)

			cmds = append(cmds, cmd)
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, AuthKeyMaps.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, AuthKeyMaps.Submit):
			// TODO instead of handling submit in the auth model handle it in the authform model and send launch a cmd from the authForm model to send a msg indicating to the auth model to start the spinner
			m.isLoading = !m.isLoading //! this is not the actual logic for authenticating
			cmds = append(cmds, m.spinner.Tick)
		}
	}
	// add the forum update function here to handle text input, cursor blinking and authentication logic launching
	updatedAuthForum, cmd := m.authForum.Update(msg)

	m.authForum = updatedAuthForum.(comps.AuthForumModel)

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

// Assuming 'styles' is the package name from the theme files
func (m AuthModel) View() string {
	// build your auth page body as before
	header := lipgloss.JoinVertical(
		lipgloss.Center,
		styles.Title.Render(" welcome to Echo "),
		styles.Subtitle.Render(" a small chat application served over ssh \n"),
		styles.Subtitle.Render(m.authForum.AuthMode.String()+"\n"),
	)

	forum := m.authForum.View()

	spinner := ""

	if m.isLoading {
		spinner = m.spinner.View() + "accessing Echo"
	}

	help := m.help.View(AuthKeyMaps)

	// Add some vertical spacing if needed (lipgloss.Height includes margin)
	// vspace := lipgloss.Height(header) +
	// 	lipgloss.Height(forum) +
	// 	lipgloss.Height(help)

	body := lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		forum,
		lipgloss.NewStyle().Height(2).Render(""), // Adjust height calculation as needed
		spinner,
		lipgloss.NewStyle().Height(2).Render(""), // Adjust height calculation as needed
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
