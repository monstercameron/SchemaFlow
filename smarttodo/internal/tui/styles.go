package tui

import "github.com/charmbracelet/lipgloss"

// Styles - all styling definitions in one place
var (
	// Colors
	primaryColor   = lipgloss.Color("#00D9FF")
	secondaryColor = lipgloss.Color("#FF00D9")
	successColor   = lipgloss.Color("#00FF88")
	warningColor   = lipgloss.Color("#FFD700")
	errorColor     = lipgloss.Color("#FF4444")
	mutedColor     = lipgloss.Color("#666666")

	// Base styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			Background(lipgloss.Color("#1a1a1a")).
			Padding(0, 1)

	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25283D"))

	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor)

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor)

	priorityHighStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(errorColor)

	priorityMediumStyle = lipgloss.NewStyle().
				Foreground(warningColor)

	priorityLowStyle = lipgloss.NewStyle().
				Foreground(successColor)

	deadlineStyle = lipgloss.NewStyle().
			Foreground(secondaryColor)

	categoryStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2)

	inputStyle = lipgloss.NewStyle().
			Foreground(primaryColor)

	focusedStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	blurredStyle = lipgloss.NewStyle().
			Foreground(mutedColor)
)

func getCategoryEmoji(category string) string {
	emojis := map[string]string{
		"work":     "💼",
		"personal": "🏠",
		"health":   "🏃",
		"urgent":   "🚨",
		"learning": "📚",
		"shopping": "🛒",
		"social":   "👥",
		"finance":  "💰",
	}
	if emoji, ok := emojis[category]; ok {
		return emoji
	}
	return "📌"
}

func getEffortEmoji(effort string) string {
	emojis := map[string]string{
		"minimal": "⚡",
		"low":     "🔹",
		"medium":  "🔸",
		"high":    "🔶",
		"massive": "🔥",
	}
	if emoji, ok := emojis[effort]; ok {
		return emoji
	}
	return "🔸"
}