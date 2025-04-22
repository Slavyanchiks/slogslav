package slogslav

import (
	"bytes"
	"context"
	"github.com/fatih/color"
	"io"
	"log/slog"
	"time"
)

// ColorFormatterPalette implements ColorFormatter
type ColorFormatterPalette struct {
	err       *color.Color
	warn      *color.Color
	debug     *color.Color
	info      *color.Color
	time      *color.Color
	separator byte
}

type ColorFormatter interface {
	Err(level slog.Level) string
	Warn(level slog.Level) string
	Debug(level slog.Level) string
	Info(level slog.Level) string
	Time(logtime string) string
	Separator() string
}

func (c *ColorFormatterPalette) Err(level slog.Level) string {
	return "level=[" + c.err.Sprint(level.String()) + "]"
}

func (c *ColorFormatterPalette) Warn(level slog.Level) string {
	return "level=[" + c.warn.Sprint(level.String()) + "] "
}

func (c *ColorFormatterPalette) Debug(level slog.Level) string {
	return "level=[" + c.debug.Sprint(level.String()) + "]"
}

func (c *ColorFormatterPalette) Info(level slog.Level) string {
	return "level=[" + c.info.Sprint(level.String()) + "] "
}

func (c *ColorFormatterPalette) Time(logtime string) string {
	return c.time.Sprint("time=" + logtime)
}

func (c *ColorFormatterPalette) Separator() string {
	return string(c.separator)
}

// ColorfulHandler implements slog.Handler
type ColorfulHandler struct {
	out       io.Writer
	formatter ColorFormatter
	handler   slog.Handler
}

func (c *ColorfulHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return c.handler.Enabled(ctx, level)
}

func (c *ColorfulHandler) Handle(ctx context.Context, r slog.Record) error {
	buf := bytes.Buffer{}

	if ctx.Err() != nil {
		r.AddAttrs(slog.Attr{Key: "ctx error", Value: slog.StringValue(ctx.Err().Error())})
	}

	if _, err := buf.Write([]byte(c.formatter.Time(r.Time.Format(time.StampMicro)))); err != nil {
		return err
	}
	if _, err := buf.Write([]byte(c.formatter.Separator())); err != nil {
		return err
	}

	level := r.Level
	switch {
	case level < slog.LevelInfo:
		if _, err := buf.Write([]byte(c.formatter.Debug(level))); err != nil {
			return err
		}
	case level < slog.LevelWarn:
		if _, err := buf.Write([]byte(c.formatter.Info(level))); err != nil {
			return err
		}
	case level < slog.LevelError:
		if _, err := buf.Write([]byte(c.formatter.Warn(level))); err != nil {
			return err
		}
	default:
		if _, err := buf.Write([]byte(c.formatter.Err(level))); err != nil {
			return err
		}
	}
	if _, err := buf.Write([]byte(c.formatter.Separator())); err != nil {
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
	return newColorfulHandler(c.out, c.formatter, c.handler.WithAttrs(attrs))
}

func (c *ColorfulHandler) WithGroup(name string) slog.Handler {
	return newColorfulHandler(c.out, c.formatter, c.handler.WithGroup(name))
}

func newColorfulHandler(out io.Writer, formatter ColorFormatter, handler slog.Handler) *ColorfulHandler {
	if formatter == nil {
		formatter = &ColorFormatterPalette{
			err:       color.New(color.FgRed),
			warn:      color.New(color.FgYellow),
			debug:     color.New(color.FgBlue),
			info:      color.New(color.FgCyan),
			time:      color.RGB(112, 128, 144),
			separator: ' ',
		}
	}
	return &ColorfulHandler{
		out:       out,
		formatter: formatter,
		handler:   handler,
	}
}

func NewColorfulTextHandler(out io.Writer, formatter ColorFormatter, handlerOpts *slog.HandlerOptions) *ColorfulHandler {
	return newColorfulHandler(out, formatter, slog.NewTextHandler(out, handlerOpts))
}
