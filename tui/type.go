package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type Renderer = lipgloss.Renderer

type mode int

const (
	Nav mode = iota
	Ins
)

// String implements the fmt.Stringer interface for the mode type.
func (m mode) String() string {
	switch m {
	case Nav:
		return "Nav" // Return the string representation for Nav
	case Ins:
		return "Ins" // Return the string representation for Ins
	default:
		// Handle potential invalid values gracefully
		return fmt.Sprintf("mode(%d)", m)
	}
}

type route int

const (
	Auth route = iota
	Chat
	MaxRoute
)

// String implements the fmt.Stringer interface for the route type.
func (r route) String() string {
	// You can also use an array/slice lookup if values are contiguous
	// and you handle bounds checks, but switch is often clearer.
	switch r {
	case Auth:
		return "Auth"
	case Chat:
		return "Chat"
	case MaxRoute:
		return "MaxRoute"
	default:
		// Handle potential invalid values
		return fmt.Sprintf("route(%d)", r)
	}
}