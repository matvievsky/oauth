package log

import (
	"log/slog"
	"os"
	"strings"
)

func init() {
	var level slog.Level
	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "d", "debug":
		level = slog.LevelDebug
	case "i", "info":
		level = slog.LevelInfo
	case "w", "warn", "warning":
		level = slog.LevelWarn
	case "e", "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})))
}
