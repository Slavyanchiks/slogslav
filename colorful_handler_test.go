package slogslav

import (
	"bytes"
	"github.com/fatih/color"
	"log/slog"
	"testing"
)

const (
	withColor = "color"
	noColor   = "nocolor"
)

// Let's use for example red
var (
	colorOff = []byte("\033[0m")
	colorRed = []byte("\033[0;31m")
)

type MockColorWriter struct {
	*bytes.Buffer
}

func (m *MockColorWriter) Check() bool {
	if m.Buffer == nil {
		return false
	}
	return bytes.Index(m.Bytes(), colorRed) != -1 || bytes.Index(m.Bytes(), colorOff) != -1 || bytes.Index(m.Bytes(), colorOff) > bytes.Index(m.Bytes(), colorRed)
}

func Test_coloring(t *testing.T) {
	for _, test := range []struct {
		name    string
		log     func(logger *slog.Logger)
		handler string
		result  bool
	}{
		{
			"color applied ansi found",
			func(logger *slog.Logger) {
				logger.Debug("some log")
				logger.Info("some log")
			},
			"color",
			true,
		},
		{
			"color not applied ansi not found",
			func(logger *slog.Logger) {
				logger.Debug("some log")
				logger.Info("some log")
			},
			"nocolor",
			false,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			var mcw = MockColorWriter{&bytes.Buffer{}}
			logColor := slog.New(NewColorfulTextHandler(
				mcw,
				&ColorFormatterPalette{
					err:       color.New(color.BgRed),
					debug:     color.New(color.BgRed),
					info:      color.New(color.BgRed),
					warn:      color.New(color.BgRed),
					time:      color.New(color.BgRed),
					separator: ' '},
				nil))
			logNoColor := slog.New(slog.NewTextHandler(mcw, nil))

			switch test.handler {
			case withColor:
				test.log(logColor)
			case noColor:
				test.log(logNoColor)
			}

			if mcw.Check() != test.result {
				t.Errorf("got %v, want %v", mcw.Check(), test.result)
			}
		})
	}
}
