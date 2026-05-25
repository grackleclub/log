package log

import (
	"io"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

const ISO8601 = "2006-01-02T15:04:05.000Z"

// New creates a new slog.Logger with a tint handler that defaults to stderr.
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
func New(opts slog.HandlerOptions, writers ...io.Writer) (*slog.Logger, error) {
	var out io.Writer
	if len(writers) == 0 {
		out = os.Stderr
	} else {
		out = io.MultiWriter(writers...)
	}

	_, envNoColor := os.LookupEnv("NO_COLOR")
	if opts.Level == nil {
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
		tint.NewHandler(out, &tint.Options{
			Level:      opts.Level,
			AddSource:  opts.AddSource,
			NoColor:    envNoColor,
			TimeFormat: ISO8601,
		}),
	)
	logger.Debug(
		"new logger",
		"source", opts.AddSource,
		"level", opts.Level,
		"NO_COLOR", envNoColor,
		"DEBUG", envDebug,
	)
	return logger, nil
}

// NewWithHandlers creates a logger that fans out to the default tint
// handler (stderr) plus any additional slog.Handlers, such as an
// OpenTelemetry log handler. Equivalent to New when no extras are given.
func NewWithHandlers(opts slog.HandlerOptions, extra ...slog.Handler) (*slog.Logger, error) {
	if len(extra) == 0 {
		return New(opts)
	}
	_, envNoColor := os.LookupEnv("NO_COLOR")
	if opts.Level == nil {
		opts.Level = slog.LevelInfo
	}
	_, envDebug := os.LookupEnv("DEBUG")
	if envDebug {
		opts.Level = slog.LevelDebug
		opts.AddSource = true
	}
	tintH := tint.NewHandler(os.Stderr, &tint.Options{
		Level:      opts.Level,
		AddSource:  opts.AddSource,
		NoColor:    envNoColor,
		TimeFormat: ISO8601,
	})
	all := make([]slog.Handler, 0, len(extra)+1)
	all = append(all, tintH)
	all = append(all, extra...)
	logger := slog.New(&fanoutHandler{handlers: all})
	logger.Debug(
		"new logger (fanout)",
		"source", opts.AddSource,
		"level", opts.Level,
		"handlers", len(all),
	)
	return logger, nil
}
