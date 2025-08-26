package appmatcher

import (
	"fmt"

	"github.com/pltanton/autobrowser/common/pkg/matchers"
	"github.com/pltanton/autobrowser/macos/internal/macevents"
)

type macAppMatcher struct {
	sourceApp macevents.AppInfo
}

var _ matchers.Matcher = &macAppMatcher{}

type macAppMatcherConfig struct {
	DisplayName    string `toml:"display_name,omitempty"`
	BundleID       string `toml:"bundle_id,omitempty"`
	BundlePath     string `toml:"bundle_path,omitempty"`
	ExecutablePath string `toml:"executable_path,omitempty"`
}

// Match implements matchers.Matcher.
func (h *macAppMatcher) Match(configProvider matchers.MatcherConfigProvider) (bool, error) {
	var c macAppMatcherConfig
	if err := configProvider(&c); err != nil {
		return false, fmt.Errorf("failed to load mac app matcher config: %w", err)
	}

	if c.DisplayName != "" && h.sourceApp.LocalizedName != c.DisplayName {
		return false, nil
	}
	if c.BundleID != "" && h.sourceApp.BundleID != c.BundleID {
		return false, nil
	}
	if c.BundlePath != "" && h.sourceApp.BundleURL != c.BundlePath {
		return false, nil
	}
	if c.ExecutablePath != "" && h.sourceApp.ExecutableURL != c.ExecutablePath {
		return false, nil
	}

	return true, nil
}

func New(ppid int) matchers.Matcher {
	return &macAppMatcher{
		sourceApp: macevents.GetRunningAppInfo(ppid),
	}
}
