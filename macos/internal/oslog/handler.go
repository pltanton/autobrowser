package oslog

import (
	"context"
	"log/slog"
)

type Handler struct {
	osLogger *Logger
}

func NewHandler() *Handler {
	return &Handler{osLogger: NewLogger("dev.pltanton.autobrowser", "AppLog")}
}

func (h *Handler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

func (h *Handler) Handle(_ context.Context, r slog.Record) error {
	// Format log entry
	msg := r.Message
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
	// naive: just log attrs inline
	return h
}

func (h *Handler) WithGroup(name string) slog.Handler {
	// groups not handled in this simple version
	return h
}
