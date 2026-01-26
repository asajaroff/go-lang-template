package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/asajaroff/go-lang-template/internal/config"
)

// Logger wraps slog.Logger with additional functionality
type Logger struct {
	*slog.Logger
}

// New creates a new logger based on the provided configuration
// It outputs to stdout for info/debug and stderr for warn/error
func New(cfg config.LogConfig) *Logger {
	level := parseLevel(cfg.Level)

	var handler slog.Handler

	switch strings.ToLower(cfg.Format) {
	case "json":
		handler = newJSONHandler(level)
	case "text":
		handler = newTextHandler(level)
	default:
		handler = newJSONHandler(level)
	}

	return &Logger{
		Logger: slog.New(handler),
	}
}

// newJSONHandler creates a JSON handler that routes logs to stdout/stderr
func newJSONHandler(level slog.Level) slog.Handler {
	return slog.NewJSONHandler(newLevelWriter(level), &slog.HandlerOptions{
		Level: level,
	})
}

// newTextHandler creates a text handler that routes logs to stdout/stderr
func newTextHandler(level slog.Level) slog.Handler {
	return slog.NewTextHandler(newLevelWriter(level), &slog.HandlerOptions{
		Level: level,
	})
}

// parseLevel converts a string log level to slog.Level
func parseLevel(levelStr string) slog.Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// levelWriter is an io.Writer that routes logs to stdout or stderr based on level
type levelWriter struct {
	level slog.Level
}

// newLevelWriter creates a new level-based writer
func newLevelWriter(level slog.Level) io.Writer {
	return &levelWriter{level: level}
}

// Write implements io.Writer
// Routes Info/Debug to stdout and Warn/Error to stderr
// This allows for proper stream separation in containerized environments
func (w *levelWriter) Write(p []byte) (n int, err error) {
	// Parse the log level from the message by checking for level indicators
	// We support both JSON format ("level":"ERROR") and text format (level=ERROR)
	// This string-based approach works because slog formats the level field before calling Write
	msg := string(p)

	// Check if this is a warning or error message that should go to stderr
	// All other levels (INFO, DEBUG) go to stdout
	if strings.Contains(msg, `"level":"ERROR"`) ||
		strings.Contains(msg, `"level":"WARN"`) ||
		strings.Contains(msg, "level=ERROR") ||
		strings.Contains(msg, "level=WARN") {
		return os.Stderr.Write(p)
	}

	return os.Stdout.Write(p)
}

// Default creates a default logger with info level and JSON format
func Default() *Logger {
	return New(config.LogConfig{
		Level:  "info",
		Format: "json",
	})
}
