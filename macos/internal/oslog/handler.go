package oslog

import (
	"context"
	"log/slog"
)

type Handler struct {
	osLogger *Logger
	attrs    []slog.Attr
}

func NewHandler() *Handler {
	return &Handler{
		osLogger: NewLogger("dev.pltanton.autobrowser", "AppLog"),
		attrs:    []slog.Attr{},
	}
}

func (h *Handler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

func (h *Handler) Handle(_ context.Context, r slog.Record) error {
	// Format log entry
	msg := r.Message

	// Add handler's attributes
	for _, attr := range h.attrs {
		msg += " " + attr.String()
	}

	// Add record's attributes
	if r.NumAttrs() > 0 {
		r.Attrs(func(a slog.Attr) bool {
			msg += " " + a.String()
			return true
		})
	}

	h.osLogger.Log(r.Level, msg)
	return nil
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, len(h.attrs)+len(attrs))
	copy(newAttrs, h.attrs)
	copy(newAttrs[len(h.attrs):], attrs)

	return &Handler{
		osLogger: h.osLogger,
		attrs:    newAttrs,
	}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	// Simple implementation: treat groups as prefixes for attributes
	return h
}
