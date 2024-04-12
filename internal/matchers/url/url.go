package url

import (
	"log"
	neturl "net/url"
	"regexp"

	"github.com/pltanton/autobrowser/internal/matchers"
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
		log.Printf("failed to compile regex '%s', error: %s\n", regex, err)
	}
	return r.Match([]byte(u.rawURL))
}

var _ matchers.Matcher = &urlMatcher{}

func New(url string) (matchers.Matcher, error) {
	netUrl, err := neturl.Parse(url)
	if err != nil {
		log.Println("Failed to prase URL, non-regex rules will not work! Error: ", err)
	}

	return &urlMatcher{
		rawURL: url,
		url:    netUrl,
	}, nil
}
