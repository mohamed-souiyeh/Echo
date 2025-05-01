package keymaps

import "github.com/charmbracelet/bubbles/key"

// TODO: use this "https://github.com/charmbracelet/bubbletea/tree/main/examples/help" to use key binding properly in the help menu
type GlobalKeyMap struct {
	Quit key.Binding
	Help key.Binding
}

var GlobalKeyMaps = GlobalKeyMap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
}

type AuthKeyMap struct {
	GlobalKeyMap
	AuthMode key.Binding
	Submit key.Binding
	Up     key.Binding
	Down   key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k AuthKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit, k.AuthMode, k.Submit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k AuthKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Help, k.Submit},   // first column
		{k.Quit, k.AuthMode}, // second column
		{k.Up, k.Down},
	}
}

var AuthKeyMaps = AuthKeyMap{
	GlobalKeyMap: GlobalKeyMaps,
	Submit: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("⏎", "sign-In/Up"),
	),
	AuthMode: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "toggle sign-In/Up"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "shift+tab"),
		key.WithHelp("↑", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "tab"),
		key.WithHelp("↓", "mode down"),
	),
}
