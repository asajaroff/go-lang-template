package logger

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/asajaroff/go-lang-template/internal/config"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name   string
		cfg    config.LogConfig
		verify func(t *testing.T, logger *Logger)
	}{
		{
			name: "JSON format with info level",
			cfg: config.LogConfig{
				Level:  "info",
				Format: "json",
			},
			verify: func(t *testing.T, logger *Logger) {
				if logger == nil {
					t.Fatal("expected non-nil logger")
				}
				if logger.Logger == nil {
					t.Fatal("expected non-nil slog.Logger")
				}
			},
		},
		{
			name: "Text format with debug level",
			cfg: config.LogConfig{
				Level:  "debug",
				Format: "text",
			},
			verify: func(t *testing.T, logger *Logger) {
				if logger == nil {
					t.Fatal("expected non-nil logger")
				}
			},
		},
		{
			name: "JSON format with error level",
			cfg: config.LogConfig{
				Level:  "error",
				Format: "json",
			},
			verify: func(t *testing.T, logger *Logger) {
				if logger == nil {
					t.Fatal("expected non-nil logger")
				}
			},
		},
		{
			name: "Default to JSON when format is unknown",
			cfg: config.LogConfig{
				Level:  "info",
				Format: "unknown",
			},
			verify: func(t *testing.T, logger *Logger) {
				if logger == nil {
					t.Fatal("expected non-nil logger")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := New(tt.cfg)
			tt.verify(t, logger)
		})
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected slog.Level
	}{
		{"debug", slog.LevelDebug},
		{"DEBUG", slog.LevelDebug},
		{"info", slog.LevelInfo},
		{"INFO", slog.LevelInfo},
		{"warn", slog.LevelWarn},
		{"WARN", slog.LevelWarn},
		{"warning", slog.LevelWarn},
		{"error", slog.LevelError},
		{"ERROR", slog.LevelError},
		{"unknown", slog.LevelInfo}, // defaults to info
		{"", slog.LevelInfo},        // defaults to info
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseLevel(tt.input)
			if result != tt.expected {
				t.Errorf("parseLevel(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestDefault(t *testing.T) {
	logger := Default()
	if logger == nil {
		t.Fatal("expected non-nil logger")
	}
	if logger.Logger == nil {
		t.Fatal("expected non-nil slog.Logger")
	}
}

func TestLoggerOutput(t *testing.T) {
	tests := []struct {
		name         string
		cfg          config.LogConfig
		logFunc      func(*Logger)
		expectStdout bool
		expectStderr bool
		checkContent func(t *testing.T, output string)
	}{
		{
			name: "Info level logs to stdout in JSON format",
			cfg: config.LogConfig{
				Level:  "info",
				Format: "json",
			},
			logFunc: func(l *Logger) {
				l.Info("test info message", "key", "value")
			},
			expectStdout: true,
			expectStderr: false,
			checkContent: func(t *testing.T, output string) {
				if !strings.Contains(output, "test info message") {
					t.Errorf("expected output to contain 'test info message', got: %s", output)
				}
				if !strings.Contains(output, `"level":"INFO"`) {
					t.Errorf("expected JSON level field, got: %s", output)
				}
				// Verify it's valid JSON
				var data map[string]interface{}
				if err := json.Unmarshal([]byte(output), &data); err != nil {
					t.Errorf("expected valid JSON, got error: %v", err)
				}
			},
		},
		{
			name: "Debug level logs to stdout in JSON format",
			cfg: config.LogConfig{
				Level:  "debug",
				Format: "json",
			},
			logFunc: func(l *Logger) {
				l.Debug("test debug message", "debug", true)
			},
			expectStdout: true,
			expectStderr: false,
			checkContent: func(t *testing.T, output string) {
				if !strings.Contains(output, "test debug message") {
					t.Errorf("expected output to contain 'test debug message', got: %s", output)
				}
			},
		},
		{
			name: "Error level logs to stderr in JSON format",
			cfg: config.LogConfig{
				Level:  "info",
				Format: "json",
			},
			logFunc: func(l *Logger) {
				l.Error("test error message", "error", "something went wrong")
			},
			expectStdout: false,
			expectStderr: true,
			checkContent: func(t *testing.T, output string) {
				if !strings.Contains(output, "test error message") {
					t.Errorf("expected output to contain 'test error message', got: %s", output)
				}
				if !strings.Contains(output, `"level":"ERROR"`) {
					t.Errorf("expected JSON error level, got: %s", output)
				}
			},
		},
		{
			name: "Warn level logs to stderr in JSON format",
			cfg: config.LogConfig{
				Level:  "info",
				Format: "json",
			},
			logFunc: func(l *Logger) {
				l.Warn("test warn message", "warning", "be careful")
			},
			expectStdout: false,
			expectStderr: true,
			checkContent: func(t *testing.T, output string) {
				if !strings.Contains(output, "test warn message") {
					t.Errorf("expected output to contain 'test warn message', got: %s", output)
				}
				if !strings.Contains(output, `"level":"WARN"`) {
					t.Errorf("expected JSON warn level, got: %s", output)
				}
			},
		},
		{
			name: "Text format info logs to stdout",
			cfg: config.LogConfig{
				Level:  "info",
				Format: "text",
			},
			logFunc: func(l *Logger) {
				l.Info("text format test", "format", "text")
			},
			expectStdout: true,
			expectStderr: false,
			checkContent: func(t *testing.T, output string) {
				if !strings.Contains(output, "text format test") {
					t.Errorf("expected output to contain 'text format test', got: %s", output)
				}
				if !strings.Contains(output, "level=INFO") {
					t.Errorf("expected text level field, got: %s", output)
				}
			},
		},
		{
			name: "Text format error logs to stderr",
			cfg: config.LogConfig{
				Level:  "info",
				Format: "text",
			},
			logFunc: func(l *Logger) {
				l.Error("text format error", "type", "critical")
			},
			expectStdout: false,
			expectStderr: true,
			checkContent: func(t *testing.T, output string) {
				if !strings.Contains(output, "text format error") {
					t.Errorf("expected output to contain 'text format error', got: %s", output)
				}
				if !strings.Contains(output, "level=ERROR") {
					t.Errorf("expected text error level, got: %s", output)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout and stderr
			oldStdout := os.Stdout
			oldStderr := os.Stderr
			defer func() {
				os.Stdout = oldStdout
				os.Stderr = oldStderr
			}()

			stdoutR, stdoutW, _ := os.Pipe()
			stderrR, stderrW, _ := os.Pipe()
			os.Stdout = stdoutW
			os.Stderr = stderrW

			// Create logger and execute log function
			logger := New(tt.cfg)
			tt.logFunc(logger)

			// Close write ends to flush
			stdoutW.Close()
			stderrW.Close()

			// Read captured output
			stdoutBuf := new(bytes.Buffer)
			stderrBuf := new(bytes.Buffer)
			io.Copy(stdoutBuf, stdoutR)
			io.Copy(stderrBuf, stderrR)

			stdoutOutput := stdoutBuf.String()
			stderrOutput := stderrBuf.String()

			// Verify output routing
			if tt.expectStdout && stdoutOutput == "" {
				t.Error("expected stdout output, got none")
			}
			if !tt.expectStdout && stdoutOutput != "" {
				t.Errorf("expected no stdout output, got: %s", stdoutOutput)
			}
			if tt.expectStderr && stderrOutput == "" {
				t.Error("expected stderr output, got none")
			}
			if !tt.expectStderr && stderrOutput != "" {
				t.Errorf("expected no stderr output, got: %s", stderrOutput)
			}

			// Check content
			if tt.expectStdout {
				tt.checkContent(t, stdoutOutput)
			}
			if tt.expectStderr {
				tt.checkContent(t, stderrOutput)
			}
		})
	}
}

func TestLoggerLevelFiltering(t *testing.T) {
	tests := []struct {
		name         string
		configLevel  string
		logFunc      func(*Logger)
		shouldOutput bool
	}{
		{
			name:        "Info level should not log debug messages",
			configLevel: "info",
			logFunc: func(l *Logger) {
				l.Debug("debug message")
			},
			shouldOutput: false,
		},
		{
			name:        "Info level should log info messages",
			configLevel: "info",
			logFunc: func(l *Logger) {
				l.Info("info message")
			},
			shouldOutput: true,
		},
		{
			name:        "Error level should not log info messages",
			configLevel: "error",
			logFunc: func(l *Logger) {
				l.Info("info message")
			},
			shouldOutput: false,
		},
		{
			name:        "Error level should log error messages",
			configLevel: "error",
			logFunc: func(l *Logger) {
				l.Error("error message")
			},
			shouldOutput: true,
		},
		{
			name:        "Debug level should log all messages",
			configLevel: "debug",
			logFunc: func(l *Logger) {
				l.Debug("debug message")
			},
			shouldOutput: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout and stderr
			oldStdout := os.Stdout
			oldStderr := os.Stderr
			defer func() {
				os.Stdout = oldStdout
				os.Stderr = oldStderr
			}()

			stdoutR, stdoutW, _ := os.Pipe()
			stderrR, stderrW, _ := os.Pipe()
			os.Stdout = stdoutW
			os.Stderr = stderrW

			// Create logger and execute log function
			logger := New(config.LogConfig{
				Level:  tt.configLevel,
				Format: "json",
			})
			tt.logFunc(logger)

			// Close write ends to flush
			stdoutW.Close()
			stderrW.Close()

			// Read captured output
			stdoutBuf := new(bytes.Buffer)
			stderrBuf := new(bytes.Buffer)
			io.Copy(stdoutBuf, stdoutR)
			io.Copy(stderrBuf, stderrR)

			hasOutput := stdoutBuf.Len() > 0 || stderrBuf.Len() > 0

			if tt.shouldOutput && !hasOutput {
				t.Error("expected output, got none")
			}
			if !tt.shouldOutput && hasOutput {
				t.Errorf("expected no output, got stdout: %s, stderr: %s",
					stdoutBuf.String(), stderrBuf.String())
			}
		})
	}
}
