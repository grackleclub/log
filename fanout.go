package log

import (
	"context"
	"errors"
	"log/slog"
)

// fanoutHandler distributes log records to multiple slog.Handlers.
type fanoutHandler struct {
	handlers []slog.Handler
}

// Enabled reports whether any underlying handler is enabled at the given level.
func (f *fanoutHandler) Enabled(ctx context.Context, l slog.Level) bool {
	for _, h := range f.handlers {
		if h.Enabled(ctx, l) {
			return true
		}
	}
	return false
}

// Handle sends the record to every enabled handler, collecting all errors.
func (f *fanoutHandler) Handle(ctx context.Context, r slog.Record) error {
	var errs []error
	for _, h := range f.handlers {
		if h.Enabled(ctx, r.Level) {
			if err := h.Handle(ctx, r); err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errors.Join(errs...)
}

// WithAttrs returns a new fanoutHandler with attrs appended to each handler.
func (f *fanoutHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	cloned := make([]slog.Handler, len(f.handlers))
	for i, h := range f.handlers {
		cloned[i] = h.WithAttrs(attrs)
	}
	return &fanoutHandler{handlers: cloned}
}

// WithGroup returns a new fanoutHandler with the group applied to each handler.
func (f *fanoutHandler) WithGroup(name string) slog.Handler {
	cloned := make([]slog.Handler, len(f.handlers))
	for i, h := range f.handlers {
		cloned[i] = h.WithGroup(name)
	}
	return &fanoutHandler{handlers: cloned}
}
