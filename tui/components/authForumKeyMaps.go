package components

import "github.com/charmbracelet/bubbles/key"

type authForumKeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Submit   key.Binding
	AuthMode key.Binding
}

var authForumKeyMaps = authForumKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "shift+tab"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "tab"),
	),
	Submit: key.NewBinding(
		key.WithKeys("enter"),
	),
	AuthMode: key.NewBinding(
		key.WithKeys("ctrl+s"),
	),
}
