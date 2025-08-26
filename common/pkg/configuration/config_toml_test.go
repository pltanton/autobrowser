// Package configuration provides TOML-based configuration parsing
package configuration

import (
	"testing"
)

// TestParseConfig tests the configuration parser with various inputs
func TestParseConfig(t *testing.T) {
	// Test basic config with command
	t.Run("basic config", func(t *testing.T) {
		input := `
default_command = "open"

[command.open]
cmd = "firefox {{url}}"
placeholder = "{{url}}"
`
		config, err := ParseConfig(input)
		if err != nil {
			t.Fatalf("ParseConfig() error = %v", err)
		}

		// Check default command
		if config.DefaultCommand != "open" {
			t.Errorf("DefaultCommand = %q, want %q", config.DefaultCommand, "open")
		}

		// Check commands
		if len(config.Commands) != 1 {
			t.Errorf("Commands count = %d, want 1", len(config.Commands))
			return
		}

		// Check command exists
		if _, exists := config.Commands["open"]; !exists {
			t.Errorf("Command 'open' not found in config")
			return
		}

		// Check placeholder
		cmd := config.Commands["open"]
		if cmd.Placeholder != "{{url}}" {
			t.Errorf("Placeholder = %q, want %q", cmd.Placeholder, "{{url}}")
		}
	})

	// Test config with rules and matchers
	t.Run("config with rules", func(t *testing.T) {
		input := `
[command.test]
cmd = "echo test"

[[rules]]
command = "test"
matchers = [
    { type = "url", pattern = "example.com" },
    { type = "title", pattern = "Example" }
]
`
		config, err := ParseConfig(input)
		if err != nil {
			t.Fatalf("ParseConfig() error = %v", err)
		}

		// Check rules
		if len(config.Rules) != 1 {
			t.Errorf("Rules count = %d, want 1", len(config.Rules))
			return
		}

		// Check rule properties
		rule := config.Rules[0]
		if rule.Command != "test" {
			t.Errorf("Rule command = %q, want %q", rule.Command, "test")
		}

		// Check matchers
		if len(rule.Matchers) != 2 {
			t.Errorf("Matchers count = %d, want 2", len(rule.Matchers))
			return
		}

		// Check matcher types
		if rule.Matchers[0].Type != "url" {
			t.Errorf("First matcher type = %q, want %q", rule.Matchers[0].Type, "url")
		}
		if rule.Matchers[1].Type != "title" {
			t.Errorf("Second matcher type = %q, want %q", rule.Matchers[1].Type, "title")
		}
	})

	// Test invalid configuration
	t.Run("invalid config", func(t *testing.T) {
		input := `
[command.invalid]
cmd = 123  # This should cause an error - cmd must be string or array
`
		c, err := ParseConfig(input)
		if err == nil {
			t.Errorf("ParseConfig() did not return error for invalid command: %v", c.Commands)
		}
	})
}
