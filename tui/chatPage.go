package tui

import (
	db "echo/db/repository"
	"echo/tui/commands"
	"echo/tui/keymaps"
	"echo/tui/messages"
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type ChatModel struct {
	width       int
	height      int
	CurrentUser db.User
	help        help.Model
}

func InitChatModel() ChatModel {
	return ChatModel{
		help: help.New(),
	}
}

func (m ChatModel) Init() tea.Cmd {
	return nil
}

func (m ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Handle Window Size Changes
	case tea.WindowSizeMsg:
		m.width = msg.Width   // Update the width
		m.height = msg.Height // Update the height
		m.help.Width = msg.Width
	case messages.AccessChatMsg:
		m.CurrentUser = msg.User
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keymaps.ChatKeyMaps.Logout):
			return m, commands.LogoutCmd(m.CurrentUser.Username)
		}
	}
	return m, nil
}

func (m ChatModel) View() string {
	help := m.help.View(keymaps.ChatKeyMaps)
	return fmt.Sprintf("chat current user is => %s", m.CurrentUser.Username) + "\n" + help
}
