package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/monstercameron/schemaflow/examples/smarttodo/internal/models"
)

// Task view for managing subtasks - modal overlay
func (m Model) taskViewRender() string {
	if m.selectedTodo == nil {
		return m.listViewRender()
	}
	
	// Render the background (main list view)
	background := m.listViewRender()
	
	todo := m.selectedTodo
	
	// Modal dimensions
	modalWidth := 70
	if modalWidth > m.width-10 {
		modalWidth = m.width - 10
	}
	modalHeight := 24  // Increased to accommodate input field
	if modalHeight > m.height-6 {
		modalHeight = m.height - 6
	}
	
	// Create modal header
	header := lipgloss.NewStyle().
		Background(lipgloss.Color("#1a1a2e")).
		Foreground(primaryColor).
		Bold(true).
		Padding(0, 2).
		Width(modalWidth - 2).
		Align(lipgloss.Center).
		Render(fmt.Sprintf("Subtasks • %s", todo.Title))
	
	// Calculate completion and progress
	var progressSection string
	var tasksSection string
	
	if len(todo.Tasks) > 0 {
		completedCount := 0
		for _, task := range todo.Tasks {
			if task.Completed {
				completedCount++
			}
		}
		
		// Progress bar
		progressBar := m.renderProgressBar((completedCount * 100) / len(todo.Tasks))
		progressSection = fmt.Sprintf("%s\n%s\n",
			lipgloss.NewStyle().Foreground(primaryColor).Bold(true).Render(fmt.Sprintf("Progress: %d/%d tasks completed", completedCount, len(todo.Tasks))),
			progressBar,
		)
		
		// Sort tasks: uncompleted first, then completed
		type indexedTask struct {
			task  models.Task
			index int
		}
		
		var uncompletedTasks []indexedTask
		var completedTasks []indexedTask
		
		for i, task := range todo.Tasks {
			if task.Completed {
				completedTasks = append(completedTasks, indexedTask{task: task, index: i})
			} else {
				uncompletedTasks = append(uncompletedTasks, indexedTask{task: task, index: i})
			}
		}
		
		// Combine uncompleted first, then completed
		sortedTasks := append(uncompletedTasks, completedTasks...)
		
		// Task list with scroll window
		maxVisibleTasks := 8
		taskList := []string{}
		
		// Calculate scroll window based on sorted position
		// Find selected task's position in sorted list
		selectedSortedIdx := 0
		for i, t := range sortedTasks {
			if t.index == m.selectedTask {
				selectedSortedIdx = i
				break
			}
		}
		
		startIdx := 0
		endIdx := len(sortedTasks)
		
		if len(sortedTasks) > maxVisibleTasks {
			// Implement scrolling based on sorted position
			if selectedSortedIdx >= maxVisibleTasks/2 {
				startIdx = selectedSortedIdx - maxVisibleTasks/2
				if startIdx+maxVisibleTasks > len(sortedTasks) {
					startIdx = len(sortedTasks) - maxVisibleTasks
				}
			}
			endIdx = startIdx + maxVisibleTasks
			if endIdx > len(sortedTasks) {
				endIdx = len(sortedTasks)
			}
		}
		
		// Show scroll indicators if needed
		if startIdx > 0 {
			taskList = append(taskList, lipgloss.NewStyle().Foreground(mutedColor).Render("  ↑ more above..."))
		}
		
		// Track if we need to show separator
		lastWasUncompleted := false
		
		for i := startIdx; i < endIdx; i++ {
			t := sortedTasks[i]
			checkbox := "☐"
			taskStyle := lipgloss.NewStyle()
			
			if t.task.Completed {
				checkbox = "☑"
				taskStyle = taskStyle.Foreground(mutedColor).Strikethrough(true)
			}
			
			// Add separator before first completed task
			if t.task.Completed && lastWasUncompleted {
				taskList = append(taskList, lipgloss.NewStyle().Foreground(mutedColor).Render("  ─────────────────"))
			}
			
			if !t.task.Completed {
				lastWasUncompleted = true
			} else {
				lastWasUncompleted = false
			}
			
			// Highlight selected task (using original index)
			prefix := "  "
			if t.index == m.selectedTask {
				prefix = "▶ "
				taskStyle = taskStyle.Bold(true).Foreground(primaryColor)
			}
			
			// Make task text full width in modal
			taskText := t.task.Text
			maxTaskWidth := modalWidth - 12 // Account for checkbox, prefix, and padding
			if len(taskText) > maxTaskWidth {
				taskText = taskText[:maxTaskWidth-3] + "..."
			}
			
			taskList = append(taskList, fmt.Sprintf("%s%s %s", prefix, checkbox, taskStyle.Render(taskText)))
		}
		
		if endIdx < len(sortedTasks) {
			taskList = append(taskList, lipgloss.NewStyle().Foreground(mutedColor).Render("  ↓ more below..."))
		}
		
		tasksSection = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(0, 1).
			Width(modalWidth - 6). // Use most of modal width
			Height(maxVisibleTasks + 2).
			Render(strings.Join(taskList, "\n"))
	} else {
		// No tasks yet - show helpful message
		progressSection = lipgloss.NewStyle().
			Foreground(mutedColor).
			Render("No tasks added yet\n")
		
		tasksSection = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(mutedColor).
			Padding(2, 4).
			Width(modalWidth - 6). // Use most of modal width
			Align(lipgloss.Center).
			Render(lipgloss.JoinVertical(
				lipgloss.Center,
				lipgloss.NewStyle().Foreground(mutedColor).Render("📝 No subtasks yet"),
				"",
				lipgloss.NewStyle().Foreground(primaryColor).Bold(true).Render("Press Ctrl+A to add your first task!"),
			))
	}
	
	// Input field section (shown when in input mode)
	var inputSection string
	if m.taskInputMode {
		// Update input width to fit modal
		m.taskInput.Width = modalWidth - 10
		
		inputSection = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(successColor).
			Padding(0, 1).
			Width(modalWidth - 4).
			Render(lipgloss.JoinVertical(
				lipgloss.Left,
				lipgloss.NewStyle().Foreground(successColor).Bold(true).Render("➕ Add New Task:"),
				m.taskInput.View(),
				lipgloss.NewStyle().Foreground(mutedColor).Italic(true).Render("Enter: Add • Esc: Cancel"),
			))
	}
	
	// Compact instructions
	var instructions string
	if m.taskInputMode {
		instructions = lipgloss.NewStyle().
			Foreground(mutedColor).
			Render("Adding tasks... Press Enter to add each task, Esc when done")
	} else {
		instructions = lipgloss.NewStyle().
			Foreground(mutedColor).
			Render("↑/↓ Navigate • Space: Toggle • Ctrl+A: Add • Ctrl+D: Delete • Esc: Close")
	}
	
	// Build modal content
	var modalContent string
	if m.taskInputMode {
		modalContent = lipgloss.JoinVertical(
			lipgloss.Left,
			header,
			progressSection,
			tasksSection,
			"",
			inputSection,
			"",
			instructions,
		)
	} else {
		modalContent = lipgloss.JoinVertical(
			lipgloss.Left,
			header,
			progressSection,
			tasksSection,
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
	
	// Only set max height if not in input mode to ensure input is visible
	if !m.taskInputMode {
		modalStyle = modalStyle.MaxHeight(modalHeight)
	}
	
	modal := modalStyle.Render(modalContent)
	
	// Create an overlay effect by first rendering the background
	// then placing the modal on top
	modalOverlay := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		modal,
	)
	
	// Combine background and modal
	// Split both into lines and overlay them
	bgLines := strings.Split(background, "\n")
	overlayLines := strings.Split(modalOverlay, "\n")
	
	// Simple overlay: modal takes precedence where it has content
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
