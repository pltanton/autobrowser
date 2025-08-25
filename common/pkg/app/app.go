package app

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/pltanton/autobrowser/common/pkg/config"
	"github.com/pltanton/autobrowser/common/pkg/matchers"
)

func SetupAndRun(configPath string, urlString string, registry *matchers.MatchersRegistry) {
	// Load and parse TOML configuration
	tomlParser, err := config.LoadTomlConfig(configPath)
	if err != nil {
		slog.Error("Failed to load TOML configuration", "err", err)
		os.Exit(1)
	}

	// Convert TOML config to instructions
	instructions := tomlParser.ConvertToInstructions()

	variables := make(map[string][]string)

	// Process all instructions
	for _, instruction := range instructions {
		if assignment, ok := instruction.Assignment(); ok {
			variables[assignment.Key] = assignment.Value
		} else if rule, ok := instruction.Rule(); ok {
			matches, err := registry.EvalRule(rule)
			if err != nil {
				slog.Error("Failed to evaluate rule", "err", err)
				os.Exit(1)
			}

			if matches {
				// Replace all placeholders in command to url
				command := rule.Command
				// Try find command in variables if it's single word command
				if len(command) == 1 {
					if newCommand, ok := variables[command[0]]; ok {
						command = newCommand
					}
				}

				urlEscaped := url.QueryEscape(urlString)

				for i := range command {
					command[i] = strings.Replace(command[i], "{}", urlEscaped, 1)
					command[i] = strings.Replace(command[i], "{escape}", urlEscaped, 1)
				}

				slog.Info("Launching CMD", "command", command)

				out, err := exec.Command(command[0], command[1:]...).CombinedOutput()
				if err != nil {
					slog.Error("Failed to run command", "err", err, "output", string(out))
					os.Exit(1)
				}

				slog.Debug("Command executed successfully", "output", string(out))
				return
			}
		} else {
			slog.Error(fmt.Sprintf("Unknown instruction type %+v", instruction))
		}
	}

	slog.Info("Nothing matched, please specify a default rule or add a fallback rule in your configuration!")
}
