package slogslav

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
)

type MockWriter struct {
	*bytes.Buffer
}

func (m *MockWriter) Check(count int) bool {
	res := 0
	defer m.Reset()

	for _, s := range strings.Split(m.String(), "\n") {
		if strings.Contains(s, "level=DEBUG") {
			res++
		}
	}

	return res == count
}

func Test_filtering(t *testing.T) {
	for _, test := range []struct {
		name   string
		log    func(logger *slog.Logger)
		result int
	}{
		{
			"stated once written once",
			func(logger *slog.Logger) {
				logger.Debug("some log")
			},
			1,
		},
		{
			"stated several written several",
			func(logger *slog.Logger) {
				logger.Debug("some log")
				logger.Debug("some log")
			},
			2,
		},
		{
			"stated once many imposters written once",
			func(logger *slog.Logger) {
				logger.Info("some log")
				logger.Error("some log")
				logger.Warn("some log")
				logger.Debug("some log")
				logger.Info("some log")
			},
			1,
		},
		{
			"stated zero written zero",
			func(logger *slog.Logger) {
				logger.Info("some log")
			},
			0,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			var mwr = MockWriter{&bytes.Buffer{}}
			lg := slog.New(NewFilterTextHandler(mwr, nil, slog.LevelDebug))
			test.log(lg)
			if !mwr.Check(test.result) {
				t.Errorf("Log check failed: expected %d", test.result)
			}
		})
	}
}
