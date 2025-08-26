package configuration

import (
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/pltanton/autobrowser/common/pkg/matchers"
)

type Config struct {
	DefaultCommand string             `toml:"default_command"`
	Commands       map[string]Command `toml:"command"`
	Rules          []Rule             `toml:"rules"`

	md toml.MetaData
}

type Command struct {
	CMD          []string       `toml:"-"`
	CMDPrimitive toml.Primitive `toml:"cmd"`
	Placeholder  string         `toml:"placeholder,omitempty"`
	QueryEscape  bool           `toml:"query_escape,omitempty"`
}

type Rule struct {
	Command           string           `toml:"command"`
	MatchersPrimitive []toml.Primitive `toml:"matchers"`
	Matchers          []TypedMatcher   `toml:"-"`
}

type TypedMatcher struct {
	Type      string
	Primitive toml.Primitive
}

type matcherType struct {
	Type string `toml:"type"`
}

func ParseConfigFile(path string) (*Config, error) {
	var config Config
	md, err := toml.DecodeFile(path, &config)
	if err != nil {
		return nil, err
	}

	config.md = md
	if err := parseConfig(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

func ParseConfig(str string) (*Config, error) {
	var config Config
	md, err := toml.Decode(str, &config)
	if err != nil {
		return nil, err
	}

	config.md = md
	if err := parseConfig(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

func parseConfig(config *Config) error {
	for name, command := range config.Commands {
		if command.Placeholder == "" {
			command.Placeholder = "{}"
		}

		var sliceCommand []string
		err := config.md.PrimitiveDecode(command.CMDPrimitive, &sliceCommand)
		if err == nil {
			command.CMD = sliceCommand
			config.Commands[name] = command
			continue
		}

		var stringCommand string
		err = config.md.PrimitiveDecode(command.CMDPrimitive, &stringCommand)
		if err == nil {
			command.CMD = splitQuoted(stringCommand)
			config.Commands[name] = command
			continue
		}

		return fmt.Errorf("Failed to parse command cmd %s, cmd=%v", name, command.CMDPrimitive)
	}

	for i, rule := range config.Rules {
		config.Rules[i].Matchers = make([]TypedMatcher, len(rule.MatchersPrimitive))

		for j, matcher := range rule.MatchersPrimitive {
			var matcherType matcherType
			err := config.md.PrimitiveDecode(matcher, &matcherType)
			if err != nil {
				return fmt.Errorf("Failed to parse matcher type for rule %d, matcher %d", i, j)
			}

			config.Rules[i].Matchers[j].Type = matcherType.Type
			config.Rules[i].Matchers[j].Primitive = matcher
		}
	}

	return nil
}

func splitQuoted(s string) []string {
	var result []string
	var buf strings.Builder
	var quote rune

	flush := func() {
		if buf.Len() > 0 {
			result = append(result, buf.String())
			buf.Reset()
		}
	}

	for _, r := range s {
		switch {
		case (r == '"' || r == '\''):
			if quote == r {
				quote = 0
			} else if quote == 0 {
				quote = r
			} else {
				buf.WriteRune(r)
			}
		case r == ' ' || r == '\t':
			if quote != 0 {
				buf.WriteRune(r)
			} else {
				flush()
			}
		default:
			buf.WriteRune(r)
		}
	}
	flush()
	return result
}

func NewDefaultCommand(cmdString string) Command {
	return Command{
		CMD:         splitQuoted(cmdString),
		Placeholder: "{}",
		QueryEscape: false,
	}
}

func (c *Config) ConfigProvider(matcher TypedMatcher) matchers.MatcherConfigProvider {
	return func(v any) error { return c.md.PrimitiveDecode(matcher.Primitive, v) }
}
