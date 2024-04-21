package app

import (
	"errors"
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
		slog.Error("Failed to open cofig file: %s", err)
		os.Exit(1)
	}

	parser := config.NewParser(configFile)

	for rule, err := parser.ParseRule(); !errors.Is(err, io.EOF); rule, err = parser.ParseRule() {
		if err != nil {
			slog.Info("Failed to parse rule", "err", err)
			os.Exit(1)
		}

		matches, err := registry.EvalRule(rule)
		if err != nil {
			slog.Info("Failed to evaluate rule", "err", err)
			os.Exit(1)
		}

		if matches {
			// Replace all placeholders in command to url
			command := rule.Command
			for i := range command {
				command[i] = strings.Replace(command[i], "{}", url, 1)
			}

			cmd := exec.Command(command[0], command[1:]...)
			if err := cmd.Run(); err != nil {
				slog.Error("Failed to run command", "err", err)
				os.Exit(1)
			}
			return
		}
	}

	slog.Info("Nothing matched, please specify 'fallback' rule to setup default browser!")
}
