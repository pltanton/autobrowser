package config

// RuleMatcher matches url by specific rule
type RuleMatcher interface {
	Match(arg string) bool
}

type MatcherProps map[string]string

// Rule single parsed row of configured rules
type Rule struct {
	// Prepeared rule matcher with parsed argument
	Matchers map[string]MatcherProps

	// Target browser command
	Target string
}
