package macevents

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>
#include "browser.h"

*/
import "C"
import (
	"fmt"
	"time"
)

type URLEvent struct {
	URL string
	PID int
}

func WaitForURL(timeout time.Duration) (URLEvent, error) {
	var eventListener chan URLEvent = make(chan URLEvent)

	cancel := time.After(timeout)

	C.RunApp()

	select {
	case e := <-eventListener:
		return e
	case <-cancel:
		return URLEvent{}, fmt.Errorf("failed to get url event, timeout reached")
	}
}

var urlEventChan = make(chan URLEvent)

//export handleURL
func handleURL(u *C.char, i C.int) {
	urlEventChan <- URLEvent{URL: C.GoString(u), PID: int(i)}
}

type AppInfo struct {
	LocalizedName  string
	BundleId       string
	BundlePath     string
	ExecutablePath string
}

func GetRunningAppInfo(pid int) AppInfo {
	return AppInfo{
		LocalizedName:  C.GoString(C.GetLocalizedName(runningApp)),
		BundleId:       C.GoString(C.GetBundleIdentifier(runningApp)),
		BundlePath:     C.GoString(C.GetBundleURL(runningApp)),
		ExecutablePath: C.GoString(C.GetExecutableURL(runningApp)),
	}
}
