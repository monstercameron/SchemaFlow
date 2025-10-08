package tui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Custom item delegate for better highlighting
type itemDelegate struct{}

func (d itemDelegate) Height() int {
	// Dynamic height based on content - increased for tasks
	return 6
}
func (d itemDelegate) Spacing() int                            { return 1 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(todoItem)
	if !ok {
		return
	}

	str := i.Title()
	desc := i.Description()

	var itemStyle lipgloss.Style
	var textStyle lipgloss.Style
	
	// Apply styles based on completion status
	if i.todo.Completed {
		// Completed items are dimmed
		textStyle = lipgloss.NewStyle().
			Foreground(mutedColor)
	} else {
		textStyle = lipgloss.NewStyle()
	}
	
	if index == m.Index() {
		// Selected item - highlight with background and border
		itemStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#1a1a2e")).
			BorderStyle(lipgloss.NormalBorder()).
			BorderLeft(true).
			BorderForeground(primaryColor).
			PaddingLeft(1).
			Width(m.Width() - 2).
			Bold(true)
	} else {
		// Non-selected item
		itemStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			Width(m.Width() - 2)
	}

	// Render the item - title on first line, description below
	// Apply strikethrough to completed items using ANSI codes
	if i.todo.Completed {
		str = formatCompletedText(str, true)
	}
	
	fullContent := textStyle.Render(str)
	if desc != "" {
		fullContent += "\n" + textStyle.Render(desc)
	}
	fmt.Fprint(w, itemStyle.Render(fullContent))
}