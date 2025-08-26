package urlmatcher

import (
	"fmt"
	"log/slog"
	neturl "net/url"
	"regexp"

	"github.com/pltanton/autobrowser/common/pkg/matchers"
)

type urlMatcher struct {
	rawURL string
	url    *neturl.URL
}

type urlMatcherConfig struct {
	Regex  string `toml:"regex,omitempty"`
	Host   string `toml:"host,omitempty"`
	Scheme string `toml:"scheme,omitempty"`
}

// Match implements matchers.Matcher.
func (u *urlMatcher) Match(configProvider matchers.MatcherConfigProvider) (bool, error) {
	var c urlMatcherConfig
	if err := configProvider(&c); err != nil {
		return false, fmt.Errorf("failed to load url matcher config %w", err)
	}

	if c.Regex != "" && !u.matchByRegex(c.Regex) {
		return false, nil
	}

	if c.Host != "" && !u.matchByHost(c.Host) {
		return false, nil
	}

	if c.Scheme != "" && !u.matchByScheme(c.Scheme) {
		return false, nil
	}

	return true, nil
}

func (u *urlMatcher) matchByHost(host string) bool {
	return u.url.Host == host
}

func (u *urlMatcher) matchByScheme(scheme string) bool {
	return u.url.Scheme == scheme
}

func (u *urlMatcher) matchByRegex(regex string) bool {
	r, err := regexp.Compile(regex)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to compile regex '%s'", regex), "err", err)
	}
	return r.Match([]byte(u.rawURL))
}

var _ matchers.Matcher = &urlMatcher{}

func New(url string) matchers.Matcher {
	netUrl, err := neturl.Parse(url)
	if err != nil {
		slog.Error("Failed to prase URL, non-regex rules will not work!", "err", err)
	}

	return &urlMatcher{
		rawURL: url,
		url:    netUrl,
	}
}
