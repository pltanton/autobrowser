package args

import "flag"

type Args struct {
	ConfigPath string
	Url        string
}

func Parse() Args {
	result := Args{}

	flag.StringVar(&result.ConfigPath, "config", "", "configuration file path")
	flag.StringVar(&result.Url, "url", "", "url to open")

	flag.Parse()

	return result
}

func ParseConfig() string {
	var result string

	flag.StringVar(&result, "config", "", "configuration file path")

	flag.Parse()

	return result
}
