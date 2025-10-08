package processor

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/monstercameron/schemaflow"
	"github.com/monstercameron/schemaflow/examples/smarttodo/internal/database"
	"github.com/monstercameron/schemaflow/examples/smarttodo/internal/models"
)

// TodoProcessor handles AI-powered todo processing
type TodoProcessor struct {
	TotalCost    float64 // Track total API costs for session
	LastCost     float64 // Last operation cost
	FastCalls    int     // Number of fast model calls
	SmartCalls   int     // Number of smart model calls
}

// NewTodoProcessor creates a new TodoProcessor instance
func NewTodoProcessor() *TodoProcessor {
	return &TodoProcessor{}
}

// RefineTaskText cleans up and improves task text without full AI processing
func (tp *TodoProcessor) RefineTaskText(taskText string) string {
	// Simple refinement - capitalize, trim, add punctuation if needed
	text := strings.TrimSpace(taskText)
	if text == "" {
		return text
	}
	
	// Capitalize first letter
	runes := []rune(text)
	if len(runes) > 0 {
		runes[0] = unicode.ToUpper(runes[0])
	}
	text = string(runes)
	
	// Remove duplicate spaces
	text = strings.Join(strings.Fields(text), " ")
	
	// Add period if missing and doesn't end with punctuation
	lastChar := text[len(text)-1]
	if lastChar != '.' && lastChar != '!' && lastChar != '?' {
		text += "."
	}
	
	return text
}

// FixTaskGrammar uses AI to fix grammatical errors in task text
func (tp *TodoProcessor) FixTaskGrammar(taskText string) (string, error) {
	// Track cost
	tp.FastCalls++
	tp.LastCost = 0.0001 // Minimal cost for fast model
	tp.TotalCost += tp.LastCost
	
	fixed, err := schemaflow.Transform[string, string](taskText, schemaflow.OpOptions{
		Intelligence: schemaflow.Fast,
		Mode:         schemaflow.TransformMode,
		Steering:     "Fix any grammatical errors and improve clarity. Keep it concise.",
	})
	
	if err != nil {
		// Fallback to simple refinement
		return tp.RefineTaskText(taskText), nil
	}
	
	// Ensure it's not too long
	if len(fixed) > 100 {
		fixed = fixed[:97] + "..."
	}
	
	return fixed, nil
}

func (tp *TodoProcessor) ProcessNote(input string) (*models.SmartTodo, error) {
	// Track cost - estimate for fast model
	tp.FastCalls++
	tp.LastCost = 0.0002 // Approximate cost for GPT-3.5 call
	tp.TotalCost += tp.LastCost
	
	// Try to extract todo from natural language
	rawTodo, err := schemaflow.Extract[models.SmartTodo](input, schemaflow.OpOptions{
		Mode:         schemaflow.TransformMode,
		Intelligence: schemaflow.Quick,
		Steering: `Extract todo information from the input. Infer missing fields intelligently:
		- Title: Create a clear, actionable title (max 50 chars)
		- Description: Expand with helpful details
		- Priority: Determine as high/medium/low based on urgency words
		- Category: Choose from: work, personal, urgent, health, learning, shopping, social, finance
		- Location: Extract location (home, office, gym, store, outside, etc) or default to "home"
		- Tasks: Break down into subtasks if the input mentions multiple steps, bullet points, or a list of things to do
		- Deadline: Parse any date/time mentions (today, tomorrow, next week, etc)
		- Effort: Estimate as minimal/low/medium/high/massive based on complexity
		- Context: Add helpful context or tips for completing the task
		
		Be intelligent about parsing natural language. Examples:
		- "urgent" -> priority: high
		- "quick" -> effort: minimal
		- "buy groceries" -> category: shopping, location: store
		- "workout" -> category: health, location: gym`,
	})
	
	if err != nil {
		// Fallback to basic parsing
		rawTodo = models.SmartTodo{
			Title:       input,
			Description: "",
			Priority:    "medium",
			Category:    "personal",
			Location:    "home",
			Effort:      "medium",
			CreatedAt:   time.Now(),
			Cost:        tp.LastCost,
		}
		
		if len(input) > 50 {
			rawTodo.Title = input[:47] + "..."
			rawTodo.Description = input
		}
	}
	
	// Clean up the title
	if len(rawTodo.Title) > 50 {
		rawTodo.Title = rawTodo.Title[:47] + "..."
	}
	
	// Set location if empty
	if rawTodo.Location == "" {
		rawTodo.Location = "home"
	}
	
	// Ensure valid priority
	if rawTodo.Priority != "high" && rawTodo.Priority != "medium" && rawTodo.Priority != "low" {
		rawTodo.Priority = "medium"
	}
	
	// Ensure valid effort
	validEfforts := map[string]bool{"minimal": true, "low": true, "medium": true, "high": true, "massive": true}
	if !validEfforts[strings.ToLower(rawTodo.Effort)] {
		rawTodo.Effort = "medium"
	}
	
	// Generate context advice if missing
	if rawTodo.Context == "" && rawTodo.Priority == "high" {
		contextPrompt := fmt.Sprintf("Task: %s, Category: %s, Effort: %s", 
			rawTodo.Title, rawTodo.Category, rawTodo.Effort)
		
		context, err := schemaflow.Generate[string](fmt.Sprintf(
			"Provide a brief (20 words max) helpful tip for completing this task: %s",
			contextPrompt,
		), schemaflow.OpOptions{
			Intelligence: schemaflow.Fast,
		})
		
		if err == nil && context != "" {
			rawTodo.Context = context
		}
	}
	
	// Set created time
	rawTodo.CreatedAt = time.Now()
	
	// Track cost
	rawTodo.Cost = tp.LastCost
	
	// Break down into tasks if description is long but no tasks specified
	if len(rawTodo.Tasks) == 0 && len(rawTodo.Description) > 100 {
		// Try to extract subtasks from description
		type TasksResult struct {
			Tasks []models.Task `json:"tasks"`
		}
		
		tasksResult, err := schemaflow.Extract[TasksResult](rawTodo.Description, schemaflow.OpOptions{
			Intelligence: schemaflow.Fast,
			Steering:     "Extract 2-5 subtasks from this description. Each task should be a clear action item.",
		})
		
		if err == nil && len(tasksResult.Tasks) > 0 {
			rawTodo.Tasks = tasksResult.Tasks
		}
	}
	
	return &rawTodo, nil
}

func (tp *TodoProcessor) SuggestNext(todos []*models.SmartTodo) (*models.SmartTodo, error) {
	if len(todos) == 0 {
		return nil, fmt.Errorf("no todos available")
	}
	
	// Single todo, return it
	if len(todos) == 1 {
		return todos[0], nil
	}
	
	// Track cost
	tp.FastCalls++
	tp.LastCost = 0.0002
	tp.TotalCost += tp.LastCost
	
	// Get current time context
	now := time.Now()
	timeContext := fmt.Sprintf("Current time: %s (%s)", 
		now.Format("3:04 PM"), now.Format("Monday"))
	
	best, err := schemaflow.Choose(todos, schemaflow.OpOptions{
		Intelligence: schemaflow.Fast,
		Steering: fmt.Sprintf(`%s
		
		Select the best task to work on right now based on:
		- Priority (high priority first)
		- Deadlines (urgent deadlines first)
		- Current time and day
		- Effort required vs available time
		- Dependencies (tasks blocking others go first)
		- Category fit (work tasks during work hours)
		
		Return the most suitable task for immediate action.`, timeContext),
	})
	
	if err != nil {
		// Fallback: return highest priority or earliest deadline
		var selected *models.SmartTodo
		for _, todo := range todos {
			if selected == nil {
				selected = todo
				continue
			}
			
			// Prefer high priority
			if todo.Priority == "high" && selected.Priority != "high" {
				selected = todo
			} else if todo.Deadline != nil && selected.Deadline != nil {
				// Prefer earlier deadline
				if todo.Deadline.Before(*selected.Deadline) {
					selected = todo
				}
			} else if todo.Deadline != nil && selected.Deadline == nil {
				// Prefer task with deadline over one without
				selected = todo
			}
		}
		
		return selected, nil
	}
	
	// The Choose function returns the selected item
	return best, nil
}

func (tp *TodoProcessor) FilterByContext(todos []*models.SmartTodo, context string) ([]*models.SmartTodo, error) {
	if len(todos) <= 1 {
		return todos, nil
	}
	
	filtered, err := schemaflow.Filter(todos, schemaflow.OpOptions{
		Intelligence: schemaflow.Quick,
		Steering: fmt.Sprintf(`Filter tasks suitable for this context: "%s"
		
		Consider:
		- Location (home, office, gym, etc)
		- Available time
		- Energy level (morning = high energy, evening = low energy)
		- Category relevance
		
		Include tasks that match the current context.`, context),
	})
	
	if err != nil {
		// Return all todos if filtering fails
		return todos, nil
	}
	
	// Return the filtered list directly
	return filtered, nil
}

func (tp *TodoProcessor) GroupByCategory(todos []*models.SmartTodo) map[string][]*models.SmartTodo {
	groups := make(map[string][]*models.SmartTodo)
	
	for _, todo := range todos {
		category := todo.Category
		if category == "" {
			category = "uncategorized"
		}
		groups[category] = append(groups[category], todo)
	}
	
	return groups
}

func (tp *TodoProcessor) SortByPriority(todos []*models.SmartTodo) []*models.SmartTodo {
	// Manual sort since we need specific priority ordering
	high := []*models.SmartTodo{}
	medium := []*models.SmartTodo{}
	low := []*models.SmartTodo{}
	
	for _, todo := range todos {
		switch todo.Priority {
		case "high":
			high = append(high, todo)
		case "medium":
			medium = append(medium, todo)
		default:
			low = append(low, todo)
		}
	}
	
	// Within each priority, sort by deadline
	sortByDeadline := func(todos []*models.SmartTodo) {
		for i := 0; i < len(todos)-1; i++ {
			for j := i + 1; j < len(todos); j++ {
				// Both have deadlines
				if todos[i].Deadline != nil && todos[j].Deadline != nil {
					if todos[i].Deadline.After(*todos[j].Deadline) {
						todos[i], todos[j] = todos[j], todos[i]
					}
				} else if todos[i].Deadline == nil && todos[j].Deadline != nil {
					// Task with deadline comes first
					todos[i], todos[j] = todos[j], todos[i]
				}
			}
		}
	}
	
	sortByDeadline(high)
	sortByDeadline(medium)
	sortByDeadline(low)
	
	// Combine
	result := make([]*models.SmartTodo, 0, len(todos))
	result = append(result, high...)
	result = append(result, medium...)
	result = append(result, low...)
	
	return result
}

// ProcessEditContext merges edit context with existing todo using AI
func (tp *TodoProcessor) ProcessEditContext(todo *models.SmartTodo, context string) (*models.SmartTodo, error) {
	// Track cost
	tp.FastCalls++
	tp.LastCost = 0.0003 // Approximate cost
	tp.TotalCost += tp.LastCost
	
	prompt := fmt.Sprintf(`Current todo:
Title: %s
Description: %s
Priority: %s
Category: %s
Effort: %s
Context: %s

Edit context: %s

Apply the edit context to update the todo. Rules:
- If it mentions location, update location
- If it's a priority change (urgent, asap, etc), update priority
- If it mentions deadline, update deadline
- If it mentions effort/time, update effort
- Keep original information unless explicitly overridden
- Be smart about what the user intends`, 
		todo.Title, todo.Description, todo.Priority, 
		todo.Category, todo.Effort, todo.Context, context)
	
	updatedTodo, err := schemaflow.Transform[string, models.SmartTodo](prompt, schemaflow.OpOptions{
		Intelligence: schemaflow.Quick,
		Mode:         schemaflow.TransformMode,
	})
	
	if err != nil {
		return todo, fmt.Errorf("failed to process edit: %w", err)
	}
	
	// Preserve original ID and timestamps
	updatedTodo.ID = todo.ID
	updatedTodo.CreatedAt = todo.CreatedAt
	updatedTodo.Completed = todo.Completed
	updatedTodo.CompletedAt = todo.CompletedAt
	
	// Update cost
	updatedTodo.Cost = todo.Cost + tp.LastCost
	
	return &updatedTodo, nil
}

func (tp *TodoProcessor) EstimateTimeToComplete(todos []*models.SmartTodo) time.Duration {
	var total time.Duration
	
	for _, todo := range todos {
		switch strings.ToLower(todo.Effort) {
		case "minimal":
			total += 5 * time.Minute
		case "low":
			total += 30 * time.Minute
		case "medium":
			total += 90 * time.Minute
		case "high":
			total += 3 * time.Hour
		case "massive":
			total += 5 * time.Hour
		default:
			total += 1 * time.Hour
		}
	}
	
	return total
}

// SmartPrioritize reorders todos using AI based on multiple factors
func (tp *TodoProcessor) SmartPrioritize(todos []*models.SmartTodo) ([]*models.SmartTodo, error) {
	if len(todos) <= 1 {
		return todos, nil
	}
	
	// Track cost for smart model
	tp.SmartCalls++
	tp.LastCost = 0.001 // Approximate cost for smart model
	tp.TotalCost += tp.LastCost
	
	// Filter only incomplete todos for prioritization
	var incompleteTodos []*models.SmartTodo
	var completedTodos []*models.SmartTodo
	
	for _, todo := range todos {
		if todo.Completed {
			completedTodos = append(completedTodos, todo)
		} else {
			incompleteTodos = append(incompleteTodos, todo)
		}
	}
	
	if len(incompleteTodos) <= 1 {
		return todos, nil
	}
	
	// Get current context
	now := time.Now()
	timeContext := fmt.Sprintf("Current time: %s on %s", 
		now.Format("3:04 PM"), now.Format("Monday, January 2"))
	
	sorted, err := schemaflow.Sort(incompleteTodos, schemaflow.OpOptions{
		Intelligence: schemaflow.Smart,
		Steering: fmt.Sprintf(`Sort these tasks by priority considering:
		
		%s
		
		Factors to consider (in order of importance):
		1. Overdue deadlines (most overdue first)
		2. Priority level (high > medium > low)
		3. Today's deadlines
		4. Dependencies (tasks that block others go first)
		5. Time of day appropriateness (work tasks during work hours)
		6. Effort vs available time
		7. Category clustering (group similar tasks)
		
		Return tasks ordered from most important to least important.`, timeContext),
	})
	
	if err != nil {
		// Fallback to simple priority sort
		return tp.SortByPriority(todos), nil
	}
	
	// Append completed todos at the end
	return append(sorted, completedTodos...), nil
}

// SemanticFilter filters todos using natural language queries
func (tp *TodoProcessor) SemanticFilter(todos []*models.SmartTodo, query string) ([]*models.SmartTodo, error) {
	if len(todos) == 0 || query == "" {
		return todos, nil
	}
	
	// Track cost
	tp.FastCalls++
	tp.LastCost = 0.0002
	tp.TotalCost += tp.LastCost
	
	filtered, err := schemaflow.Filter(todos, schemaflow.OpOptions{
		Intelligence: schemaflow.Fast,
		Steering: fmt.Sprintf(`Filter tasks based on this natural language query: "%s"

Examples of queries to understand:
- "urgent" - high priority tasks or approaching deadlines
- "quick tasks" - minimal or low effort tasks
- "home" - tasks with location=home
- "today" - tasks due today or without deadline
- "work stuff" - work category tasks
- "overdue" - past deadline tasks
- "almost done" - tasks with high completion percentage
- "meeting" or other keywords - search in title and description

Be intelligent about understanding the user's intent.
Include tasks that match the query semantics.`, query),
	})
	
	if err != nil {
		// Fallback to simple substring matching
		var fallbackFiltered []*models.SmartTodo
		queryLower := strings.ToLower(query)
		
		for _, todo := range todos {
			if strings.Contains(strings.ToLower(todo.Title), queryLower) ||
			   strings.Contains(strings.ToLower(todo.Description), queryLower) ||
			   strings.Contains(strings.ToLower(todo.Category), queryLower) ||
			   strings.Contains(strings.ToLower(todo.Priority), queryLower) {
				fallbackFiltered = append(fallbackFiltered, todo)
			}
		}
		
		return fallbackFiltered, nil
	}
	
	return filtered, nil
}

// UpdateDeadlines checks and updates deadline-based priorities
func (tp *TodoProcessor) UpdateDeadlines(db *database.Database) error {
	todos, err := db.GetAllTodos()
	if err != nil {
		return err
	}
	
	now := time.Now()
	updated := false
	
	for _, todo := range todos {
		if todo.Completed || todo.Deadline == nil {
			continue
		}
		
		// Check if deadline is approaching
		daysUntil := int(todo.Deadline.Sub(now).Hours() / 24)
		
		// Auto-escalate priority if deadline is near
		if daysUntil <= 1 && todo.Priority != "high" {
			todo.Priority = "high"
			if err := db.UpdateTodo(todo); err == nil {
				updated = true
			}
		} else if daysUntil <= 3 && todo.Priority == "low" {
			todo.Priority = "medium"
			if err := db.UpdateTodo(todo); err == nil {
				updated = true
			}
		}
	}
	
	if updated {
		return nil
	}
	
	return nil
}