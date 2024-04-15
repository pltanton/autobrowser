package mac_opener

import (
	"github.com/pltanton/autobrowser/internal/matchers"
)

type MacOpenerMatcher struct {
	DisplayName    string
	BundleId       string
	BundlePath     string
	ExecutablePath string
}

// Match implements matchers.Matcher.
func (h *MacOpenerMatcher) Match(args map[string]string) bool {
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

var _ matchers.Matcher = &MacOpenerMatcher{}
