package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// TODO: add athorization logic to routes
// and implement the authorization toggling as a cmd that is returned from the chiled route that return a msg to indicate to the root model that it need to toggel the authorized status.
type RootModel struct {
	activeRoute route
	currentMode mode
	Routes      []tea.Model
	authorized  bool
	quiting     bool
}

func InitialRootModel() RootModel {
	return RootModel{
		activeRoute: Auth,
		currentMode: Nav,
		Routes:      []tea.Model{InitialAuthModel(), InitChatModel()},
		authorized:  false,
		quiting:     false,
	}
}

// TODO: make it call the childs init and return there returned cmds as a batch => DONE
func (m RootModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, route := range m.Routes {
		cmds = append(cmds, route.Init())
	}
	return tea.Batch(cmds...)
}

// TODO: handel global keybindings in the root model update. => DONE
// TODO: on windowsizemsg call all child updates and return there cmds as a batch => DONE
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
	case tea.KeyMsg:
		switch msg.String() {
		// ! remove this switching key binding
		case "s":
			m.activeRoute = (m.activeRoute + 1) % MaxRoute
		case "q", "ctrl+c":
			m.quiting = true
			return m, tea.Quit
		default:
			m.Routes[m.activeRoute], cmd = m.Routes[m.activeRoute].Update(msg)
		}
	}
	return m, cmd
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
