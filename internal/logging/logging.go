// Package logging configures structured logging.
//
// Factor XI (Logs): logs are an unbuffered event stream written to stdout
// as JSON, one line per event, via the standard library's log/slog. Nothing
// here writes to a file or rotates anything — aggregating, routing, and
// retaining logs is the runtime environment's job (e.g. the cluster's log
// collector), not the app's.
package logging

import (
	"log/slog"
	"os"
	"strings"
)

// New builds a JSON slog.Logger writing to stdout at the given level
// ("debug", "info", "warn", "error" — case-insensitive, defaults to info
// on an unrecognized value).
func New(level string) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: parseLevel(level),
	})
	return slog.New(handler)
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
