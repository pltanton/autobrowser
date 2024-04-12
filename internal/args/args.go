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
