package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) addViewRender() string {
	// Render the background (main list view)
	background := m.listViewRender()
	
	// Modal dimensions
	modalWidth := 60
	if modalWidth > m.width-10 {
		modalWidth = m.width - 10
	}
	
	// Create modal header
	header := lipgloss.NewStyle().
		Background(lipgloss.Color("#1a1a2e")).
		Foreground(primaryColor).
		Bold(true).
		Padding(0, 2).
		Width(modalWidth - 2).
		Align(lipgloss.Center).
		Render("New Task")
	
	// Convert textarea to single-line input for modal
	// Use the input field instead of textarea for consistency
	if !m.input.Focused() {
		m.input.Focus()
	}
	m.input.Width = modalWidth - 6
	
	// AI capabilities hint
	aiHint := lipgloss.NewStyle().
		Foreground(mutedColor).
		Italic(true).
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			"🤖 AI will extract:",
			"• Priority, deadline, category",
			"• Subtasks if mentioned",
			"• Location context",
		))
	
	// Input section
	inputSection := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(successColor).
		Padding(0, 1).
		Width(modalWidth - 4).
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.NewStyle().Foreground(successColor).Bold(true).Render("📝 What needs to be done?"),
			m.input.View(),
		))
	
	// Instructions
	instructions := lipgloss.NewStyle().
		Foreground(mutedColor).
		Render("Enter: Add • Esc: Cancel")
	
	// Build modal content
	modalContent := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		inputSection,
		"",
		aiHint,
		"",
		instructions,
	)
	
	// Create modal with border and shadow effect
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(primaryColor).
		Background(lipgloss.Color("#0a0a0a")).
		Padding(1, 2).
		Width(modalWidth)
	
	modal := modalStyle.Render(modalContent)
	
	// Create an overlay effect
	modalOverlay := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		modal,
	)
	
	// Combine background and modal
	bgLines := strings.Split(background, "\n")
	overlayLines := strings.Split(modalOverlay, "\n")
	
	result := make([]string, len(bgLines))
	for i := range bgLines {
		if i < len(overlayLines) && strings.TrimSpace(overlayLines[i]) != "" {
			result[i] = overlayLines[i]
		} else {
			result[i] = bgLines[i]
		}
	}
	
	return strings.Join(result, "\n")
}

func (m Model) editViewRender() string {
	if m.selectedTodo == nil {
		return m.listViewRender()
	}
	
	// Render the background (main list view)
	background := m.listViewRender()
	
	// Modal dimensions
	modalWidth := 70
	if modalWidth > m.width-10 {
		modalWidth = m.width - 10
	}
	
	// Create modal header
	header := lipgloss.NewStyle().
		Background(lipgloss.Color("#1a1a2e")).
		Foreground(primaryColor).
		Bold(true).
		Padding(0, 2).
		Width(modalWidth - 2).
		Align(lipgloss.Center).
		Render("Edit Task")
	
	// Current todo info
	todoInfo := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(mutedColor).
		Padding(0, 1).
		Width(modalWidth - 4).
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.NewStyle().Foreground(primaryColor).Bold(true).Render("Current Todo:"),
			m.selectedTodo.Title,
			"",
			lipgloss.NewStyle().Foreground(mutedColor).Render(fmt.Sprintf("Priority: %s | Category: %s", m.selectedTodo.Priority, m.selectedTodo.Category)),
		))
	
	// Context input section
	if !m.editInput.Focused() {
		m.editInput.Focus()
	}
	m.editInput.Width = modalWidth - 6
	
	inputSection := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(successColor).
		Padding(0, 1).
		Width(modalWidth - 4).
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.NewStyle().Foreground(successColor).Bold(true).Render("📝 Add Context:"),
			m.editInput.View(),
			"",
			lipgloss.NewStyle().Foreground(mutedColor).Italic(true).Render(
				"💡 Examples:\n"+
				"• \"due Monday at 3pm\"\n"+
				"• \"high priority, needs review from John\"\n"+
				"• \"add location: conference room B\"",
			),
		))
	
	// AI hint
	aiHint := lipgloss.NewStyle().
		Foreground(secondaryColor).
		Italic(true).
		Render("🤖 AI will intelligently merge your context with the existing todo")
	
	// Instructions
	instructions := lipgloss.NewStyle().
		Foreground(mutedColor).
		Render("Enter: Save • Esc: Cancel")
	
	// Build modal content based on processing state
	var modalContent string
	if m.editProcessing {
		// Show processing animation
		spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		frame := spinner[m.loadingFrame%len(spinner)]
		
		processingSteps := []string{
			"🧠 Analyzing your context...",
			"🔍 Understanding intent...",
			"✨ Merging changes intelligently...",
			"💾 Updating todo...",
		}
		step := processingSteps[(m.loadingFrame/10)%len(processingSteps)]
		
		processingContent := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(primaryColor).
			Padding(2, 3).
			Width(modalWidth - 4).
			Render(lipgloss.JoinVertical(
				lipgloss.Center,
				lipgloss.NewStyle().
					Foreground(primaryColor).
					Bold(true).
					Render(fmt.Sprintf("%s Processing Edit", frame)),
				"",
				lipgloss.NewStyle().
					Foreground(secondaryColor).
					Render(step),
				"",
				lipgloss.NewStyle().
					Foreground(mutedColor).
					Render("Please wait..."),
			))
		
		modalContent = lipgloss.JoinVertical(
			lipgloss.Left,
			header,
			"",
			todoInfo,
			"",
			processingContent,
		)
	} else {
		modalContent = lipgloss.JoinVertical(
			lipgloss.Left,
			header,
			"",
			todoInfo,
			"",
			inputSection,
			"",
			aiHint,
			"",
			instructions,
		)
	}
	
	// Create modal with border and shadow effect
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(primaryColor).
		Background(lipgloss.Color("#0a0a0a")).
		Padding(1, 2).
		Width(modalWidth)
	
	modal := modalStyle.Render(modalContent)
	
	// Create an overlay effect
	modalOverlay := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		modal,
	)
	
	// Combine background and modal
	bgLines := strings.Split(background, "\n")
	overlayLines := strings.Split(modalOverlay, "\n")
	
	result := make([]string, len(bgLines))
	for i := range bgLines {
		if i < len(overlayLines) && strings.TrimSpace(overlayLines[i]) != "" {
			result[i] = overlayLines[i]
		} else {
			result[i] = bgLines[i]
		}
	}
	
	return strings.Join(result, "\n")
}

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