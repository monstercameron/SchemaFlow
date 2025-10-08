package tui

import (
	"strings"
	"github.com/charmbracelet/lipgloss"
)

// renderModalOverlay creates a proper modal overlay with dimmed background
func renderModalOverlay(background, modal string, width, height int) string {
	// Create a semi-transparent overlay effect by dimming the background
	dimmedBg := dimBackground(background)
	
	// Split into lines
	bgLines := strings.Split(dimmedBg, "\n")
	modalLines := strings.Split(modal, "\n")
	
	// Calculate modal position (centered)
	modalHeight := len(modalLines)
	modalWidth := 0
	for _, line := range modalLines {
		w := lipgloss.Width(line)
		if w > modalWidth {
			modalWidth = w
		}
	}
	
	// Calculate starting positions
	startY := (height - modalHeight) / 2
	startX := (width - modalWidth) / 2
	
	if startY < 0 {
		startY = 0
	}
	if startX < 0 {
		startX = 0
	}
	
	// Ensure we have enough background lines
	for len(bgLines) < height {
		bgLines = append(bgLines, strings.Repeat(" ", width))
	}
	
	// Overlay modal on background
	for i, modalLine := range modalLines {
		y := startY + i
		if y >= 0 && y < len(bgLines) {
			// Get the background line
			bgLine := bgLines[y]
			bgRunes := []rune(bgLine)
			
			// Pad background line if needed
			for len(bgRunes) < width {
				bgRunes = append(bgRunes, ' ')
			}
			
			// Calculate the actual width of the modal line
			modalRunes := []rune(modalLine)
			modalLineWidth := lipgloss.Width(modalLine)
			
			// Replace the center portion with the modal line
			if startX < len(bgRunes) {
				// Create the new line with modal content
				newLine := make([]rune, 0, width)
				
				// Add left padding (dimmed background)
				if startX > 0 {
					newLine = append(newLine, bgRunes[:startX]...)
				}
				
				// Add modal content
				newLine = append(newLine, modalRunes...)
				
				// Add right padding (dimmed background)
				endX := startX + modalLineWidth
				if endX < len(bgRunes) {
					newLine = append(newLine, bgRunes[endX:]...)
				}
				
				bgLines[y] = string(newLine)
			}
		}
	}
	
	// Trim to height
	if len(bgLines) > height {
		bgLines = bgLines[:height]
	}
	
	return strings.Join(bgLines, "\n")
}

// dimBackground applies a dimming effect to the background
func dimBackground(content string) string {
	// Apply a subtle dimming by using muted colors
	lines := strings.Split(content, "\n")
	dimmedLines := make([]string, len(lines))
	
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#4a4a4a"))
	
	for i, line := range lines {
		// Preserve empty lines
		if strings.TrimSpace(line) == "" {
			dimmedLines[i] = line
		} else {
			// Apply dimming to non-empty lines
			// This is a simple approach - more sophisticated dimming could be added
			dimmedLines[i] = dimStyle.Render(stripANSI(line))
		}
	}
	
	return strings.Join(dimmedLines, "\n")
}

// stripANSI removes ANSI escape codes from a string
func stripANSI(str string) string {
	// Simple ANSI stripping - removes color codes
	result := ""
	inEscape := false
	
	for _, r := range str {
		if r == '\x1b' {
			inEscape = true
		} else if inEscape {
			if r == 'm' {
				inEscape = false
			}
		} else {
			result += string(r)
		}
	}
	
	return result
}

// createModalBox creates a consistent modal box with proper styling
func createModalBox(title, content string, width int, borderColor lipgloss.Color) string {
	// Create header
	headerStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#1a1a2e")).
		Foreground(borderColor).
		Bold(true).
		Padding(0, 2).
		Width(width - 4).
		Align(lipgloss.Center)
	
	header := headerStyle.Render(title)
	
	// Create modal box
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(borderColor).
		Background(lipgloss.Color("#0a0a0a")).
		Padding(1, 2).
		Width(width).
		MaxWidth(width)
	
	// Combine header and content
	fullContent := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		content,
	)
	
	return modalStyle.Render(fullContent)
}