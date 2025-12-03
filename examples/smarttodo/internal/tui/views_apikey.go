package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/monstercameron/schemaflow"
)

// apiKeyViewRender renders the API key setup view
func (m Model) apiKeyViewRender() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		MarginBottom(2)

	warningStyle := lipgloss.NewStyle().
		Foreground(warningColor).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(warningColor).
		Padding(1, 2).
		MarginBottom(2)

	instructionStyle := lipgloss.NewStyle().
		Foreground(secondaryColor).
		MarginBottom(2)

	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(0, 1).
		Width(60)

	helpStyle := lipgloss.NewStyle().
		Foreground(mutedColor).
		MarginTop(2)

	errorStyle := lipgloss.NewStyle().
		Foreground(errorColor).
		Bold(true).
		MarginTop(1)

	title := titleStyle.Render("🔑 API Key Setup")

	warning := warningStyle.Render(`⚠️  No OpenAI API key found!
Smart Todo uses AI to process your tasks intelligently.`)

	instructions := instructionStyle.Render(`Please enter your OpenAI API key:
(Get one at https://platform.openai.com/api-keys)`)

	// Mask the API key input for security
	maskedInput := m.setupInput.Value()
	if len(maskedInput) > 8 {
		maskedInput = maskedInput[:4] + strings.Repeat("•", len(maskedInput)-8) + maskedInput[len(maskedInput)-4:]
	}
	m.setupInput.SetValue(maskedInput)
	inputBox := inputStyle.Render(m.setupInput.View())
	m.setupInput.SetValue(m.setupInput.Value()) // Reset to actual value

	help := helpStyle.Render(`Press Enter to validate and save
Press Esc to quit`)

	content := []string{
		title,
		warning,
		instructions,
		inputBox,
		help,
	}

	// Add error message if validation failed
	if m.statusMsg != "" && m.statusType == "error" {
		content = append(content, errorStyle.Render("❌ "+m.statusMsg))
	}

	joined := lipgloss.JoinVertical(
		lipgloss.Center,
		content...,
	)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		joined,
	)
}

// validateAPIKey tests if the API key works
func validateAPIKey(apiKey string) error {
	// Temporarily set the API key for validation
	os.Setenv("OPENAI_API_KEY", apiKey)

	// Try a simple API call to validate
	_, err := schemaflow.Generate[string]("Say 'test'", schemaflow.OpOptions{
		Intelligence: schemaflow.Fast,
	})

	if err != nil {
		// Clear the invalid key
		os.Unsetenv("OPENAI_API_KEY")
		return fmt.Errorf("invalid API key: %v", err)
	}

	return nil
}

// saveAPIKey saves the API key to .env file
func saveAPIKey(apiKey string) error {
	// Create or append to .env file
	envPath := ".env"

	// Read existing .env content
	content, err := os.ReadFile(envPath)
	existingLines := []string{}
	if err == nil {
		existingLines = strings.Split(string(content), "\n")
	}

	// Update or add OPENAI_API_KEY
	found := false
	for i, line := range existingLines {
		if strings.HasPrefix(line, "OPENAI_API_KEY=") {
			existingLines[i] = fmt.Sprintf("OPENAI_API_KEY=%s", apiKey)
			found = true
			break
		}
	}

	if !found {
		existingLines = append(existingLines, fmt.Sprintf("OPENAI_API_KEY=%s", apiKey))
	}

	// Write back to file
	newContent := strings.Join(existingLines, "\n")
	return os.WriteFile(envPath, []byte(newContent), 0600)
}

