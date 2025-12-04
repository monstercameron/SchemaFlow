// package ops - Negotiate operation for reconciling competing constraints
package ops

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/monstercameron/SchemaFlow/internal/config"
	"github.com/monstercameron/SchemaFlow/internal/logger"
	"github.com/monstercameron/SchemaFlow/internal/types"
)

// NegotiateOptions configures the Negotiate operation
type NegotiateOptions struct {
	// Priorities maps constraint names to their importance weights (0.0-1.0)
	Priorities map[string]float64

	// Constraints as natural language requirements
	Constraints []string

	// MinSatisfaction is the minimum acceptable satisfaction score (0.0-1.0)
	MinSatisfaction float64

	// MaxAlternatives limits the number of alternative solutions returned
	MaxAlternatives int

	// Strategy guides the negotiation approach ("balanced", "maximize_primary", "pareto")
	Strategy string

	// Common options
	Steering     string
	Mode         types.Mode
	Intelligence types.Speed
	Context      context.Context
}

// Tradeoff describes what was sacrificed to gain something else
type Tradeoff struct {
	// Sacrificed is what was reduced or given up
	Sacrificed string `json:"sacrificed"`

	// Gained is what was improved or obtained
	Gained string `json:"gained"`

	// Impact describes the magnitude of the tradeoff
	Impact string `json:"impact"`

	// Reasoning explains why this tradeoff was made
	Reasoning string `json:"reasoning,omitempty"`
}

// NegotiateResult contains the negotiated solution and metadata
type NegotiateResult[T any] struct {
	// Solution is the optimal result balancing all constraints
	Solution T `json:"solution"`

	// Satisfaction maps constraint names to how well they were met (0.0-1.0)
	Satisfaction map[string]float64 `json:"satisfaction"`

	// OverallSatisfaction is the weighted average satisfaction
	OverallSatisfaction float64 `json:"overall_satisfaction"`

	// Tradeoffs describes what was sacrificed for what
	Tradeoffs []Tradeoff `json:"tradeoffs,omitempty"`

	// Alternatives are other Pareto-optimal solutions
	Alternatives []T `json:"alternatives,omitempty"`

	// Reasoning explains the negotiation process
	Reasoning string `json:"reasoning,omitempty"`

	// Confidence in the solution quality (0.0-1.0)
	Confidence float64 `json:"confidence"`

	// Metadata contains additional operation information
	Metadata map[string]any `json:"metadata,omitempty"`
}

// Negotiate reconciles competing constraints to find an optimal typed solution.
//
// Type parameter T specifies the output type that balances the constraints.
//
// Examples:
//
//	// Example 1: Project planning with competing constraints
//	type ProjectPlan struct {
//	    Duration int      `json:"duration_weeks"`
//	    Budget   int      `json:"budget"`
//	    Features []string `json:"features"`
//	}
//	constraints := map[string]any{
//	    "max_budget": 100000,
//	    "deadline": "2024-06-01",
//	    "required_features": []string{"auth", "dashboard"},
//	    "desired_features": []string{"analytics", "mobile"},
//	}
//	result, err := Negotiate[ProjectPlan](constraints)
//	fmt.Printf("Plan: %d weeks, $%d\n", result.Solution.Duration, result.Solution.Budget)
//	for _, t := range result.Tradeoffs {
//	    fmt.Printf("Tradeoff: gave up %s for %s\n", t.Sacrificed, t.Gained)
//	}
//
//	// Example 2: Resource allocation with priorities
//	type Allocation struct {
//	    CPUCores  int `json:"cpu_cores"`
//	    MemoryGB  int `json:"memory_gb"`
//	    StorageGB int `json:"storage_gb"`
//	}
//	result, err := Negotiate[Allocation](demands, NegotiateOptions{
//	    Strategy: "balanced",
//	    Priorities: map[string]float64{
//	        "performance": 0.5,
//	        "cost": 0.3,
//	        "reliability": 0.2,
//	    },
//	})
//
//	// Example 3: Salary negotiation
//	type Offer struct {
//	    Salary     int    `json:"salary"`
//	    RemoteDays int    `json:"remote_days"`
//	    Equity     string `json:"equity"`
//	}
//	result, err := Negotiate[Offer](map[string]any{
//	    "candidate_min": 150000,
//	    "company_max": 140000,
//	    "remote_preference": "3-5 days",
//	}, NegotiateOptions{Steering: "Find creative compensation solutions"})
func Negotiate[T any](constraints any, opts ...NegotiateOptions) (NegotiateResult[T], error) {
	log := logger.GetLogger()
	log.Debug("Starting negotiate operation")

	var result NegotiateResult[T]
	result.Satisfaction = make(map[string]float64)
	result.Metadata = make(map[string]any)

	// Apply defaults
	opt := NegotiateOptions{
		MinSatisfaction: 0.6,
		MaxAlternatives: 3,
		Strategy:        "balanced",
		Mode:            types.TransformMode,
		Intelligence:    types.Fast,
	}
	if len(opts) > 0 {
		opt = mergeNegotiateOptions(opt, opts[0])
	}

	// Get context
	ctx := opt.Context
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel := context.WithTimeout(ctx, config.GetTimeout())
	defer cancel()

	// Convert constraints to JSON
	constraintsJSON, err := json.Marshal(constraints)
	if err != nil {
		log.Error("Negotiate operation failed: marshal error", "error", err)
		return result, fmt.Errorf("failed to marshal constraints: %w", err)
	}

	// Get target type schema
	var zero T
	typeSchema := GenerateTypeSchema(reflect.TypeOf(zero))

	// Build priorities description
	prioritiesDesc := ""
	if len(opt.Priorities) > 0 {
		var parts []string
		for name, weight := range opt.Priorities {
			parts = append(parts, fmt.Sprintf("%s (weight: %.2f)", name, weight))
		}
		prioritiesDesc = fmt.Sprintf("\n\nPriorities:\n%s", strings.Join(parts, "\n"))
	}

	// Build constraints description
	constraintsDesc := ""
	if len(opt.Constraints) > 0 {
		constraintsDesc = fmt.Sprintf("\n\nAdditional constraints:\n- %s", strings.Join(opt.Constraints, "\n- "))
	}

	systemPrompt := fmt.Sprintf(`You are a negotiation and optimization expert. Find the best solution that balances competing constraints.

Strategy: %s
Minimum acceptable satisfaction: %.0f%%%s%s

Return a JSON object matching this schema:
{
  "solution": %s,
  "satisfaction": {"constraint_name": 0.0-1.0, ...},
  "overall_satisfaction": 0.0-1.0,
  "tradeoffs": [{"sacrificed": "what was reduced", "gained": "what was improved", "impact": "low/medium/high", "reasoning": "why"}],
  "alternatives": [alternative solutions matching the schema above],
  "reasoning": "explanation of negotiation process",
  "confidence": 0.0-1.0
}

Rules:
- "solution" must be a valid instance of the target schema
- "satisfaction" shows how well each constraint was satisfied (1.0 = fully satisfied)
- "tradeoffs" explains what compromises were made
- "alternatives" provides up to %d Pareto-optimal alternatives
- Balance all constraints according to their priorities`,
		opt.Strategy, opt.MinSatisfaction*100, prioritiesDesc, constraintsDesc, typeSchema, opt.MaxAlternatives)

	steeringNote := ""
	if opt.Steering != "" {
		steeringNote = fmt.Sprintf("\n\nAdditional guidance: %s", opt.Steering)
	}

	userPrompt := fmt.Sprintf(`Find the optimal solution that balances these constraints:

%s%s`, string(constraintsJSON), steeringNote)

	// Build OpOptions for LLM call
	opOpts := types.OpOptions{
		Mode:         opt.Mode,
		Intelligence: opt.Intelligence,
		Context:      ctx,
	}

	response, err := callLLM(ctx, systemPrompt, userPrompt, opOpts)
	if err != nil {
		log.Error("Negotiate operation LLM call failed", "error", err)
		return result, fmt.Errorf("negotiation failed: %w", err)
	}

	// Clean up response
	response = strings.TrimSpace(response)
	if strings.HasPrefix(response, "```json") {
		response = strings.TrimPrefix(response, "```json")
		response = strings.TrimSuffix(response, "```")
		response = strings.TrimSpace(response)
	} else if strings.HasPrefix(response, "```") {
		response = strings.TrimPrefix(response, "```")
		response = strings.TrimSuffix(response, "```")
		response = strings.TrimSpace(response)
	}

	// Parse response
	var parsed struct {
		Solution            json.RawMessage    `json:"solution"`
		Satisfaction        map[string]float64 `json:"satisfaction"`
		OverallSatisfaction float64            `json:"overall_satisfaction"`
		Tradeoffs           []Tradeoff         `json:"tradeoffs"`
		Alternatives        []json.RawMessage  `json:"alternatives"`
		Reasoning           string             `json:"reasoning"`
		Confidence          float64            `json:"confidence"`
	}

	if err := json.Unmarshal([]byte(response), &parsed); err != nil {
		log.Error("Negotiate operation failed: parse error", "error", err, "response", response)
		return result, fmt.Errorf("failed to parse negotiation result: %w", err)
	}

	// Parse solution
	if len(parsed.Solution) > 0 {
		if err := json.Unmarshal(parsed.Solution, &result.Solution); err != nil {
			log.Error("Negotiate operation failed: solution parse error", "error", err)
			return result, fmt.Errorf("failed to parse solution: %w", err)
		}
	}

	// Parse alternatives
	for _, alt := range parsed.Alternatives {
		var alternative T
		if err := json.Unmarshal(alt, &alternative); err == nil {
			result.Alternatives = append(result.Alternatives, alternative)
		}
	}

	result.Satisfaction = parsed.Satisfaction
	result.OverallSatisfaction = parsed.OverallSatisfaction
	result.Tradeoffs = parsed.Tradeoffs
	result.Reasoning = parsed.Reasoning
	result.Confidence = parsed.Confidence

	log.Debug("Negotiate operation succeeded",
		"overallSatisfaction", result.OverallSatisfaction,
		"tradeoffs", len(result.Tradeoffs),
		"alternatives", len(result.Alternatives))

	return result, nil
}

// mergeNegotiateOptions merges user options with defaults
func mergeNegotiateOptions(defaults, user NegotiateOptions) NegotiateOptions {
	if user.Priorities != nil {
		defaults.Priorities = user.Priorities
	}
	if user.Constraints != nil {
		defaults.Constraints = user.Constraints
	}
	if user.MinSatisfaction > 0 {
		defaults.MinSatisfaction = user.MinSatisfaction
	}
	if user.MaxAlternatives > 0 {
		defaults.MaxAlternatives = user.MaxAlternatives
	}
	if user.Strategy != "" {
		defaults.Strategy = user.Strategy
	}
	if user.Steering != "" {
		defaults.Steering = user.Steering
	}
	if user.Mode != 0 {
		defaults.Mode = user.Mode
	}
	if user.Intelligence != 0 {
		defaults.Intelligence = user.Intelligence
	}
	if user.Context != nil {
		defaults.Context = user.Context
	}
	return defaults
}
