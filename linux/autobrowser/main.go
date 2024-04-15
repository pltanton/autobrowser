package main

import (
	"github.com/pltanton/autobrowser/common/pkg/app"
	"github.com/pltanton/autobrowser/common/pkg/args"
	"github.com/pltanton/autobrowser/common/pkg/matchers"
	"github.com/pltanton/autobrowser/common/pkg/matchers/fallback"
	"github.com/pltanton/autobrowser/common/pkg/matchers/url"
	"github.com/pltanton/autobrowser/linux/internal/matchers/hyprland"
)

func main() {
	cfg := args.Parse()

	registry := matchers.NewMatcherRegistry()

	registry.RegisterLazyRule("url", func() (matchers.Matcher, error) {
		return url.New(cfg.Url)
	})

	registry.RegisterLazyRule("app", hyprland.New)
	registry.RegisterLazyRule("fallback", fallback.New)

	app.SetupAndRun(cfg, registry)
}
