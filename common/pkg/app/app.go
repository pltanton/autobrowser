package app

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"strings"

	"github.com/pltanton/autobrowser/common/pkg/config"
	"github.com/pltanton/autobrowser/common/pkg/matchers"
)

func SetupAndRun(configPath string, url string, registry *matchers.MatchersRegistry) {
	configFile, err := os.Open(configPath)
	if err != nil {
		slog.Error("Failed to open cofig file", "err", err)
		os.Exit(1)
	}

	parser := config.NewParser(configFile)

	variables := make(map[string][]string)

	for instruction, err := parser.ParseInstruction(); !errors.Is(err, io.EOF); instruction, err = parser.ParseInstruction() {
		if err != nil {
			slog.Error("Failed to parse instruction", "err", err)
			os.Exit(1)
		}

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

				for i := range command {
					command[i] = strings.Replace(command[i], "{}", url, 1)
				}

				cmd := exec.Command(command[0], command[1:]...)
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout

				if err := cmd.Run(); err != nil {
					slog.Error("Failed to run command", "err", err)
					os.Exit(1)
				}
				return
			}
		} else {
			slog.Error(fmt.Sprintf("Unknown instruction type %+v", instruction))
		}
	}

	slog.Info("Nothing matched, please specify 'fallback' rule to setup default browser!")
}
