package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/monstercameron/schemaflow/examples/smarttodo/internal/localization"
)

func (m Model) listViewRender() string {
	// Build professional header with user info
	// Create header sections
	leftSection := localization.T(localization.AppName)
	if m.userName != "" {
		leftSection = fmt.Sprintf("%s • %s", localization.T(localization.AppName), m.userName)
	}
	
	// Add filter indicator if active
	if m.list.SettingFilter() {
		filterIndicator := lipgloss.NewStyle().
			Foreground(warningColor).
			Bold(true).
			Render(" 🔍 FILTER MODE")
		leftSection += filterIndicator
	}
	
	// Add today's date
	today := time.Now().Format("Monday, January 2, 2006")
	rightSection := today
	
	
	// Create professional bordered header
	// Top line of header
	// Ensure header width is safe
	headerWidth := m.width - 2
	if headerWidth < 40 {
		headerWidth = 40
	}
	headerTop := "╭" + strings.Repeat("─", headerWidth) + "╮"
	
	// Main header content line
	leftPart := lipgloss.NewStyle().Foreground(primaryColor).Bold(true).Render(leftSection)
	rightPart := lipgloss.NewStyle().Foreground(secondaryColor).Render(rightSection)
	
	// Calculate spacing safely
	leftWidth := lipgloss.Width(leftPart)
	rightWidth := lipgloss.Width(rightPart)
	spacerWidth := headerWidth - leftWidth - rightWidth - 2 // Account for side borders
	if spacerWidth < 1 {
		spacerWidth = 1
	}
	spacer := strings.Repeat(" ", spacerWidth)
	
	headerMainLine := fmt.Sprintf("│ %s%s%s │", leftPart, spacer, rightPart)
	
	// Status line with location and progress info
	var locationDisplay string
	if m.selectedTodo != nil && m.selectedTodo.Context != "" {
		locationDisplay = localization.T(localization.StatusLocation, m.selectedTodo.Context)
	} else {
		locationDisplay = localization.T(localization.StatusLocation, localization.T(localization.LocationHome))
	}
	
	// Count today's tasks
	completedToday := 0
	totalToday := 0
	for _, todo := range m.todos {
		if todo.CreatedAt.Day() == time.Now().Day() {
			totalToday++
			if todo.Completed {
				completedToday++
			}
		}
	}
	
	tasksDisplay := localization.T(localization.StatusTasks, completedToday, totalToday)
	streakDisplay := localization.T(localization.StatusStreak, 3) // TODO: Calculate actual streak
	
	statusParts := fmt.Sprintf("%s   %s   %s", locationDisplay, tasksDisplay, streakDisplay)
	statusLine := fmt.Sprintf("│ %s%*s │", 
		lipgloss.NewStyle().Foreground(mutedColor).Render(statusParts),
		m.width-lipgloss.Width(statusParts)-4, "")
	
	// Bottom line of header
	headerBottom := "╰" + strings.Repeat("─", headerWidth) + "╯"
	
	headerContent := lipgloss.NewStyle().Foreground(mutedColor).Render(headerTop) + "\n" +
		headerMainLine + "\n" +
		statusLine + "\n" +
		lipgloss.NewStyle().Foreground(mutedColor).Render(headerBottom)
	
	// Stats bar with beautiful formatting
	statsBar := m.renderEnhancedStatsBar()
	
	// Todo list with adjusted height
	// Calculate available height more accurately
	headerLines := 4    // Header box
	statsLines := 2     // Stats bar
	bottomPanelLines := 9 // Keyboard shortcuts and activity log
	statusLines := 1    // Status message
	spacing := 2         // Vertical spacing
	
	listHeight := m.height - (headerLines + statsLines + bottomPanelLines + statusLines + spacing)
	if listHeight < 5 {
		listHeight = 5
	}
	if listHeight > 50 {
		listHeight = 50 // Cap max height for readability
	}
	
	// Ensure list width doesn't cause overflow
	listWidth := m.width - 4
	if listWidth < 60 {
		listWidth = 60
	}
	m.list.SetSize(listWidth, listHeight)
	listContent := m.list.View()
	
	// Create inline keyboard shortcuts and activity log (50/50 split)
	// Ensure minimum width for both panels
	halfWidth := (m.width - 6) / 2
	if halfWidth < 35 {
		halfWidth = 35
	}
	
	// Keyboard shortcuts (left side)
	keyboardShortcuts := []string{
		lipgloss.NewStyle().Foreground(primaryColor).Bold(true).Render(localization.T(localization.ShortcutsTitle)),
		"",
		lipgloss.NewStyle().Foreground(successColor).Render(localization.T(localization.ShortcutAdd)) + " " + localization.T(localization.ActionAdd) + "  " + 
			lipgloss.NewStyle().Foreground(secondaryColor).Render(localization.T(localization.ShortcutSuggest)) + " " + localization.T(localization.ActionSuggest),
		lipgloss.NewStyle().Foreground(successColor).Render(localization.T(localization.ShortcutComplete)) + "  " + localization.T(localization.ActionComplete) + "  " +
			lipgloss.NewStyle().Foreground(secondaryColor).Render(localization.T(localization.ShortcutStats)) + " " + localization.T(localization.ActionStats), 
		lipgloss.NewStyle().Foreground(successColor).Render(localization.T(localization.ShortcutSubtasks)) + "      " + localization.T(localization.ActionSubtasks) + "  " +
			lipgloss.NewStyle().Foreground(mutedColor).Render(localization.T(localization.ShortcutEdit)) + " " + localization.T(localization.ActionEdit),
		lipgloss.NewStyle().Foreground(mutedColor).Render(localization.T(localization.ShortcutNavigate)) + "    " + localization.T(localization.ActionNavigate) + "  " +
			lipgloss.NewStyle().Foreground(errorColor).Render(localization.T(localization.ShortcutDelete)) + " " + localization.T(localization.ActionDelete),
		lipgloss.NewStyle().Foreground(mutedColor).Render(localization.T(localization.ShortcutDetails)) + "      " + localization.T(localization.ActionDetails) + "   " +
			lipgloss.NewStyle().Foreground(errorColor).Render(localization.T(localization.ShortcutQuit)) + "    " + localization.T(localization.ActionQuit),
	}
	
	keyHelpBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(mutedColor).
		Padding(0, 1).
		Width(halfWidth).
		Height(9).
		Render(strings.Join(keyboardShortcuts, "\n"))
	
	// Activity log (right side) 
	activityLogs := []string{
		lipgloss.NewStyle().Foreground(primaryColor).Bold(true).Render(localization.T(localization.ActivityLogTitle)),
		"",
	}
	
	// Add cost tracking at the top if available
	if m.processor != nil && m.processor.TotalCost > 0 {
		costLine := fmt.Sprintf("💰 $%.4f | Fast:%d | Smart:%d",
			m.processor.TotalCost,
			m.processor.FastCalls,
			m.processor.SmartCalls)
		activityLogs = append(activityLogs, 
			lipgloss.NewStyle().Foreground(warningColor).Render(costLine))
	}
	
	// Add the actual logs, truncating long lines
	if len(m.consoleLogs) == 0 {
		activityLogs = append(activityLogs, lipgloss.NewStyle().Foreground(mutedColor).Render(localization.T(localization.ActivityNoRecent)))
	} else {
		// Truncate each log line to fit within the box
		maxLogWidth := halfWidth - 4 // Account for padding and borders
		for _, log := range m.consoleLogs {
			if lipgloss.Width(log) > maxLogWidth {
				// Truncate with ellipsis
				runes := []rune(log)
				if len(runes) > maxLogWidth-3 {
					log = string(runes[:maxLogWidth-3]) + "..."
				}
			}
			activityLogs = append(activityLogs, log)
		}
	}
	
	// Ensure we don't exceed the height
	maxLogLines := 7
	if len(activityLogs) > maxLogLines {
		activityLogs = append(activityLogs[:2], activityLogs[len(activityLogs)-(maxLogLines-2):]...)
	}
	
	activityBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(mutedColor).
		Padding(0, 1).
		Width(halfWidth).
		Height(9).
		Render(strings.Join(activityLogs, "\n"))
	
	// Join keyboard and activity side by side (no gap)
	bottomPanel := lipgloss.JoinHorizontal(
		lipgloss.Top,
		keyHelpBox,
		activityBox,
	)
	
	// Status message with animation
	var statusMsgLine string
	if m.statusMsg != "" {
		var statusStyle lipgloss.Style
		switch m.statusType {
		case "success":
			statusStyle = lipgloss.NewStyle().Foreground(successColor).Bold(true)
		case "error":
			statusStyle = lipgloss.NewStyle().Foreground(errorColor).Bold(true)
		case "warning":
			statusStyle = lipgloss.NewStyle().Foreground(warningColor).Bold(true)
		default:
			statusStyle = lipgloss.NewStyle().Foreground(primaryColor)
		}
		statusMsgLine = lipgloss.NewStyle().
			Padding(0, 2).
			Render(statusStyle.Render("▶ " + m.statusMsg))
	}
	
	// Compose the full view with proper spacing
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		headerContent,
		statsBar,
		listContent,
		statusMsgLine,
		bottomPanel,
	)
	
	// Use Top alignment for better layout
	return lipgloss.Place(m.width, m.height, lipgloss.Left, lipgloss.Top, content)
}
