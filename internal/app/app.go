package app

import (
	"fmt"
	"log"
	"os"

	"github.com/pltanton/autobrowser/internal/args"
	"github.com/pltanton/autobrowser/internal/config"
	"github.com/pltanton/autobrowser/internal/matchers"
)

func SetupAndRun(cfg args.Args, registry *matchers.MatchersRegistry) {
	configFile, err := os.Open(cfg.ConfigPath)
	if err != nil {
		log.Fatalf("Failed to open cofig file: %s", err)
	}

	parser := config.NewParser(configFile)

	for rule, over, err := parser.ParseRule(); !over; rule, over, err = parser.ParseRule() {
		if err != nil {
			log.Fatalf("Failed to parse rule: %s", err)
		}

		matches, err := registry.EvalRule(rule)
		if err != nil {
			log.Fatalf("Failed to evaluate rule: %v", err)
		}

		if matches {
			// TODO: do actual execution here
			fmt.Printf("Great success rule matched! Target: %s\n", rule.Target)
			return
		}
	}

	fmt.Println("No rules matches")
}
