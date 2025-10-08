// package core - JSON parsing utilities
package core

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// ParseJSON attempts to parse JSON from LLM response, handling common formatting issues
func ParseJSON[T any](response string, target *T) error {
	response = strings.TrimSpace(response)

	// Remove markdown code blocks if present
	if strings.HasPrefix(response, "```json") {
		response = strings.TrimPrefix(response, "```json")
		response = strings.TrimSuffix(response, "```")
		response = strings.TrimSpace(response)
	} else if strings.HasPrefix(response, "```") {
		response = strings.TrimPrefix(response, "```")
		response = strings.TrimSuffix(response, "```")
		response = strings.TrimSpace(response)
	}

	// Try standard unmarshaling first
	if err := json.Unmarshal([]byte(response), target); err != nil {
		// Try with a decoder for better error messages
		decoder := json.NewDecoder(strings.NewReader(response))
		decoder.DisallowUnknownFields()
		if decodeErr := decoder.Decode(target); decodeErr != nil {
			GetLogger().Error("JSON parsing failed",
				"error", decodeErr,
				"response", response[:min(len(response), 200)], // Log first 200 chars
			)
			return fmt.Errorf("JSON decode error: %w", decodeErr)
		}
	}

	return nil
}

// GetTypeDescription generates a string representation of a type for LLM prompts.
func GetTypeDescription(t reflect.Type) string {
	// For built-in types, use their standard name.
	// This also handles interface types like `error`.
	if t.PkgPath() == "" && t.Name() != "" {
		return t.Name()
	}

	switch t.Kind() {
	case reflect.Ptr:
		return GetTypeDescription(t.Elem()) + " (optional)"
	case reflect.Struct:
		// To avoid infinite recursion with self-referential structs, we can use the type name
		if t.Name() != "" {
			pkgPathParts := strings.Split(t.PkgPath(), "/")
			pkgName := pkgPathParts[len(pkgPathParts)-1]
			return pkgName + "." + t.Name()
		}
		var builder strings.Builder
		builder.WriteString("struct { ")
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if i > 0 {
				builder.WriteString("; ")
			}
			builder.WriteString(field.Name + " " + GetTypeDescription(field.Type))
		}
		builder.WriteString(" }")
		return builder.String()
	case reflect.Slice:
		return "[]" + GetTypeDescription(t.Elem())
	case reflect.Map:
		return "map[" + GetTypeDescription(t.Key()) + "]" + GetTypeDescription(t.Elem())
	default:
		return t.String()
	}
}

// ValidateExtractedData checks if the extracted data meets a quality threshold.
// It returns an error if the data is considered invalid (e.g., nil, empty, or zero value).
func ValidateExtractedData(data interface{}, threshold float64) error {
	if data == nil {
		return fmt.Errorf("extracted data is nil")
	}

	val := reflect.ValueOf(data)

	// At a low threshold, even zero values are acceptable.
	if threshold < 0.2 {
		return nil
	}

	if val.IsZero() {
		// For structs, IsZero is a good check for uninitialized data.
		if val.Kind() == reflect.Struct {
			// Allow zero-value structs only if the threshold is low
			if threshold < 0.5 {
				return nil
			}
			return fmt.Errorf("extracted data is the zero value for its type")
		}
	}

	// Additional checks for specific types
	switch val.Kind() {
	case reflect.String:
		// An empty string is only an error if the threshold is high
		if val.String() == "" && threshold >= 0.8 {
			return fmt.Errorf("extracted string is empty")
		}
	case reflect.Slice, reflect.Map:
		if val.Len() == 0 && threshold >= 0.8 {
			return fmt.Errorf("extracted slice or map is empty")
		}
	case reflect.Struct:
		jsonData, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to marshal struct to JSON for validation: %w", err)
		}
		// An empty JSON object is only an error at high thresholds
		if string(jsonData) == "{}" && threshold >= 0.9 {
			return fmt.Errorf("extracted struct is empty")
		}
	}

	return nil
}
