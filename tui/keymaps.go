package tui

import "github.com/charmbracelet/bubbles/key"

// TODO: use this "https://github.com/charmbracelet/bubbletea/tree/main/examples/help" to use key binding properly in the help menu
type GlobalKeyMap struct {
	Enter key.Binding
	Quit key.Binding
	Info key.Binding
}


var GlobalKeyMaps = GlobalKeyMap {
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "submit"),
	),
}