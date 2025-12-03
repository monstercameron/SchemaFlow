package tui

import (
	"github.com/charmbracelet/lipgloss"
)

func (m Model) quitConfirmViewRender() string {
	// Create confirmation modal
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(errorColor).
		Padding(2, 4).
		Background(lipgloss.Color("#1a1a1a"))

	title := lipgloss.NewStyle().
		Foreground(errorColor).
		Bold(true).
		MarginBottom(1).
		Render("Quit Smart Todo?")

	message := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffffff")).
		MarginBottom(2).
		Render("Are you sure you want to quit?")

	// Create button-like options
	enterOption := lipgloss.NewStyle().
		Foreground(errorColor).
		Bold(true).
		Padding(0, 2).
		Background(lipgloss.Color("#330000")).
		Render("[ Enter ] Yes, Quit")

	escOption := lipgloss.NewStyle().
		Foreground(successColor).
		Bold(true).
		Padding(0, 2).
		Background(lipgloss.Color("#003300")).
		MarginLeft(2).
		Render("[ Esc ] Cancel")

	buttons := lipgloss.JoinHorizontal(
		lipgloss.Center,
		enterOption,
		escOption,
	)

	modal := modalStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			title,
			message,
			buttons,
		),
	)

	// Create overlay with modal centered
	overlay := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		modal,
	)

	return overlay
}

