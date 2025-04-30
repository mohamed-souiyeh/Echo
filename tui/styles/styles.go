package styles

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

var ClientRenderer *lipgloss.Renderer

var EchoSpinner spinner.Spinner = spinner.Spinner {
	Frames: []string {
		"   ( )   ",
		"  (( ))  ",
		" ((( ))) ",
		" ((   )) ",
		" (     ) ",
		" )     ( ",
		" ))   (( ",
		"  )) ((  ",
		"   ) (   ",
		"         ",
	},
	FPS: time.Second / 5,
}

var (
	// Base colors
	ColorBg      lipgloss.Style
	ColorFg      lipgloss.Style
	ColorAccent1 lipgloss.Style
	ColorAccent2 lipgloss.Style
	ColorSuccess lipgloss.Style
	ColorError   lipgloss.Style

	NoStyle = ClientRenderer.NewStyle()

	// authForm
	FocusedStyle = ClientRenderer.NewStyle().Foreground(lipgloss.Color("205"))
	BlurredStyle = ClientRenderer.NewStyle().Foreground(lipgloss.Color("240"))
	AuthFormFocusedButton = FocusedStyle.Render("[ lemme in ]")
	AuthFormBlurredButton = fmt.Sprintf("[ %s ]", BlurredStyle.Render("lemme in"))
	
	// Container for the entire page
	Container lipgloss.Style

	// Title + subtitle
	Title lipgloss.Style = ClientRenderer.NewStyle().Bold(true).
		Padding(0, 1)  // Add some horizontal padding

	Subtitle lipgloss.Style = ClientRenderer.NewStyle().
			Padding(0, 1)

	// Input boxes
	Input lipgloss.Style                    = ClientRenderer.NewStyle().
		Border(lipgloss.RoundedBorder()). // Keep the rounded border
		Padding(0, 1).
		Align(lipgloss.Center)

	// Buttons / actionable text
	Button lipgloss.Style                    = ClientRenderer.NewStyle().
		Border(lipgloss.RoundedBorder()). // Keep the rounded border
		Align(lipgloss.Center).
		Padding(0, 2). // More padding for button feel
		Margin(0, 0).  // Add vertical margin
		Width(14)

	// Help bar at bottom
	HelpBar lipgloss.Style = ClientRenderer.NewStyle().
		Padding(1, 1)  // Padding top/bottom and sides
)
