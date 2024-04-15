package opener

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

NSRunningApplication *GetById(int pid) {
    return [NSRunningApplication runningApplicationWithProcessIdentifier:pid];
}
char* GetLocalizedName(NSRunningApplication* runningApp) {
  return [[runningApp localizedName] UTF8String];
}
char* GetBundleIdentifier(NSRunningApplication* runningApp) {
  return [[runningApp bundleIdentifier] UTF8String];
}
char* GetBundleURL(NSRunningApplication* runningApp) {
  return [[[runningApp bundleURL] absoluteString] UTF8String];
}
char* GetExecutableURL(NSRunningApplication* runningApp) {
  return [[[runningApp executableURL] absoluteString] UTF8String];
}
*/
import "C"

import (
	"github.com/pltanton/autobrowser/common/pkg/matchers"
)

type macOpenerMatcher struct {
	DisplayName    string
	BundleId       string
	BundlePath     string
	ExecutablePath string
}

var _ matchers.Matcher = &macOpenerMatcher{}

// Match implements matchers.Matcher.
func (h *macOpenerMatcher) Match(args map[string]string) bool {
	if displayName, ok := args["display_name"]; ok && h.DisplayName != displayName {
		return false
	}
	if bundleId, ok := args["bundle_id"]; ok && h.BundleId != bundleId {
		return false
	}
	if bundlePath, ok := args["bundle_path"]; ok && h.BundlePath != bundlePath {
		return false
	}
	if executablePath, ok := args["executable_path"]; ok && h.ExecutablePath != executablePath {
		return false
	}

	return true
}

func New(pid int) (matchers.Matcher, error) {
	runningApp := C.GetById(C.int(pid))
	return &macOpenerMatcher{
		DisplayName:    C.GoString(C.GetLocalizedName(runningApp)),
		BundleId:       C.GoString(C.GetBundleIdentifier(runningApp)),
		BundlePath:     C.GoString(C.GetBundleURL(runningApp)),
		ExecutablePath: C.GoString(C.GetExecutableURL(runningApp)),
	}, nil
}
