package appmatcher

import (
	"fmt"
	"log/slog"
	"regexp"

	"github.com/pltanton/autobrowser/common/pkg/matchers"
	"github.com/pltanton/autobrowser/linux/internal/deinfo"
)

type appMatcher struct {
	provider *deinfo.DeInfoProvider
}

type appMatcherConfig struct {
	Class string `toml:"class,omitempty"`
	Title string `toml:"title,omitempty"`
}

// Match implements matchers.Matcher.
func (m *appMatcher) Match(configProvider matchers.MatcherConfigProvider) (bool, error) {
	var c appMatcherConfig
	if err := configProvider(&c); err != nil {
		return false, fmt.Errorf("failed to load app matcher config: %w", err)
	}

	if c.Class != "" && !m.matchByClass(c.Class) {
		return false, nil
	}

	if c.Title != "" && !m.matchByTitle(c.Title) {
		return false, nil
	}

	return true, nil
}

func (m *appMatcher) matchByTitle(regex string) bool {
	r, err := regexp.Compile(regex)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to compile regex '%s'", regex), "err", err)
	}

	return r.Match([]byte(m.provider.GetActiveApp().Title))
}

func (h *appMatcher) matchByClass(class string) bool {
	return h.provider.GetActiveApp().Class == class
}

var _ matchers.Matcher = &appMatcher{}

func New(provider *deinfo.DeInfoProvider) matchers.Matcher {
	return &appMatcher{
		provider: provider,
	}
}
