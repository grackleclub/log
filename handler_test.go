package log

import (
	"log/slog"
	"os"
	"testing"
)

func TestSlogHandler(t *testing.T) {

	t.Run("debug color", func(t *testing.T) {
		Init(slog.LevelDebug)
		slog.Debug("debug message", "key", "value", "key2", "value2")
		slog.Info("info message", "key", "value", "key2", "value2")
		slog.Warn("warn message", "key", "value", "key2", "value2")
		slog.Error("error message", "key", "value", "key2", "value2")
	})

	t.Run("info color", func(t *testing.T) {
		Init(slog.LevelInfo)
		slog.Debug("debug message", "key", "value", "key2", "value2")
		slog.Info("info message", "key", "value", "key2", "value2")
		slog.Warn("warn message", "key", "value", "key2", "value2")
		slog.Error("error message", "key", "value", "key2", "value2")
	})

	t.Run("debug no color", func(t *testing.T) {
		os.Setenv("NO_COLOR", "true")
		Init(slog.LevelDebug)
		slog.Debug("debug message", "key", "value", "key2", "value2")
		slog.Info("info message", "key", "value", "key2", "value2")
		slog.Warn("warn message", "key", "value", "key2", "value2")
		slog.Error("error message", "key", "value", "key2", "value2")
	})

	t.Run("info no color", func(t *testing.T) {
		os.Setenv("NO_COLOR", "true")
		Init(slog.LevelInfo)
		slog.Debug("debug message", "key", "value", "key2", "value2")
		slog.Info("info message", "key", "value", "key2", "value2")
		slog.Warn("warn message", "key", "value", "key2", "value2")
		slog.Error("error message", "key", "value", "key2", "value2")
	})

}
