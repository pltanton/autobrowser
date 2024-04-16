package main

import (
	"flag"
	"os"

	"github.com/pltanton/autobrowser/common/pkg/app"
	"github.com/pltanton/autobrowser/common/pkg/matchers"
	"github.com/pltanton/autobrowser/common/pkg/matchers/fallback"
	"github.com/pltanton/autobrowser/common/pkg/matchers/url"
	"github.com/pltanton/autobrowser/linux/internal/matchers/hyprland"
)

type args struct {
	ConfigPath string
	Url        string
}

func parseArgs() args {
	result := args{}
	dir, _ := os.UserHomeDir()
	flag.StringVar(&result.ConfigPath, "config", dir+"/.config/autobrowser.config", "configuration file path")
	flag.StringVar(&result.Url, "url", "", "url to open")

	flag.Parse()

	return result
}

func main() {
	cfg := parseArgs()

	registry := matchers.NewMatcherRegistry()

	registry.RegisterLazyRule("url", func() (matchers.Matcher, error) {
		return url.New(cfg.Url)
	})

	registry.RegisterLazyRule("app", hyprland.New)
	registry.RegisterLazyRule("fallback", fallback.New)

	app.SetupAndRun(cfg.ConfigPath, cfg.Url, registry)
}
