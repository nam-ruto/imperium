package tui

import "github.com/charmbracelet/lipgloss"

// Styles for the TUI
var (
	// Title style
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205"))

	// Prompt style for the input line
	PromptStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205"))

	// Input box style with border
	InputStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("205")).
			Padding(0, 1)

	// Header style for table headers (no background for transparency)
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("39"))

	// Selected item style (reversed for visibility)
	SelectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")).
			Bold(true)

	// Normal item style
	NormalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	// Dim text for secondary information
	DimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	// Success message style
	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")). // Green
			Bold(true)

	// Error message style
	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("203")) // Red

	// Pointer/cursor indicator
	PointerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true)

	// Column separator
	SeparatorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)
