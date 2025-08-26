package fallback

import "github.com/pltanton/autobrowser/common/pkg/matchers"

type fallbackMatcher struct {
}

type fallbackMatcherConfig struct {
}

// Match implements matchers.Matcher.
func (*fallbackMatcher) Match(configProvider matchers.MatcherConfigProvider) (bool, error) {
	var c fallbackMatcherConfig
	if err := configProvider(&c); err != nil {
		// Even if config parsing fails, fallback matcher should still return true
		return true, nil
	}
	return true, nil
}

var _ matchers.Matcher = &fallbackMatcher{}

func New() matchers.Matcher {
	return &fallbackMatcher{}
}
