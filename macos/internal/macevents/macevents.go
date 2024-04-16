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

var urlChan = make(chan URLEvent)

//export handleURL
func handleURL(u *C.char, i C.int) {
	urlChan <- URLEvent{URL: C.GoString(u), PID: int(i)}
}

// WaitForURL wait for URL event at URL chan or return error
func WaitForURL(timeout time.Duration) (URLEvent, error) {
	cancel := time.After(timeout)

	go C.RunApp()
	defer C.StopApp()

	select {
	case e := <-urlChan:
		return e, nil
	case <-cancel:
		return URLEvent{}, fmt.Errorf("failed to get event, timeout reached")
	}
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
