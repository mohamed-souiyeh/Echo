package tui

import (
	db "echo/db/repository"
	"echo/services"
	"echo/tui/keymaps"
	"echo/tui/messages"
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type RootModel struct {
	activeRoute     route
	currentMode     mode
	Routes          []tea.Model
	isAuthenticated bool
	userService     *services.UserService
	quiting         bool
}

func InitialRootModel(userRepo db.UserRepository) RootModel {
	return RootModel{
		activeRoute:     Auth,
		currentMode:     Nav,
		Routes:          []tea.Model{InitialAuthModel(), InitChatModel()},
		isAuthenticated: false,
		userService:     services.NewUserService(userRepo),
		quiting:         false,
	}
}

func (m RootModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, route := range m.Routes {
		cmds = append(cmds, route.Init())
	}
	return tea.Batch(cmds...)
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd = nil
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		var cmds []tea.Cmd = nil
		for routeidx, route := range m.Routes {
			m.Routes[routeidx], cmd = route.Update(msg)
			cmds = append(cmds, cmd)
		}
		debug.print(func() {
			fmt.Println("window size event")
		})
		return m, tea.Batch(cmds...)
	case messages.SignUpAttemptMsg:
		cmd := m.signUp(msg.Username, msg.Password)
		return m, cmd
	case messages.SignInAttemptMsg:
		cmd := m.signIn(msg.Username, msg.Password)
		return m, cmd
	case messages.AuthFailedMsg:
		if m.activeRoute == Auth {
			m.Routes[m.activeRoute], cmd = m.Routes[m.activeRoute].Update(msg)
		}
		return m, cmd
	case messages.AuthSuccessMsg:
		if m.activeRoute == Auth {
			m.isAuthenticated = true
			m.Routes[m.activeRoute], cmd = m.Routes[m.activeRoute].Update(msg)
		}
		return m, cmd
	case messages.AccessChatMsg:
		if m.isAuthenticated {
			m.activeRoute = Chat
			keymaps.AuthKeyMaps.Deactivate()
			keymaps.ChatKeyMaps.Activate()
			m.Routes[m.activeRoute], cmd = m.Routes[m.activeRoute].Update(msg)
		}
		return m, cmd
	case messages.LogoutMsg:
		m.isAuthenticated = false
		m.activeRoute = Auth
		keymaps.ChatKeyMaps.Deactivate()
		keymaps.AuthKeyMaps.Activate()
		m.Routes[m.activeRoute], cmd = m.Routes[m.activeRoute].Update(msg)
		return m, cmd
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keymaps.GlobalKeyMaps.Quit):
			m.quiting = true
			return m, tea.Quit
		default:
			m.Routes[m.activeRoute], cmd = m.Routes[m.activeRoute].Update(msg)
			return m, cmd
		}
	default:
		m.Routes[m.activeRoute], cmd = m.Routes[m.activeRoute].Update(msg)
		return m, cmd
	}
}

func (m RootModel) signUp(username string, password string) tea.Cmd {
	return func() tea.Msg {

		if len(username) < 1 {
			return messages.AuthFailedMsg{
				Reason:      "You can't be nameless, put a username already -_-",
				DebugReason: "username too short",
			}
		}
		if len(password) < 6 {
			return messages.AuthFailedMsg{
				Reason:      "Size does matter in passwords, it need to be at least 13 characters long",
				DebugReason: "password too short",
			}
		}

		return m.userService.SignUp(username, password)
	}
}

func (m RootModel) signIn(username, password string) tea.Cmd {
	return func() tea.Msg {

		if len(username) < 1 || len(password) < 1 {
			return messages.AuthFailedMsg{
				Reason: "You aren't expecting to get access without crediantials, are u?",
			}
		}

		return m.userService.SignIn(username, password)
	}
}

func (m RootModel) View() string {
	if m.quiting {
		return "byeee!\n"
	}
	debug.print(func() {
		fmt.Println("current view:", "'", m.Routes[m.activeRoute].View(), "'")
		fmt.Println("\nactive route: ", m.activeRoute)
	})
	return m.Routes[m.activeRoute].View() + fmt.Sprint("\nactive route: ", m.activeRoute)
}
