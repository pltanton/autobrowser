package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#include "browser.h"
*/
import "C"
import (
	"os"
	"time"

	"github.com/pltanton/autobrowser/common/pkg/app"
	"github.com/pltanton/autobrowser/common/pkg/args"
	"github.com/pltanton/autobrowser/common/pkg/matchers"
	"github.com/pltanton/autobrowser/common/pkg/matchers/fallback"
	"github.com/pltanton/autobrowser/common/pkg/matchers/url"
	"github.com/pltanton/autobrowser/macos/internal/matchers/opener"
)

type event struct {
	url string
	pid int
}

var eventListener chan event = make(chan event)

func main() {
	cfg := args.ParseConfig()
	if cfg == "" {
		panic("Config path is not provided")
	}

	go func() {
		timeout := time.After(4 * time.Second)
		select {
		case e := <-eventListener:
			urlStr := e.url
			pid := e.pid
			registry := matchers.NewMatcherRegistry()

			registry.RegisterLazyRule("url", func() (matchers.Matcher, error) {
				return url.New(urlStr)
			})
			registry.RegisterLazyRule("app", func() (matchers.Matcher, error) {
				return opener.New(pid)
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
func HandleURL(u *C.char, i C.int) {
	eventListener <- event{url: C.GoString(u), pid: int(i)}
}
