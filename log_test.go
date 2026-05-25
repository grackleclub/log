package log

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestNewWriters(t *testing.T) {
	t.Run("single writer", func(t *testing.T) {
		var buf bytes.Buffer
		opts := slog.HandlerOptions{Level: slog.LevelInfo}
		log, err := New(opts, &buf)
		require.NoError(t, err)
		log.Info("hello")
		assert.Contains(t, buf.String(), "hello")
	})
	t.Run("multi writer", func(t *testing.T) {
		var buf1, buf2 bytes.Buffer
		opts := slog.HandlerOptions{Level: slog.LevelInfo}
		log, err := New(opts, &buf1, &buf2)
		require.NoError(t, err)
		log.Info("hello")
		assert.Contains(t, buf1.String(), "hello")
		assert.Contains(t, buf2.String(), "hello")
	})
}

func TestNewWithHandlers(t *testing.T) {
	t.Run("extra handler receives records", func(t *testing.T) {
		h := newCaptureHandler(slog.LevelInfo)
		log, err := NewWithHandlers(slog.HandlerOptions{Level: slog.LevelInfo}, h)
		require.NoError(t, err)
		log.Info("test msg", "k", "v")
		require.Len(t, h.state.records, 1)
		assert.Equal(t, "test msg", h.state.records[0].Message)
	})
	t.Run("level filtering", func(t *testing.T) {
		h := newCaptureHandler(slog.LevelWarn)
		log, err := NewWithHandlers(slog.HandlerOptions{Level: slog.LevelInfo}, h)
		require.NoError(t, err)
		log.Info("below warn")
		assert.Empty(t, h.state.records)
		log.Warn("at warn")
		require.Len(t, h.state.records, 1)
	})
	t.Run("WithAttrs propagation", func(t *testing.T) {
		h := newCaptureHandler(slog.LevelInfo)
		log, err := NewWithHandlers(slog.HandlerOptions{Level: slog.LevelInfo}, h)
		require.NoError(t, err)
		log = log.With("service", "test")
		log.Info("msg")
		require.Len(t, h.state.records, 1)
		assert.True(t, h.state.hadAttrs, "WithAttrs should propagate")
	})
	t.Run("WithGroup propagation", func(t *testing.T) {
		h := newCaptureHandler(slog.LevelInfo)
		log, err := NewWithHandlers(slog.HandlerOptions{Level: slog.LevelInfo}, h)
		require.NoError(t, err)
		log = log.WithGroup("grp")
		log.Info("msg")
		require.Len(t, h.state.records, 1)
		assert.True(t, h.state.hadGroup, "WithGroup should propagate")
	})
	t.Run("error collection", func(t *testing.T) {
		errA := errors.New("handler A failed")
		errB := errors.New("handler B failed")
		hA := newCaptureHandler(slog.LevelInfo)
		hA.err = errA
		hB := newCaptureHandler(slog.LevelInfo)
		hB.err = errB
		f := &fanoutHandler{handlers: []slog.Handler{hA, hB}}
		err := f.Handle(context.Background(), slog.Record{})
		require.Error(t, err)
		assert.ErrorIs(t, err, errA)
		assert.ErrorIs(t, err, errB)
	})
}

type captureState struct {
	mu       sync.Mutex
	records  []slog.Record
	hadAttrs bool
	hadGroup bool
}

type captureHandler struct {
	state *captureState
	level slog.Level
	err   error
}

func newCaptureHandler(level slog.Level) *captureHandler {
	return &captureHandler{
		state: &captureState{},
		level: level,
	}
}

func (c *captureHandler) Enabled(_ context.Context, l slog.Level) bool {
	return l >= c.level
}

func (c *captureHandler) Handle(_ context.Context, r slog.Record) error {
	c.state.mu.Lock()
	defer c.state.mu.Unlock()
	c.state.records = append(c.state.records, r)
	return c.err
}

func (c *captureHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	c.state.mu.Lock()
	c.state.hadAttrs = true
	c.state.mu.Unlock()
	return &captureHandler{state: c.state, level: c.level, err: c.err}
}

func (c *captureHandler) WithGroup(_ string) slog.Handler {
	c.state.mu.Lock()
	c.state.hadGroup = true
	c.state.mu.Unlock()
	return &captureHandler{state: c.state, level: c.level, err: c.err}
}

func printAll(log *slog.Logger) {
	log = log.With("service", "cool_service")
	log.Debug("debug message test", "key", "value")
	log.Info("info message test", "key", "value")
	log.Warn("warn message test", "key", "value")
	log.Error("error message test", "key", "value")
}
