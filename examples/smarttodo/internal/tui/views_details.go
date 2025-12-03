package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) detailViewRender() string {
	if m.selectedTodo == nil {
		return m.listViewRender()
	}

	todo := m.selectedTodo
	// Create professional header for detail view
	headerStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(0, 2).
		Width(60)
	
	title := headerStyle.Render(
		lipgloss.NewStyle().Foreground(primaryColor).Bold(true).Render(fmt.Sprintf("Task Details • %s", todo.Title)),
	)

	// Build details
	details := []string{}

	// Description
	if todo.Description != "" {
		details = append(details, fmt.Sprintf("%s\n%s", 
			focusedStyle.Render("Description:"),
			todo.Description))
	}
	
	// Tasks/Subtasks
	if len(todo.Tasks) > 0 {
		completedCount := 0
		for _, task := range todo.Tasks {
			if task.Completed {
				completedCount++
			}
		}
		
		progressBar := m.renderProgressBar((completedCount * 100) / len(todo.Tasks))
		tasksStr := fmt.Sprintf("%s (%d/%d)\n%s\n",
			focusedStyle.Render("Tasks:"),
			completedCount,
			len(todo.Tasks),
			progressBar,
		)
		
		for i, task := range todo.Tasks {
			checkbox := "☐"
			taskStyle := lipgloss.NewStyle()
			if task.Completed {
				checkbox = "☑"
				taskStyle = taskStyle.Foreground(mutedColor).Strikethrough(true)
			}
			tasksStr += fmt.Sprintf("  %d. %s %s\n", i+1, checkbox, taskStyle.Render(task.Text))
		}
		
		details = append(details, tasksStr)
	}

	// Priority
	var priorityStr string
	switch todo.Priority {
	case "high":
		priorityStr = priorityHighStyle.Render("HIGH PRIORITY")
	case "medium":
		priorityStr = priorityMediumStyle.Render("MEDIUM PRIORITY")
	case "low":
		priorityStr = priorityLowStyle.Render("LOW PRIORITY")
	}
	details = append(details, fmt.Sprintf("%s %s",
		focusedStyle.Render("Priority:"),
		priorityStr))

	// Category
	details = append(details, fmt.Sprintf("%s %s %s",
		focusedStyle.Render("Category:"),
		getCategoryEmoji(todo.Category),
		todo.Category))

	// Effort
	details = append(details, fmt.Sprintf("%s %s %s",
		focusedStyle.Render("Effort:"),
		getEffortEmoji(todo.Effort),
		todo.Effort))

	// Deadline
	if todo.Deadline != nil {
		daysUntil := int(time.Until(*todo.Deadline).Hours() / 24)
		deadlineStr := todo.Deadline.Format("Monday, January 2, 2006")
		if daysUntil < 0 {
			deadlineStr = lipgloss.NewStyle().Foreground(errorColor).Render(fmt.Sprintf("%s (OVERDUE by %d days)", deadlineStr, -daysUntil))
		} else if daysUntil == 0 {
			deadlineStr = lipgloss.NewStyle().Foreground(warningColor).Render(fmt.Sprintf("%s (TODAY!)", deadlineStr))
		} else if daysUntil == 1 {
			deadlineStr = lipgloss.NewStyle().Foreground(warningColor).Render(fmt.Sprintf("%s (Tomorrow)", deadlineStr))
		} else {
			deadlineStr = fmt.Sprintf("%s (%d days)", deadlineStr, daysUntil)
		}
		details = append(details, fmt.Sprintf("%s %s",
			focusedStyle.Render("Deadline:"),
			deadlineStr))
	}

	// Dependencies
	if len(todo.Dependencies) > 0 {
		details = append(details, fmt.Sprintf("%s\n  • %s",
			focusedStyle.Render("Dependencies:"),
			strings.Join(todo.Dependencies, "\n  • ")))
	}

	// Context advice
	if todo.Context != "" && todo.Context != "No specific context required" {
		details = append(details, fmt.Sprintf("%s\n%s",
			focusedStyle.Render("💡 Best Context:"),
			lipgloss.NewStyle().Foreground(secondaryColor).Render(todo.Context)))
	}

	// Created
	details = append(details, fmt.Sprintf("%s %s",
		focusedStyle.Render("Created:"),
		todo.CreatedAt.Format("Jan 2, 2006 at 3:04 PM")))

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		borderStyle.Render(strings.Join(details, "\n\n")),
		"",
		helpStyle.Render("Press Esc or Enter to go back"),
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}
