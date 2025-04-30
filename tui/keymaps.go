package tui

import "github.com/charmbracelet/bubbles/key"

// TODO: use this "https://github.com/charmbracelet/bubbletea/tree/main/examples/help" to use key binding properly in the help menu
type GlobalKeyMap struct {
	Quit key.Binding
	Mode key.Binding
	Help key.Binding
}


var GlobalKeyMaps = GlobalKeyMap{
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit"),
		),
		Mode: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("ctrl+s", "toggle sign-In/Up"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "expand help"),
		),
	}

type AuthKeyMap struct {
	GlobalKeyMap
	Submit key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k AuthKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit, k.Mode, k.Submit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k AuthKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Help, k.Quit}, // first column
		{k.Mode, k.Submit}, // second column
	}
}

var AuthKeyMaps = AuthKeyMap{
	GlobalKeyMap: GlobalKeyMaps,
	Submit: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("⏎", "sign-in/up"),
	),
}
