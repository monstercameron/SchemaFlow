package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/monstercameron/schemaflow/examples/smarttodo/internal/models"
)

// List item for todo display
type todoItem struct {
	todo *models.SmartTodo
}

func (t todoItem) Title() string {
	// Format: [checkbox] [urgency] Title
	checkbox := "☐"
	if t.todo.Completed {
		checkbox = "☑"
	}
	
	// Urgency indicator
	var urgency string
	switch t.todo.Priority {
	case "high":
		urgency = lipgloss.NewStyle().Foreground(errorColor).Bold(true).Render("🔴")
	case "medium":
		urgency = lipgloss.NewStyle().Foreground(warningColor).Render("🟡")
	case "low":
		urgency = lipgloss.NewStyle().Foreground(successColor).Render("🟢")
	default:
		urgency = "⚪"
	}
	
	return fmt.Sprintf("%s %s %s", checkbox, urgency, t.todo.Title)
}

func (t todoItem) Description() string {
	var lines []string
	
	// First line: metadata (category, effort, deadline, dependencies)
	details := []string{}
	
	// Category
	if t.todo.Category != "" && t.todo.Category != "pending" {
		icon := getCategoryEmoji(t.todo.Category)
		details = append(details, fmt.Sprintf("%s %s", icon, t.todo.Category))
	}
	
	// Effort
	if t.todo.Effort != "" {
		effortEmoji := getEffortEmoji(t.todo.Effort)
		details = append(details, fmt.Sprintf("%s %s", effortEmoji, t.todo.Effort))
	}
	
	// Deadline
	if t.todo.Deadline != nil {
		daysUntil := int(time.Until(*t.todo.Deadline).Hours() / 24)
		var deadlineStr string
		if daysUntil < 0 {
			deadlineStr = lipgloss.NewStyle().Foreground(errorColor).Bold(true).Render(fmt.Sprintf("📅 OVERDUE %dd!", -daysUntil))
		} else if daysUntil == 0 {
			deadlineStr = lipgloss.NewStyle().Foreground(warningColor).Bold(true).Render("📅 TODAY!")
		} else if daysUntil == 1 {
			deadlineStr = lipgloss.NewStyle().Foreground(warningColor).Render("📅 Tomorrow")
		} else if daysUntil <= 7 {
			deadlineStr = fmt.Sprintf("📅 %dd", daysUntil)
		} else {
			deadlineStr = fmt.Sprintf("📅 %s", t.todo.Deadline.Format("Jan 2"))
		}
		details = append(details, deadlineStr)
	}
	
	// Dependencies
	if len(t.todo.Dependencies) > 0 {
		details = append(details, fmt.Sprintf("🔗 %d deps", len(t.todo.Dependencies)))
	}
	
	// Add metadata line if we have any details
	if len(details) > 0 {
		lines = append(lines, " "+strings.Join(details, " • "))
	}
	
	// Tasks/Subtasks - show each on its own line
	if len(t.todo.Tasks) > 0 {
		// Show completion progress
		completedCount := 0
		var uncompletedTasks []models.Task
		var completedTasks []models.Task
		
		for _, task := range t.todo.Tasks {
			if task.Completed {
				completedCount++
				completedTasks = append(completedTasks, task)
			} else {
				uncompletedTasks = append(uncompletedTasks, task)
			}
		}
		
		// Combine uncompleted first, then completed
		sortedTasks := append(uncompletedTasks, completedTasks...)
		
		// Add progress indicator
		progressStr := fmt.Sprintf(" [%d/%d]", completedCount, len(t.todo.Tasks))
		if completedCount == len(t.todo.Tasks) {
			progressStr = lipgloss.NewStyle().Foreground(successColor).Render(progressStr + " ✓")
		} else if completedCount > 0 {
			progressStr = lipgloss.NewStyle().Foreground(warningColor).Render(progressStr)
		} else {
			progressStr = lipgloss.NewStyle().Foreground(mutedColor).Render(progressStr)
		}
		
		lines = append(lines, fmt.Sprintf(" Tasks%s:", progressStr))
		
		// Show up to 3 tasks (sorted order)
		showCount := len(sortedTasks)
		if showCount > 3 {
			showCount = 3
		}
		
		for i := 0; i < showCount; i++ {
			task := sortedTasks[i]
			checkbox := "☐"
			if task.Completed {
				checkbox = "☑"
			}
			
			taskText := task.Text
			// Allow much wider text (60-80% of typical terminal width)
			if len(taskText) > 70 {
				taskText = taskText[:67] + "..."
			}
			
			if task.Completed {
				// Use ANSI codes for proper strikethrough
				taskText = formatCompletedText(taskText, true)
				taskText = lipgloss.NewStyle().Foreground(mutedColor).Render(taskText)
			}
			
			lines = append(lines, fmt.Sprintf("   %s %s", checkbox, taskText))
		}
		
		if len(sortedTasks) > 3 {
			lines = append(lines, fmt.Sprintf("   ... +%d more", len(sortedTasks)-3))
		}
	}
	
	if len(lines) == 0 {
		return ""
	}
	
	return strings.Join(lines, "\n")
}

func (t todoItem) FilterValue() string {
	return t.todo.Title
}
