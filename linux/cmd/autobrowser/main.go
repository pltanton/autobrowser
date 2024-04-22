package main

import (
	"github.com/pltanton/autobrowser/common/pkg/app"
	"github.com/pltanton/autobrowser/common/pkg/matchers"
	"github.com/pltanton/autobrowser/common/pkg/matchers/fallback"
	"github.com/pltanton/autobrowser/common/pkg/matchers/urlmatcher"
	"github.com/pltanton/autobrowser/common/pkg/utils"
	"github.com/pltanton/autobrowser/linux/internal/deinfo"
	"github.com/pltanton/autobrowser/linux/internal/envx"
	"github.com/pltanton/autobrowser/linux/internal/matchers/appmatcher"
)

func main() {
	options := envx.GetOptions()
	utils.SetLogLevel(options.LogLevel)

	registry := matchers.NewMatcherRegistry()

	// Might be reused to fetch other stuff for other providers
	deInfoProvider := deinfo.New(options.Mode)

	registry.RegisterMatcher("url", urlmatcher.New(options.Url))
	registry.RegisterMatcher("app", appmatcher.New(deInfoProvider))
	registry.RegisterMatcher("fallback", fallback.New())

	app.SetupAndRun(options.ConfigPath, options.Url, registry)
}
