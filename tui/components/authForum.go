package components

import (
	"echo/tui/commands"
	"echo/tui/keymaps"
	"echo/tui/messages"
	"echo/tui/styles"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AuthForumModel struct {
	width        int
	height       int
	focusIndex   int
	inputs       []textinput.Model
	isLoading    bool
	AuthMode     AuthMode
	spinner      spinner.Model
	responseMsg  string
	showResponse bool
}

func InitialAuthForumModel() AuthForumModel {

	spin := spinner.New()

	spin.Spinner = styles.EchoSpinner

	m := AuthForumModel{
		focusIndex:   -1,
		inputs:       make([]textinput.Model, 2),
		isLoading:    false,
		AuthMode:     SignUp,
		spinner:      spin,
		responseMsg:  "",
		showResponse: false,
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Prompt = ""
		t.Cursor.Style = styles.FocusedStyle
		t.CharLimit = 20
		t.Width = 20

		switch i {
		case 0:
			t.Placeholder = "Username"
		case 1:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		m.inputs[i] = t

		//TODO: removing this or delaying it or showing the use how to focus on a field can hide the annoying ancii problem "]11;rgb:1e1e/1e1e/1e1e"
		// if m.focusIndex == i {
		// 	m.inputs[i].Focus()
		// 	m.inputs[i].PromptStyle = styles.FocusedStyle
		// 	m.inputs[i].TextStyle = styles.FocusedStyle
		// }
	}

	return m
}

func (m AuthForumModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m AuthForumModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width   // Update the width
		m.height = msg.Height // Update the height
		return m, nil
	case spinner.TickMsg:
		if m.isLoading {
			spinner, cmd := m.spinner.Update(msg)

			m.spinner = spinner
			return m, cmd
		}
	case messages.AuthFailedMsg:
		m.isLoading = false
		m.showResponse = true
		m.responseMsg = msg.Reason
		return m, nil
	case messages.AuthSuccessMsg:
		m.inputs[0].Reset()
		m.inputs[1].Reset()
		m.isLoading = false
		m.showResponse = false
		m.responseMsg = ""
		m.focusIndex = 0
		return m, commands.AccessChatCmd(msg.User)
	case messages.LogoutMsg:
		var cmd tea.Cmd = nil

		m.showResponse = true
		m.responseMsg = "See you soon ^_^ " + msg.Username
		m.focusIndex = 0 // i know, i am assigning here just out of paranoia, maybe a bad reason but i will keep it for now.
		cmd = m.inputs[m.focusIndex].Focus()
		m.inputs[m.focusIndex].PromptStyle = styles.FocusedStyle
		m.inputs[m.focusIndex].TextStyle = styles.FocusedStyle
		return m, tea.Batch(textinput.Blink, cmd)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keymaps.AuthKeyMaps.AuthMode):
			m.AuthMode = (m.AuthMode + 1) % MaxMode
			return m, nil
		case key.Matches(msg, keymaps.AuthKeyMaps.Down, keymaps.AuthKeyMaps.Up):

			if m.focusIndex < len(m.inputs) && m.focusIndex >= 0 {
				m.inputs[m.focusIndex].Blur()
				m.inputs[m.focusIndex].PromptStyle = styles.NoStyle
				m.inputs[m.focusIndex].TextStyle = styles.NoStyle
			}

			if key.Matches(msg, keymaps.AuthKeyMaps.Up) {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			var cmd tea.Cmd
			if m.focusIndex < len(m.inputs) && m.focusIndex >= 0 {
				cmd = m.inputs[m.focusIndex].Focus()
				m.inputs[m.focusIndex].PromptStyle = styles.FocusedStyle
				m.inputs[m.focusIndex].TextStyle = styles.FocusedStyle
			}

			return m, cmd
		case key.Matches(msg, keymaps.AuthKeyMaps.Submit):
			m.isLoading = true
			m.showResponse = false
			var cmds []tea.Cmd

			if m.focusIndex < len(m.inputs) && m.focusIndex >= 0 {
				m.inputs[m.focusIndex].Blur()
				m.inputs[m.focusIndex].PromptStyle = styles.NoStyle
				m.inputs[m.focusIndex].TextStyle = styles.NoStyle
			}

			m.focusIndex = len(m.inputs)
			cmds = append(cmds, m.spinner.Tick)

			if m.AuthMode == SignUp {
				cmds = append(cmds, commands.SignUpAttemptCmd(m.inputs[0].Value(), m.inputs[1].Value()))
			} else if m.AuthMode == SignIn {
				cmds = append(cmds, commands.SignInAttemptCmd(m.inputs[0].Value(), m.inputs[1].Value()))
			}

			return m, tea.Batch(cmds...)
		}
	}

	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *AuthForumModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m AuthForumModel) View() string {
	var b strings.Builder

	for i := 0; i < len(m.inputs); i++ {
		b.WriteString(styles.Input.Render(m.inputs[i].View()))
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &styles.AuthFormBlurredButton
	if m.focusIndex == len(m.inputs) {
		button = &styles.AuthFormFocusedButton
	}

	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	spinner := ""

	if m.isLoading {
		spinner = styles.Subtitle.Width(m.width).Align(lipgloss.Center).Render(m.spinner.View() + "accessing Echo")
	} else if m.showResponse {
		spinner = styles.Subtitle.Width(m.width).Align(lipgloss.Center).Render(m.responseMsg)
	}

	forum := lipgloss.JoinVertical(
		lipgloss.Center,
		b.String(),
		lipgloss.NewStyle().Height(2).Render(""), // Adjust height calculation as needed
		spinner,
	)

	return forum
}
