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

// Match implements matchers.Matcher.
func (u *urlMatcher) Match(args map[string]string) bool {
	if regex, ok := args["regex"]; ok && !u.matchByRegex(regex) {
		return false
	}

	if host, ok := args["host"]; ok && !u.matchByHost(host) {
		return false
	}

	if scheme, ok := args["scheme"]; ok && !u.matchByScheme(scheme) {
		return false
	}

	return true
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
