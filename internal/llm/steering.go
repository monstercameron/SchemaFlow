package llm

import (
	"fmt"
	"strings"
)

type SteeringPresets struct{}

var Steering = SteeringPresets{}

func (SteeringPresets) BusinessTone(additionalContext ...string) string {
	base := `Use professional, clear business language:
- Formal but approachable tone
- Action-oriented phrasing
- Avoid jargon and buzzwords
- Focus on outcomes and value`
	return appendContext(base, additionalContext...)
}

func (SteeringPresets) CasualTone(additionalContext ...string) string {
	base := `Use friendly, conversational language:
- Natural, relaxed tone
- Use contractions and common phrases
- Be helpful and encouraging
- Keep it simple and clear`
	return appendContext(base, additionalContext...)
}

func (SteeringPresets) TechnicalTone(additionalContext ...string) string {
	base := `Use precise, technical language:
- Accurate terminology
- Structured and logical flow
- Include relevant details
- Be concise but complete`
	return appendContext(base, additionalContext...)
}

func (SteeringPresets) UrgencyScore(additionalContext ...string) string {
	base := `Rate urgency on scale 0.0-1.0 based on:
- Time sensitivity and deadlines
- Impact of delays
- Dependency blocking potential
- Keywords: urgent, ASAP, critical, deadline`
	return appendContext(base, additionalContext...)
}

func (SteeringPresets) ImportanceScore(additionalContext ...string) string {
	base := `Rate importance on scale 0.0-1.0 based on:
- Strategic value and long-term impact
- Alignment with key goals
- Resource investment required
- Stakeholder priority level`
	return appendContext(base, additionalContext...)
}

func (SteeringPresets) QualityScore(additionalContext ...string) string {
	base := `Rate quality on scale 0.0-1.0 based on:
- Completeness and accuracy
- Clarity and coherence
- Professional standards
- User experience impact`
	return appendContext(base, additionalContext...)
}

func (SteeringPresets) PrioritySort(additionalContext ...string) string {
	base := `Sort by priority considering:
1. High priority items first
2. Urgent deadlines next
3. Quick wins (high value, low effort)
4. Dependencies that block others
5. Strategic alignment`
	return appendContext(base, additionalContext...)
}

func (SteeringPresets) EffortSort(additionalContext ...string) string {
	base := `Sort by effort level considering:
- Time required to complete
- Complexity and skill needed
- Resources and tools required
- Risk and uncertainty factors`
	return appendContext(base, additionalContext...)
}

func (SteeringPresets) DeadlineSort(additionalContext ...string) string {
	base := `Sort by deadline urgency:
- Overdue items first
- Due today/this week
- Due this month
- Future deadlines
- No deadline (lowest priority)`
	return appendContext(base, additionalContext...)
}

func (SteeringPresets) WorkContext(additionalContext ...string) string {
	base := `Filter for work-appropriate tasks:
- Professional environment suitable
- Business hours compatible
- Office tools and resources available
- Collaboration requirements met`
	return appendContext(base, additionalContext...)
}

func (SteeringPresets) HomeContext(additionalContext ...string) string {
	base := `Filter for home environment tasks:
- Personal time suitable
- Home tools and space available
- Family-friendly timing
- Relaxed environment appropriate`
	return appendContext(base, additionalContext...)
}

func (SteeringPresets) MobileContext(additionalContext ...string) string {
	base := `Filter for mobile/on-the-go tasks:
- Can be done on phone/tablet
- No desktop software required
- Short time windows suitable
- Location independent`
	return appendContext(base, additionalContext...)
}

func (SteeringPresets) StrictExtraction(additionalContext ...string) string {
	base := `Extract data with strict validation:
- All required fields must be present
- Format validation enforced
- Type constraints respected
- Fail if ambiguous or incomplete`
	return appendContext(base, additionalContext...)
}

func (SteeringPresets) FlexibleExtraction(additionalContext ...string) string {
	base := `Extract data with flexible interpretation:
- Infer missing fields intelligently
- Use reasonable defaults
- Handle various input formats
- Be forgiving with incomplete data`
	return appendContext(base, additionalContext...)
}

func (SteeringPresets) DetailedExtraction(additionalContext ...string) string {
	base := `Extract comprehensive details:
- Capture all available information
- Infer implicit relationships
- Add contextual metadata
- Preserve original nuances`
	return appendContext(base, additionalContext...)
}

func CustomerServiceTone(specificContext ...string) string {
	return Steering.BusinessTone(
		"Customer service guidelines:",
		"- Empathetic and understanding",
		"- Solution-focused responses",
		"- Acknowledge concerns explicitly",
		"- Provide clear next steps",
		strings.Join(specificContext, "\n"),
	)
}

func UrgentWorkTasks(deadline string, resources string) string {
	return Steering.UrgencyScore("Deadline: "+deadline) +
		"\n\n" + Steering.WorkContext("Available resources: "+resources)
}

func ProjectSpecificSort(projectType, timeline, teamSize string) string {
	return Steering.PrioritySort(
		fmt.Sprintf("Project type: %s", projectType),
		fmt.Sprintf("Timeline: %s", timeline),
		fmt.Sprintf("Team size: %s people", teamSize),
	)
}

func appendContext(base string, additional ...string) string {
	if len(additional) == 0 {
		return base
	}

	result := base
	for _, context := range additional {
		if context != "" {
			result += "\n\nAdditional Context:\n" + context
		}
	}
	return result
}
