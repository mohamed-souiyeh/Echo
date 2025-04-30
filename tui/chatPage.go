package tui

import (
	db "echo/db/repository"
	"echo/tui/styles"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type ChatModel struct {
	currentUser db.User
	count       int
}

func InitChatModel() ChatModel {
	return ChatModel{}
}

func (m ChatModel) Init() tea.Cmd {
	return nil
}

func (m ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Handle Window Size Changes
	// case tea.WindowSizeMsg:
	// 	mc.width = msg.Width   // Update the width
	// 	mc.height = msg.Height // Update the height
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

func (m ChatModel) View() string {
	debug.print(func() {
		fmt.Println("we are return the chat view")
	})
	help := styles.HelpBar.Render("?:Toggle Help • ⏎: submit • ctrl+c: quit")
	return fmt.Sprintf("chat count is => %d", m.count) + "\n" + help
}
