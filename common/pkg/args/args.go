package args

import (
	"flag"
	"os"
)

type Args struct {
	ConfigPath string
	Url        string
}

func Parse() Args {
	result := Args{}
	dir, _ := os.UserHomeDir()
	flag.StringVar(&result.ConfigPath, "config", dir+"/.config/autobrowser.config", "configuration file path")
	flag.StringVar(&result.Url, "url", "", "url to open")

	flag.Parse()

	return result
}

func ParseConfig() string {
	var result string

	dir, _ := os.UserHomeDir()
	flag.StringVar(&result, "config", dir+"/.config/autobrowser.config", "configuration file path")

	flag.Parse()

	return result
}
