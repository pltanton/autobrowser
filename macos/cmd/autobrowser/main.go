package main

import (
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/pltanton/autobrowser/common/pkg/app"
	"github.com/pltanton/autobrowser/common/pkg/matchers"
	"github.com/pltanton/autobrowser/common/pkg/matchers/fallback"
	"github.com/pltanton/autobrowser/common/pkg/matchers/urlmatcher"
	"github.com/pltanton/autobrowser/macos/internal/macevents"
	"github.com/pltanton/autobrowser/macos/internal/matchers/appmatcher"
	"github.com/pltanton/autobrowser/macos/internal/oslog"
)

func parseConfig() string {
	var result string

	dir, _ := os.UserHomeDir()
	flag.StringVar(&result, "config", dir+"/.config/autobrowser.config", "configuration file path")

	flag.Parse()

	return result
}

func main() {
	if os.Getenv("TERM") == "" {
		slog.Info("Current runtime considered as non-terminal, redirecting logs to OSLog")
		slog.SetDefault(slog.New(oslog.NewHandler()))
	}

	slog.Debug("Autobrowser launched")
	cfg := parseConfig()
	if cfg == "" {
		slog.Error("Please provide config by -config parameter")
		os.Exit(1)
	}

	urlEvent, err := macevents.WaitForURL(4 * time.Second)
	if err != nil {
		slog.Error("Failed to receive url event", "err", err)
		os.Exit(1)
	}

	registry := matchers.NewMatcherRegistry()

	registry.RegisterMatcher("url", urlmatcher.New(urlEvent.URL))
	registry.RegisterMatcher("app", appmatcher.New(urlEvent.PID))
	registry.RegisterMatcher("fallback", fallback.New())

	app.SetupAndRun(cfg, urlEvent.URL, registry)
}
