package config

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/BurntSushi/toml"
)

// For testing - will be replaced by runtime.GOOS if empty
var currentOS = ""

// TomlConfig represents the structure of the TOML configuration file
type TomlConfig struct {
	// Variables are named browser commands that can be referenced in rules
	Variables map[string]string `toml:"variables"`

	// Rules define the matchers and corresponding browser commands
	Rules []TomlRule `toml:"rules"`

	// Default is the fallback browser command to use if no rules match
	Default string `toml:"default"`

	// Linux-specific configuration
	Linux *OSSpecificConfig `toml:"linux"`

	// Darwin (macOS) specific configuration
	Darwin *OSSpecificConfig `toml:"darwin"`
}

// OSSpecificConfig holds OS-specific configuration
type OSSpecificConfig struct {
	// Variables specific to this OS
	Variables map[string]string `toml:"variables"`

	// Rules specific to this OS
	Rules []TomlRule `toml:"rules"`

	// Default browser for this OS
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

	// Process OS-specific configurations based on current OS
	osConfig := p.getOSSpecificConfig()

	// Add variable assignments (global + OS-specific)
	variables := make(map[string]string)

	// First add global variables
	for k, v := range p.config.Variables {
		variables[k] = v
	}

	// Then add OS-specific variables (overriding globals if needed)
	if osConfig != nil && osConfig.Variables != nil {
		for k, v := range osConfig.Variables {
			variables[k] = v
		}
	}

	// Convert variables to assignments
	for key, value := range variables {
		parts := strings.Fields(value)
		assignment := Assignment{
			Key:   key,
			Value: parts,
		}
		instructions = append(instructions, FromAssignment(assignment))
	}

	// Process rules - first global rules
	rules := p.config.Rules

	// Then append OS-specific rules if any
	if osConfig != nil && osConfig.Rules != nil {
		rules = append(rules, osConfig.Rules...)
	}

	// Convert all rules
	for _, tomlRule := range rules {
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

	// Add default rule if specified - OS-specific default overrides global default
	defaultValue := p.config.Default
	if osConfig != nil && osConfig.Default != "" {
		defaultValue = osConfig.Default
	}

	if defaultValue != "" {
		defaultRule := Rule{
			Command:  strings.Fields(defaultValue),
			Matchers: map[string]MatcherProps{"fallback": {}},
		}
		instructions = append(instructions, FromRule(defaultRule))
	} else if osConfig != nil && osConfig.Default != "" {
		// Ensure OS-specific default is always added when present
		osDefaultRule := Rule{
			Command:  strings.Fields(osConfig.Default),
			Matchers: map[string]MatcherProps{"fallback": {}},
		}
		instructions = append(instructions, FromRule(osDefaultRule))
	}

	return instructions
}

// getOSSpecificConfig returns the OS-specific configuration based on runtime OS
func (p *TomlParser) getOSSpecificConfig() *OSSpecificConfig {
	// Use the mock OS value for testing if set, otherwise use the real runtime.GOOS
	os := currentOS
	if os == "" {
		os = runtime.GOOS
	}

	if os == "darwin" {
		return p.config.Darwin
	} else if os == "linux" {
		return p.config.Linux
	}
	return nil
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
