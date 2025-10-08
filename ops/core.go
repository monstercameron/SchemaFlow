// package ops - Data operations for structured data extraction and transformation
package ops

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// Extract converts unstructured data into strongly-typed Go structs using LLM interpretation.
// It handles various input formats (string, JSON, structs) and maps them to the target type.
//
// Type parameter T specifies the target Go type to extract into.
//
// Examples:
//
//	// Basic extraction
//	person, err := Extract[Person]("John Smith, 28 years old", NewExtractOptions())
//
//	// Extraction with schema hints
//	data, err := Extract[Invoice](jsonString, NewExtractOptions().
//	    WithStrictSchema(true).
//	    WithSchemaHints(map[string]string{
//	        "date": "ISO 8601 format",
//	        "amount": "USD currency",
//	    }))
//
//	// Extraction with examples
//	product, err := Extract[Product](description, NewExtractOptions().
//	    WithExamples(exampleProduct1, exampleProduct2).
//	    WithSteering("Focus on technical specifications"))
//
// The operation uses schema inference to guide the LLM in structured extraction.
// In Strict mode, all required fields must be present. In Transform mode (default),
// the LLM will intelligently infer missing fields.
func Extract[T any](input any, opts ExtractOptions) (T, error) {
	// Use default client if available, otherwise use global config
	if defaultClient != nil {
		return ClientExtract[T](defaultClient, input, opts)
	}
	return extractImpl[T](input, opts)
}

// ClientExtract is a client method that converts unstructured data into strongly-typed Go structs.
// This is the client-based version of the Extract operation.
// Usage: result, err := ClientExtract[Person](client, input, NewExtractOptions())
func ClientExtract[T any](c *Client, input any, opts ExtractOptions) (T, error) {
	// Temporarily set global variables to client values for compatibility
	oldClient := client
	oldTimeout := timeout
	oldMaxRetries := maxRetries
	oldLogger := logger
	
	c.mu.RLock()
	client = c.openaiClient
	timeout = c.timeout
	maxRetries = c.maxRetries
	logger = c.logger
	c.mu.RUnlock()
	
	defer func() {
		client = oldClient
		timeout = oldTimeout
		maxRetries = oldMaxRetries
		logger = oldLogger
	}()
	
	return extractImpl[T](input, opts)
}

// extractImpl contains the actual implementation of Extract
func extractImpl[T any](input any, opts ExtractOptions) (T, error) {
	var result T
	
	// Validate options
	if err := opts.Validate(); err != nil {
		return result, fmt.Errorf("invalid options: %w", err)
	}
	
	// Convert to legacy OpOptions for internal use
	opt := opts.toOpOptions()
	
	// Enhance steering with extraction-specific options
	if opts.SchemaHints != nil || opts.Examples != nil || opts.FieldRules != nil {
		var steeringParts []string
		if opts.Steering != "" {
			steeringParts = append(steeringParts, opts.Steering)
		}
		
		if opts.StrictSchema {
			steeringParts = append(steeringParts, "Enforce strict schema validation. All fields must be present and valid.")
		}
		
		if opts.AllowPartial {
			steeringParts = append(steeringParts, "Allow partial extraction if some fields are missing.")
		}
		
		if len(opts.SchemaHints) > 0 {
			hints := "Schema hints: "
			for field, hint := range opts.SchemaHints {
				hints += fmt.Sprintf("%s (%s), ", field, hint)
			}
			steeringParts = append(steeringParts, strings.TrimSuffix(hints, ", "))
		}
		
		if len(opts.FieldRules) > 0 {
			rules := "Field rules: "
			for field, rule := range opts.FieldRules {
				rules += fmt.Sprintf("%s: %s; ", field, rule)
			}
			steeringParts = append(steeringParts, strings.TrimSuffix(rules, "; "))
		}
		
		if len(opts.Examples) > 0 {
			examplesJSON, _ := json.Marshal(opts.Examples)
			steeringParts = append(steeringParts, fmt.Sprintf("Follow these examples: %s", string(examplesJSON)))
		}
		
		opt.Steering = strings.Join(steeringParts, ". ")
	}
	
	// Start operation timing
	startTime := time.Now()
	defer func() {
		if metricsEnabled {
			recordMetric("extract_duration", time.Since(startTime).Milliseconds(), map[string]string{
				"type": reflect.TypeOf(result).String(),
				"mode": opt.Mode.String(),
			})
		}
	}()
	
	// Log operation start
	logger.Info("Extract operation started",
		"requestID", opt.requestID,
		"targetType", reflect.TypeOf(result).String(),
		"mode", opt.Mode.String(),
		"intelligence", opt.Intelligence.String(),
	)
	
	// Validate input
	if input == nil {
		err := ExtractError{
			Input:      input,
			TargetType: reflect.TypeOf(result).String(),
			Reason:     "input cannot be nil",
			RequestID:  opt.requestID,
			Timestamp:  time.Now(),
		}
		logger.Error("Extract failed: nil input", "requestID", opt.requestID, "error", err)
		return result, err
	}
	
	// Get context with timeout
	ctx := opt.context
	if ctx == nil {
		ctx = context.Background()
	}
	
	// Generate type schema for the target type
	targetType := reflect.TypeOf(result)
	typeInfo := generateTypeSchema(targetType)
	
	// Convert input to string format for LLM processing
	inputStr, err := normalizeInput(input)
	if err != nil {
		extractErr := ExtractError{
			Input:      input,
			TargetType: targetType.String(),
			Reason:     fmt.Sprintf("failed to normalize input: %v", err),
			RequestID:  opt.requestID,
			Timestamp:  time.Now(),
		}
		logger.Error("Extract failed: input normalization error",
			"requestID", opt.requestID,
			"error", extractErr,
		)
		return result, extractErr
	}
	
	// Log input details in debug mode
	if debugMode || false {
		logger.Debug("Extract input normalized",
			"requestID", opt.requestID,
			"inputLength", len(inputStr),
			"inputPreview", inputStr[:min(len(inputStr), 100)],
		)
	}
	
	// Build system prompt based on mode
	systemPrompt := buildExtractSystemPrompt(typeInfo, opt.Mode)
	
	// Build user prompt
	userPrompt := fmt.Sprintf("Extract structured data from this input:\n%s", inputStr)
	
	// Call LLM for extraction
	response, err := callLLM(ctx, systemPrompt, userPrompt, opt)
	if err != nil {
		extractErr := ExtractError{
			Input:      input,
			TargetType: targetType.String(),
			Reason:     err.Error(),
			Confidence: 0,
			RequestID:  opt.requestID,
			Timestamp:  time.Now(),
		}
		logger.Error("Extract failed: LLM error",
			"requestID", opt.requestID,
			"error", extractErr,
		)
		return result, extractErr
	}
	
	// Parse JSON response into target type
	if err := parseJSON(response, &result); err != nil {
		// Calculate partial confidence based on parsing attempt
		confidence := calculateParsingConfidence(response, targetType)
		
		extractErr := ExtractError{
			Input:      input,
			TargetType: targetType.String(),
			Reason:     fmt.Sprintf("failed to parse response: %v", err),
			Confidence: confidence,
			RequestID:  opt.requestID,
			Timestamp:  time.Now(),
		}
		
		logger.Error("Extract failed: JSON parsing error",
			"requestID", opt.requestID,
			"confidence", confidence,
			"error", extractErr,
		)
		return result, extractErr
	}
	
	// Validate extracted data if in Strict mode
	if opt.Mode == Strict {
		if err := validateExtractedData(result, opt.Threshold); err != nil {
			extractErr := ExtractError{
				Input:      input,
				TargetType: targetType.String(),
				Reason:     fmt.Sprintf("validation failed: %v", err),
				Confidence: opt.Threshold - 0.1, // Just below threshold
				RequestID:  opt.requestID,
				Timestamp:  time.Now(),
			}
			logger.Error("Extract failed: validation error",
				"requestID", opt.requestID,
				"error", extractErr,
			)
			return result, extractErr
		}
	}
	
	logger.Info("Extract operation completed",
		"requestID", opt.requestID,
		"duration", time.Since(startTime),
	)
	
	return result, nil
}

// Transform converts data from one type to another using semantic mapping.
// It understands relationships between different structures and maps fields intelligently.
//
// Type parameters:
//   - T: Source type
//   - U: Target type
//
// Examples:
//
//	// Basic transformation
//	employee, err := Transform[Person, Employee](person, NewTransformOptions())
//
//	// Transformation with mapping rules
//	output, err := Transform[OldFormat, NewFormat](input, NewTransformOptions().
//	    WithMappingRules(map[string]string{
//	        "fullName": "firstName + ' ' + lastName",
//	        "age": "calculate from birthDate",
//	    }).
//	    WithPreserveFields([]string{"id", "createdAt"}))
//
// The operation uses semantic understanding to map between related but structurally
// different types. It can handle field renaming, type conversion, and derived fields.
func Transform[T any, U any](input T, opts TransformOptions) (U, error) {
	if defaultClient != nil {
		return ClientTransform[T, U](defaultClient, input, opts)
	}
	return transformImpl[T, U](input, opts)
}

// ClientTransform is a client method that converts data from one type to another.
// Usage: result, err := ClientTransform[Person, Employee](client, input)
func ClientTransform[T any, U any](c *Client, input T, opts TransformOptions) (U, error) {
	// Temporarily set global variables to client values
	oldClient := client
	oldTimeout := timeout
	oldMaxRetries := maxRetries
	oldLogger := logger
	
	c.mu.RLock()
	client = c.openaiClient
	timeout = c.timeout
	maxRetries = c.maxRetries
	logger = c.logger
	c.mu.RUnlock()
	
	defer func() {
		client = oldClient
		timeout = oldTimeout
		maxRetries = oldMaxRetries
		logger = oldLogger
	}()
	
	return transformImpl[T, U](input, opts)
}


// transformImpl contains the actual implementation
func transformImpl[T any, U any](input T, opts TransformOptions) (U, error) {
	var result U
	
	// Validate options
	if err := opts.Validate(); err != nil {
		return result, fmt.Errorf("invalid options: %w", err)
	}
	
	// Convert to legacy OpOptions
	opt := opts.toOpOptions()
	
	// Enhance steering with transformation-specific options
	var steeringParts []string
	if opts.Steering != "" {
		steeringParts = append(steeringParts, opts.Steering)
	}
	
	if opts.TransformLogic != "" {
		steeringParts = append(steeringParts, fmt.Sprintf("Apply this transformation: %s", opts.TransformLogic))
	}
	
	if len(opts.MappingRules) > 0 {
		rules := "Field mappings: "
		for target, source := range opts.MappingRules {
			rules += fmt.Sprintf("%s <- %s; ", target, source)
		}
		steeringParts = append(steeringParts, strings.TrimSuffix(rules, "; "))
	}
	
	if len(opts.PreserveFields) > 0 {
		steeringParts = append(steeringParts, fmt.Sprintf("Preserve these fields: %s", strings.Join(opts.PreserveFields, ", ")))
	}
	
	if opts.MergeStrategy != "" {
		steeringParts = append(steeringParts, fmt.Sprintf("Use %s merge strategy", opts.MergeStrategy))
	}
	
	if len(steeringParts) > 0 {
		opt.Steering = strings.Join(steeringParts, ". ")
	}
	
	startTime := time.Now()
	defer func() {
		if metricsEnabled {
			recordMetric("transform_duration", time.Since(startTime).Milliseconds(), map[string]string{
				"from_type": reflect.TypeOf(input).String(),
				"to_type":   reflect.TypeOf(result).String(),
				"mode":      opt.Mode.String(),
			})
		}
	}()
	
	logger.Info("Transform operation started",
		"requestID", opt.requestID,
		"fromType", reflect.TypeOf(input).String(),
		"toType", reflect.TypeOf(result).String(),
	)
	
	// Validate input
	if reflect.ValueOf(input).IsZero() {
		logger.Warn("Transform received zero value input",
			"requestID", opt.requestID,
			"type", reflect.TypeOf(input).String(),
		)
	}
	
	ctx := opt.context
	if ctx == nil {
		ctx = context.Background()
	}
	
	// Get type information
	fromType := reflect.TypeOf(input)
	toType := reflect.TypeOf(result)
	
	fromSchema := generateTypeSchema(fromType)
	toSchema := generateTypeSchema(toType)
	
	// Marshal input to JSON
	inputJSON, err := json.Marshal(input)
	if err != nil {
		transformErr := TransformError{
			Input:     input,
			FromType:  fromType.String(),
			ToType:    toType.String(),
			Reason:    fmt.Sprintf("failed to marshal input: %v", err),
			RequestID: opt.requestID,
			Timestamp: time.Now(),
		}
		logger.Error("Transform failed: marshaling error",
			"requestID", opt.requestID,
			"error", transformErr,
		)
		return result, transformErr
	}
	
	// Build transformation prompt
	systemPrompt := fmt.Sprintf(`You are a data transformation expert. Transform data from one type to another using semantic mapping.

Source schema:
%s

Target schema:
%s

Transformation rules:
- Map semantically related fields even if names differ
- Infer and calculate derived fields when possible
- Handle type conversions intelligently (e.g., string to number)
- Combine or split fields as needed (e.g., fullName <-> firstName+lastName)
- Use reasonable defaults for missing required fields
- Preserve data integrity and meaning
- Return ONLY valid JSON matching the target schema`, fromSchema, toSchema)
	
	userPrompt := fmt.Sprintf("Transform this data:\n%s", string(inputJSON))
	
	// Log transformation details in debug mode
	if debugMode || false {
		logger.Debug("Transform schemas",
			"requestID", opt.requestID,
			"fromSchema", fromSchema,
			"toSchema", toSchema,
		)
	}
	
	// Call LLM for transformation
	response, err := callLLM(ctx, systemPrompt, userPrompt, opt)
	if err != nil {
		transformErr := TransformError{
			Input:     input,
			FromType:  fromType.String(),
			ToType:    toType.String(),
			Reason:    err.Error(),
			RequestID: opt.requestID,
			Timestamp: time.Now(),
		}
		logger.Error("Transform failed: LLM error",
			"requestID", opt.requestID,
			"error", transformErr,
		)
		return result, transformErr
	}
	
	// Parse transformed data
	if err := parseJSON(response, &result); err != nil {
		transformErr := TransformError{
			Input:      input,
			FromType:   fromType.String(),
			ToType:     toType.String(),
			Reason:     fmt.Sprintf("failed to parse response: %v", err),
			Confidence: 0.5,
			RequestID:  opt.requestID,
			Timestamp:  time.Now(),
		}
		logger.Error("Transform failed: parsing error",
			"requestID", opt.requestID,
			"error", transformErr,
		)
		return result, transformErr
	}
	
	logger.Info("Transform operation completed",
		"requestID", opt.requestID,
		"duration", time.Since(startTime),
	)
	
	return result, nil
}

// Generate creates structured data from natural language prompts.
// It can generate both simple strings and complex structured types.
//
// Type parameter T specifies the type to generate.
//
// Examples:
//
//	// Basic generation
//	product, err := Generate[Product]("Create a laptop listing", NewGenerateOptions())
//
//	// Generation with constraints
//	users, err := Generate[[]User]("Generate test users", NewGenerateOptions().
//	    WithCount(10).
//	    WithConstraints(map[string]interface{}{
//	        "age": "18-65",
//	        "country": "US or Canada",
//	    }).
//	    WithEnsureUnique(true))
//
//	// Generation with template and examples
//	content, err := Generate[BlogPost]("Tech article", NewGenerateOptions().
//	    WithTemplate(articleTemplate).
//	    WithExamples(example1, example2).
//	    WithStyle("technical but accessible"))
//
// The operation understands the target type structure and generates
// appropriate data that conforms to the schema.
func Generate[T any](prompt string, opts GenerateOptions) (T, error) {
	var result T
	
	// Validate options
	if err := opts.Validate(); err != nil {
		return result, fmt.Errorf("invalid options: %w", err)
	}
	
	// Handle batch generation if Count > 1
	if opts.Count > 1 {
		// This would need special handling for slice types
		return result, fmt.Errorf("batch generation not yet supported - use Count=1")
	}
	
	// Convert to legacy OpOptions
	opt := opts.toOpOptions()
	
	// Build enhanced prompt
	var promptParts []string
	promptParts = append(promptParts, prompt)
	
	if opts.Template != "" {
		promptParts = append(promptParts, fmt.Sprintf("Use this template: %s", opts.Template))
	}
	
	if opts.Style != "" {
		promptParts = append(promptParts, fmt.Sprintf("Style: %s", opts.Style))
	}
	
	if len(opts.Constraints) > 0 {
		constraints := "Constraints: "
		for key, value := range opts.Constraints {
			constraints += fmt.Sprintf("%s=%v; ", key, value)
		}
		promptParts = append(promptParts, strings.TrimSuffix(constraints, "; "))
	}
	
	if opts.SeedData != nil {
		seedJSON, _ := json.Marshal(opts.SeedData)
		promptParts = append(promptParts, fmt.Sprintf("Base on this seed data: %s", string(seedJSON)))
	}
	
	if len(opts.Examples) > 0 {
		examplesJSON, _ := json.Marshal(opts.Examples)
		promptParts = append(promptParts, fmt.Sprintf("Follow these examples: %s", string(examplesJSON)))
	}
	
	prompt = strings.Join(promptParts, ". ")
	if opts.Steering != "" {
		opt.Steering = opts.Steering
	}
	
	startTime := time.Now()
	defer func() {
		if metricsEnabled {
			recordMetric("generate_duration", time.Since(startTime).Milliseconds(), map[string]string{
				"type": reflect.TypeOf(result).String(),
				"mode": opt.Mode.String(),
			})
		}
	}()
	
	logger.Info("Generate operation started",
		"requestID", opt.requestID,
		"targetType", reflect.TypeOf(result).String(),
		"promptLength", len(prompt),
	)
	
	// Validate prompt
	prompt = strings.TrimSpace(prompt)
	if prompt == "" {
		err := GenerateError{
			Prompt:     prompt,
			TargetType: reflect.TypeOf(result).String(),
			Reason:     "prompt cannot be empty",
			RequestID:  opt.requestID,
			Timestamp:  time.Now(),
		}
		logger.Error("Generate failed: empty prompt",
			"requestID", opt.requestID,
			"error", err,
		)
		return result, err
	}
	
	// Sanitize prompt length
	const maxPromptLength = 10000
	if len(prompt) > maxPromptLength {
		logger.Warn("Generate prompt truncated",
			"requestID", opt.requestID,
			"originalLength", len(prompt),
			"maxLength", maxPromptLength,
		)
		prompt = prompt[:maxPromptLength]
	}
	
	ctx := opt.context
	if ctx == nil {
		ctx = context.Background()
	}
	
	targetType := reflect.TypeOf(result)
	
	// Handle string generation differently (simpler)
	if targetType.Kind() == reflect.String {
		systemPrompt := buildGenerateStringPrompt(opt.Mode)
		
		response, err := callLLM(ctx, systemPrompt, prompt, opt)
		if err != nil {
			genErr := GenerateError{
				Prompt:     prompt,
				TargetType: targetType.String(),
				Reason:     err.Error(),
				RequestID:  opt.requestID,
				Timestamp:  time.Now(),
			}
			logger.Error("Generate failed: LLM error",
				"requestID", opt.requestID,
				"error", genErr,
			)
			return result, genErr
		}
		
		// Set string result using reflection
		reflect.ValueOf(&result).Elem().SetString(response)
		
		logger.Info("Generate operation completed (string)",
			"requestID", opt.requestID,
			"duration", time.Since(startTime),
			"responseLength", len(response),
		)
		return result, nil
	}
	
	// Handle structured type generation
	typeSchema := generateTypeSchema(targetType)
	
	systemPrompt := fmt.Sprintf(`You are a data generation expert. Generate structured data based on the prompt.

Target schema:
%s

Generation rules:
- Generate realistic and coherent data
- Follow the prompt requirements precisely
- Ensure all required fields are populated with appropriate values
- Use sensible defaults where not specified
- Maintain internal consistency (e.g., related fields should make sense together)
- Return ONLY valid JSON matching the schema, no explanations`, typeSchema)
	
	// Log generation details in debug mode
	if debugMode || false {
		logger.Debug("Generate schema",
			"requestID", opt.requestID,
			"schema", typeSchema,
			"prompt", prompt[:min(len(prompt), 200)],
		)
	}
	
	response, err := callLLM(ctx, systemPrompt, prompt, opt)
	if err != nil {
		genErr := GenerateError{
			Prompt:     prompt,
			TargetType: targetType.String(),
			Reason:     err.Error(),
			RequestID:  opt.requestID,
			Timestamp:  time.Now(),
		}
		logger.Error("Generate failed: LLM error",
			"requestID", opt.requestID,
			"error", genErr,
		)
		return result, genErr
	}
	
	// Parse generated data
	if err := parseJSON(response, &result); err != nil {
		genErr := GenerateError{
			Prompt:     prompt,
			TargetType: targetType.String(),
			Reason:     fmt.Sprintf("failed to parse response: %v", err),
			RequestID:  opt.requestID,
			Timestamp:  time.Now(),
		}
		logger.Error("Generate failed: parsing error",
			"requestID", opt.requestID,
			"error", genErr,
			"response", response[:min(len(response), 200)],
		)
		return result, genErr
	}
	
	logger.Info("Generate operation completed",
		"requestID", opt.requestID,
		"duration", time.Since(startTime),
	)
	
	return result, nil
}

// Helper functions for data operations

// generateTypeSchema creates a human-readable schema description for a Go type
func generateTypeSchema(targetType reflect.Type) string {
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
			fieldType := getTypeDescription(field.Type)
			
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
		return fmt.Sprintf("[]%s", generateTypeSchema(elemType))
		
	case reflect.Map:
		keyType := targetType.Key()
		valueType := targetType.Elem()
		return fmt.Sprintf("map[%s]%s", keyType.String(), generateTypeSchema(valueType))
		
	default:
		return getTypeDescription(targetType)
	}
}

// getTypeDescription returns a simple description of a type
func getTypeDescription(targetType reflect.Type) string {
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
		return getTypeDescription(targetType.Elem()) + " (optional)"
	default:
		return targetType.String()
	}
}


// normalizeInput converts various input types to string for LLM processing
func normalizeInput(input any) (string, error) {
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

// buildExtractSystemPrompt creates the system prompt for extraction based on mode
func buildExtractSystemPrompt(typeSchema string, mode Mode) string {
	base := fmt.Sprintf(`You are a data extraction expert. Extract structured data from the input and return it as JSON.
Target schema:
%s

Rules:
- Extract all relevant information that maps to the schema
- Return ONLY valid JSON, no explanations or markdown`, typeSchema)
	
	switch mode {
	case Strict:
		return base + `
- All required fields MUST be present and valid
- Fail if any required field cannot be extracted
- Use null only for explicitly optional fields
- Validate data types strictly`
		
	case TransformMode:
		return base + `
- Infer missing fields intelligently when possible
- Use reasonable defaults for missing data
- Be flexible with type conversions
- Preserve as much information as possible`
		
	case Creative:
		return base + `
- Creatively interpret ambiguous data
- Generate plausible values for missing fields
- Use context to enrich extracted data
- Prioritize completeness over strict accuracy`
		
	default:
		return base
	}
}

// buildGenerateStringPrompt creates the system prompt for string generation
func buildGenerateStringPrompt(mode Mode) string {
	switch mode {
	case Strict:
		return "You are a precise content generator. Generate exactly what is requested, following all specifications strictly."
	case TransformMode:
		return "You are a creative content generator. Generate the requested content while maintaining quality and relevance."
	case Creative:
		return "You are a highly creative content generator. Generate engaging, original content based on the prompt."
	default:
		return "You are a content generator. Generate the requested content based on the prompt."
	}
}

// calculateParsingConfidence estimates confidence when parsing partially fails
func calculateParsingConfidence(response string, targetType reflect.Type) float64 {
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

// validateExtractedData validates extracted data meets requirements
func validateExtractedData(data any, threshold float64) error {
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
