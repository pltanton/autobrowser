package config

import (
	"os"
	"strings"
	"testing"
)

func TestTomlParser_ParseBasic(t *testing.T) {
	configStr := `
[variables]
work = "firefox 'ext+container:name=Work&url={}'"
personal = "firefox {}"

[[rules]]
command = "work"
[rules.matchers]
app_class = "Slack"

[[rules]]
command = "personal"
[rules.matchers]
fallback = true
`
	parser, err := NewTomlParser(strings.NewReader(configStr))
	if err != nil {
		t.Fatalf("Failed to parse TOML: %v", err)
	}

	instructions := parser.ConvertToInstructions()

	// Check that we have at least the minimum number of instructions
	if len(instructions) < 3 {
		t.Fatalf("Expected at least 3 instructions, got %d", len(instructions))
	}

	// Check that we have variables
	varCount := 0
	for _, instr := range instructions {
		if _, ok := instr.Assignment(); ok {
			varCount++
		}
	}
	if varCount != 2 {
		t.Errorf("Expected 2 variable assignments, got %d", varCount)
	}

	// Check that we have rules with the right matchers
	appMatcherFound := false
	fallbackFound := false

	for _, instr := range instructions {
		rule, ok := instr.Rule()
		if !ok {
			continue
		}

		// Check for app matcher
		if props, ok := rule.Matchers["app"]; ok {
			if props["class"] == "Slack" {
				appMatcherFound = true
			}
		}

		// Check for fallback
		if _, ok := rule.Matchers["fallback"]; ok {
			fallbackFound = true
		}
	}

	if !appMatcherFound {
		t.Errorf("Expected to find an app matcher with class=Slack")
	}

	if !fallbackFound {
		t.Errorf("Expected to find a fallback matcher")
	}
}

func TestTomlParser_OSSpecificConfigDetection(t *testing.T) {
	// This is a simpler test that just checks if the OS-specific config is detected
	parser := &TomlParser{
		config: TomlConfig{
			Variables: map[string]string{"global": "value"},
			Linux: &OSSpecificConfig{
				Variables: map[string]string{"linux_var": "linux_value"},
			},
			Darwin: &OSSpecificConfig{
				Variables: map[string]string{"darwin_var": "darwin_value"},
			},
		},
	}

	// Save original mocked OS and restore it after the test
	originalOS := currentOS
	defer func() { currentOS = originalOS }()

	// Test Linux detection
	currentOS = "linux"
	osConfig := parser.getOSSpecificConfig()
	if osConfig == nil {
		t.Fatal("Linux config not detected")
	}
	if _, ok := osConfig.Variables["linux_var"]; !ok {
		t.Error("Linux variable not found in config")
	}

	// Test Darwin detection
	currentOS = "darwin"
	osConfig = parser.getOSSpecificConfig()
	if osConfig == nil {
		t.Fatal("Darwin config not detected")
	}
	if _, ok := osConfig.Variables["darwin_var"]; !ok {
		t.Error("Darwin variable not found in config")
	}
}

func TestTomlParser_ParseWithOSSpecific(t *testing.T) {
	configStr := `
[variables]
work = "firefox 'ext+container:name=Work&url={}'"
personal = "firefox {}"

[[rules]]
command = "work"
[rules.matchers]
url_regex = ".*jira.*"

default = "personal"

[linux]
variables = { work = "firefox -p work {}" }

[[linux.rules]]
command = "work"
[linux.rules.matchers]
app_class = "Slack"

[darwin]
variables = { personal = "open -a 'Safari' '{}'" }

[[darwin.rules]]
command = "personal"
[darwin.rules.matchers]
app_bundle_id = "com.apple.safari"

darwin.default = "open -a 'Safari' '{}'"
`
	parser, err := NewTomlParser(strings.NewReader(configStr))
	if err != nil {
		t.Fatalf("Failed to parse TOML: %v", err)
	}

	// Just verify the TOML was parsed correctly
	if len(parser.config.Variables) != 2 {
		t.Errorf("Expected 2 global variables, got %d", len(parser.config.Variables))
	}

	if parser.config.Linux == nil {
		t.Error("Linux config section missing")
	} else if len(parser.config.Linux.Variables) != 1 {
		t.Errorf("Expected 1 Linux variable, got %d", len(parser.config.Linux.Variables))
	}

	if parser.config.Darwin == nil {
		t.Error("Darwin config section missing")
	} else if len(parser.config.Darwin.Variables) != 1 {
		t.Errorf("Expected 1 Darwin variable, got %d", len(parser.config.Darwin.Variables))
	}
}

func TestTomlParser_LoadFromFile(t *testing.T) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "config.*.toml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write test configuration
	configStr := `
[variables]
test = "test {}"

[[rules]]
command = "test"
[rules.matchers]
fallback = true
`
	if _, err := tempFile.WriteString(configStr); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	// Load from file
	parser, err := LoadTomlConfig(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to load TOML from file: %v", err)
	}

	// Just verify the parser loaded the TOML correctly
	if _, ok := parser.config.Variables["test"]; !ok {
		t.Error("Expected 'test' variable but not found")
	}

	if len(parser.config.Rules) != 1 {
		t.Errorf("Expected 1 rule, got %d", len(parser.config.Rules))
	}
}
