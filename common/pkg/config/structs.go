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

	// Command browser command
	Command []string
}

// Assignment is variable definition for config
type Assignment struct {
	Key   string
	Value []string
}

// Insturtion is isngle line in config
type Instruction struct {
	instruction any
}

func FromRule(r Rule) Instruction {
	return Instruction{r}
}

func FromAssignment(a Assignment) Instruction {
	return Instruction{a}
}

func (i Instruction) Rule() (Rule, bool) {
	rule, ok := i.instruction.(Rule)
	return rule, ok
}

func (i Instruction) Assignment() (Assignment, bool) {
	assignment, ok := i.instruction.(Assignment)
	return assignment, ok
}
