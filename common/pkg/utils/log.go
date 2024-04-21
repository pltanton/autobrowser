package utils

import (
	"log/slog"
	"strings"
)

func SetLogLevel(levelStr string) {
	var logLevel = slog.LevelInfo
	switch strings.ToUpper(levelStr) {
	case "DEBUG":
		logLevel = slog.LevelDebug
	case "INFO":
		logLevel = slog.LevelInfo
	case "WARN":
		logLevel = slog.LevelWarn
	case "ERROR":
		logLevel = slog.LevelError
	}

	slog.SetLogLoggerLevel(logLevel)
}
