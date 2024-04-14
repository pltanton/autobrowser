package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#include "browser.h"
*/
import "C"

import (
	"github.com/pltanton/autobrowser/internal/app"
	"github.com/pltanton/autobrowser/internal/args"
	"github.com/pltanton/autobrowser/internal/matchers"
	"github.com/pltanton/autobrowser/internal/matchers/fallback"
	"github.com/pltanton/autobrowser/internal/matchers/url"
	"os"
	"time"
)

var urlListener chan string = make(chan string)

func main() {
	cfg := args.ParseConfig()
	if cfg == "" {
		panic("Config path is not provided")
	}

	go func() {
		timeout := time.After(4 * time.Second)
		select {
		case urlStr := <-urlListener:
			registry := matchers.NewMatcherRegistry()

			registry.RegisterLazyRule("url", func() (matchers.Matcher, error) {
				return url.New(urlStr)
			})
			registry.RegisterLazyRule("fallback", fallback.New)

			app.SetupAndRun(args.Args{ConfigPath: cfg, Url: urlStr}, registry)
			os.Exit(0)
		case <-timeout:
			os.Exit(1)
		}
	}()

	C.RunApp()
}

//export HandleURL
func HandleURL(u *C.char) {
	urlListener <- C.GoString(u)
}