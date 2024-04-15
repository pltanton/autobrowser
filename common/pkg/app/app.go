package app

import (
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/pltanton/autobrowser/common/pkg/args"
	"github.com/pltanton/autobrowser/common/pkg/config"
	"github.com/pltanton/autobrowser/common/pkg/matchers"
)

func SetupAndRun(cfg args.Args, registry *matchers.MatchersRegistry) {
	configFile, err := os.Open(cfg.ConfigPath)
	if err != nil {
		log.Fatalf("Failed to open cofig file: %s", err)
	}

	parser := config.NewParser(configFile)

	for rule, err := parser.ParseRule(); !errors.Is(err, io.EOF); rule, err = parser.ParseRule() {
		if err != nil {
			log.Fatalf("Failed to parse rule: %s", err)
		}

		matches, err := registry.EvalRule(rule)
		if err != nil {
			log.Fatalf("Failed to evaluate rule: %v", err)
		}

		if matches {
			// Replace all placeholders in command to url
			command := rule.Command
			for i := range command {
				command[i] = strings.Replace(command[i], "{}", cfg.Url, 1)
			}

			cmd := exec.Command(command[0], command[1:]...)
			if err := cmd.Run(); err != nil {
				log.Fatalln("Failed to run command: ", err)
			}
			return
		}
	}

	log.Println("Nothing matched, please specify 'fallback' rule to setup default browser!")
}
