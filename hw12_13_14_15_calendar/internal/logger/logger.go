package logger

import (
	"os"

	"golang.org/x/exp/slog"
)

func Setup(logLevel string, logFile *os.File) *slog.Logger {
	var log *slog.Logger
	var slogLevel slog.Level

	switch logLevel {
	case "debug":
		slogLevel = slog.LevelDebug
	case "info":
		slogLevel = slog.LevelInfo
	case "warn":
		slogLevel = slog.LevelWarn
	case "error":
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelDebug
	}

	log = slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{
		Level: slogLevel,
	}))

	return log
}
