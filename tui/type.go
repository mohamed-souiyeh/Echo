package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// chat types
const msgTimeoutDuration = time.Millisecond * 5000

type Window struct {
	Width  int
	Height int
}

type Renderer = lipgloss.Renderer

type mode int

const (
	Nav mode = iota
	Ins
)

func (m mode) String() string {
	switch m {
	case Nav:
		return "Nav"
	case Ins:
		return "Ins"
	default:
		return fmt.Sprintf("mode(%d)", m)
	}
}

type route int

const (
	Auth route = iota
	Chat
	MaxRoute
)

func (r route) String() string {
	switch r {
	case Auth:
		return "Auth"
	case Chat:
		return "Chat"
	case MaxRoute:
		return "MaxRoute"
	default:
		return fmt.Sprintf("route(%d)", r)
	}
}
