package logger

import (
	"log/slog"
	"testing"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger()

	if logger == nil {
		t.Errorf("expected logger to not be nil")
	}

	_, ok := logger.Handler().(*slog.JSONHandler)
	if !ok {
		t.Errorf("expected logger handler to be of type *slog.JSONHandler")
	}
}
