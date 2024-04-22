package envx

import (
	"flag"
	"os"
)

type Options struct {
	LogLevel   string
	ConfigPath string
	Url        string
	Mode       AppMode
}

var options Options

func GetOptions() Options {
	return options
}

func init() {
	flags := struct {
		ConfigPath string
		Url        string

		HyprlandMode bool
		GnomeMode    bool

		LogLevel string
	}{}

	dir, _ := os.UserHomeDir()
	flag.StringVar(&flags.ConfigPath, "config", dir+"/.config/autobrowser.config", "configuration file path")
	flag.StringVar(&flags.Url, "url", "", "url to open")
	flag.StringVar(&flags.LogLevel, "log", "INFO", "log level: DEBUG, INFO, WARN, ERROR")

	flag.BoolVar(&flags.HyprlandMode, "hyprland", false, "use hyprland for app matcher")
	flag.BoolVar(&flags.GnomeMode, "gnome", false, "use gnome for app matcher")

	flag.Parse()

	options = Options{
		ConfigPath: flags.ConfigPath,
		Url:        flags.Url,
		Mode:       getAppMode(flags.HyprlandMode, flags.GnomeMode),
		LogLevel:   flags.LogLevel,
	}
}

type AppMode int

const (
	UNKNOWN AppMode = iota
	HYPRLAND
	GNOME
)

func getAppMode(hyprlandFlag, gnomeFlag bool) AppMode {
	switch {
	case hyprlandFlag:
		return HYPRLAND
	case gnomeFlag:
		return GNOME
	}

	// Try to determine it then
	switch {
	case os.Getenv("HYPRLAND_INSTANCE_SIGNATURE") != "":
		return HYPRLAND
	case os.Getenv("DESKTOP_SESSION") == "gnome":
		return GNOME
	}

	return UNKNOWN
}
