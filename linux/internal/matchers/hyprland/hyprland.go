package hyprland

import (
	"log"
	"os"
	"regexp"

	ipc "github.com/labi-le/hyprland-ipc-client"
	"github.com/pltanton/autobrowser/common/pkg/matchers"
)

type hyprlandMatcher struct {
	title string
	class string
}

// Match implements matchers.Matcher.
func (h *hyprlandMatcher) Match(args map[string]string) bool {
	if class, ok := args["class"]; ok && !h.matchByClass(class) {
		return false
	}

	if title, ok := args["title"]; ok && !h.matchByTitle(title) {
		return false
	}

	return true
}

func (h *hyprlandMatcher) matchByTitle(regex string) bool {
	r, err := regexp.Compile(regex)
	if err != nil {
		log.Printf("failed to compile regex '%s', error: %s\n", regex, err)
	}

	return r.Match([]byte(h.title))
}

func (h *hyprlandMatcher) matchByClass(class string) bool {
	return h.class == class
}

var _ matchers.Matcher = &hyprlandMatcher{}

func New() (matchers.Matcher, error) {
	c := ipc.NewClient(os.Getenv("HYPRLAND_INSTANCE_SIGNATURE"))

	window, err := c.ActiveWindow()
	if err != nil {
		log.Println("Failed to fetch hyprland's window, rulle will not be applied! Error: ", err)
	}

	return &hyprlandMatcher{
		title: window.Title,
		class: window.Class,
	}, nil
}
