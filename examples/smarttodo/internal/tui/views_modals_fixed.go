package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) addViewRenderFixed() string {
	// Render the background (main list view)
	background := m.listViewRender()
	
	// Modal dimensions with safe bounds
	modalWidth := 60
	if modalWidth > m.width-10 {
		modalWidth = m.width - 10
	}
	if modalWidth < 40 {
		modalWidth = 40
	}
	
	// Ensure input is focused
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
		inputSection,
		"",
		aiHint,
		"",
		instructions,
	)
	
	// Create modal using the new modal box helper
	modal := createModalBox("New Task", modalContent, modalWidth, primaryColor)
	
	// Place modal centered
	centeredModal := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		modal,
	)
	
	// Use the new overlay function for proper rendering
	return renderModalOverlay(background, centeredModal, m.width, m.height)
}

func (m Model) editViewRenderFixed() string {
	if m.selectedTodo == nil {
		return m.listViewRender()
	}
	
	// Render the background (main list view)
	background := m.listViewRender()
	
	// Modal dimensions with safe bounds
	modalWidth := 70
	if modalWidth > m.width-10 {
		modalWidth = m.width - 10
	}
	if modalWidth < 50 {
		modalWidth = 50
	}
	
	// Current todo info
	currentInfo := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(mutedColor).
		Padding(0, 1).
		Width(modalWidth - 4).
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.NewStyle().Foreground(mutedColor).Bold(true).Render("Current Todo:"),
			truncateText(m.selectedTodo.Title, modalWidth-8),
			"",
			lipgloss.NewStyle().Foreground(mutedColor).Render(fmt.Sprintf("Priority: %s | Category: %s", 
				m.selectedTodo.Priority, m.selectedTodo.Category)),
		))
	
	// Edit input
	if !m.editInput.Focused() {
		m.editInput.Focus()
	}
	m.editInput.Width = modalWidth - 6
	
	// Processing indicator
	var processingIndicator string
	if m.editProcessing {
		frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		frame := frames[m.loadingFrame%len(frames)]
		
		steps := []string{
			"Analyzing changes...",
			"Understanding intent...",
			"Merging with existing data...",
			"Updating todo...",
		}
		
		currentStep := steps[(m.loadingFrame/10)%len(steps)]
		processingIndicator = lipgloss.NewStyle().
			Foreground(warningColor).
			Bold(true).
			Render(fmt.Sprintf("%s %s", frame, currentStep))
	}
	
	// Input section
	inputSection := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(secondaryColor).
		Padding(0, 1).
		Width(modalWidth - 4).
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.NewStyle().Foreground(secondaryColor).Bold(true).Render("✏️ Add context or updates:"),
			lipgloss.NewStyle().Foreground(mutedColor).Italic(true).Render("AI will merge with existing todo"),
			m.editInput.View(),
		))
	
	// Instructions
	instructions := lipgloss.NewStyle().
		Foreground(mutedColor).
		Render("Enter: Save • Esc: Cancel")
	
	// Build modal content
	var modalContent string
	if m.editProcessing {
		modalContent = lipgloss.JoinVertical(
			lipgloss.Left,
			currentInfo,
			"",
			processingIndicator,
			"",
			instructions,
		)
	} else {
		modalContent = lipgloss.JoinVertical(
			lipgloss.Left,
			currentInfo,
			"",
			inputSection,
			"",
			instructions,
		)
	}
	
	// Create modal
	modal := createModalBox("Edit Task", modalContent, modalWidth, secondaryColor)
	
	// Place modal centered
	centeredModal := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		modal,
	)
	
	// Use the new overlay function
	return renderModalOverlay(background, centeredModal, m.width, m.height)
}

func (m Model) suggestViewRenderFixed() string {
	// Render the background
	background := m.listViewRender()
	
	// Modal dimensions
	modalWidth := 60
	if modalWidth > m.width-10 {
		modalWidth = m.width - 10
	}
	if modalWidth < 40 {
		modalWidth = 40
	}
	
	var content string
	
	if m.loading {
		// Loading animation
		frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		frame := frames[m.loadingFrame%len(frames)]
		
		loadingText := lipgloss.NewStyle().
			Foreground(warningColor).
			Bold(true).
			Padding(2, 0).
			Align(lipgloss.Center).
			Render(fmt.Sprintf("%s Analyzing your tasks...\n\nFinding the best next action...", frame))
		
		content = loadingText
	} else if m.selectedTodo != nil && m.mode == suggestView {
		// Show suggestion
		suggestionBox := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(successColor).
			Padding(0, 1).
			Width(modalWidth - 4).
			Render(lipgloss.JoinVertical(
				lipgloss.Left,
				lipgloss.NewStyle().Foreground(successColor).Bold(true).Render("💡 Suggested Next Task:"),
				"",
				truncateText(m.selectedTodo.Title, modalWidth-8),
				"",
				lipgloss.NewStyle().Foreground(mutedColor).Render(
					fmt.Sprintf("Priority: %s | Category: %s", 
						m.selectedTodo.Priority, m.selectedTodo.Category)),
			))
		
		instructions := lipgloss.NewStyle().
			Foreground(mutedColor).
			Render("Enter: Start Task • Esc: Dismiss")
		
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			suggestionBox,
			"",
			instructions,
		)
	} else {
		content = lipgloss.NewStyle().
			Foreground(mutedColor).
			Padding(2, 0).
			Align(lipgloss.Center).
			Render("No suggestions available")
	}
	
	// Create modal
	modal := createModalBox("AI Suggestion", content, modalWidth, primaryColor)
	
	// Place modal centered
	centeredModal := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		modal,
	)
	
	return renderModalOverlay(background, centeredModal, m.width, m.height)
}
