package appmatcher

import (
	"github.com/pltanton/autobrowser/common/pkg/matchers"
	"github.com/pltanton/autobrowser/macos/internal/macevents"
)

type macAppMatcher struct {
	sourceApp macevents.AppInfo
}

var _ matchers.Matcher = &macAppMatcher{}

// Match implements matchers.Matcher.
func (h *macAppMatcher) Match(args map[string]string) bool {
	if displayName, ok := args["display_name"]; ok && h.sourceApp.LocalizedName != displayName {
		return false
	}
	if bundleId, ok := args["bundle_id"]; ok && h.sourceApp.BundleID != bundleId {
		return false
	}
	if bundlePath, ok := args["bundle_path"]; ok && h.sourceApp.BundleURL != bundlePath {
		return false
	}
	if executablePath, ok := args["executable_path"]; ok && h.sourceApp.ExecutableURL != executablePath {
		return false
	}

	return true
}

func New(ppid int) matchers.Matcher {
	return &macAppMatcher{
		sourceApp: macevents.GetRunningAppInfo(ppid),
	}
}
