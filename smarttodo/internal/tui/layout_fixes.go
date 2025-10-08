package tui

import (
	"github.com/charmbracelet/lipgloss"
	"strings"
)

// Layout constants to prevent rendering issues
const (
	// Minimum dimensions to prevent crashes
	MinWidth  = 80
	MinHeight = 24
	
	// Fixed component heights for consistent layout
	HeaderHeight       = 4  // Header with border
	StatsBarHeight     = 2  // Stats bar with margin
	BottomPanelHeight  = 9  // Keyboard shortcuts and activity log
	StatusLineHeight   = 1  // Status message line
	SpacingHeight      = 3  // Total spacing between components
	
	// Margins and padding
	ListHorizontalPadding = 4
	ModalMinMargin        = 10
	
	// Maximum widths to prevent overflow
	MaxTodoTitleWidth     = 80
	MaxTaskTextWidth      = 70
	MaxProgressBarWidth   = 40
)

// safeWidth returns a width that won't cause rendering issues
func (m Model) safeWidth() int {
	if m.width < MinWidth {
		return MinWidth
	}
	return m.width
}

// safeHeight returns a height that won't cause rendering issues
func (m Model) safeHeight() int {
	if m.height < MinHeight {
		return MinHeight
	}
	return m.height
}

// calculateListHeight computes the available height for the todo list
func (m Model) calculateListHeight() int {
	totalFixedHeight := HeaderHeight + StatsBarHeight + BottomPanelHeight + StatusLineHeight + SpacingHeight
	availableHeight := m.safeHeight() - totalFixedHeight
	
	// Ensure minimum list height
	if availableHeight < 8 {
		return 8
	}
	
	// Cap maximum to prevent excessive scrolling
	if availableHeight > 50 {
		return 50
	}
	
	return availableHeight
}

// calculateModalDimensions returns safe dimensions for modal dialogs
func (m Model) calculateModalDimensions(preferredWidth, preferredHeight int) (width, height int) {
	safeW := m.safeWidth()
	safeH := m.safeHeight()
	
	// Calculate maximum modal dimensions
	maxWidth := safeW - ModalMinMargin
	maxHeight := safeH - 6
	
	// Use preferred dimensions but cap at maximums
	width = preferredWidth
	if width > maxWidth {
		width = maxWidth
	}
	if width < 40 {
		width = 40
	}
	
	height = preferredHeight
	if height > maxHeight {
		height = maxHeight
	}
	if height < 10 {
		height = 10
	}
	
	return width, height
}

// truncateText safely truncates text to fit within width
func truncateText(text string, maxWidth int) string {
	if lipgloss.Width(text) <= maxWidth {
		return text
	}
	
	// Account for ellipsis
	if maxWidth < 4 {
		return "..."
	}
	
	// Binary search for the right truncation point
	runes := []rune(text)
	left, right := 0, len(runes)
	result := ""
	
	for left < right {
		mid := (left + right + 1) / 2
		candidate := string(runes[:mid]) + "..."
		if lipgloss.Width(candidate) <= maxWidth {
			result = candidate
			left = mid
		} else {
			right = mid - 1
		}
	}
	
	if result == "" && maxWidth >= 3 {
		return "..."
	}
	
	return result
}

// renderSafeBox creates a box with safe dimensions
func renderSafeBox(content string, width, height int, borderStyle lipgloss.Border, borderColor lipgloss.Color) string {
	// Ensure minimum dimensions
	if width < 10 {
		width = 10
	}
	if height < 3 {
		height = 3
	}
	
	style := lipgloss.NewStyle().
		Border(borderStyle).
		BorderForeground(borderColor).
		Width(width).
		MaxWidth(width)
	
	// Only set height if content might overflow
	contentLines := strings.Count(content, "\n") + 1
	if contentLines > height-2 { // Account for borders
		style = style.Height(height).MaxHeight(height)
	}
	
	return style.Render(content)
}

// calculateSplitWidth returns widths for side-by-side panels
func (m Model) calculateSplitWidth(gapSize int) (leftWidth, rightWidth int) {
	availableWidth := m.safeWidth() - gapSize
	
	// Equal split
	leftWidth = availableWidth / 2
	rightWidth = availableWidth - leftWidth
	
	// Ensure minimum widths
	minPanelWidth := 30
	if leftWidth < minPanelWidth {
		leftWidth = minPanelWidth
	}
	if rightWidth < minPanelWidth {
		rightWidth = minPanelWidth
	}
	
	return leftWidth, rightWidth
}

// renderSafeProgressBar creates a progress bar that fits the available width
func (m Model) renderSafeProgressBar(progress int) string {
	// Calculate available width for the bar
	availableWidth := m.safeWidth() - 30 // Leave room for text
	if availableWidth > MaxProgressBarWidth {
		availableWidth = MaxProgressBarWidth
	}
	if availableWidth < 10 {
		availableWidth = 10
	}
	
	filled := (progress * availableWidth) / 100
	if filled < 0 {
		filled = 0
	}
	if filled > availableWidth {
		filled = availableWidth
	}
	
	bar := strings.Repeat("█", filled) + strings.Repeat("░", availableWidth-filled)
	
	return lipgloss.NewStyle().Foreground(successColor).Render(bar)
}

// ensureVisibleArea makes sure content is visible in the viewport
func ensureVisibleArea(content string, width, height int) string {
	lines := strings.Split(content, "\n")
	
	// Truncate lines that are too long
	for i, line := range lines {
		if lipgloss.Width(line) > width {
			lines[i] = truncateText(line, width)
		}
	}
	
	// Limit number of lines
	if len(lines) > height {
		lines = lines[:height]
	}
	
	return strings.Join(lines, "\n")
}

// wrapText wraps text to fit within the specified width
func wrapText(text string, width int) string {
	if width < 10 {
		width = 10
	}
	
	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}
	
	var lines []string
	currentLine := ""
	
	for _, word := range words {
		testLine := currentLine
		if testLine != "" {
			testLine += " "
		}
		testLine += word
		
		if lipgloss.Width(testLine) > width {
			if currentLine != "" {
				lines = append(lines, currentLine)
				currentLine = word
			} else {
				// Word is too long, truncate it
				lines = append(lines, truncateText(word, width))
				currentLine = ""
			}
		} else {
			currentLine = testLine
		}
	}
	
	if currentLine != "" {
		lines = append(lines, currentLine)
	}
	
	return strings.Join(lines, "\n")
}