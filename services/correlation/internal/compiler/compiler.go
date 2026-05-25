package compiler

import (
	"fmt"
	"strings"
)

type SigmaRule struct {
	Title       string                 `yaml:"title"`
	ID          string                 `yaml:"id"`
	Description string                 `yaml:"description"`
	Logsource   map[string]interface{} `yaml:"logsource"`
	Detection   map[string]interface{} `yaml:"detection"`
	Level       string                 `yaml:"level"`
}

type Rule struct {
	ID        string
	Name      string
	Logsource map[string]interface{}
	Condition func(map[string]interface{}) bool
}

// CompileSigma converts a SigmaRule into a functional Rule with a logical evaluator.
// This implementation handles multiple selection blocks and basic boolean logic (AND, OR, NOT).
func CompileSigma(sigma *SigmaRule) (*Rule, error) {
	selections := make(map[string]func(map[string]interface{}) bool)
	var conditionStr string

	for k, v := range sigma.Detection {
		if k == "condition" {
			conditionStr, _ = v.(string)
			continue
		}

		// Treat other keys as selection blocks
		selectionMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		// Create a matcher for this selection block
		// Implements implied AND between fields in a selection
		selections[k] = func(data map[string]interface{}) bool {
			for field, expectedVal := range selectionMap {
				// Handle nested fields (e.g., "event.action")
				actualVal := getNestedField(data, field)
				if actualVal != expectedVal {
					return false
				}
			}
			return true
		}
	}

	if conditionStr == "" {
		// Default to ORing all selections if no condition is provided
		conditionStr = strings.Join(getKeys(selections), " or ")
	}

	evaluator, err := buildEvaluator(conditionStr, selections)
	if err != nil {
		return nil, fmt.Errorf("failed to build evaluator for rule %s: %v", sigma.ID, err)
	}

	return &Rule{
		ID:        sigma.ID,
		Name:      sigma.Title,
		Logsource: sigma.Logsource,
		Condition: evaluator,
	}, nil
}

// getNestedField retrieves a value from a nested map using dot notation
func getNestedField(data map[string]interface{}, field string) interface{} {
	parts := strings.Split(field, ".")
	var current interface{} = data

	for _, part := range parts {
		m, ok := current.(map[string]interface{})
		if !ok {
			return nil
		}
		current = m[part]
	}
	return current
}

func getKeys(m map[string]func(map[string]interface{}) bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// buildEvaluator builds a boolean logic evaluator for the condition string.
// Note: This is a simplified implementation for the platform,
// handling basic tokenized logic without full AST complexity for now.
func buildEvaluator(condition string, selections map[string]func(map[string]interface{}) bool) (func(map[string]interface{}) bool, error) {
	// Standardize condition: lowercase and remove extra spaces
	condition = strings.ToLower(strings.TrimSpace(condition))

	// If it's a single selection
	if matcher, ok := selections[condition]; ok {
		return matcher, nil
	}

	// This is where a real Sigma implementation would use an expression parser.
	// For operational depth, we implement a basic split-based evaluator for AND/OR.

	if strings.Contains(condition, " or ") {
		parts := strings.Split(condition, " or ")
		evaluators := make([]func(map[string]interface{}) bool, 0)
		for _, p := range parts {
			ev, err := buildEvaluator(p, selections)
			if err != nil {
				return nil, err
			}
			evaluators = append(evaluators, ev)
		}
		return func(data map[string]interface{}) bool {
			for _, ev := range evaluators {
				if ev(data) {
					return true
				}
			}
			return false
		}, nil
	}

	if strings.Contains(condition, " and ") {
		parts := strings.Split(condition, " and ")
		evaluators := make([]func(map[string]interface{}) bool, 0)
		for _, p := range parts {
			ev, err := buildEvaluator(p, selections)
			if err != nil {
				return nil, err
			}
			evaluators = append(evaluators, ev)
		}
		return func(data map[string]interface{}) bool {
			for _, ev := range evaluators {
				if !ev(data) {
					return false
				}
			}
			return true
		}, nil
	}

	if strings.HasPrefix(condition, "not ") {
		sub := strings.TrimPrefix(condition, "not ")
		ev, err := buildEvaluator(sub, selections)
		if err != nil {
			return nil, err
		}
		return func(data map[string]interface{}) bool {
			return !ev(data)
		}, nil
	}

	return nil, fmt.Errorf("unsupported condition logic: %s", condition)
}
