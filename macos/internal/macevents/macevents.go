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

func start() {
	C.RunApp()
}

var urlChan = make(chan URLEvent)

//export handleURL
func handleURL(u *C.char, i C.int) {
	urlChan <- URLEvent{URL: C.GoString(u), PID: int(i)}
}

// StartListenNCEvents starts goroutine with Cocoa listener. Unsafe.
func stop() {
	C.StopApp()
}

// WaitForURL wait for URL event at URL chan or return error
func WaitForURL(timeout time.Duration) (URLEvent, error) {
	cancel := time.After(timeout)

	var closure URLEvent
	var ok bool

	go func() {
		select {
		case e := <-urlChan:
			ok = true
			closure = e
			stop()
		case <-cancel:
			stop()
		}
	}()

	start()

	if !ok {
		return closure, fmt.Errorf("failed to get event, timeout reached")
	}
	return closure, nil
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
