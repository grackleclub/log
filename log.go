package log

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

const ISO8601 = "2006-01-02T15:04:05.000Z"

// New creates a new slog.Logger with a tint handler that outputs to stderr.
// If empty slog.HandlerOptions are provided, env var DEBUG is checked.
// If DEBUG is unset:
//   - Level: slog.LevelInfo
//   - AddSource: false.
//
// If DEBUG is set:
//   - Level: slog.LevelDebug
//   - AddSource: true.
//
// NO_COLOR environment variable disables color output.
func New(opts slog.HandlerOptions) (*slog.Logger, error) {
	w := os.Stderr
	_, noColor := os.LookupEnv("NO_COLOR")
	if opts.Level == slog.Level(0) {
		opts.Level = slog.LevelInfo
	}
	if !opts.AddSource {
		opts.AddSource = false
	}
	_, envDebug := os.LookupEnv("DEBUG")
	if envDebug {
		opts.Level = slog.LevelDebug
		opts.AddSource = true
	}
	logger := slog.New(
		tint.NewHandler(w, &tint.Options{
			Level:      opts.Level,
			AddSource:  opts.AddSource,
			NoColor:    noColor,
			TimeFormat: ISO8601,
		}),
	)
	logger.Info(
		"new logger",
		"source", opts.AddSource,
		"level", opts.Level,
		"NO_COLOR", os.Getenv("NO_COLOR"),
		"DEBUG", os.Getenv("DEBUG"),
	)
	return logger, nil
}
