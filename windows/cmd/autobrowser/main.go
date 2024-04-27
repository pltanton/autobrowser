package main

import (
	"flag"
	"github.com/pltanton/autobrowser/common/pkg/app"
	"github.com/pltanton/autobrowser/common/pkg/matchers"
	"github.com/pltanton/autobrowser/common/pkg/matchers/fallback"
	"github.com/pltanton/autobrowser/common/pkg/matchers/urlmatcher"
	"log/slog"
	"os"
)

type config struct {
	url        string
	configPath string
}

func parseConfig() config {
	var result config

	dir, _ := os.UserHomeDir()
	flag.StringVar(&result.configPath, "config", dir+"/.config/autobrowser.config", "configuration file path")
	flag.StringVar(&result.url, "url", "", "url to open")

	flag.Parse()

	return result
}

func main() {
	slog.Debug("Autobrowser launched")
	cfg := parseConfig()
	if cfg.configPath == "" {
		slog.Error("Please provide config by -config parameter")
		os.Exit(1)
	}

	registry := matchers.NewMatcherRegistry()

	registry.RegisterMatcher("url", urlmatcher.New(cfg.url))
	registry.RegisterMatcher("fallback", fallback.New())

	app.SetupAndRun(cfg.configPath, cfg.url, registry)
}
