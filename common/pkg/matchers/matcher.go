package matchers

import (
	"fmt"
	"log/slog"

	"github.com/pltanton/autobrowser/common/pkg/config"
)

type Matcher interface {
	Match(args map[string]string) bool
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

func (r *MatchersRegistry) EvalRule(rule config.Rule) (bool, error) {
	for name, args := range rule.Matchers {
		matcher, err := r.getMatcher(name)
		if err != nil {
			return false, err
		}

		if !matcher.Match(args) {
			slog.Debug("Matcher not matched", "name", name, "matcher", matcher, "command", rule.Command)
			return false, nil
		}
	}

	slog.Debug("Rule matched", "command", rule)
	return true, nil
}

func (r *MatchersRegistry) getMatcher(name string) (Matcher, error) {
	matcher, ok := r.matchers[name]
	if !ok {
		return nil, fmt.Errorf("unknown matcher %s", name)
	}

	return matcher, nil
}
