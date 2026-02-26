package logger

import (
	"log/slog"
	"os"
)

// NewLogger initializes a structured JSON logger using slog.
func NewLogger() *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	return slog.New(handler)
}
