package compiler

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
	Condition func(map[string]interface{}) bool
}

// Simple rule compiler logic
func CompileSigma(sigma *SigmaRule) (*Rule, error) {
	// In a real implementation, we would parse the 'detection' section
	// and build a nested logical evaluator.
	// For this platform, we implement a mock compiler that handles field matches.

	return &Rule{
		ID:   sigma.ID,
		Name: sigma.Title,
		Condition: func(data map[string]interface{}) bool {
			// Basic match logic for the sake of the platform structure
			for k, v := range sigma.Detection {
				if k == "selection" {
					selection := v.(map[string]interface{})
					for field, val := range selection {
						if data[field] != val {
							return false
						}
					}
					return true
				}
			}
			return false
		},
	}, nil
}
