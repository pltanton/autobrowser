package hyprland

import (
	"log"
	"os"
	"regexp"

	client "github.com/labi-le/hyprland-ipc-client"
	"github.com/pltanton/autobrowser/internal/matchers"
)

type hyprlandMatcher struct {
	window client.Window
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

	return r.Match([]byte(h.window.Title))
}

func (h *hyprlandMatcher) matchByClass(class string) bool {
	return h.window.Class == class
}

var _ matchers.Matcher = &hyprlandMatcher{}

func New() (matchers.Matcher, error) {
	c := client.NewClient(os.Getenv("HYPRLAND_INSTANCE_SIGNATURE"))
	window, err := c.ActiveWindow()
	if err != nil {
		log.Println("Failed to fetch hyprland's window, rulle will not be applied! Error: ", err)
	}

	return &hyprlandMatcher{window: window}, nil
}
