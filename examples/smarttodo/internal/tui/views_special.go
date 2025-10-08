package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) splashViewRender() string {
	// Different splash for new vs returning users
	if m.userName == "" {
		// New user splash
		titleStyle := lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			MarginBottom(2)
		
		welcomeStyle := lipgloss.NewStyle().
			Foreground(secondaryColor).
			MarginBottom(2)
		
		featuresStyle := lipgloss.NewStyle().
			Foreground(mutedColor).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 3)
		
		instructionStyle := lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			MarginTop(2)
		
		content := lipgloss.JoinVertical(
			lipgloss.Center,
			titleStyle.Render("Smart Todo"),
			welcomeStyle.Render("AI-Powered Task Management"),
			"",
			featuresStyle.Render(lipgloss.JoinVertical(
				lipgloss.Left,
				"• Intelligent task processing",
				"• Smart deadline management",
				"• Context-aware suggestions",
				"• Automatic task organization",
			)),
			"",
			instructionStyle.Render("Press Enter to get started"),
		)
		
		// Add decorative border
		modalStyle := lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(primaryColor).
			Padding(2, 4)
		
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			modalStyle.Render(content),
		)
	}
	
	// Returning user splash
	titleStyle := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		MarginBottom(2)
	
	greetingStyle := lipgloss.NewStyle().
		Foreground(secondaryColor).
		MarginBottom(1)
	
	dateStyle := lipgloss.NewStyle().
		Foreground(mutedColor).
		MarginBottom(2)
	
	statsStyle := lipgloss.NewStyle().
		Foreground(successColor).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(successColor).
		Padding(1, 3)
	
	// Calculate stats for today
	completedToday := 0
	pendingCount := 0
	streakDays := 7 // Placeholder
	
	for _, todo := range m.todos {
		if todo.Completed && todo.CreatedAt.Day() == time.Now().Day() {
			completedToday++
		}
		if !todo.Completed {
			pendingCount++
		}
	}
	
	// Due today count
	dueTodayCount := 0
	for _, todo := range m.todos {
		if todo.Deadline != nil && todo.Deadline.Day() == time.Now().Day() && !todo.Completed {
			dueTodayCount++
		}
	}
	
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		titleStyle.Render(fmt.Sprintf("Welcome back, %s", m.userName)),
		greetingStyle.Render(m.listTitle),
		dateStyle.Render(time.Now().Format("Monday, January 2")),
		"",
		statsStyle.Render(lipgloss.JoinVertical(
			lipgloss.Left,
			"📊 Your Progress:",
			fmt.Sprintf("• %d tasks completed today", completedToday),
			fmt.Sprintf("• %d tasks due today", dueTodayCount),
			fmt.Sprintf("• %d tasks pending", pendingCount),
			fmt.Sprintf("• Current streak: %d days 🔥", streakDays),
		)),
		"",
		lipgloss.NewStyle().Foreground(primaryColor).Bold(true).Render("Press Enter to continue"),
	)
	
	// Add decorative border
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(primaryColor).
		Padding(2, 4)
	
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		modalStyle.Render(content),
	)
}

func (m Model) setupViewRender() string {
	title := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		MarginBottom(2).
		Render("Smart Todo Setup")
	
	var prompt string
	if m.userName == "" {
		prompt = "What's your name?"
	} else {
		prompt = fmt.Sprintf("Hi %s! What would you like to call your todo list?", m.userName)
	}
	
	promptStyle := lipgloss.NewStyle().
		Foreground(secondaryColor).
		MarginBottom(1).
		Render(prompt)
	
	inputBox := borderStyle.Render(m.setupInput.View())
	
	helpText := blurredStyle.Render("Press Enter to continue")
	
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		promptStyle,
		inputBox,
		helpText,
	)
	
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func (m Model) closingViewRender() string {
	// Create a beautiful closing animation (frame animation removed for cleaner look)
	
	// Progress bar for closing
	progressBar := m.renderProgressBar(m.closingProgress)
	
	// Goodbye messages
	messages := []string{
		"Saving your progress...",
		"Closing database connections...",
		"Tidying up...",
		"Almost done...",
		"Goodbye! 👋",
	}
	
	messageIdx := m.closingProgress / 20
	if messageIdx >= len(messages) {
		messageIdx = len(messages) - 1
	}
	currentMessage := messages[messageIdx]
	
	// Stats summary
	completedToday := 0
	pendingCount := 0
	for _, todo := range m.todos {
		if todo.Completed {
			completedToday++
		} else {
			pendingCount++
		}
	}
	
	// Build the closing screen
	titleStyle := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		MarginBottom(2)
		
	messageStyle := lipgloss.NewStyle().
		Foreground(secondaryColor).
		MarginBottom(1)
		
	statsStyle := lipgloss.NewStyle().
		Foreground(mutedColor).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(mutedColor).
		Padding(1, 2).
		MarginTop(2).
		MarginBottom(2)
	
	stats := fmt.Sprintf(`Today's Summary:
  ✅ Completed: %d tasks
  ⏳ Remaining: %d tasks
  
  Great work, %s! See you next time!`,
		completedToday,
		pendingCount,
		m.userName,
	)
	
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		titleStyle.Render("Smart Todo"),
		messageStyle.Render(currentMessage),
		progressBar,
		statsStyle.Render(stats),
	)
	
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

func (m Model) statsViewRender() string {
	// Create professional header for stats view
	headerStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(0, 2).
		Width(50)
	
	title := headerStyle.Render(
		lipgloss.NewStyle().Foreground(primaryColor).Bold(true).Render("Smart Todo • Statistics"),
	)
	
	stats := []string{
		focusedStyle.Render("Task Statistics"),
		"",
		fmt.Sprintf("📝 Total Tasks: %s", lipgloss.NewStyle().Foreground(primaryColor).Render(fmt.Sprintf("%d", m.stats["total"]))),
		fmt.Sprintf("✅ Completed: %s", lipgloss.NewStyle().Foreground(successColor).Render(fmt.Sprintf("%d", m.stats["completed"]))),
		fmt.Sprintf("⏳ Pending: %s", lipgloss.NewStyle().Foreground(warningColor).Render(fmt.Sprintf("%d", m.stats["pending"]))),
		"",
		focusedStyle.Render("Priority Breakdown"),
		fmt.Sprintf("🔴 High: %d", m.stats["high"]),
		fmt.Sprintf("🟡 Medium: %d", m.stats["medium"]),
		fmt.Sprintf("🟢 Low: %d", m.stats["low"]),
		"",
	}
	
	if m.stats["overdue"] > 0 {
		stats = append(stats, lipgloss.NewStyle().Foreground(errorColor).Render(fmt.Sprintf("⚠️ OVERDUE: %d tasks", m.stats["overdue"])))
	}
	
	completionRate := 0
	if m.stats["total"] > 0 {
		completionRate = (m.stats["completed"] * 100) / m.stats["total"]
	}
	
	stats = append(stats, "", focusedStyle.Render("Completion Rate"))
	stats = append(stats, m.renderProgressBar(completionRate))
	
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		borderStyle.Render(strings.Join(stats, "\n")),
		"",
		helpStyle.Render("Press Esc to go back"),
	)
	
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func (m Model) suggestViewRender() string {
	if m.selectedTodo == nil {
		return m.listViewRender()
	}
	
	todo := m.selectedTodo
	
	// Create professional header for suggest view
	headerStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(0, 2).
		Width(50)
	
	title := headerStyle.Render(
		lipgloss.NewStyle().Foreground(primaryColor).Bold(true).Render("Smart Todo • AI Suggestion"),
	)
	
	suggestion := fmt.Sprintf(`
%s

%s

%s %s

This task is recommended because:
• %s priority with %s effort required
• %s
• Best suited for your current context

%s
`,
		focusedStyle.Render("Suggested Task:"),
		lipgloss.NewStyle().Foreground(primaryColor).Bold(true).Render(todo.Title),
		getCategoryEmoji(todo.Category),
		todo.Description,
		todo.Priority,
		todo.Effort,
		todo.Category,
		helpStyle.Render("Press Enter to view details, Esc to go back"),
	)
	
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		borderStyle.Render(suggestion),
	)
	
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}