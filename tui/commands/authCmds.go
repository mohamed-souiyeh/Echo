package commands

import (
	db "echo/db/repository"
	"echo/tui/messages"
	msgs "echo/tui/messages"

	tea "github.com/charmbracelet/bubbletea"
)

func SignUpAttemptCmd(username string, password string) tea.Cmd {
	return func() tea.Msg {
		return msgs.SignUpAttemptMsg{
			Username: username,
			Password: password,
		}
	}
}

func SignInAttemptCmd(username string, password string) tea.Cmd {
	return func() tea.Msg {
		return msgs.SignInAttemptMsg{
			Username: username,
			Password: password,
		}
	}
}

func AuthSucessCmd(user db.User) tea.Cmd {
	return func() tea.Msg {
		return msgs.AuthSuccessMsg{
			User: user,
		}
	}
}

func AutFailedCmd(reason string, debugReason string) tea.Cmd {
	return func() tea.Msg {
		return msgs.AuthFailedMsg{
			Reason:      reason,
			DebugReason: debugReason,
		}
	}
}

func AccessChatCmd(user db.User) tea.Cmd {
	return func() tea.Msg {
		return msgs.AccessChatMsg{
			User: user,
		}
	}
}

func LogoutCmd(username string) tea.Cmd {
	return func() tea.Msg {
		return messages.LogoutMsg{
			Username: username,
		}
	}
}
