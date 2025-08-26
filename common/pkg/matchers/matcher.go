package matchers

import (
	"fmt"
)

type MatcherConfigProvider func(v any) error

type Matcher interface {
	Match(configProvider MatcherConfigProvider) (bool, error)
}

type MatchersRegistry struct {
	matchers map[string]Matcher
}

func NewMatcherRegistry() *MatchersRegistry {
	return &MatchersRegistry{
		matchers: map[string]Matcher{},
	}
}

func (r *MatchersRegistry) RegisterMatcher(name string, matcher Matcher) {
	r.matchers[name] = matcher
}

func (r *MatchersRegistry) GetMatcher(name string) (Matcher, error) {
	matcher, ok := r.matchers[name]
	if !ok {
		return nil, fmt.Errorf("unknown matcher %s", name)
	}

	return matcher, nil
}
