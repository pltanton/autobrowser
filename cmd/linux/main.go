package main

import (
	"github.com/pltanton/autobrowser/internal/app"
	"github.com/pltanton/autobrowser/internal/args"
	"github.com/pltanton/autobrowser/internal/matchers"
	"github.com/pltanton/autobrowser/internal/matchers/url"
)

func main() {
	cfg := args.Parse()

	registry := matchers.NewMatcherRegistry()
	registry.RegisterLazyRule("url", func() (matchers.Matcher, error) {
		return url.New(cfg.Url)
	})

	app.SetupAndRun(cfg, registry)
}
