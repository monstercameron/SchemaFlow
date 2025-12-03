package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/monstercameron/schemaflow/examples/smarttodo/internal/models"
)

// Calendar view shows todos in a daily time-based view
func (m Model) calendarViewRender() string {
	// Create time slots from 6 AM to 11 PM (30-minute intervals)
	startHour := 6
	endHour := 23
	slotDuration := 30 // minutes

	// Create slots
	type timeSlot struct {
		time  time.Time
		todos []*models.SmartTodo
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	
	// Initialize slots
	var slots []timeSlot
	for hour := startHour; hour <= endHour; hour++ {
		for minute := 0; minute < 60; minute += slotDuration {
			slotTime := today.Add(time.Duration(hour)*time.Hour + time.Duration(minute)*time.Minute)
			slots = append(slots, timeSlot{time: slotTime, todos: []*models.SmartTodo{}})
		}
	}

	// Place todos in appropriate slots based on deadline
	for _, todo := range m.todos {
		if todo.Completed {
			continue // Skip completed todos
		}
		
		if todo.Deadline != nil {
			// Check if deadline is today
			if todo.Deadline.Year() == today.Year() && 
			   todo.Deadline.Month() == today.Month() && 
			   todo.Deadline.Day() == today.Day() {
				// Find the appropriate slot
				for i := range slots {
					if i < len(slots)-1 {
						if todo.Deadline.After(slots[i].time) && todo.Deadline.Before(slots[i+1].time) {
							slots[i].todos = append(slots[i].todos, todo)
							break
						}
					} else {
						// Last slot
						if todo.Deadline.After(slots[i].time) {
							slots[i].todos = append(slots[i].todos, todo)
						}
					}
				}
			}
		}
	}

	// Build the calendar view
	headerStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(0, 2).
		Width(80)

	titleContent := lipgloss.JoinHorizontal(
		lipgloss.Left,
		lipgloss.NewStyle().Foreground(primaryColor).Bold(true).Render("📅 Daily Calendar"),
		lipgloss.NewStyle().Foreground(mutedColor).Render(" • "),
		lipgloss.NewStyle().Foreground(secondaryColor).Render(now.Format("Monday, January 2, 2006")),
	)

	header := headerStyle.Render(titleContent)

	// Build slot display
	var calendarLines []string
	currentSlotIndex := -1

	for i, slot := range slots {
		// Check if this is the current time slot
		if now.After(slot.time) && (i == len(slots)-1 || now.Before(slots[i+1].time)) {
			currentSlotIndex = i
		}

		timeStr := slot.time.Format("3:04 PM")
		
		// Style based on current time
		var timeStyle lipgloss.Style
		if i == currentSlotIndex {
			timeStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true).
				Background(lipgloss.Color("#1a1a2e"))
		} else if slot.time.Before(now) {
			timeStyle = lipgloss.NewStyle().Foreground(mutedColor)
		} else {
			timeStyle = lipgloss.NewStyle().Foreground(secondaryColor)
		}

		line := timeStyle.Render(fmt.Sprintf("%-8s", timeStr))

		// Add todos for this slot
		if len(slot.todos) > 0 {
			for _, todo := range slot.todos {
				// Color code by priority
				var todoStyle lipgloss.Style
				switch todo.Priority {
				case "high":
					todoStyle = lipgloss.NewStyle().Foreground(errorColor)
				case "medium":
					todoStyle = lipgloss.NewStyle().Foreground(warningColor)
				default:
					todoStyle = lipgloss.NewStyle().Foreground(successColor)
				}

				// Truncate title if too long
				title := todo.Title
				if len(title) > 50 {
					title = title[:47] + "..."
				}

				line += " │ " + todoStyle.Render(title)
				
				// Add category icon
				if todo.Category != "" {
					line += " " + lipgloss.NewStyle().Foreground(mutedColor).Render(fmt.Sprintf("[%s]", todo.Category))
				}
			}
		} else if i == currentSlotIndex {
			line += " │ " + lipgloss.NewStyle().Foreground(primaryColor).Italic(true).Render("← Now")
		} else {
			line += " │ " + lipgloss.NewStyle().Foreground(mutedColor).Render("—")
		}

		calendarLines = append(calendarLines, line)

		// Add separator every hour
		if slot.time.Minute() == 30 && i < len(slots)-1 {
			calendarLines = append(calendarLines, lipgloss.NewStyle().Foreground(mutedColor).Render("        ├" + strings.Repeat("─", 70)))
		}
	}

	// Create scrollable content (show 20 slots at a time)
	visibleSlots := 20
	startIndex := 0
	
	// Auto-scroll to current time
	if currentSlotIndex >= 0 {
		startIndex = currentSlotIndex - visibleSlots/2
		if startIndex < 0 {
			startIndex = 0
		}
		if startIndex+visibleSlots > len(calendarLines) {
			startIndex = len(calendarLines) - visibleSlots
			if startIndex < 0 {
				startIndex = 0
			}
		}
	}

	endIndex := startIndex + visibleSlots
	if endIndex > len(calendarLines) {
		endIndex = len(calendarLines)
	}

	visibleLines := calendarLines[startIndex:endIndex]
	
	calendarContent := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(mutedColor).
		Padding(1, 2).
		Width(80).
		Render(strings.Join(visibleLines, "\n"))

	// Statistics
	todayCount := 0
	overdueCount := 0
	for _, todo := range m.todos {
		if !todo.Completed && todo.Deadline != nil {
			if todo.Deadline.Year() == today.Year() && 
			   todo.Deadline.Month() == today.Month() && 
			   todo.Deadline.Day() == today.Day() {
				todayCount++
			}
			if todo.Deadline.Before(now) {
				overdueCount++
			}
		}
	}

	statsLine := lipgloss.NewStyle().
		Foreground(mutedColor).
		Render(fmt.Sprintf("Tasks Today: %d | Overdue: %d | Time: %s", 
			todayCount, overdueCount, now.Format("3:04 PM")))

	// Instructions
	instructions := lipgloss.NewStyle().
		Foreground(primaryColor).
		Border(lipgloss.NormalBorder()).
		BorderForeground(primaryColor).
		Padding(0, 2).
		Render("↑/↓: Scroll • Enter: View Todo • Esc: Back to List • Ctrl+P: Exit Calendar")

	// Compose the view
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		calendarContent,
		"",
		statsLine,
		"",
		instructions,
	)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}
