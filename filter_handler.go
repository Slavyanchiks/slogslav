package slogslav

import (
	"context"
	"io"
	"log/slog"
)

// FilterHandler implements slog.Handler
type FilterHandler struct {
	level   slog.Level
	handler slog.Handler
}

func (fh *FilterHandler) Handler() slog.Handler {
	return fh.handler
}

func (fh *FilterHandler) Enabled(_ context.Context, level slog.Level) bool {
	return fh.level.Level() == level
}

func (fh *FilterHandler) Handle(ctx context.Context, r slog.Record) error {
	return fh.handler.Handle(ctx, r)
}

func (fh *FilterHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NewFilterHandler(fh.level, fh.handler.WithAttrs(attrs))
}

func (fh *FilterHandler) WithGroup(name string) slog.Handler {
	return NewFilterHandler(fh.level, fh.handler.WithGroup(name))
}

func NewFilterHandler(level slog.Level, handler slog.Handler) *FilterHandler {
	if fh, ok := handler.(*FilterHandler); ok {
		handler = fh.Handler()
	}
	return &FilterHandler{level: level, handler: handler}
}

func NewFilterTextHandler(w io.Writer, opts *slog.HandlerOptions, level slog.Level) *FilterHandler {
	th := slog.NewTextHandler(w, opts)
	return NewFilterHandler(level, th)
}

func NewFilterJSONHandler(w io.Writer, opts *slog.HandlerOptions, level slog.Level) *FilterHandler {
	th := slog.NewJSONHandler(w, opts)
	return NewFilterHandler(level, th)
}
