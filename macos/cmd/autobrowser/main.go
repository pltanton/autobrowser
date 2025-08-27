package main

import (
	"flag"
	"log/slog"
	"os"
	"os/user"
	"time"

	"github.com/pltanton/autobrowser/common/pkg/app"
	"github.com/pltanton/autobrowser/common/pkg/matchers"
	"github.com/pltanton/autobrowser/common/pkg/matchers/urlmatcher"
	"github.com/pltanton/autobrowser/macos/internal/macevents"
	"github.com/pltanton/autobrowser/macos/internal/matchers/appmatcher"
	"github.com/pltanton/autobrowser/macos/internal/oslog"
)

func parseConfig() string {
	var result string

	curUser, err := user.Current()
	if err != nil {
		slog.Error("Failed to get current user", "err", err)
		os.Exit(1)
	}

	flag.StringVar(&result, "config", curUser.HomeDir+"/.config/autobrowser/config.toml", "configuration file path")

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

	app.SetupAndRun(cfg, urlEvent.URL, registry)
}
