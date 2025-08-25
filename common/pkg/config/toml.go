package config

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

// TomlConfig represents the structure of the TOML configuration file
type TomlConfig struct {
	// Variables are named browser commands that can be referenced in rules
	Variables map[string]string `toml:"variables"`

	// Rules define the matchers and corresponding browser commands
	Rules []TomlRule `toml:"rules"`

	// Default is the fallback browser command to use if no rules match
	Default string `toml:"default"`
}

// TomlRule defines a single rule in the TOML configuration
type TomlRule struct {
	// Command is either a direct browser command or a reference to a variable
	Command string `toml:"command"`

	// Matchers are the conditions that must be met for this rule to apply
	Matchers TomlMatchers `toml:"matchers"`
}

// TomlMatchers holds all possible matchers
type TomlMatchers struct {
	// App matchers
	AppClass       string `toml:"app_class,omitempty"`
	AppTitle       string `toml:"app_title,omitempty"`
	AppDisplayName string `toml:"app_display_name,omitempty"`
	AppBundleID    string `toml:"app_bundle_id,omitempty"`
	AppBundlePath  string `toml:"app_bundle_path,omitempty"`
	AppExecPath    string `toml:"app_executable_path,omitempty"`

	// URL matchers
	URLHost   string `toml:"url_host,omitempty"`
	URLScheme string `toml:"url_scheme,omitempty"`
	URLRegex  string `toml:"url_regex,omitempty"`

	// Fallback matcher (always matches)
	Fallback bool `toml:"fallback,omitempty"`
}

// TomlParser handles parsing TOML configuration files
type TomlParser struct {
	config TomlConfig
}

// NewTomlParser creates a new TOML configuration parser
func NewTomlParser(r io.Reader) (*TomlParser, error) {
	var config TomlConfig

	// Read all content from the reader
	content, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration: %w", err)
	}

	// Parse TOML
	if err := toml.Unmarshal(content, &config); err != nil {
		return nil, fmt.Errorf("failed to parse TOML configuration: %w", err)
	}

	return &TomlParser{config: config}, nil
}

// ConvertToInstructions converts the TOML configuration to a sequence of Instructions
// compatible with the existing system
func (p *TomlParser) ConvertToInstructions() []Instruction {
	var instructions []Instruction

	// Add variable assignments
	for key, value := range p.config.Variables {
		parts := strings.Fields(value)
		assignment := Assignment{
			Key:   key,
			Value: parts,
		}
		instructions = append(instructions, FromAssignment(assignment))
	}

	// Add rules
	for _, tomlRule := range p.config.Rules {
		rule := Rule{
			Command:  strings.Fields(tomlRule.Command),
			Matchers: make(map[string]MatcherProps),
		}

		// Convert app matchers
		appProps := make(MatcherProps)
		if tomlRule.Matchers.AppClass != "" {
			appProps["class"] = tomlRule.Matchers.AppClass
		}
		if tomlRule.Matchers.AppTitle != "" {
			appProps["title"] = tomlRule.Matchers.AppTitle
		}
		if tomlRule.Matchers.AppDisplayName != "" {
			appProps["display_name"] = tomlRule.Matchers.AppDisplayName
		}
		if tomlRule.Matchers.AppBundleID != "" {
			appProps["bundle_id"] = tomlRule.Matchers.AppBundleID
		}
		if tomlRule.Matchers.AppBundlePath != "" {
			appProps["bundle_path"] = tomlRule.Matchers.AppBundlePath
		}
		if tomlRule.Matchers.AppExecPath != "" {
			appProps["executable_path"] = tomlRule.Matchers.AppExecPath
		}
		if len(appProps) > 0 {
			rule.Matchers["app"] = appProps
		}

		// Convert URL matchers
		urlProps := make(MatcherProps)
		if tomlRule.Matchers.URLHost != "" {
			urlProps["host"] = tomlRule.Matchers.URLHost
		}
		if tomlRule.Matchers.URLScheme != "" {
			urlProps["scheme"] = tomlRule.Matchers.URLScheme
		}
		if tomlRule.Matchers.URLRegex != "" {
			urlProps["regex"] = tomlRule.Matchers.URLRegex
		}
		if len(urlProps) > 0 {
			rule.Matchers["url"] = urlProps
		}

		// Add fallback matcher if specified
		if tomlRule.Matchers.Fallback {
			rule.Matchers["fallback"] = MatcherProps{}
		}

		instructions = append(instructions, FromRule(rule))
	}

	// Add default rule if specified
	if p.config.Default != "" {
		defaultRule := Rule{
			Command:  strings.Fields(p.config.Default),
			Matchers: map[string]MatcherProps{"fallback": {}},
		}
		instructions = append(instructions, FromRule(defaultRule))
	}

	return instructions
}

// LoadTomlConfig loads TOML configuration from a file
func LoadTomlConfig(path string) (*TomlParser, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open configuration file: %w", err)
	}
	defer file.Close()

	return NewTomlParser(file)
}
