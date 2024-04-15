package matchers

import (
	"fmt"

	"github.com/pltanton/autobrowser/common/pkg/config"
)

type Matcher interface {
	Match(args map[string]string) bool
}

type MatcherConstructor func() (Matcher, error)

type MatchersRegistry struct {
	constructors map[string]MatcherConstructor
	matchers     map[string]Matcher
}

func NewMatcherRegistry() *MatchersRegistry {
	return &MatchersRegistry{
		constructors: map[string]MatcherConstructor{},
		matchers:     map[string]Matcher{},
	}
}

func (r *MatchersRegistry) RegisterLazyRule(name string, constructor MatcherConstructor) {
	r.constructors[name] = constructor
}

func (r *MatchersRegistry) EvalRule(rule config.Rule) (bool, error) {
	for name, args := range rule.Matchers {
		matcher, err := r.getMatcher(name)
		if err != nil {
			return false, err
		}

		if !matcher.Match(args) {
			return false, nil
		}
	}

	return true, nil
}

func (r *MatchersRegistry) getMatcher(name string) (Matcher, error) {
	matcher, ok := r.matchers[name]
	if !ok {
		constructor, ok := r.constructors[name]
		if !ok {
			return nil, fmt.Errorf("unknown matcher %s", name)
		}
		var err error
		matcher, err = constructor()
		if err != nil {
			return nil, fmt.Errorf("failed to construct matcher %s: %w", name, err)
		}

		r.matchers[name] = matcher
	}

	return matcher, nil
}
