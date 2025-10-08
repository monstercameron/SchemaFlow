package schemaflow

import (
	"context"
	"fmt"
	"strings"
)

func Summarize(input string, opts SummarizeOptions) (string, error) {
	// Validate options
	if err := opts.Validate(); err != nil {
		return "", fmt.Errorf("invalid options: %w", err)
	}
	
	// Build summarization instructions
	var instructions []string
	
	if opts.TargetLength > 0 {
		instructions = append(instructions, fmt.Sprintf("Target length: %d %s", opts.TargetLength, opts.LengthUnit))
	}
	
	if opts.BulletPoints {
		instructions = append(instructions, "Format as bullet points")
	} else if opts.Style != "" {
		instructions = append(instructions, fmt.Sprintf("Style: %s", opts.Style))
	}
	
	if len(opts.FocusAreas) > 0 {
		instructions = append(instructions, fmt.Sprintf("Focus on: %s", strings.Join(opts.FocusAreas, ", ")))
	}
	
	if len(opts.PreserveInfo) > 0 {
		instructions = append(instructions, fmt.Sprintf("Must preserve: %s", strings.Join(opts.PreserveInfo, ", ")))
	}
	
	opt := opts.toOpOptions()
	if len(instructions) > 0 {
		steering := strings.Join(instructions, ". ")
		if opts.Steering != "" {
			steering = opts.Steering + ". " + steering
		}
		opt.Steering = steering
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	
	systemPrompt := `You are a text summarization expert. Create concise summaries that preserve key information.

Rules:
- Maintain the most important points
- Use clear, concise language
- Preserve critical details and context
- Keep the original tone when appropriate`
	
	userPrompt := fmt.Sprintf("Summarize this text:\n%s", input)
	
	response, err := callLLM(ctx, systemPrompt, userPrompt, opt)
	if err != nil {
		return "", SummarizeError{
			Input:  input,
			Length: len(input),
			Reason: err.Error(),
		}
	}
	
	return strings.TrimSpace(response), nil
}

func Rewrite(input string, opts RewriteOptions) (string, error) {
	// Validate options
	if err := opts.Validate(); err != nil {
		return "", fmt.Errorf("invalid options: %w", err)
	}
	
	// Build rewrite instructions
	var instructions []string
	
	if opts.TargetTone != "" {
		instructions = append(instructions, fmt.Sprintf("Target tone: %s", opts.TargetTone))
	}
	
	if opts.FormalityLevel != 5 {
		instructions = append(instructions, fmt.Sprintf("Formality level: %d/10", opts.FormalityLevel))
	}
	
	if opts.Audience != "" {
		instructions = append(instructions, fmt.Sprintf("Target audience: %s", opts.Audience))
	}
	
	if opts.StyleGuide != "" {
		instructions = append(instructions, fmt.Sprintf("Follow style: %s", opts.StyleGuide))
	}
	
	if len(opts.Changes) > 0 {
		instructions = append(instructions, fmt.Sprintf("Make these changes: %s", strings.Join(opts.Changes, ", ")))
	}
	
	if len(opts.AvoidWords) > 0 {
		instructions = append(instructions, fmt.Sprintf("Avoid: %s", strings.Join(opts.AvoidWords, ", ")))
	}
	
	if len(opts.IncludeWords) > 0 {
		instructions = append(instructions, fmt.Sprintf("Include: %s", strings.Join(opts.IncludeWords, ", ")))
	}
	
	if opts.PreserveFacts {
		instructions = append(instructions, "Preserve all factual information")
	}
	
	opt := opts.toOpOptions()
	if len(instructions) > 0 {
		steering := strings.Join(instructions, ". ")
		if opts.Steering != "" {
			steering = opts.Steering + ". " + steering
		}
		opt.Steering = steering
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	
	systemPrompt := `You are a text rewriting expert. Modify text while preserving its core meaning.

Rules:
- Maintain the original message and intent
- Improve clarity and readability
- Adapt style as requested
- Fix grammar and spelling errors`
	
	userPrompt := fmt.Sprintf("Rewrite this text:\n%s", input)
	
	response, err := callLLM(ctx, systemPrompt, userPrompt, opt)
	if err != nil {
		return "", RewriteError{
			Input:  input,
			Reason: err.Error(),
		}
	}
	
	return strings.TrimSpace(response), nil
}

func Translate(input string, opts TranslateOptions) (string, error) {
	// Validate options
	if err := opts.Validate(); err != nil {
		return "", fmt.Errorf("invalid options: %w", err)
	}
	
	// Build translation instructions
	var instructions []string
	
	instructions = append(instructions, fmt.Sprintf("Translate to %s", opts.TargetLanguage))
	
	if opts.SourceLanguage != "" {
		instructions = append(instructions, fmt.Sprintf("From %s", opts.SourceLanguage))
	}
	
	if opts.Dialect != "" {
		instructions = append(instructions, fmt.Sprintf("Use %s dialect", opts.Dialect))
	}
	
	if opts.Formality != "neutral" {
		instructions = append(instructions, fmt.Sprintf("Formality: %s", opts.Formality))
	}
	
	if opts.CulturalAdaptation != 5 {
		instructions = append(instructions, fmt.Sprintf("Cultural adaptation level: %d/10", opts.CulturalAdaptation))
	}
	
	if opts.PreserveFormatting {
		instructions = append(instructions, "Preserve formatting")
	}
	
	if len(opts.Glossary) > 0 {
		glossary := "Use glossary: "
		for term, translation := range opts.Glossary {
			glossary += fmt.Sprintf("%s=%s, ", term, translation)
		}
		instructions = append(instructions, strings.TrimSuffix(glossary, ", "))
	}
	
	opt := opts.toOpOptions()
	steering := strings.Join(instructions, ". ")
	if opts.Steering != "" {
		steering = opts.Steering + ". " + steering
	}
	opt.Steering = steering
	
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	
	systemPrompt := `You are a translation expert. Translate text accurately between languages.

Rules:
- Preserve meaning and nuance
- Maintain appropriate formality level
- Handle idioms and cultural references appropriately
- Keep technical terms accurate`
	
	userPrompt := fmt.Sprintf("Translate this text:\n%s", input)
	
	response, err := callLLM(ctx, systemPrompt, userPrompt, opt)
	if err != nil {
		return "", TranslateError{
			Input:  input,
			Reason: err.Error(),
		}
	}
	
	return strings.TrimSpace(response), nil
}

func Expand(input string, opts ExpandOptions) (string, error) {
	// Validate options
	if err := opts.Validate(); err != nil {
		return "", fmt.Errorf("invalid options: %w", err)
	}
	
	// Build expansion instructions
	var instructions []string
	
	if opts.ExpansionFactor != 2.0 {
		instructions = append(instructions, fmt.Sprintf("Expand by %.1fx", opts.ExpansionFactor))
	}
	
	instructions = append(instructions, fmt.Sprintf("Detail level: %d/10", opts.DetailLevel))
	
	if opts.ExpansionStyle != "" {
		instructions = append(instructions, fmt.Sprintf("Style: %s", opts.ExpansionStyle))
	}
	
	if opts.IncludeExamples {
		instructions = append(instructions, "Include relevant examples")
	}
	
	if len(opts.ElaborateOn) > 0 {
		instructions = append(instructions, fmt.Sprintf("Elaborate on: %s", strings.Join(opts.ElaborateOn, ", ")))
	}
	
	if len(opts.AddContext) > 0 {
		instructions = append(instructions, fmt.Sprintf("Add context about: %s", strings.Join(opts.AddContext, ", ")))
	}
	
	opt := opts.toOpOptions()
	if len(instructions) > 0 {
		steering := strings.Join(instructions, ". ")
		if opts.Steering != "" {
			steering = opts.Steering + ". " + steering
		}
		opt.Steering = steering
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	
	systemPrompt := `You are a content expansion expert. Elaborate on text with additional detail and context.

Rules:
- Add relevant details and examples
- Maintain consistency with the original
- Provide useful elaboration
- Keep the expanded content coherent`
	
	userPrompt := fmt.Sprintf("Expand on this text:\n%s", input)
	
	response, err := callLLM(ctx, systemPrompt, userPrompt, opt)
	if err != nil {
		return "", ExpandError{
			Input:  input,
			Reason: err.Error(),
		}
	}
	
	return strings.TrimSpace(response), nil
}