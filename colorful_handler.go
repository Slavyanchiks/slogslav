package slogslav

import (
	"bytes"
	"context"
	"github.com/fatih/color"
	"io"
	"log/slog"
)

type ColorfulHandlerOptions struct {
	Err       *color.Color
	Warn      *color.Color
	Debug     *color.Color
	Info      *color.Color
	Time      *color.Color
	Separator byte
}

// ColorfulHandler implements slog.Handler
type ColorfulHandler struct {
	out     io.Writer
	opts    *ColorfulHandlerOptions
	handler slog.Handler
}

func (c *ColorfulHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return c.handler.Enabled(ctx, level)
}

func (c *ColorfulHandler) Handle(ctx context.Context, r slog.Record) error {
	buf := bytes.Buffer{}

	if ctx.Err() != nil {
		r.AddAttrs(slog.Attr{Key: "ctx error", Value: slog.StringValue(ctx.Err().Error())})
	}

	if _, err := buf.Write([]byte(c.opts.Time.Sprint("time=" + r.Time.String()))); err != nil {
		return err
	}
	if err := buf.WriteByte(c.opts.Separator); err != nil {
		return err
	}

	level := r.Level
	switch {
	case level < slog.LevelInfo:
		if _, err := buf.Write([]byte("level=" + c.opts.Debug.Sprint("["+r.Level.String()+"]"))); err != nil {
			return err
		}
	case level < slog.LevelWarn:
		if _, err := buf.Write([]byte("level=" + c.opts.Info.Sprint("["+r.Level.String()+"]"))); err != nil {
			return err
		}
	case level < slog.LevelError:
		if _, err := buf.Write([]byte("level=" + c.opts.Warn.Sprint("["+r.Level.String()+"]"))); err != nil {
			return err
		}
	default:
		if _, err := buf.Write([]byte("level=" + c.opts.Err.Sprint("["+r.Level.String()+"]"))); err != nil {
			return err
		}
	}
	if err := buf.WriteByte(c.opts.Separator); err != nil {
		return err
	}

	if _, err := buf.Write([]byte("msg=" + r.Message)); err != nil {
		return err
	}

	if err := buf.WriteByte('\n'); err != nil {
		return err
	}

	_, err := c.out.Write(buf.Bytes())
	return err
}

func (c *ColorfulHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return newColorfulHandler(c.out, c.opts, c.handler.WithAttrs(attrs))
}

func (c *ColorfulHandler) WithGroup(name string) slog.Handler {
	return newColorfulHandler(c.out, c.opts, c.handler.WithGroup(name))
}

func newColorfulHandler(out io.Writer, opts *ColorfulHandlerOptions, handler slog.Handler) *ColorfulHandler {
	if opts == nil {
		opts = &ColorfulHandlerOptions{
			Err:       color.New(color.FgRed),
			Warn:      color.New(color.FgYellow),
			Debug:     color.New(color.FgBlue),
			Info:      color.New(color.FgCyan),
			Time:      color.RGB(112, 128, 144),
			Separator: ' ',
		}
	}
	return &ColorfulHandler{
		out:     out,
		opts:    opts,
		handler: handler,
	}
}

func NewColourfulTextHandler(out io.Writer, opts *ColorfulHandlerOptions, handlerOpts *slog.HandlerOptions) *ColorfulHandler {
	return newColorfulHandler(out, opts, slog.NewTextHandler(out, handlerOpts))
}
