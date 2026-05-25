package log

import (
	"context"
	"log/slog"
)

type FanoutHandler struct {
	handlers []slog.Handler
}

func (f *FanoutHandler) Enabled(ctx context.Context, l slog.Level) bool {
	for _, h := range f.handlers {
		if h.Enabled(ctx, l) {
			return true
		}
	}
	return false
}

func (f *FanoutHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range f.handlers {
		if h.Enabled(ctx, r.Level) {
			if err := h.Handle(ctx, r); err != nil {
				return err
			}
		}
	}
	return nil
}

func (f *FanoutHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	cloned := make([]slog.Handler, len(f.handlers))
	for i, h := range f.handlers {
		cloned[i] = h.WithAttrs(attrs)
	}
	return &FanoutHandler{handlers: cloned}
}

func (f *FanoutHandler) WithGroup(name string) slog.Handler {
	cloned := make([]slog.Handler, len(f.handlers))
	for i, h := range f.handlers {
		cloned[i] = h.WithGroup(name)
	}
	return &FanoutHandler{handlers: cloned}
}
