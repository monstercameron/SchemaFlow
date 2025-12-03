package ops

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ParseJSON parses a JSON string into a target struct, handling common LLM output issues
// like markdown code blocks and whitespace.
func ParseJSON(input string, target any) error {
	// Clean up the input
	cleaned := cleanJSON(input)

	// Try to unmarshal
	err := json.Unmarshal([]byte(cleaned), target)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w\nInput: %s", err, cleaned)
	}

	return nil
}

// cleanJSON removes markdown code blocks and extra whitespace
func cleanJSON(input string) string {
	input = strings.TrimSpace(input)

	// Remove markdown code blocks if present
	if strings.HasPrefix(input, "```json") {
		input = strings.TrimPrefix(input, "```json")
		input = strings.TrimSuffix(input, "```")
	} else if strings.HasPrefix(input, "```") {
		input = strings.TrimPrefix(input, "```")
		input = strings.TrimSuffix(input, "```")
	}

	return strings.TrimSpace(input)
}
