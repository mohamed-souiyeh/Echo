package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var ClientRenderer *lipgloss.Renderer

var (
	// Base colors
	ColorBg      lipgloss.Style
	ColorFg      lipgloss.Style
	ColorAccent1 lipgloss.Style
	ColorAccent2 lipgloss.Style
	ColorSuccess lipgloss.Style
	ColorError   lipgloss.Style

	// Container for the entire page
	Container lipgloss.Style

	// Title + subtitle
	Title lipgloss.Style = ClientRenderer.NewStyle().Bold(true).
		Padding(0, 1)  // Add some horizontal padding

	Subtitle lipgloss.Style = lipgloss.NewStyle().
			Padding(0, 1)

	// Input boxes
	Input lipgloss.Style                    = ClientRenderer.NewStyle().
		Border(lipgloss.RoundedBorder()). // Keep the rounded border
		Padding(0, 1).
		Align(lipgloss.Center).
		Width(24) // Example fixed width

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
