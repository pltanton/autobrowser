package fallback

import "github.com/pltanton/autobrowser/common/pkg/matchers"

type fallbackMatcher struct {
}

// Match implements matchers.Matcher.
func (*fallbackMatcher) Match(args map[string]string) bool {
	return true
}

var _ matchers.Matcher = &fallbackMatcher{}

func New() matchers.Matcher {
	return &fallbackMatcher{}
}
