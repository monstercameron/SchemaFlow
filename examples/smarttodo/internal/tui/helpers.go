package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/monstercameron/schemaflow/examples/smarttodo/internal/models"
)

// Helper functions for views

// formatCompletedText applies strikethrough effect to completed items
// Uses ANSI escape codes for better terminal compatibility
func formatCompletedText(text string, completed bool) string {
	if !completed {
		return text
	}

	// Use ANSI escape codes for strikethrough
	// \x1b[9m starts strikethrough, \x1b[0m resets
	return fmt.Sprintf("\x1b[9m%s\x1b[0m", text)
}

func (m *Model) addLog(msg string) {
	timestamp := time.Now().Format("15:04:05")
	logEntry := fmt.Sprintf("[%s] %s", timestamp, msg)
	m.consoleLogs = append(m.consoleLogs, logEntry)

	// Keep only the last N logs
	if len(m.consoleLogs) > m.maxLogs {
		m.consoleLogs = m.consoleLogs[len(m.consoleLogs)-m.maxLogs:]
	}
}

func (m Model) renderProgressBar(percentage int) string {
	width := 30
	filled := (percentage * width) / 100
	empty := width - filled

	bar := lipgloss.NewStyle().Foreground(successColor).Render(strings.Repeat("█", filled)) +
		lipgloss.NewStyle().Foreground(mutedColor).Render(strings.Repeat("░", empty))

	return fmt.Sprintf("%s %d%%", bar, percentage)
}

func (m Model) renderEnhancedStatsBar() string {
	if m.stats == nil {
		return ""
	}

	// Calculate today's todos (not tasks)
	todayCompleted := 0
	todayTotal := 0
	urgentCount := 0 // Due within 1 hour

	for _, todo := range m.todos {
		// Count all todos (not just today's)
		todayTotal++
		if todo.Completed {
			todayCompleted++
		}

		// Count urgent
		if todo.Deadline != nil && !todo.Completed {
			timeUntil := time.Until(*todo.Deadline)
			if timeUntil > 0 && timeUntil <= time.Hour {
				urgentCount++
			}
		}
	}

	// Calculate todo completion progress (not daily goal)
	todoProgress := 0
	if todayTotal > 0 {
		todoProgress = (todayCompleted * 100) / todayTotal
	}

	// Get location (default to home)
	location := "Home"
	if m.selectedTodo != nil && m.selectedTodo.Context != "" {
		location = m.selectedTodo.Context
	}

	// Calculate streak (placeholder - would need persistent storage)
	streakDays := 3

	// Build single-line progress bar
	progressBar := m.renderProgressBar(todoProgress)

	// Create single-line stats display
	statsContent := fmt.Sprintf(
		"Today: %d/%d done │ ",
		todayCompleted, todayTotal,
	)

	if urgentCount > 0 {
		statsContent += fmt.Sprintf("🔥 %d urgent │ ", urgentCount)
	}

	statsContent += fmt.Sprintf("📍 %s │ ", location)
	statsContent += fmt.Sprintf("🔥 %d day streak │ ", streakDays)
	statsContent += progressBar
	statsContent += fmt.Sprintf(" %d%%", todoProgress)

	// Create single-line display with full width
	return lipgloss.NewStyle().
		Foreground(mutedColor).
		Width(m.width).
		Padding(0, 2).
		MarginBottom(1).
		Render(statsContent)
}

func tickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return models.TickMsg(t)
	})
}

func closingTickCmd() tea.Cmd {
	return tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
		return models.ClosingTickMsg{}
	})
}

func splashTimerCmd() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return models.SplashDismissMsg{}
	})
}

func idleCheckCmd() tea.Cmd {
	return tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
		return models.IdleCheckMsg{}
	})
}
