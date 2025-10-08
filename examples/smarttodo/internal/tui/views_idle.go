package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/monstercameron/schemaflow/examples/smarttodo/internal/models"
)

func (m Model) idleViewRender() string {
	// Calculate stats
	completedToday := 0
	pendingCount := 0
	urgentCount := 0
	overdueCount := 0
	locationMap := make(map[string]int)
	
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	
	for _, todo := range m.todos {
		if !todo.Completed {
			pendingCount++
			
			// Track locations
			if todo.Location != "" {
				locationMap[todo.Location]++
			}
			
			if todo.Deadline != nil {
				if todo.Deadline.Before(now) {
					overdueCount++
				} else if todo.Deadline.Sub(now) <= time.Hour {
					urgentCount++
				}
			}
		} else if todo.Completed && todo.CreatedAt.After(todayStart) {
			completedToday++
		}
	}
	
	// Create professional idle header
	headerStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(mutedColor).
		Padding(0, 2).
		Width(60)
	
	// Calculate idle duration for display
	idleDuration := time.Since(m.lastActivity)
	idleMinutesDisplay := int(idleDuration.Minutes())
	
	headerContent := lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.NewStyle().Foreground(primaryColor).Bold(true).Render("Smart Todo • Taking a Break"),
		lipgloss.NewStyle().Foreground(mutedColor).Render(fmt.Sprintf("💤 Idle for %d minutes • Press any key to continue", idleMinutesDisplay)),
	)
	
	idleHeader := headerStyle.Render(headerContent)
	
	// Summary stats
	summaryStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(mutedColor).
		Padding(1, 3).
		Width(60)
	
	var summaryLines []string
	
	// Header
	summaryLines = append(summaryLines,
		lipgloss.NewStyle().Foreground(secondaryColor).Bold(true).Render("📊 Task Summary"),
		"",
	)
	
	// Today's progress
	if completedToday > 0 {
		summaryLines = append(summaryLines,
			lipgloss.NewStyle().Foreground(successColor).Render(
				fmt.Sprintf("✅ Completed Today: %d tasks", completedToday)),
		)
	}
	
	// Pending tasks
	if pendingCount > 0 {
		summaryLines = append(summaryLines,
			lipgloss.NewStyle().Foreground(warningColor).Render(
				fmt.Sprintf("⏳ Pending: %d tasks", pendingCount)),
		)
	}
	
	// Urgent/Overdue alerts
	if overdueCount > 0 {
		summaryLines = append(summaryLines,
			lipgloss.NewStyle().Foreground(errorColor).Bold(true).Render(
				fmt.Sprintf("🚨 OVERDUE: %d tasks need attention!", overdueCount)),
		)
	}
	
	if urgentCount > 0 {
		summaryLines = append(summaryLines,
			lipgloss.NewStyle().Foreground(warningColor).Bold(true).Render(
				fmt.Sprintf("⏰ Due Soon: %d tasks within the hour", urgentCount)),
		)
	}
	
	// Location breakdown if available
	if len(locationMap) > 0 {
		summaryLines = append(summaryLines, "")
		summaryLines = append(summaryLines,
			lipgloss.NewStyle().Foreground(secondaryColor).Bold(true).Render("📍 Tasks by Location:"))
		for location, count := range locationMap {
			summaryLines = append(summaryLines,
				lipgloss.NewStyle().Foreground(mutedColor).Render(
					fmt.Sprintf("  • %s: %d tasks", location, count)))
		}
	}
	
	// Time since idle already calculated above
	summaryLines = append(summaryLines,
		"",
		lipgloss.NewStyle().Foreground(mutedColor).Render(
			fmt.Sprintf("Idle for %d minutes", idleMinutesDisplay)),
	)
	
	// Next task suggestion if available
	if pendingCount > 0 {
		var nextTask *models.SmartTodo
		var highPriorityTask *models.SmartTodo
		
		for _, todo := range m.todos {
			if !todo.Completed {
				if nextTask == nil {
					nextTask = todo
				}
				if todo.Priority == "high" && highPriorityTask == nil {
					highPriorityTask = todo
				}
			}
		}
		
		suggestedTask := nextTask
		if highPriorityTask != nil {
			suggestedTask = highPriorityTask
		}
		
		if suggestedTask != nil {
			summaryLines = append(summaryLines,
				"",
				lipgloss.NewStyle().Foreground(primaryColor).Bold(true).Render("💡 Next Suggested Task:"),
				lipgloss.NewStyle().Foreground(successColor).Render(fmt.Sprintf("  → %s", suggestedTask.Title)),
			)
		}
	}
	
	summaryBox := summaryStyle.Render(strings.Join(summaryLines, "\n"))
	
	// Motivation quote - use AI-generated quote if available, otherwise use static quotes
	var quoteText string
	if m.aiQuote != "" {
		quoteText = m.aiQuote
	} else {
		// Fallback quotes if AI quote hasn't been generated yet
		quotes := []string{
			"\"The secret of getting ahead is getting started.\" - Mark Twain",
			"\"Focus on being productive instead of busy.\" - Tim Ferriss",
			"\"You don't have to be great to start, but you have to start to be great.\" - Zig Ziglar",
			"\"A year from now you may wish you had started today.\" - Karen Lamb",
			"\"The way to get started is to quit talking and begin doing.\" - Walt Disney",
		}
		quoteIndex := (m.loadingFrame / 100) % len(quotes)
		quoteText = quotes[quoteIndex]
	}
	
	quote := lipgloss.NewStyle().
		Foreground(mutedColor).
		Italic(true).
		MarginTop(2).
		Render(quoteText)
	
	// Wake instruction
	wakeInstruction := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Border(lipgloss.NormalBorder()).
		BorderForeground(primaryColor).
		Padding(0, 2).
		MarginTop(2).
		Render("Press any key to wake up and continue...")
	
	// Compose the idle view
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		idleHeader,
		summaryBox,
		quote,
		wakeInstruction,
	)
	
	// Center everything
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}