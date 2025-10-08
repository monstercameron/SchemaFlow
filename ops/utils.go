package ops

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/monstercameron/SchemaFlow/core"
)

// GenerateTypeSchema creates a human-readable schema description for a Go type
func GenerateTypeSchema(targetType reflect.Type) string {
	if targetType.Kind() == reflect.Ptr {
		targetType = targetType.Elem()
	}

	switch targetType.Kind() {
	case reflect.Struct:
		var fields []string
		for i := 0; i < targetType.NumField(); i++ {
			field := targetType.Field(i)

			// Skip unexported fields
			if !field.IsExported() {
				continue
			}

			// Get JSON tag or use field name
			jsonTag := field.Tag.Get("json")
			fieldName := field.Name
			if jsonTag != "" {
				parts := strings.Split(jsonTag, ",")
				if parts[0] != "-" {
					fieldName = parts[0]
				}
			}

			// Get field type description
			fieldType := GetTypeDescription(field.Type)

			// Check if field is required (no omitempty tag)
			required := !strings.Contains(jsonTag, "omitempty")
			requiredStr := ""
			if required {
				requiredStr = " (required)"
			}

			// Add field description
			fields = append(fields, fmt.Sprintf("  %s: %s%s", fieldName, fieldType, requiredStr))
		}
		return fmt.Sprintf("{\n%s\n}", strings.Join(fields, "\n"))

	case reflect.Slice:
		elemType := targetType.Elem()
		return fmt.Sprintf("[]%s", GenerateTypeSchema(elemType))

	case reflect.Map:
		keyType := targetType.Key()
		valueType := targetType.Elem()
		return fmt.Sprintf("map[%s]%s", keyType.String(), GenerateTypeSchema(valueType))

	default:
		return GetTypeDescription(targetType)
	}
}

// GetTypeDescription returns a simple description of a type
func GetTypeDescription(targetType reflect.Type) string {
	switch targetType.Kind() {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "integer"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "unsigned integer"
	case reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Bool:
		return "boolean"
	case reflect.Struct:
		if targetType.String() == "time.Time" {
			return "datetime (RFC3339)"
		}
		return targetType.String()
	case reflect.Ptr:
		return GetTypeDescription(targetType.Elem()) + " (optional)"
	default:
		return targetType.String()
	}
}

// NormalizeInput converts various input types to string for LLM processing
func NormalizeInput(input any) (string, error) {
	if input == nil {
		return "", fmt.Errorf("input is nil")
	}

	switch inputValue := input.(type) {
	case string:
		return inputValue, nil
	case []byte:
		return string(inputValue), nil
	case fmt.Stringer:
		return inputValue.String(), nil
	default:
		// Try JSON marshaling for complex types
		if b, err := json.Marshal(input); err == nil {
			return string(b), nil
		}
		// Fallback to fmt.Sprint
		return fmt.Sprint(input), nil
	}
}

// BuildExtractSystemPrompt creates the system prompt for extraction based on mode
func BuildExtractSystemPrompt(typeSchema string, mode core.Mode) string {
	base := fmt.Sprintf(`You are a data extraction expert. Extract structured data from the input and return it as JSON.
Target schema:
%s

Rules:
- Extract all relevant information that maps to the schema
- Return ONLY valid JSON, no explanations or markdown`, typeSchema)

	switch mode {
	case core.Strict:
		return base + `
- All required fields MUST be present and valid
- Fail if any required field cannot be extracted
- Use null only for explicitly optional fields
- Validate data types strictly`

	case core.TransformMode:
		return base + `
- Infer missing fields intelligently when possible
- Use reasonable defaults for missing data
- Be flexible with type conversions
- Preserve as much information as possible`

	case core.Creative:
		return base + `
- Creatively interpret ambiguous data
- Generate plausible values for missing fields
- Use context to enrich extracted data
- Prioritize completeness over strict accuracy`

	default:
		return base
	}
}

// BuildGenerateStringPrompt creates the system prompt for string generation
func BuildGenerateStringPrompt(mode core.Mode) string {
	switch mode {
	case core.Strict:
		return "You are a precise content generator. Generate exactly what is requested, following all specifications strictly."
	case core.TransformMode:
		return "You are a creative content generator. Generate the requested content while maintaining quality and relevance."
	case core.Creative:
		return "You are a highly creative content generator. Generate engaging, original content based on the prompt."
	default:
		return "You are a content generator. Generate the requested content based on the prompt."
	}
}

// CalculateParsingConfidence estimates confidence when parsing partially fails
func CalculateParsingConfidence(response string, targetType reflect.Type) float64 {
	// Basic heuristic: check if response looks like valid JSON
	response = strings.TrimSpace(response)
	if strings.HasPrefix(response, "{") && strings.HasSuffix(response, "}") {
		return 0.3 // Looks like JSON but failed to parse
	}
	if strings.HasPrefix(response, "[") && strings.HasSuffix(response, "]") {
		return 0.3 // Looks like JSON array but failed to parse
	}
	return 0.1 // Doesn't look like JSON at all
}

// ValidateExtractedData validates extracted data meets requirements
func ValidateExtractedData(data any, threshold float64) error {
	// Basic validation - can be extended based on needs
	if data == nil {
		return fmt.Errorf("data cannot be nil")
	}

	value := reflect.ValueOf(data)
	if !value.IsValid() {
		return fmt.Errorf("invalid data")
	}

	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return fmt.Errorf("data cannot be nil pointer")
		}
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return nil // Only validate structs for now
	}

	// Check for zero values in required fields
	t := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := t.Field(i)
		fieldValue := value.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Check if field is required (no omitempty tag)
		jsonTag := field.Tag.Get("json")
		if !strings.Contains(jsonTag, "omitempty") && fieldValue.IsZero() {
			return fmt.Errorf("required field %s is empty", field.Name)
		}
	}

	return nil
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
