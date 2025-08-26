package app

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/pltanton/autobrowser/common/pkg/configuration"
	"github.com/pltanton/autobrowser/common/pkg/matchers"
)

func SetupAndRun(configPath string, urlString string, r *matchers.MatchersRegistry) {
	c, err := configuration.ParseConfigFile(configPath)
	if err != nil {
		slog.Error("Failed to parse cofig file", "path", configPath, "err", err)
		os.Exit(1)
	}

	err = evaluate(c, r, urlString)
	if err != nil {
		slog.Error("Failed to evaluate", "err", err)
		os.Exit(1)
	}
}

func evaluate(c *configuration.Config, r *matchers.MatchersRegistry, urlString string) error {
	var command configuration.Command
	var matched bool

	for ruleN, rule := range c.Rules {
		for matcherN, matcherConfig := range rule.Matchers {
			log := slog.With("type", matcherConfig.Type, "rule id", ruleN, "matcher id", matcherN)
			log.Debug("Start matching")

			matcher, err := r.GetMatcher(matcherConfig.Type)
			if err != nil {
				return err
			}

			ok, err := matcher.Match(c.ConfigProvider(matcherConfig))
			if err != nil {
				return err
			}

			log.Debug("Matcher match result", "matched", ok)
			if !ok {
				continue
			}
			matched = true

			command, ok = c.Commands[rule.Command]
			if !ok {
				slog.Debug("Command not declared, using command as is", "command", rule.Command)
				command = configuration.NewDefaultCommand(rule.Command)
			}

			break
		}
	}

	if !matched {
		slog.Debug("None of matchers matched, using default command")
		var ok bool
		command, ok = c.Commands[c.DefaultCommand]
		if !ok {
			slog.Debug("Default command not declared, using command as is", "command", c.DefaultCommand)
			command = configuration.NewDefaultCommand(c.DefaultCommand)
		}
	}

	return runCommand(command, urlString)
}

func runCommand(cmdConfig configuration.Command, urlString string) error {
	cmd := cmdConfig.CMD[:]

	slog.Debug("Command config to execute", "placeholder", cmdConfig.Placeholder)

	if cmdConfig.QueryEscape {
		urlString = url.QueryEscape(urlString)
	}

	for i := range cmd {
		cmd[i] = strings.Replace(cmd[i], cmdConfig.Placeholder, urlString, 1)
	}

	slog.Debug("Launching CMD", "command", cmd)

	out, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
	if err != nil {
		slog.Error("Failed to run command", "err", err, "output", string(out))
		return fmt.Errorf("failed to execute command: %w", err)
	}

	slog.Debug("Command executed successfully", "output", string(out))
	return nil
}
