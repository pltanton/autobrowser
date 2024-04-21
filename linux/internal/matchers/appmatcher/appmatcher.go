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

// Match implements matchers.Matcher.
func (m *appMatcher) Match(args map[string]string) bool {
	if class, ok := args["class"]; ok && !m.matchByClass(class) {
		return false
	}

	if title, ok := args["title"]; ok && !m.matchByTitle(title) {
		return false
	}

	return true
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
