package main

import (
	"flag"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/pltanton/autobrowser/common/pkg/app"
	"github.com/pltanton/autobrowser/common/pkg/matchers"
	"github.com/pltanton/autobrowser/common/pkg/matchers/fallback"
	"github.com/pltanton/autobrowser/common/pkg/matchers/url"
	"github.com/pltanton/autobrowser/macos/internal/macevents"
	"github.com/pltanton/autobrowser/macos/internal/matchers/opener"
)

func parseConfig() string {
	var result string

	dir, _ := os.UserHomeDir()
	flag.StringVar(&result, "config", dir+"/.config/autobrowser.config", "configuration file path")

	flag.Parse()

	return result
}

func main() {
	slog.Debug("Autobrowser launched")
	cfg := parseConfig()
	if cfg == "" {
		slog.Error("Please provide config by -config parameter")
		os.Exit(1)
	}

	urlEvent, err := macevents.WaitForURL(4 * time.Second)
	if err != nil {
		slog.Error("Failed to receive url event: ", err)
		os.Exit(1)
	}

	registry := matchers.NewMatcherRegistry()

	registry.RegisterLazyRule("url", func() (matchers.Matcher, error) {
		return url.New(urlEvent.URL)
	})

	registry.RegisterLazyRule("app", func() (matchers.Matcher, error) {
		return opener.New(urlEvent.PID)
	})

	registry.RegisterLazyRule("fallback", fallback.New)

	app.SetupAndRun(cfg, urlEvent.URL, registry)
}
