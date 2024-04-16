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

// StartListenNCEvents starts goroutine with Cocoa listener. Unsafe.
func StartListenNCEvents() {
	go C.RunApp()
}

// StartListenNCEvents starts goroutine with Cocoa listener. Unsafe.
func StopListenNCEvents() {
	C.StopApp()
}

// WaitForURL wait for URL event at URL chan or return error
func WaitForURL(timeout time.Duration) (URLEvent, error) {
	cancel := time.After(timeout)

	select {
	case e := <-urlEventChan:
		return e, nil
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
	LocalizedName string
	BundleID      string
	BundleURL     string
	ExecutableURL string
}

func GetRunningAppInfo(pid int) AppInfo {
	appInfo := C.GetById(C.int(pid))
	fmt.Println(appInfo)
	return AppInfo{
		LocalizedName: C.GoString(appInfo.LocalizedName),
		BundleID:      C.GoString(appInfo.BundleID),
		BundleURL:     C.GoString(appInfo.BundleURL),
		ExecutableURL: C.GoString(appInfo.ExecutableURL),
	}
}
