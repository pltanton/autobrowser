package mac_opener

import (
	"github.com/pltanton/autobrowser/internal/matchers"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

type macOpenerMatcher struct {
	displayName    string
	bundleId       string
	bundlePath     string
	executablePath string
}

// Match implements matchers.Matcher.
func (h *macOpenerMatcher) Match(args map[string]string) bool {
	if displayName, ok := args["display_name"]; ok && h.displayName != displayName {
		return false
	}
	if bundleId, ok := args["bundle_id"]; ok && h.bundleId != bundleId {
		return false
	}
	if bundlePath, ok := args["bundle_path"]; ok && h.bundlePath != bundlePath {
		return false
	}
	if executablePath, ok := args["executable_path"]; ok && h.executablePath != executablePath {
		return false
	}

	return true
}

var _ matchers.Matcher = &macOpenerMatcher{}

func New(pid int) (matchers.Matcher, error) {
	executablePath, err := os.Executable()
	if err != nil {
		return nil, err
	}

	displayName := fetchInfo(pid, "displayname")
	bundleId := fetchInfo(pid, "bundleid")
	bundlePath := fetchInfo(pid, "bundlepath")

	return &macOpenerMatcher{
		displayName:    displayName,
		bundleId:       bundleId,
		bundlePath:     bundlePath,
		executablePath: executablePath,
	}, nil
}

func fetchInfo(pid int, param string) string {
	info := exec.Command("lsappinfo", "info", "-only", param, "#"+strconv.Itoa(pid))
	output, _ := info.CombinedOutput()
	outputString := string(output)
	if outputString == "" {
		return ""
	}
	regex := regexp.MustCompile("\".+\"=\"(?P<Value>.+)\"")
	match := regex.FindStringSubmatch(outputString)
	return match[1]
}
