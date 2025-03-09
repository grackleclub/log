package log

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var handlerLevel slog.Level

var (
	colorReset    = "\033[0m"
	colorGreyDark = "\033[90m"
)

type SlogHandler struct {
	level slog.Level
	attrs []slog.Attr
	group string
}

func (h *SlogHandler) Handle(ctx context.Context, r slog.Record) error {
	_, noColor := os.LookupEnv("NO_COLOR")

	var color string
	switch r.Level {
	case slog.LevelDebug:
		color = "\033[36m"
	case slog.LevelInfo:
		color = "\033[32m"
	case slog.LevelWarn:
		color = "\033[33m"
	case slog.LevelError:
		color = "\033[31m"
	default:
		color = colorReset
	}
	source := ""
	if handlerLevel == slog.LevelDebug {
		pc := make([]uintptr, 10)
		n := runtime.Callers(4, pc)
		frames := runtime.CallersFrames(pc[:n])
		frame, _ := frames.Next()
		if !noColor {
			source = fmt.Sprintf(" %s%s:%d%s", colorGreyDark, filepath.Base(frame.File), frame.Line, colorReset)
		} else {
			source = fmt.Sprintf(" %s:%d", filepath.Base(frame.File), frame.Line)
		}
	}
	timestamp := r.Time.UTC().Format(time.RFC3339)
	if !noColor {
		timestamp = fmt.Sprintf("%s%s%s", colorGreyDark, timestamp, colorReset)
	}

	level := fmt.Sprintf("%-5s", r.Level.String())
	msg := r.Message
	attrs := ""

	r.Attrs(func(a slog.Attr) bool {
		if !noColor {
			attrs += fmt.Sprintf(" %s%s=%s%v", colorGreyDark, a.Key, colorReset, a.Value)
		} else {
			attrs += fmt.Sprintf(" %s=%v", a.Key, a.Value)
		}
		return true
	})

	var destination *os.File
	if r.Level == slog.LevelError {
		destination = os.Stderr
	} else {
		destination = os.Stdout
	}
	if !noColor {
		fmt.Fprintf(destination, "%s %s%s%s%s %q%s\n",
			timestamp, color, level, colorReset, source, msg, attrs,
		)
	} else {
		fmt.Fprintf(destination, "%s %s %s %q%s\n",
			timestamp, level, source, msg, attrs,
		)

	}
	return nil
}

func (h *SlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &SlogHandler{
		level: h.level,
		attrs: attrs,
	}
}

// TODO this is no-op
func (h *SlogHandler) WithGroup(name string) slog.Handler {
	return h
}

func (h *SlogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level
}

// Init puts a pretty custom handler on otherwise standard log/slog
func Init(level slog.Level) {
	handlerLevel = level
	slog.SetDefault(slog.New(&SlogHandler{
		level: level,
	}))
}
