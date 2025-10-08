package schemaflow

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

func Match(input any, cases ...Case) {
	if len(cases) == 0 {
		return
	}
	
	executed := false
	
	for _, c := range cases {
		if c.action == nil {
			continue
		}
		
		switch cond := c.condition.(type) {
		case string:
			if cond == "_" || cond == "otherwise" || cond == "default" {
				if !executed {
					c.action()
					executed = true
				}
				break
			}
			
			if matchesStringCondition(input, cond) {
				c.action()
				executed = true
				break
			}
			
		case reflect.Type:
			if matchesType(input, cond) {
				c.action()
				executed = true
				break
			}
			
		case error:
			if err, ok := input.(error); ok {
				if reflect.TypeOf(err) == reflect.TypeOf(cond) {
					c.action()
					executed = true
					break
				}
			}
			
		default:
			inputType := reflect.TypeOf(input)
			condType := reflect.TypeOf(cond)
			
			if inputType != nil && condType != nil && inputType == condType {
				c.action()
				executed = true
				break
			}
		}
		
		if executed {
			break
		}
	}
}

func When(condition any, action func()) Case {
	return Case{
		condition: condition,
		action:    action,
	}
}

func Like(template string, action func()) Case {
	return Case{
		condition: template,
		action:    action,
	}
}

func Otherwise(action func()) Case {
	return Case{
		condition: "otherwise",
		action:    action,
	}
}

func matchesStringCondition(input any, condition string) bool {
	if condition == "" {
		return false
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*1000000000)
	defer cancel()
	
	inputStr := fmt.Sprintf("%v", input)
	
	systemPrompt := `You are a pattern matching expert. Determine if the input matches the condition.

Rules:
- Consider semantic meaning
- Be consistent in matching
- Return ONLY "true" or "false"`
	
	userPrompt := fmt.Sprintf("Does this input:\n%s\n\nMatch this condition:\n%s", inputStr, condition)
	
	opt := OpOptions{
		Intelligence: Quick,
		Mode:         TransformMode,
	}
	
	response, err := callLLM(ctx, systemPrompt, userPrompt, opt)
	if err != nil {
		return false
	}
	
	response = strings.ToLower(strings.TrimSpace(response))
	response = strings.Trim(response, "\"'")
	
	return response == "true" || response == "yes"
}

func matchesType(input any, targetType reflect.Type) bool {
	if input == nil {
		return false
	}
	
	inputType := reflect.TypeOf(input)
	
	if inputType == targetType {
		return true
	}
	
	if targetType.Kind() == reflect.Interface {
		return inputType.Implements(targetType)
	}
	
	if inputType.Kind() == reflect.Ptr && targetType.Kind() == reflect.Ptr {
		return inputType.Elem() == targetType.Elem()
	}
	
	return false
}