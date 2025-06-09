package log

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("color info", func(t *testing.T) {
		err := os.Unsetenv("NO_COLOR")
		require.NoError(t, err)
		opts := slog.HandlerOptions{
			Level:     slog.LevelInfo,
			AddSource: false,
		}
		log, err := New(opts)
		require.NoError(t, err)
		printAll(log)
	})
	t.Run("color debug", func(t *testing.T) {
		err := os.Unsetenv("NO_COLOR")
		require.NoError(t, err)
		opts := slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		}
		log, err := New(opts)
		require.NoError(t, err)
		printAll(log)
	})
	t.Run("no color debug", func(t *testing.T) {
		err := os.Setenv("NO_COLOR", "1")
		require.NoError(t, err)
		opts := slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		}
		log, err := New(opts)
		require.NoError(t, err)
		printAll(log)
	})
	t.Run("no color info", func(t *testing.T) {
		err := os.Setenv("NO_COLOR", "1")
		require.NoError(t, err)
		opts := slog.HandlerOptions{
			Level:     slog.LevelInfo,
			AddSource: false,
		}
		log, err := New(opts)
		require.NoError(t, err)
		printAll(log)
	})
	err := os.Unsetenv("NO_COLOR")
	require.NoError(t, err,
		"error restoring NO_COLOR environment variable to pre-test state",
	)
}

func printAll(log *slog.Logger) {
	log = log.With("service", "cool_service")
	log.Debug("debug message test", "key", "value")
	log.Info("info message test", "key", "value")
	log.Warn("warn message test", "key", "value")
	log.Error("error message test", "key", "value")
}
