package tui

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nam-ruto/imperium/internal/models"
	"github.com/nam-ruto/imperium/internal/search"
	"github.com/nam-ruto/imperium/internal/storage"
)

// SearchModel is the Bubble Tea model for search mode
type SearchModel struct {
	textInput textinput.Model
	commands  []models.CommandRecord
	filtered  []models.CommandRecord
	cursor    int
	selected  *models.CommandRecord
	storage   *storage.JSONStorage
	width     int
	height    int
	message   string
	quitting  bool
}

// NewSearchModel creates a new search model
func NewSearchModel(store *storage.JSONStorage, initialQuery string) SearchModel {
	ti := textinput.New()
	ti.Placeholder = "Type to search..."
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 50
	ti.SetValue(initialQuery)
	ti.Prompt = ""
	ti.PromptStyle = PromptStyle
	ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))

	commands, _ := store.Load()
	filtered := search.Search(initialQuery, commands)

	return SearchModel{
		textInput: ti,
		commands:  commands,
		filtered:  filtered,
		cursor:    0,
		storage:   store,
		width:     80,
		height:    24,
	}
}

// Init implements tea.Model
func (m SearchModel) Init() tea.Cmd {
	// Clear screen and start cursor blink
	return tea.Batch(tea.ClearScreen, textinput.Blink)
}

// Update implements tea.Model
func (m SearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "ctrl+p":
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil

		case "down", "ctrl+n":
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
			}
			return m, nil

		case "ctrl+u":
			// Page up - move cursor up by 5
			m.cursor -= 5
			if m.cursor < 0 {
				m.cursor = 0
			}
			return m, nil

		case "ctrl+d":
			// Page down - move cursor down by 5
			m.cursor += 5
			if m.cursor >= len(m.filtered) {
				m.cursor = len(m.filtered) - 1
			}
			if m.cursor < 0 {
				m.cursor = 0
			}
			return m, nil

		case "enter":
			if len(m.filtered) > 0 && m.cursor < len(m.filtered) {
				m.selected = &m.filtered[m.cursor]
				// Copy to clipboard
				if err := clipboard.WriteAll(m.selected.Command); err == nil {
					m.message = fmt.Sprintf("Copied: %s", m.selected.Command)
					// Update last used
					m.storage.UpdateLastUsed(m.selected.Command)
				}
				m.quitting = true
				return m, tea.Quit
			}
			return m, nil

		case "esc":
			// If input is empty, quit; otherwise clear
			if m.textInput.Value() == "" {
				m.quitting = true
				return m, tea.Quit
			}
			m.textInput.SetValue("")
			m.filtered = search.Search("", m.commands)
			m.cursor = 0
			return m, nil

		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		}
	}

	// Handle text input
	prevValue := m.textInput.Value()
	m.textInput, cmd = m.textInput.Update(msg)

	// If input changed, update search results
	if m.textInput.Value() != prevValue {
		m.filtered = search.Search(m.textInput.Value(), m.commands)
		// Reset cursor if out of bounds
		if m.cursor >= len(m.filtered) {
			m.cursor = len(m.filtered) - 1
		}
		if m.cursor < 0 {
			m.cursor = 0
		}
	}

	return m, cmd
}

// View implements tea.Model
func (m SearchModel) View() string {
	if m.quitting {
		if m.message != "" {
			return SuccessStyle.Render(m.message) + "\n"
		}
		return ""
	}

	var b strings.Builder

	// Calculate max visible rows based on terminal height
	maxVisible := 10
	if m.height > 0 {
		maxVisible = m.height - 10 // Account for header, input, help, etc.
		if maxVisible < 5 {
			maxVisible = 5
		}
		if maxVisible > 20 {
			maxVisible = 20
		}
	}

	// Calculate available width
	availableWidth := m.width - 4
	if availableWidth < 60 {
		availableWidth = 60
	}

	// Title
	b.WriteString(TitleStyle.Render("[>] imperium"))
	b.WriteString("\n\n")

	// Search input with border
	inputBox := InputStyle.Render(m.textInput.View())
	b.WriteString(inputBox)
	b.WriteString("\n")

	// Results count
	countText := fmt.Sprintf("  %d/%d commands", len(m.filtered), len(m.commands))
	b.WriteString(DimStyle.Render(countText))
	b.WriteString("\n\n")

	// Calculate responsive column widths
	cmdWidth, usageWidth, tagsWidth := m.calculateColumnWidths(availableWidth)

	// Table header
	header := fmt.Sprintf("  %-*s  %-*s  %s",
		cmdWidth, "Command",
		usageWidth, "Usage",
		"Tags",
	)
	b.WriteString(HeaderStyle.Render(header))
	b.WriteString("\n")

	// Table separator
	totalWidth := cmdWidth + usageWidth + tagsWidth + 10
	sep := strings.Repeat("─", totalWidth)
	b.WriteString(DimStyle.Render("  " + sep))
	b.WriteString("\n")

	// Calculate scroll window
	start := 0
	if m.cursor >= maxVisible {
		start = m.cursor - maxVisible + 1
	}
	end := start + maxVisible
	if end > len(m.filtered) {
		end = len(m.filtered)
	}

	// Render rows - always render maxVisible rows to prevent ghost lines
	rowsRendered := 0
	for i := start; i < end; i++ {
		record := m.filtered[i]
		row := m.formatRow(record, i == m.cursor, cmdWidth, usageWidth, tagsWidth)
		b.WriteString(row)
		b.WriteString("\n")
		rowsRendered++
	}

	// Handle empty state
	if len(m.filtered) == 0 {
		emptyMsg := "  No matching commands"
		b.WriteString(DimStyle.Render(emptyMsg))
		b.WriteString("\n")
		rowsRendered++
	}

	// Pad with empty lines to maintain consistent height (prevents ghost lines)
	for i := rowsRendered; i < maxVisible; i++ {
		b.WriteString(strings.Repeat(" ", totalWidth+4))
		b.WriteString("\n")
	}

	// Scroll indicator (or empty line to maintain consistent height)
	if len(m.filtered) > maxVisible {
		scrollInfo := fmt.Sprintf("  [%d-%d of %d]", start+1, end, len(m.filtered))
		b.WriteString(DimStyle.Render(scrollInfo))
	} else {
		b.WriteString(strings.Repeat(" ", 20))
	}
	b.WriteString("\n\n")

	// Help text
	help := "  ↑/↓: navigate | enter: copy & exit | ctrl+u/d: page | esc: clear/quit"
	b.WriteString(DimStyle.Render(help))

	return b.String()
}

// calculateColumnWidths returns responsive column widths based on available space
func (m SearchModel) calculateColumnWidths(available int) (cmd, usage, tags int) {
	if available >= 120 {
		return 45, 40, 25
	} else if available >= 100 {
		return 40, 35, 20
	} else if available >= 80 {
		return 35, 30, 15
	} else {
		return 30, 25, 10
	}
}

// formatRow formats a single result row
func (m SearchModel) formatRow(cmd models.CommandRecord, isSelected bool, cmdWidth, usageWidth, tagsWidth int) string {
	// Format columns
	cmdStr := truncate(cmd.Command, cmdWidth)
	cmdStr = padRight(cmdStr, cmdWidth)

	usageStr := truncate(cmd.Usage, usageWidth)
	usageStr = padRight(usageStr, usageWidth)

	tagsStr := ""
	if len(cmd.Tags) > 0 {
		tagsStr = truncate(strings.Join(cmd.Tags, ", "), tagsWidth)
	}
	tagsStr = padRight(tagsStr, tagsWidth)

	// Build row content
	row := fmt.Sprintf("%-*s  %-*s  %s", cmdWidth, cmdStr, usageWidth, usageStr, tagsStr)

	// Apply styles
	if isSelected {
		return SelectedStyle.Render("> " + row)
	}
	return NormalStyle.Render("  " + row)
}

// Selected returns the selected command (if any)
func (m SearchModel) Selected() *models.CommandRecord {
	return m.selected
}

// truncate shortens a string to max length with ellipsis
func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if max <= 3 {
		return s[:max]
	}
	return s[:max-3] + "..."
}

// padRight pads a string to the specified width, accounting for ANSI codes
func padRight(s string, width int) string {
	visibleLen := lipgloss.Width(s)
	if visibleLen >= width {
		return s
	}
	return s + strings.Repeat(" ", width-visibleLen)
}

// RunSearch runs the search TUI
func RunSearch(store *storage.JSONStorage, initialQuery string) (*models.CommandRecord, error) {
	model := NewSearchModel(store, initialQuery)
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(), // Use alternate screen buffer for clean rendering
	)

	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}

	if m, ok := finalModel.(SearchModel); ok {
		return m.Selected(), nil
	}
	return nil, nil
}
