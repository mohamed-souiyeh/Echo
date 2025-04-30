package components

import (
	"echo/tui/styles"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type AuthMode int

// it is very important to keep the order of (SignUp = 0) because the authForm depends on it
const (
	SignUp AuthMode = iota
	SignIn

	MaxMode
)

func (m AuthMode) String() string {
	switch m {
	case SignIn:
		return "Sign-In"
	case SignUp:
		return "Sign-Up"
	default:
		return fmt.Sprintf("AuthMode(%d)", m)
	}
}

type AuthForumModel struct {
	focusIndex int
	inputs     []textinput.Model
	AuthMode   AuthMode
}

func InitialAuthForumModel() AuthForumModel {

	m := AuthForumModel{
		focusIndex: int(SignUp),
		inputs:     make([]textinput.Model, 3),
		AuthMode:   SignUp, // also means authMode passed from the auth parent model
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
			t.Placeholder = "Nickname"
		case 1:
			t.Placeholder = "Username"
		case 2:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		m.inputs[i] = t

		//TODO: removing this or delaying it or showing the use how to focus on a field can hide the annoying ancii problem "]11;rgb:1e1e/1e1e/1e1e"
		if m.focusIndex == i {
			m.inputs[i].Focus()
			m.inputs[i].PromptStyle = styles.FocusedStyle
			m.inputs[i].TextStyle = styles.FocusedStyle
		}
	}

	return m
}

func (m AuthForumModel) Init() tea.Cmd {
	return textinput.Blink
}

func simulateDownKeyCmd() tea.Cmd {
	return func() tea.Msg {
		// Construct the specific KeyMsg for the Down arrow
		keyMsg := tea.KeyMsg{
			Type: tea.KeyDown, // This is the crucial part for the Down key
			// Runes field is empty for non-printable keys like arrows
			// Alt modifier is implicitly false
		}
		return keyMsg
	}
}

func (m AuthForumModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, authForumKeyMaps.AuthMode):
			var cmd tea.Cmd
			if m.AuthMode == SignUp && m.focusIndex == 0 {
				cmd = simulateDownKeyCmd()
			}
			m.AuthMode = (m.AuthMode + 1) % MaxMode
			return m, cmd
		case key.Matches(msg, authForumKeyMaps.Down, authForumKeyMaps.Up):

			if m.focusIndex < len(m.inputs) {
				m.inputs[m.focusIndex].Blur()
				m.inputs[m.focusIndex].PromptStyle = styles.NoStyle
				m.inputs[m.focusIndex].TextStyle = styles.NoStyle
			}

			if key.Matches(msg, authForumKeyMaps.Up) {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = int(m.AuthMode)
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			var cmd tea.Cmd
			if m.focusIndex < len(m.inputs) {
				cmd = m.inputs[m.focusIndex].Focus()
				m.inputs[m.focusIndex].PromptStyle = styles.FocusedStyle
				m.inputs[m.focusIndex].TextStyle = styles.FocusedStyle
			}

			return m, cmd
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

	for i := int(m.AuthMode); i < len(m.inputs); i++ {
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

	return b.String()
}
